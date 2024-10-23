package rabbitmq

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/krixlion/dev_forum-lib/tracing"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sony/gobreaker"
)

type RabbitMQ struct {
	consumerName string
	shutdown     context.CancelFunc
	config       Config
	url          string // Connection string to RabbitMQ broker.

	conn      *amqp.Connection
	breaker   *gobreaker.TwoStepCircuitBreaker
	connMutex sync.Mutex // Mutex protecting connection during reconnecting.

	notifyConnClose chan *amqp.Error        // Channel to watch for errors from the broker in order to renew the connection.
	publishQueue    chan Message            // Queue for messages waiting to be republished.
	getChannel      chan chan *amqp.Channel // Access channel for accessing the RabbitMQ Channel in a thread-safe way.

	opts options
}

// NewRabbitMQ returns a new initialized connection struct.
// It will manage the active connection in the background.
// Connection should be closed in order to shut it down gracefully.
//
//	func example() {
//		user := "guest"
//		pass := "guest"
//		host := "localhost"
//		port := "5672"
//		consumer := "user-service" //  Unique name for each consumer used to sign messages.
//
//		var customLogger Logger
//
//		// You can specify your own config or use rabbitmq.DefaultConfig() instead.
//		config := Config{
//			QueueSize:         100,
//			MaxWorkers:        50,
//			ReconnectInterval: time.Second * 2,
//			MaxRequests:       5,
//			ClearInterval:     time.Second * 5,
//			ClosedTimeout:     time.Second * 5,
//		}
//
//		// Logger and tracer are optional.
//		rabbit := rabbitmq.NewRabbitMQ(consumer, user, pass, host, port, config, WithLogger(customLogger))
//		defer rabbit.Close()
//	}
func NewRabbitMQ(consumer, user, pass, host, port string, config Config, opts ...Option) *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	ctx, cancel := context.WithCancel(context.Background())

	// Make sure there is at least one worker.
	if config.MaxWorkers == 0 {
		config.MaxWorkers++
	}

	mq := &RabbitMQ{
		consumerName: consumer,
		shutdown:     cancel,

		publishQueue:    make(chan Message, config.QueueSize),
		getChannel:      make(chan chan *amqp.Channel),
		notifyConnClose: make(chan *amqp.Error, 16),

		connMutex: sync.Mutex{},
		url:       url,
		config:    config,
		opts:      defaultOptions(),
		breaker: gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
			Name:        consumer,
			MaxRequests: config.MaxRequests,
			Interval:    config.ClearInterval,
			Timeout:     config.ClosedTimeout,
		}),
	}

	for _, opt := range opts {
		opt.apply(&mq.opts)
	}

	defer mq.run(ctx)
	return mq
}

// run initializes the RabbitMQ connection and manages it in
// separate goroutines while blocking the goroutine it was called from.
// You should use Close() in order to shutdown the connection.
func (mq *RabbitMQ) run(ctx context.Context) {
	mq.opts.logger.Log(ctx, "Connecting to RabbitMQ")
	mq.reDial(ctx)

	go mq.runPublishQueue(ctx)
	go mq.handleConnectionErrors(ctx)
	go mq.handleChannelPropagation(ctx)
}

// Close closes active connection gracefully.
func (mq *RabbitMQ) Close() error {
	mq.shutdown()

	if mq.conn != nil && !mq.conn.IsClosed() {
		return mq.conn.Close()
	}

	return nil
}

// runPublishQueue is meant to be run in a separate goroutine.
func (mq *RabbitMQ) runPublishQueue(ctx context.Context) {
	preparedMessages := mq.prepareExchangePipelined(ctx, mq.publishQueue)
	mq.publishPipelined(ctx, preparedMessages)
}

// handleChannelPropagation is meant to be run in a separate goroutine.
func (mq *RabbitMQ) handleChannelPropagation(ctx context.Context) {
	limiter := make(chan struct{}, mq.config.MaxWorkers)
	for {
		select {
		case req := <-mq.getChannel:
			limiter <- struct{}{}
			go func() {
				ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.handleChannelRead")
				defer span.End()
				defer func() { <-limiter }()
				done, err := mq.breaker.Allow()
				if err != nil {
					req <- nil
					tracing.SetSpanErr(span, err)
					return
				}

				channel, err := mq.conn.Channel()
				if err != nil {
					req <- nil
					mq.opts.logger.Log(ctx, "Failed to open a new channel", "err", err)
					done(false)
					tracing.SetSpanErr(span, err)
					return
				}

				done(true)
				req <- channel
			}()
		case <-ctx.Done():
			return
		}
	}
}

// handleConnectionErrors is meant to be run in a separate goroutine.
func (mq *RabbitMQ) handleConnectionErrors(ctx context.Context) {
	for {
		select {
		case e := <-mq.notifyConnClose:
			if e == nil {
				continue
			}
			mq.reDial(ctx)

		case <-ctx.Done():
			return
		}
	}
}

// reDial will keep reconecting until it succeeds.
func (mq *RabbitMQ) reDial(ctx context.Context) {
	ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.ReDial")
	defer span.End()

	for {
		if err := ctx.Err(); err != nil {
			return
		}

		err := mq.dial(ctx)
		if err == nil {
			return
		}

		tracing.SetSpanErr(span, err)
		mq.opts.logger.Log(ctx, "Failed to connect to RabbitMQ", "err", err)

		time.Sleep(mq.config.ReconnectInterval)
		mq.opts.logger.Log(ctx, "Reconnecting to RabbitMQ")
	}
}

// dial renews current TCP connection.
func (mq *RabbitMQ) dial(ctx context.Context) (err error) {
	_, span := mq.opts.tracer.Start(ctx, "rabbitmq.dial")
	defer span.End()

	done, err := mq.breaker.Allow()
	if err != nil {
		tracing.SetSpanErr(span, err)
		return err
	}

	conn, err := amqp.Dial(mq.url)
	if err != nil {
		tracing.SetSpanErr(span, err)
		done(!isConnectionError(err))
		return err
	}
	done(true)

	mq.connMutex.Lock()
	mq.notifyConnClose = conn.NotifyClose(mq.notifyConnClose)
	mq.conn = conn
	mq.connMutex.Unlock()

	return nil
}

// askForChannel returns a *amqp.Channel in a thread-safe way.
func (mq *RabbitMQ) askForChannel() *amqp.Channel {
	for {
		ask := make(chan *amqp.Channel)
		mq.getChannel <- ask

		channel := <-ask
		if channel != nil {
			return channel
		}

		time.Sleep(mq.config.ReconnectInterval)
	}
}

func isConnectionError(e error) bool {
	err, ok := e.(*amqp.Error)
	if !ok {
		return false
	}

	channelErrCodes := []int{
		amqp.ContentTooLarge,    // 311
		amqp.NoConsumers,        // 313
		amqp.AccessRefused,      // 403
		amqp.NotFound,           // 404
		amqp.ResourceLocked,     // 405
		amqp.PreconditionFailed, // 406
	}

	return !slices.Contains(channelErrCodes, err.Code) || err.Recover || !err.Server
}

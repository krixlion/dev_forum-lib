package rabbitmq

import (
	"context"
	"errors"

	"github.com/krixlion/dev_forum-lib/tracing"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
)

var ErrFullQueue = errors.New("publish queue is full")

// enqueue appends a message to the publishQueue and return a non-nil error if the queue is full.
func (mq *RabbitMQ) Enqueue(msg Message) error {
	select {
	case mq.publishQueue <- msg:
		return nil
	default:
		return ErrFullQueue
	}
}

func (mq *RabbitMQ) tryToEnqueue(ctx context.Context, message Message, err error, logErrorMessage string) {
	if err := mq.Enqueue(message); err != nil {
		mq.opts.logger.Log(ctx, "Failed to enqueue message", "err", err)
	}

	mq.opts.logger.Log(ctx, logErrorMessage, "err", err)
}

func (mq *RabbitMQ) publishPipelined(ctx context.Context, messages <-chan Message) {
	go func() {
		channel := mq.askForChannel()
		defer channel.Close()

		limiter := make(chan struct{}, mq.config.MaxWorkers)

		for {
			select {
			case message := <-messages:
				limiter <- struct{}{}
				go func() {
					ctx := tracing.InjectMetadataIntoContext(ctx, message.Headers)
					ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.publishPipelined", trace.WithSpanKind(trace.SpanKindProducer))
					defer span.End()
					defer func() { <-limiter }()

					done, err := mq.breaker.Allow()
					if err != nil {
						tracing.SetSpanErr(span, err)
						mq.tryToEnqueue(ctx, message, err, "Failed to publish msg")
						return
					}

					p := amqp.Publishing{
						ContentType: string(message.ContentType),
						Body:        message.Body,
						Timestamp:   message.Timestamp,
						Headers:     extractAMQPHeadersFromCtx(ctx),
					}

					if err := channel.PublishWithContext(ctx, message.ExchangeName, message.RoutingKey, false, false, p); err != nil {
						tracing.SetSpanErr(span, err)
						done(!isConnectionError(err))
						mq.tryToEnqueue(ctx, message, err, "Failed to publish msg")
						return
					}
					done(true)
				}()

			case <-ctx.Done():
				channel.Close()
				return
			}
		}
	}()
}

func (mq *RabbitMQ) prepareExchangePipelined(ctx context.Context, msgs <-chan Message) <-chan Message {
	preparedMessages := make(chan Message)

	go func() {
		channel := mq.askForChannel()
		defer channel.Close()
		limiter := make(chan struct{}, mq.config.MaxWorkers)

		for {
			select {
			case message := <-msgs:
				limiter <- struct{}{}
				go func() {
					ctx := tracing.InjectMetadataIntoContext(ctx, message.Headers)
					ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.prepareExchangePipelined", trace.WithSpanKind(trace.SpanKindProducer))
					defer span.End()
					defer func() { <-limiter }()

					done, err := mq.breaker.Allow()
					if err != nil {
						done(!isConnectionError(err))
						tracing.SetSpanErr(span, err)
						mq.tryToEnqueue(ctx, message, err, "Failed to prepare exchange before publishing")
						return
					}

					if err := channel.ExchangeDeclare(message.ExchangeName, message.ExchangeType, true, false, false, false, nil); err != nil {
						done(!isConnectionError(err))
						tracing.SetSpanErr(span, err)
						mq.tryToEnqueue(ctx, message, err, "Failed to declare exchange")
						return
					}
					done(true)

					preparedMessages <- message
				}()
			case <-ctx.Done():
				return
			}
		}
	}()

	return preparedMessages
}

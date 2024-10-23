package rabbitmq

import (
	"context"
	"time"

	"github.com/krixlion/dev_forum-lib/nulls"
	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Log(ctx context.Context, msg string, keyvals ...interface{})
}

type Option interface {
	apply(*options)
}

type Config struct {
	QueueSize         int           // Max number of messages internally queued for publishing.
	MaxWorkers        int           // Max number of concurrent workers per operation type.
	ReconnectInterval time.Duration // Time between reconnect attempts.

	// Settings for the internal circuit breaker.
	MaxRequests   uint32        // Number of requests allowed to half-open state.
	ClearInterval time.Duration // Time after which failed calls count is cleared.
	ClosedTimeout time.Duration // Time after which closed state becomes half-open.
}

func DefaultConfig() Config {
	return Config{
		QueueSize:         100,              // Max number of messages internally queued for publishing.
		MaxWorkers:        30,               // Max number of concurrent workers per operation type.
		ReconnectInterval: time.Second * 2,  // Time between reconnect attempts.
		MaxRequests:       10,               // Number of requests allowed to half-open state.
		ClearInterval:     time.Second * 10, // Time after which failed calls count is cleared.
		ClosedTimeout:     time.Second * 10, // Time after which closed state becomes half-open.
	}
}

func WithTracer(tracer trace.Tracer) Option {
	return optionFunc(func(opts *options) {
		opts.tracer = tracer
	})
}

func WithLogger(logger Logger) Option {
	return optionFunc(func(opts *options) {
		opts.logger = logger
	})
}

type optionFunc func(opts *options)

func (fn optionFunc) apply(opts *options) {
	fn(opts)
}

type options struct {
	tracer trace.Tracer
	logger Logger
}

func defaultOptions() options {
	return options{
		tracer: nulls.NullTracer{},
		logger: nulls.NullLogger{},
	}
}

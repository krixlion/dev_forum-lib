package broker

import (
	"context"
	"encoding/json"

	"github.com/krixlion/dev_forum-lib/event"
	"github.com/krixlion/dev_forum-lib/logging"
	"github.com/krixlion/dev_forum-lib/tracing"
	rabbitmq "github.com/krixlion/dev_forum-rabbitmq"
	"go.opentelemetry.io/otel/trace"
)

// Broker is a wrapper for rabbitmq.RabbitMQ
type Broker struct {
	messageQueue *rabbitmq.RabbitMQ
	logger       logging.Logger
	tracer       trace.Tracer
}

func NewBroker(mq *rabbitmq.RabbitMQ, logger logging.Logger, tracer trace.Tracer) *Broker {
	return &Broker{
		messageQueue: mq,
		logger:       logger,
		tracer:       tracer,
	}
}

// ResilientPublish returns an error only if the queue is full or if it failed to serialize the event.
func (b *Broker) ResilientPublish(e event.Event) error {
	return b.messageQueue.Enqueue(messageFromEvent(e))
}

func (b *Broker) Publish(ctx context.Context, e event.Event) error {
	ctx, span := b.tracer.Start(ctx, "broker.Publish", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	return b.messageQueue.Publish(ctx, messageFromEvent(e))
}

func (b *Broker) Consume(ctx context.Context, queue string, eventType event.EventType) (<-chan event.Event, error) {
	messages, err := b.messageQueue.Consume(ctx, queue, routeFromEvent(eventType))
	if err != nil {
		return nil, err
	}

	events := make(chan event.Event)
	go func() {
		ctx, span := b.tracer.Start(ctx, "broker.Consume", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		for message := range messages {
			var event event.Event
			if err := json.Unmarshal(message.Body, &event); err != nil {
				tracing.SetSpanErr(span, err)
				b.logger.Log(ctx, "Failed to process message", "err", err)
				continue
			}

			events <- event
		}
	}()

	return events, nil
}

func (b *Broker) Close() error {
	return b.messageQueue.Close()
}

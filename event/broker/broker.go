package broker

import (
	"context"
	"encoding/json"

	"github.com/krixlion/dev_forum-lib/event"
	"github.com/krixlion/dev_forum-lib/logging"
	rabbitmq "github.com/krixlion/dev_forum-lib/rabbitmq"
	"github.com/krixlion/dev_forum-lib/tracing"
	"go.opentelemetry.io/otel/trace"
)

// Broker is a wrapper for rabbitmq.RabbitMQ.
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
	msg, err := messageFromEvent(e)
	if err != nil {
		return err
	}

	return b.messageQueue.Enqueue(msg)
}

func (b *Broker) Publish(ctx context.Context, e event.Event) (err error) {
	ctx, span := b.tracer.Start(ctx, "broker.Publish", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()
	defer tracing.SetSpanErr(span, err)

	msg, err := messageFromEvent(e)
	if err != nil {
		return err
	}

	return b.messageQueue.Publish(ctx, msg)
}

func (b *Broker) Consume(ctx context.Context, queue string, eventType event.EventType) (_ <-chan event.Event, err error) {
	ctx, span := b.tracer.Start(ctx, "broker.Consume init")
	defer span.End()
	defer tracing.SetSpanErr(span, err)

	r, err := routeFromEvent(eventType)
	if err != nil {
		return nil, err
	}

	messages, err := b.messageQueue.Consume(ctx, queue, r)
	if err != nil {
		return nil, err
	}

	events := make(chan event.Event)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messages:
				func() {
					ctx, span := b.tracer.Start(tracing.InjectMetadataIntoContext(context.Background(), msg.Headers), "broker.Consume", trace.WithSpanKind(trace.SpanKindConsumer))
					defer span.End()

					e := event.Event{}
					if err := json.Unmarshal(msg.Body, &e); err != nil {
						tracing.SetSpanErr(span, err)
						b.logger.Log(ctx, "Failed to process message", "err", err)
						return
					}

					e.Metadata = tracing.ExtractMetadataFromContext(ctx)
					events <- e
				}()
			}
		}
	}()

	return events, nil
}

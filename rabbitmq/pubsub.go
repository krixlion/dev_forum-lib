package rabbitmq

import (
	"context"

	"github.com/krixlion/dev_forum-lib/chans"
	"github.com/krixlion/dev_forum-lib/tracing"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
)

func (mq *RabbitMQ) Publish(ctx context.Context, msg Message) (err error) {
	ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.Publish")
	defer span.End()
	defer tracing.SetSpanErr(span, err)

	if err := mq.prepareExchange(ctx, msg.Route); err != nil {
		return err
	}

	return mq.publish(ctx, msg)
}

// prepareExchange validates a message and declares a RabbitMQ exchange derived from the message.
func (mq *RabbitMQ) prepareExchange(ctx context.Context, route Route) (err error) {
	ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.prepareExchange")
	defer span.End()
	defer tracing.SetSpanErr(span, err)

	ch := mq.askForChannel()
	defer ch.Close()

	if err := ctx.Err(); err != nil {
		return err
	}

	done, err := mq.breaker.Allow()
	if err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(route.ExchangeName, route.ExchangeType, true, false, false, false, nil); err != nil {
		done(!isConnectionError(err))
		return err
	}
	done(true)

	return nil
}

func (mq *RabbitMQ) publish(ctx context.Context, msg Message) (err error) {
	ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.publish", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()
	defer tracing.SetSpanErr(span, err)

	ch := mq.askForChannel()
	defer ch.Close()

	if err := ctx.Err(); err != nil {
		return err
	}

	done, err := mq.breaker.Allow()
	if err != nil {
		return err
	}

	p := amqp.Publishing{
		ContentType: string(msg.ContentType),
		Body:        msg.Body,
		Timestamp:   msg.Timestamp,
		Headers:     extractAMQPHeadersFromCtx(ctx),
	}

	if err := ch.PublishWithContext(ctx, msg.ExchangeName, msg.RoutingKey, false, false, p); err != nil {
		done(!isConnectionError(err))
		return err
	}

	done(true)
	return nil
}

func (mq *RabbitMQ) Consume(ctx context.Context, command string, route Route) (_ <-chan Message, err error) {
	ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.Consume init")
	defer span.End()
	defer tracing.SetSpanErr(span, err)

	messages := make(chan Message)
	ch := mq.askForChannel()

	queue, err := mq.prepareQueue(ctx, command, route)
	if err != nil {
		return nil, err
	}

	done, err := mq.breaker.Allow()
	if err != nil {
		return nil, err
	}

	deliveries, err := ch.ConsumeWithContext(ctx, queue.Name, mq.consumerName, false, false, false, false, nil)
	if err != nil {
		done(!isConnectionError(err))
		return nil, err
	}
	done(true)

	go func() {
		for {
			select {
			case delivery := <-deliveries:
				func() {
					ctx := injectAMQPHeadersIntoCtx(context.Background(), delivery.Headers)
					ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.Consume", trace.WithSpanKind(trace.SpanKindConsumer))
					defer span.End()

					if err := delivery.Ack(false); err != nil {
						tracing.SetSpanErr(span, err)
						mq.opts.logger.Log(ctx, "Failed to acknowledge message delivery", "err", err)
						return
					}

					message := Message{
						Route:       route,
						Body:        delivery.Body,
						ContentType: ContentType(delivery.ContentType),
						Timestamp:   delivery.Timestamp,
						Headers:     tracing.ExtractMetadataFromContext(ctx),
					}

					chans.NonBlockSend(messages, message)
				}()
			case <-ctx.Done():
				close(messages)
				return
			}
		}
	}()

	return messages, nil
}

func (mq *RabbitMQ) prepareQueue(ctx context.Context, command string, route Route) (_ amqp.Queue, err error) {
	ctx, span := mq.opts.tracer.Start(ctx, "rabbitmq.prepareQueue")
	defer span.End()
	defer tracing.SetSpanErr(span, err)

	ch := mq.askForChannel()

	done, err := mq.breaker.Allow()
	if err != nil {
		return amqp.Queue{}, err
	}

	queue, err := ch.QueueDeclare(command, false, false, false, false, nil)
	if err != nil {
		done(!isConnectionError(err))
		return amqp.Queue{}, err
	}
	done(true)

	if err := ctx.Err(); err != nil {
		return amqp.Queue{}, err
	}

	if err := mq.prepareExchange(ctx, route); err != nil {
		return amqp.Queue{}, err
	}

	done, err = mq.breaker.Allow()
	if err != nil {
		return amqp.Queue{}, err
	}

	if err := ch.QueueBind(queue.Name, route.RoutingKey, route.ExchangeName, false, nil); err != nil {
		done(!isConnectionError(err))
		return amqp.Queue{}, err
	}
	done(true)

	return queue, nil
}

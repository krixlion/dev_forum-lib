package rabbitmq_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/krixlion/dev_forum-lib/internal/gentest"
	rabbitmq "github.com/krixlion/dev_forum-rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestPubSubPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Pub/Sub Pipeline integration test")
	}

	tests := []struct {
		desc    string
		msg     rabbitmq.Message
		wantErr bool
	}{
		{
			desc: "Test if a simple message is correctly published through a pipeline and consumed.",
			msg: rabbitmq.Message{
				Body:        gentest.RandomJSONArticle(2, 5),
				ContentType: rabbitmq.ContentTypeJson,
				Timestamp:   time.Now().Round(time.Second),
				Route: rabbitmq.Route{
					ExchangeName: "test",
					ExchangeType: amqp.ExchangeTopic,
					RoutingKey:   "test.event." + strings.ToLower(gentest.RandomString(5)),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			mq := setUpMQ(t)
			defer mq.Close()

			err := mq.Enqueue(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQ.Enqueue() error = %+v\n wantErr = %+v\n", err, tt.wantErr)
				return
			}

			messages, err := mq.Consume(ctx, gentest.RandomString(5), tt.msg.Route)
			if (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQ.Consume() error = %+v\n wantErr = %+v\n", err, tt.wantErr)
				return
			}

			msg := <-messages
			if !cmp.Equal(tt.msg, msg) {
				t.Errorf("Messages are not equal:\n want = %+v\n  got = %+v\n diff = %+v\n", tt.msg, msg, cmp.Diff(tt.msg, msg))
				return
			}
		})
	}
}

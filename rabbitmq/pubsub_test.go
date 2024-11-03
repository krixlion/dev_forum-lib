package rabbitmq_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/krixlion/dev_forum-lib/env"
	"github.com/krixlion/dev_forum-lib/internal/gentest"
	"github.com/krixlion/dev_forum-lib/logging"
	rabbitmq "github.com/krixlion/dev_forum-lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func setUpMQ(t *testing.T) *rabbitmq.RabbitMQ {
	const consumer = "TESTING"
	port := "5672"
	host := "localhost"
	user := "guest"
	pass := "guest"

	amqp.SetLogger(log.Default())
	logger, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}

	if err := env.Load("dev_forum-lib"); err != nil {
		t.Logf("Failed to load env file, using default settings, err: %s", err)
	} else {
		port = os.Getenv("MQ_PORT")
		host = os.Getenv("MQ_HOST")
		user = os.Getenv("MQ_USER")
		pass = os.Getenv("MQ_PASS")
	}

	config := rabbitmq.Config{
		QueueSize:         100,
		ReconnectInterval: time.Millisecond * 100,
		MaxRequests:       300000,
		ClearInterval:     time.Millisecond,
		ClosedTimeout:     time.Millisecond,
		MaxWorkers:        10,
	}

	return rabbitmq.NewRabbitMQ(consumer, user, pass, host, port, config, rabbitmq.WithLogger(logger))
}

func TestPubSub(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Pub/Sub integration test...")
	}

	tests := []struct {
		desc    string
		msg     rabbitmq.Message
		wantErr bool
	}{
		{
			desc: "Test if a simple message is correctly published and consumed.",
			msg: rabbitmq.Message{
				Body:        gentest.RandomJSONArticle(2, 5),
				ContentType: rabbitmq.ContentTypeJson,
				Timestamp:   time.Now().Round(time.Second),
				Route: rabbitmq.Route{
					ExchangeName: "test",
					ExchangeType: amqp.ExchangeTopic,
					RoutingKey:   "test.pubsub",
				},
				Headers: make(map[string]string),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mq := setUpMQ(t)
			defer mq.Close()

			err := mq.Publish(ctx, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQ.Publish() error = %+v\n, wantErr = %+v\n", err, tt.wantErr)
				return
			}

			time.Sleep(time.Second * 2)

			msgs, err := mq.Consume(ctx, "test_pubsub", tt.msg.Route)
			if (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQ.Consume() error = %+v\n, wantErr = %+v\n", err, tt.wantErr)
				return
			}

			msg := <-msgs
			if !cmp.Equal(tt.msg, msg) {
				t.Errorf("Messages are not equal:\n want = %+v\n got = %+v\n diff = %+v\n", tt.msg, msg, cmp.Diff(tt.msg, msg))
			}
		})
	}
}

func TestIfExchangeIsCreatedBeforeBindingQueue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	mq := setUpMQ(t)
	defer mq.Close()

	route := rabbitmq.Route{
		ExchangeName: gentest.RandomString(7),
		ExchangeType: amqp.ExchangeTopic,
		RoutingKey:   gentest.RandomString(6),
	}

	if _, err := mq.Consume(ctx, gentest.RandomString(5), route); err != nil {
		t.Errorf("RabbitMQ.Consume() error = %+v\n", err)
	}
}

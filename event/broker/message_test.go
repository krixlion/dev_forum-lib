package broker

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/krixlion/dev_forum-lib/event"
	"github.com/krixlion/dev_forum-lib/internal/gentest"
	rabbitmq "github.com/krixlion/dev_forum-lib/rabbitmq"
)

func Test_messageFromEvent(t *testing.T) {
	jsonArticle := gentest.RandomJSONArticle(3, 5)
	e := event.Event{
		AggregateId: "article",
		Type:        event.ArticleCreated,
		Body:        jsonArticle,
		Timestamp:   time.Now(),
		Metadata:    map[string]string{},
	}
	jsonEvent, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		desc    string
		arg     event.Event
		want    rabbitmq.Message
		wantErr bool
	}{
		{
			desc: "Test if message is correctly processed from simple event",
			arg:  e,
			want: rabbitmq.Message{
				Body:        jsonEvent,
				ContentType: "application/json",
				Timestamp:   e.Timestamp,
				Headers:     map[string]string{},
				Route: rabbitmq.Route{
					ExchangeName: "article",
					ExchangeType: "topic",
					RoutingKey:   "article.event.created",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := messageFromEvent(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("messageFromEvent():\n error = %v\n wantErr = %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("messageFromEvent():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}

func Test_routeFromEvent(t *testing.T) {
	type args struct {
		Type event.EventType
	}
	tests := []struct {
		desc    string
		args    args
		want    rabbitmq.Route
		wantErr bool
	}{
		{
			desc: "Test if returns correct route with simple data.",
			args: args{
				Type: event.ArticleCreated,
			},
			want: rabbitmq.Route{
				ExchangeName: "article",
				ExchangeType: "topic",
				RoutingKey:   "article.event.created",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := routeFromEvent(tt.args.Type)
			if (err != nil) != tt.wantErr {
				t.Errorf("routeFromEvent():\n error = %v\n wantErr = %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("routeFromEvent():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}

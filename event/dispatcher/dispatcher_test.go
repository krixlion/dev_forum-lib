package dispatcher

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/krixlion/dev_forum-lib/event"
	"github.com/krixlion/dev_forum-lib/internal/gentest"
	"github.com/krixlion/dev_forum-lib/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/sync/errgroup"
)

func containsHandler(handlers []event.Handler, target event.Handler) bool {
	for _, handler := range handlers {
		if cmp.Equal(handler, target, cmpopts.IgnoreUnexported(mock.Mock{})) {
			return true
		}
	}
	return false
}

func TestDispatcher_Subscribe(t *testing.T) {
	tests := []struct {
		name     string
		handlers []event.Handler
		eType    event.EventType
	}{
		{
			name:     "Check if simple handler is subscribed succesfully",
			handlers: []event.Handler{mocks.Handler{Mock: new(mock.Mock)}},
			eType:    event.UserCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher := NewDispatcher(10)
			dispatcher.Subscribe(tt.eType, tt.handlers...)

			for _, handler := range tt.handlers {
				if !containsHandler(dispatcher.handlers[tt.eType], handler) {
					t.Errorf("event.Handler was not registered succesfully")
				}
			}
		})
	}
}

func Test_mergeChans(t *testing.T) {
	tests := []struct {
		name string
		want []event.Event
	}{
		{
			name: "Test if receives all events from multiple channels",
			want: []event.Event{
				{
					AggregateId: event.AggregateId(gentest.RandomString(5)),
				},
				{
					AggregateId: event.AggregateId(gentest.RandomString(5)),
					Type:        event.ArticleDeleted,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			chans := func() (chans []<-chan event.Event) {
				for _, e := range tt.want {
					v := make(chan event.Event, 1)
					v <- e
					chans = append(chans, v)
				}
				return
			}()

			out := mergeChans(chans...)
			var got []event.Event
			for i := 0; i < len(tt.want); i++ {
				got = append(got, <-out)
			}

			if !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("Events are not equal:\n got = %+v\n want = %+v\n", got, tt.want)
				return
			}
		})
	}
}

func TestDispatcher_Dispatch(t *testing.T) {
	tests := []struct {
		name    string
		arg     event.Event
		handler mocks.Handler
		broker  mocks.Broker
	}{
		{
			name: "Test if handler is called on simple event",
			arg: event.Event{
				Type:        event.ArticleCreated,
				AggregateId: "article",
			},
			handler: func() mocks.Handler {
				m := mocks.Handler{Mock: new(mock.Mock)}
				m.On("Handle", mock.AnythingOfType("Event")).Return().Once()
				return m
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDispatcher(2)
			d.Subscribe(tt.arg.Type, tt.handler)
			d.Dispatch(tt.arg)

			// Wait for the handler to get invoked in a seperate goroutine.
			time.Sleep(time.Millisecond * 5)

			tt.handler.AssertCalled(t, "Handle", tt.arg)
			tt.handler.AssertNumberOfCalls(t, "Handle", 1)
		})
	}
}

func TestDispatcher_Run(t *testing.T) {
	t.Run("Test if Run() returns on context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		errg, ctx := errgroup.WithContext(ctx)

		d := NewDispatcher(20)
		errg.Go(func() error {
			d.Run(ctx)
			return nil
		})

		before := time.Now()

		cancel()
		errg.Wait() //nolint:errcheck // Err is always nil

		stopTime := time.Since(before)

		// Since dispatcher is not doing any work, shutdown should happen near instantly.
		if stopTime > time.Millisecond {
			t.Errorf("Run did not return on context cancellation\n Time needed for func to return: %v", stopTime.Seconds())
			return
		}
	})
}

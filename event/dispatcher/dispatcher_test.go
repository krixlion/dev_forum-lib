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

func Test_Subscribe(t *testing.T) {
	testCases := []struct {
		desc     string
		handlers []event.Handler
		eType    event.EventType
	}{
		{
			desc:     "Check if simple handler is subscribed succesfully",
			handlers: []event.Handler{mocks.Handler{Mock: new(mock.Mock)}},
			eType:    event.UserCreated,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			dispatcher := NewDispatcher(nil, 10)
			dispatcher.Subscribe(tC.eType, tC.handlers...)

			for _, handler := range tC.handlers {
				if !containsHandler(dispatcher.handlers[tC.eType], handler) {
					t.Errorf("event.Handler was not registered succesfully")
				}
			}
		})
	}
}

func Test_mergeChans(t *testing.T) {
	testCases := []struct {
		desc string
		want []event.Event
	}{
		{
			desc: "Test if receives all events from multiple channels",
			want: []event.Event{
				{
					AggregateId: gentest.RandomString(5),
				},
				{
					AggregateId: gentest.RandomString(5),
					Type:        event.ArticleDeleted,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			chans := func() (chans []<-chan event.Event) {
				for _, e := range tC.want {
					v := make(chan event.Event, 1)
					v <- e
					chans = append(chans, v)
				}
				return
			}()

			out := mergeChans(chans...)
			var got []event.Event
			for i := 0; i < len(tC.want); i++ {
				got = append(got, <-out)
			}

			if !assert.ElementsMatch(t, got, tC.want) {
				t.Errorf("Events are not equal:\n got = %+v\n want = %+v\n", got, tC.want)
				return
			}
		})
	}
}

func Test_Publish(t *testing.T) {
	testCases := []struct {
		desc   string
		arg    event.Event
		broker mocks.Broker
	}{
		{
			desc: "",
			arg:  event.MakeEvent(event.ArticleDeleted, gentest.RandomString(5)),
			broker: func() mocks.Broker {
				m := mocks.Broker{Mock: new(mock.Mock)}
				m.On("ResilientPublish", mock.AnythingOfType("event.Event")).Return(nil).Once()
				return m
			}(),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			d := NewDispatcher(tC.broker, 2)
			d.Publish(tC.arg)

			tC.broker.AssertCalled(t, "ResilientPublish", tC.arg)
			tC.broker.AssertExpectations(t)
			tC.broker.AssertNumberOfCalls(t, "ResilientPublish", 1)
		})
	}
}

func Test_Dispatch(t *testing.T) {
	testCases := []struct {
		desc    string
		arg     event.Event
		handler mocks.Handler
		broker  mocks.Broker
	}{
		{
			desc: "Test if handler is called on simple event",
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
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			d := NewDispatcher(tC.broker, 2)
			d.Subscribe(tC.arg.Type, tC.handler)
			d.Dispatch(tC.arg)

			// Wait for the handler to get invoked in a seperate goroutine.
			time.Sleep(time.Millisecond * 5)

			tC.handler.AssertCalled(t, "Handle", tC.arg)
			tC.handler.AssertNumberOfCalls(t, "Handle", 1)
		})
	}
}

func Test_Run(t *testing.T) {
	t.Run("Test if Run() returns on context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		errg, ctx := errgroup.WithContext(ctx)

		d := NewDispatcher(nil, 20)
		errg.Go(func() error {
			d.Run(ctx)
			return nil
		})

		before := time.Now()
		cancel()
		errg.Wait()
		after := time.Now()
		stopTime := after.Sub(before)

		// If time needed for Run to return was longer than a millisecond or unexpected error was returned.
		if stopTime > time.Millisecond {
			t.Errorf("Run did not stop on context cancellation\n Time needed for func to return: %v", stopTime.Seconds())
			return
		}
	})
}

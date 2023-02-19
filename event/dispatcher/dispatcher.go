package dispatcher

import (
	"context"
	"sync"

	"github.com/krixlion/dev_forum-lib/event"
)

type Dispatcher struct {
	maxWorkers  int
	mu          sync.Mutex
	handlers    map[event.EventType][]event.Handler
	events      <-chan event.Event
	broker      event.Broker
	syncHandler event.Handler
	syncEvents  <-chan event.Event
}

func NewDispatcher(broker event.Broker, maxWorkers int) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		broker:     broker,
		handlers:   make(map[event.EventType][]event.Handler),
	}
}

// AddEventProviders registers provided channels as an event source.
// Events from these providers will be parsed by the subscribed handlers
func (d *Dispatcher) AddEventProviders(providers ...<-chan event.Event) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = mergeChans(providers...)
}

// AddSyncEventProviders registers provided channels as an sync event source.
// Events from these providers will be parsed by the syncHandler.
func (d *Dispatcher) AddSyncEventProviders(providers ...<-chan event.Event) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.syncEvents = mergeChans(providers...)
}

// SetSyncHandler registers provided handler to be used as CatchUp handler
// applying sync events, eg. updating the read model.
func (d *Dispatcher) SetSyncHandler(handler event.Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.syncHandler = handler
}

func (d *Dispatcher) Run(ctx context.Context) {
	for {
		select {
		case event := <-d.syncEvents:
			d.syncHandler.Handle(event)
		case event := <-d.events:
			d.Dispatch(event)
		case <-ctx.Done():
			return
		}
	}
}

func (d *Dispatcher) Subscribe(eType event.EventType, handlers ...event.Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eType] = append(d.handlers[eType], handlers...)
}

func (d *Dispatcher) Publish(e event.Event) {
	if err := d.broker.ResilientPublish(e); err != nil {
		panic(err)
	}
	d.Dispatch(e)
}

func (d *Dispatcher) Dispatch(e event.Event) {
	limit := make(chan struct{}, d.maxWorkers)

	for _, handler := range d.handlers[e.Type] {
		limit <- struct{}{}
		go func(handler event.Handler) {
			handler.Handle(e)
			<-limit
		}(handler)
	}
}

func mergeChans(cs ...<-chan event.Event) <-chan event.Event {
	out := make(chan event.Event)

	wg := sync.WaitGroup{}
	wg.Add(len(cs))

	for _, c := range cs {
		go func(c <-chan event.Event) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	return out
}

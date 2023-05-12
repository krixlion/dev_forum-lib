package dispatcher

import (
	"context"
	"sync"

	"github.com/krixlion/dev_forum-lib/event"
)

type Listener interface {
	// EventHandlers returns all event handlers specific to an implementation
	// which are registered at the composition root.
	// These handlers should not be used to sync read and write models
	// and should be separated from them, applying other domain events.
	EventHandlers() map[event.EventType][]event.Handler
}

type Dispatcher struct {
	maxWorkers int
	events     <-chan event.Event
	mu         sync.Mutex
	handlers   map[event.EventType][]event.Handler
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
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

// Run blocks until the context is cancelled.
// Run starts the dispatcher to listen for events from its event providers and dispatch those events.
func (d *Dispatcher) Run(ctx context.Context) {
	for {
		select {
		case event := <-d.events:
			d.Dispatch(event)
		case <-ctx.Done():
			return
		}
	}
}

// Subscribe registers handlers for specified event type.
// They will be invoked when an according event is dispatched.
func (d *Dispatcher) Subscribe(eType event.EventType, handlers ...event.Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eType] = append(d.handlers[eType], handlers...)
}

// Register is a helper method allowing to subscribe multiple event listeners at once.
func (d *Dispatcher) Register(h ...Listener) {
	for _, v := range h {
		for eType, handlers := range v.EventHandlers() {
			d.Subscribe(eType, handlers...)
		}
	}
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

func mergeChans(channels ...<-chan event.Event) <-chan event.Event {
	out := make(chan event.Event)

	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, c := range channels {
		go func(c <-chan event.Event) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	return out
}

package event

import (
	"context"
)

type Broker interface {
	Consumer
	Publisher
}

type Consumer interface {
	Consume(ctx context.Context, queue string, eventType EventType) (<-chan Event, error)
}

type Publisher interface {

	// Exchanges and queues are maintained internally depending on the type of event.
	Publish(context.Context, Event) error

	// Resilient publish returns only parsing error and on any other error retries each event until it succeeds.
	ResilientPublish(Event) error
}

type Subscriber interface {
	// Subscribe registers an event handler for sepcified types of events.
	Subscribe(Handler, ...EventType)
}

type Handler interface {
	Handle(Event)
}

type HandlerFunc func(Event)

func (fn HandlerFunc) Handle(event Event) {
	fn(event)
}

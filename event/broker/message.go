package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/krixlion/dev_forum-lib/event"
	rabbitmq "github.com/krixlion/dev_forum-lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// messageFromEvent returns a message suitable for pub/sub methods and
// a non-nil error if the event could not be marshaled into JSON.
func messageFromEvent(e event.Event) (rabbitmq.Message, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return rabbitmq.Message{}, fmt.Errorf("invalid JSON tags on event.Event, err: %v", err)
	}

	r, err := routeFromEvent(e.Type)
	if err != nil {
		return rabbitmq.Message{}, err
	}

	return rabbitmq.Message{
		Body:        body,
		ContentType: rabbitmq.ContentTypeJson,
		Route:       r,
		Timestamp:   e.Timestamp,
		Headers:     e.Metadata,
	}, nil
}

func routeFromEvent(eType event.EventType) (rabbitmq.Route, error) {
	noun, action, found := strings.Cut(string(eType), "-")
	if !found {
		return rabbitmq.Route{}, errors.New("event type does not follow {noun}-{action} format")
	}

	return rabbitmq.Route{
		ExchangeName: noun,
		ExchangeType: amqp.ExchangeTopic,
		RoutingKey:   noun + ".event." + action,
	}, nil
}

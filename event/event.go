package event

import (
	"encoding/json"
	"time"
)

// Events are sent to the queue in JSON format.
type Event struct {
	AggregateId AggregateId `json:"aggregate_id,omitempty"`
	Type        EventType   `json:"type,omitempty"`
	Body        []byte      `json:"body,omitempty"` // Must be marshaled to JSON.
	Timestamp   time.Time   `json:"timestamp,omitempty"`
}

// MakeEvent returns an event serialized for general use.
// Panics when data cannot be marshaled into json.
func MakeEvent(aggregateId AggregateId, eType EventType, data interface{}) Event {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return Event{
		AggregateId: aggregateId,
		Type:        eType,
		Body:        jsonData,
		Timestamp:   time.Now(),
	}
}

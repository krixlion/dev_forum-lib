package event

import (
	"encoding/json"
	"time"
)

// Events are sent to the queue in JSON format.
type Event struct {
	AggregateId AggregateId       `json:"aggregate_id,omitempty"`
	Type        EventType         `json:"type,omitempty"`
	Body        []byte            `json:"body,omitempty"` // Must be marshaled to JSON.
	Timestamp   time.Time         `json:"timestamp,omitempty"`
	Metadata    map[string]string // TraceID etc.
}

// MakeEvent returns an event serialized for general use.
// Returns an error when bodycannot be marshaled into json.
func MakeEvent(aggregateId AggregateId, eType EventType, body interface{}, metadata map[string]string) (Event, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return Event{}, err
	}

	return Event{
		AggregateId: aggregateId,
		Type:        eType,
		Body:        jsonBody,
		Metadata:    metadata,
		Timestamp:   time.Now(),
	}, nil
}

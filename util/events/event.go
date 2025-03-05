package events

import (
	"errors"
	"time"
)

type Event struct {
	RoutingKey `json:"routingKey"`
	Msg        any       `json:"msg"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewEvent(routingKey RoutingKey, msg any) Event {
	return Event{
		RoutingKey: routingKey,
		Msg:        msg,
		Timestamp:  time.Now(),
	}
}

func (e Event) IsValid() (bool, error) {
	if !e.RoutingKey.IsValid() {
		return false, errors.New("invalid routing key: " + e.RoutingKey.String())
	}

	if e.Msg == nil {
		return false, errors.New("Event.Msg cannot be a nil pointer")
	}

	if e.Timestamp.After(time.Now()) {
		return false, errors.New("Event.Timestamp cannot be in the future")
	}

	return true, nil
}

package entity

import (
	"encoding/json"
	"errors"
	"time"
)

type EventType string

type EventInfo struct {
	ID        int       `json:"id"`
	Type      EventType `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Event struct {
	Info    EventInfo `json:"event_info"`
	Payload []byte    `json:"payload"`
}

func NewEvent[T any](info EventInfo, jsonPayload T) (Event, error) {

	if info.Type == "" {
		return Event{}, errors.New("empty event type")
	}

	payload, err := json.Marshal(jsonPayload)
	if err != nil {
		return Event{}, err
	}

	return Event{
		Info:    info,
		Payload: payload,
	}, nil
}

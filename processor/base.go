package processor

import "time"

type Event struct {
	ID           int       `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	Type         int       `json:"type"`
	Points       *float64  `json:"points"`
	Payload      *string   `json:"payload"`
	Minute       int       `json:"minute"`
	Second       int       `json:"second" `
	X            float64   `json:"x"`
	Y            float64   `json:"y"`
	PlayerID     *string   `json:"player_id"`
	TeamID       *string   `json:"team_id"`
	MatchID      string    `json:"match_id"`
	MatchEventID int       `json:"match_event_id"`
}

type EventSource interface {
	GetEvents(matchID string) ([]*Event, error)
}

type EventHandler interface {
	Reset() error
	Handle(event *Event, alreadyProcessed bool) error
}

type EventProcessor interface {
	NewEvent(newEvent *Event) error
}

type EventProcessorFactory interface {
	NewEventProcessor(matchID string) (EventProcessor, error)
}

type StateStore interface {
	GetCurrentID(processorID string) (int, error)
	SetCurrentID(processorID string, currentID int) error
}

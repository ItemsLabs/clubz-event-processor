package processor

import (
	"fmt"

	"github.com/gameon-app-inc/fanclash-event-processor/database"

	"github.com/palantir/stacktrace"
)

type StoreEventSource struct {
	store database.Store
}

func NewStoreEventSource(store database.Store) *StoreEventSource {
	return &StoreEventSource{
		store: store,
	}
}

func (s *StoreEventSource) GetEvents(matchID string) ([]*Event, error) {
	events, err := s.store.GetMatchEvents(matchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("cannot get events for match %s", matchID))
	}

	// convert to *Event
	result := make([]*Event, 0, len(events))
	for _, ev := range events {
		var payload, playerID, teamID *string
		if ev.Payload.Valid {
			payload = &ev.Payload.String
		}
		if ev.PlayerID.Valid {
			playerID = &ev.PlayerID.String
		}
		if ev.TeamID.Valid {
			teamID = &ev.TeamID.String
		}

		result = append(result, &Event{
			ID:           ev.ID,
			Timestamp:    ev.Timestamp,
			Type:         ev.Type,
			Points:       &ev.Points.Float64,
			Payload:      payload,
			Minute:       ev.Minute,
			Second:       ev.Second,
			X:            ev.X,
			Y:            ev.Y,
			PlayerID:     playerID,
			TeamID:       teamID,
			MatchID:      ev.MatchID,
			MatchEventID: ev.MatchEventID,
		})
	}

	return result, nil
}

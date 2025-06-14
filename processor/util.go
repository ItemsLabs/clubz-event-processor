package processor

import "github.com/gameon-app-inc/fanclash-event-processor/database/schema"

func FromMatchEvent(ev *schema.MatchEvent) *Event {
	return &Event{
		ID:           ev.ID,
		Timestamp:    ev.Timestamp,
		Type:         ev.Type,
		Points:       ev.Points.Ptr(),
		Payload:      ev.Payload.Ptr(),
		Minute:       ev.Minute,
		Second:       ev.Second,
		X:            ev.X,
		Y:            ev.Y,
		PlayerID:     ev.PlayerID.Ptr(),
		TeamID:       ev.TeamID.Ptr(),
		MatchID:      ev.MatchID,
		MatchEventID: ev.MatchEventID,
	}
}

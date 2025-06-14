package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) InsertGameEvent(ev *schema.GameEvent) error {
	return ev.Insert(s.db, boil.Infer())
}

func (s *DBStore) DeleteGameEventsByMatchEventID(id int) error {
	_, err := schema.GameEvents(
		qm.Where("match_event_id = ?", id),
	).DeleteAll(s.db)
	return err
}

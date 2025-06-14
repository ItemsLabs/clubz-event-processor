package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetMatchEvents(matchID string) (schema.MatchEventSlice, error) {
	return schema.MatchEvents(
		qm.Where("match_id = ?", matchID),
		qm.OrderBy("match_event_id"),
	).All(s.db)
}

func (s *DBStore) GetMatchEventByID(id int) (*schema.MatchEvent, error) {
	return schema.MatchEvents(
		qm.Where("id = ?", id),
	).One(s.db)
}

func (s *DBStore) GetMatchEventByMatchEventID(matchID string, matchEventID int) (*schema.MatchEvent, error) {
	return schema.MatchEvents(
		qm.Where("match_id = ?", matchID),
		qm.Where("match_event_id = ?", matchEventID),
	).One(s.db)
}

func (s *DBStore) UpdateMatchEventStatus(ev *schema.MatchEvent) error {
	_, err := ev.Update(s.db, boil.Whitelist("status"))
	return err
}

func (s *DBStore) GetLatestMatchEvent(matchID string) (*schema.MatchEvent, error) {
	return schema.MatchEvents(
		qm.Where("match_id = ?", matchID),
		qm.OrderBy("match_event_id desc"),
		qm.Limit(1),
	).One(s.db)
}

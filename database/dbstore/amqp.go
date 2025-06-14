package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (s *DBStore) InsertAMQPEvent(ev *schema.AmqpEvent) (*schema.AmqpEvent, error) {
	return ev, ev.Insert(s.db, boil.Infer())
}

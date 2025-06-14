package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetPlayerByID(playerID string) (*schema.Player, error) {
	return schema.Players(
		qm.Where("id = ?", playerID),
	).One(s.db)
}

package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) InsertUserTransaction(t *schema.Transaction) (*schema.Transaction, error) {
	return t, t.Insert(s.db, boil.Infer())
}

func (s *DBStore) GetUserTransaction(userID, matchID string) (*schema.Transaction, error) {
	return schema.Transactions(
		qm.Where("user_id = ?", userID),
		qm.Where("match_id = ?", matchID),
	).One(s.db)
}

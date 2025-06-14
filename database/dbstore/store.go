package dbstore

import (
	"database/sql"

	"github.com/gameon-app-inc/fanclash-event-processor/database"

	"github.com/gsamokovarov/sx"
)

var _ database.Store = (*DBStore)(nil)

type DBStore struct {
	db sx.Transactor
}

func New(db *sql.DB) *DBStore {
	return &DBStore{
		db: sx.NewTransactor(db),
	}
}

func (s *DBStore) Transaction(fn func(database.Store) error) (err error) {
	return sx.Transaction(s.db, func(tx sx.Transactor) error {
		return fn(&DBStore{db: tx})
	})
}

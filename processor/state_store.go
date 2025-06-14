package processor

import "github.com/gameon-app-inc/fanclash-event-processor/database"

var _ StateStore = (*DBStateStore)(nil)

type DBStateStore struct {
	store database.Store
}

func NewDBStateStore(store database.Store) *DBStateStore {
	return &DBStateStore{
		store: store,
	}
}

func (s *DBStateStore) GetCurrentID(processorID string) (int, error) {
	return s.store.GetMatchEventProcessorLastProcessedID(processorID)
}

func (s *DBStateStore) SetCurrentID(processorID string, currentID int) error {
	return s.store.SetMatchEventProcessorLastProcessedID(processorID, currentID)
}

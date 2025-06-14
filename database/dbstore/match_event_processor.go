package dbstore

import (
	"database/sql"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetOrCreateMatchEventProcessor(matchID string, typ int) (*schema.MatchEventProcessor, error) {
	// first get
	proc, err := schema.MatchEventProcessors(
		qm.Where("match_id = ?", matchID),
		qm.Where("type = ?", typ),
	).One(s.db)

	if err == nil {
		return proc, nil
	}

	if err == sql.ErrNoRows {
		// create new one
		proc := &schema.MatchEventProcessor{
			ID:      uuid.New().String(),
			MatchID: matchID,
			Type:    typ,
		}

		return proc, proc.Insert(s.db, boil.Infer())
	}
	return nil, err
}

func (s *DBStore) GetMatchEventProcessorLastProcessedID(processorID string) (int, error) {
	proc, err := schema.MatchEventProcessors(
		qm.Where("id = ?", processorID),
	).One(s.db)
	if err != nil {
		return 0, err
	}

	return proc.LastProcessedID, nil
}

func (s *DBStore) SetMatchEventProcessorLastProcessedID(processorID string, lastProcessedID int) error {
	proc := &schema.MatchEventProcessor{ID: processorID, LastProcessedID: lastProcessedID}
	_, err := proc.Update(s.db, boil.Whitelist("last_processed_id"))
	return err
}

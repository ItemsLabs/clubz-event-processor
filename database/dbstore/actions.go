package dbstore

import (
	"database/sql"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetActions() (schema.ActionSlice, error) {
	return schema.Actions(
		qm.OrderBy("score desc"),
	).All(s.db)
}

// GetActionNameByActionID returns the name of an action given its ID.
// It returns an empty string and an error if no action is found for the given ID.
func (s *DBStore) GetActionNameByActionID(actionID int) (string, error) {
	// Ensure s.db is properly initialized and connected to your database
	action, err := schema.FindAction(s.db, actionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // or return an error indicating that no action was found
		}
		return "", err
	}

	return action.Name, nil
}

// GetActionsByTypeID retrieves all actions that match the given type ID.
func (s *DBStore) GetActionsByTypeID(typeID int) (schema.ActionSlice, error) {
	return schema.Actions(
		qm.Where("id = ?", typeID),
		qm.OrderBy("score desc"),
	).All(s.db)
}

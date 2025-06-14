package dbstore

import "github.com/gameon-app-inc/fanclash-event-processor/database/schema"

func (s *DBStore) GetPowerUpActions() (schema.PowerupActionSlice, error) {
	return schema.PowerupActions().All(s.db)
}

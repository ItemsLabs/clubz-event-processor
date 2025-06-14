package dbstore

import (
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetActiveGamePowerUpsAtTime(matchID string, t time.Time) (schema.GamePowerupSlice, error) {
	return schema.GamePowerups(
		qm.InnerJoin("games on game_powerups.game_id = games.id"),
		qm.Where("games.match_id = ?", matchID),
		qm.Where("game_powerups.position > ?", 0),
		qm.Where("game_powerups.created_at <= ?", t),
		qm.Where("(game_powerups.ended_at is null or game_powerups.ended_at > ?)", t),
		qm.Load(schema.GamePowerupRels.Game),
		qm.Load(schema.GamePowerupRels.Powerup),
	).All(s.db)
}

func (s *DBStore) UpdateEndedAt(pu *schema.GamePowerup) error {
	_, err := pu.Update(s.db, boil.Whitelist("ended_at"))
	return err
}

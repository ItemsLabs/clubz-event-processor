package dbstore

import (
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetActivePicksAtTime(matchID string, playerID string, timestamp time.Time) (schema.GamePickSlice, error) {
	return schema.GamePicks(
		qm.InnerJoin("games on game_picks.game_id = games.id"),
		qm.Where("games.match_id = ?", matchID),
		qm.Where("game_picks.player_id = ?", playerID),
		qm.Where("game_picks.created_at <= ?", timestamp),
		qm.Where("(game_picks.ended_at is null or game_picks.ended_at > ?)", timestamp),
		qm.Load("Game"),
		qm.Load("Player"),
	).All(s.db)
}

func (s *DBStore) GetActivePicksAtMinSec(matchID string, playerID string, min, sec int) (schema.GamePickSlice, error) {
	t := min*60 + sec
	return schema.GamePicks(
		qm.InnerJoin("games on game_picks.game_id = games.id"),
		qm.Where("games.match_id = ?", matchID),
		qm.Where("game_picks.player_id = ?", playerID),
		qm.Where("(game_picks.minute * 60 + game_picks.second) <= ?", t),
		qm.Where("(game_picks.ended_minute is null or ((game_picks.ended_minute * 60 + game_picks.ended_second) > ?))", t),
		qm.Load("Game"),
	).All(s.db)
}

func (s *DBStore) UpdatePickScore(pick *schema.GamePick) error {
	_, err := pick.Update(s.db, boil.Whitelist("score"))
	return err
}

func (s *DBStore) UpdatePickEndedAt(pick *schema.GamePick) error {
	_, err := pick.Update(s.db, boil.Whitelist("ended_at"))
	return err
}

func (s *DBStore) InsertPick(pick *schema.GamePick) (*schema.GamePick, error) {
	return pick, pick.Insert(s.db, boil.Infer())
}

func (s *DBStore) GetGamePicks(gameID string) (schema.GamePickSlice, error) {
	return schema.GamePicks(
		qm.Where("game_id = ?", gameID),
	).All(s.db)
}

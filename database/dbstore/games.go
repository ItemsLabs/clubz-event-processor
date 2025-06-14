package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) UpdateMatchGamesStatus(matchID string, status string) error {
	_, err := schema.Games(
		qm.Where("match_id = ?", matchID),
		qm.Where("status != ?", status),
	).UpdateAll(s.db, schema.M{"status": status})
	return err
}

func (s *DBStore) SyncGamePremiumFlags(matchID string) error {
	_, err := s.db.Exec(`update games
set premium = (select users.premium from users where users.id = games.user_id),
    subscription_tier = (select users.subscription_tier from users where users.id = games.user_id) 
where match_id = $1`, matchID)
	return err
}

func (s *DBStore) IncGameVersion(gameID string) error {
	_, err := s.db.Exec(`update games set version = version + 1 where id = $1`, gameID)
	return err
}

func (s *DBStore) GetGamesWithoutPowerUps(matchID string) (schema.GameSlice, error) {
	return schema.Games(
		qm.Where("match_id = ?", matchID),
		qm.Where("not exists (select null from game_powerups pu where pu.game_id = games.id)"),
	).All(s.db)
}

func (s *DBStore) GetPowerUpCountForGame(matchID string) (map[string]int, error) {
	query := `
select gp.game_id,
       count(1) power_up_count
  from games g,
       game_powerups gp
 where g.match_id = $1
   and g.id = gp.game_id 
 group by gp.game_id`

	var records []struct {
		GameID       string
		PowerUpCount int
	}

	if err := queries.Raw(query, matchID).Bind(nil, s.db, &records); err != nil {
		return nil, err
	}

	var result = map[string]int{}
	for _, rec := range records {
		result[rec.GameID] = rec.PowerUpCount
	}

	return result, nil
}

// GetUserIDsByMatchID retrieves an array of user IDs for a given match ID.
func (s *DBStore) GetUserIDsByMatchID(matchID string) ([]string, error) {
	games, err := schema.Games(
		qm.Select(schema.GameColumns.UserID),
		qm.Where(schema.GameColumns.MatchID+" = ?", matchID),
	).All(s.db)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, len(games))
	for i, game := range games {
		userIDs[i] = game.UserID
	}

	return userIDs, nil
}

func (s *DBStore) GetGameByUserIDMatchID(userID, matchID string) (*schema.Game, error) {
	return schema.Games(
		qm.Where("user_id = ?", userID),
		qm.Where("match_id = ?", matchID),
		qm.Load("Match.MatchPlayers"),
		qm.Load("GamePicks"),
		qm.Load("GamePicks.Player"),
		qm.Load("GamePowerups.Powerup"),
		qm.Load("User"),
	).One(s.db)
}

func (s *DBStore) GetAllUserIDs() ([]string, error) {
	users, err := schema.Users(qm.Select(schema.UserColumns.ID)).All(s.db)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	return userIDs, nil
}

func (s *DBStore) GetUserIDsNotInMatch(matchID string) ([]string, error) {
	// Get user IDs that are part of the match
	matchUserIDs, err := s.GetUserIDsByMatchID(matchID)
	if err != nil {
		return nil, err
	}

	// Get all user IDs
	allUserIDs, err := s.GetAllUserIDs()
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	matchUserIDMap := make(map[string]struct{}, len(matchUserIDs))
	for _, id := range matchUserIDs {
		matchUserIDMap[id] = struct{}{}
	}

	// Find user IDs that are not in the match
	var notInMatchUserIDs []string
	for _, id := range allUserIDs {
		if _, found := matchUserIDMap[id]; !found {
			notInMatchUserIDs = append(notInMatchUserIDs, id)
		}
	}

	return notInMatchUserIDs, nil
}

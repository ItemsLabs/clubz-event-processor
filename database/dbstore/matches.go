package dbstore

import (
	"database/sql"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) UpdateMatch(match *schema.Match, updatedFields []string) (*schema.Match, error) {
	if _, err := match.Update(s.db, boil.Whitelist(updatedFields...)); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *DBStore) UpdateMatchTime(matchID string, minute, second int) error {
	query := `
update matches
   set minute = $1,
       second = $2
 where id = $3
   and (minute*60 + second) <= ($1*60 + $2)
`

	if _, err := s.db.Exec(query, minute, second, matchID); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) GetMatchByIDWithTeams(matchID string) (*schema.Match, error) {
	return schema.Matches(
		qm.Where("id = ?", matchID),
		qm.Load("HomeTeam"),
		qm.Load("AwayTeam"),
	).One(s.db)
}

func (s *DBStore) GetMatchByID(matchID string) (*schema.Match, error) {
	return schema.Matches(
		qm.Where("id = ?", matchID),
	).One(s.db)
}

func (s *DBStore) IncHomeScore(match *schema.Match) (*schema.Match, error) {
	match.HomeScore += 1
	if _, err := match.Update(s.db, boil.Whitelist("home_score")); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *DBStore) IncAwayScore(match *schema.Match) (*schema.Match, error) {
	match.AwayScore += 1
	if _, err := match.Update(s.db, boil.Whitelist("away_score")); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *DBStore) DecHomeScore(match *schema.Match) (*schema.Match, error) {
	match.HomeScore -= 1
	if _, err := match.Update(s.db, boil.Whitelist("home_score")); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *DBStore) DecAwayScore(match *schema.Match) (*schema.Match, error) {
	match.AwayScore -= 1
	if _, err := match.Update(s.db, boil.Whitelist("away_score")); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *DBStore) GetMatchPlayer(matchID string, teamID string, playerID string) (*schema.MatchPlayer, error) {
	return schema.MatchPlayers(
		qm.Where("match_id = ?", matchID),
		qm.Where("team_id = ?", teamID),
		qm.Where("player_id = ?", playerID),
	).One(s.db)
}

func (s *DBStore) GetMatchPlayerWithoutTeam(matchID string, playerID string) (*schema.MatchPlayer, error) {
	return schema.MatchPlayers(
		qm.Where("match_id = ?", matchID),
		qm.Where("player_id = ?", playerID),
	).One(s.db)
}

func (s *DBStore) GetOrCreateMatchPlayer(matchID string, teamID string, playerID string) (*schema.MatchPlayer, error) {
	mp, err := s.GetMatchPlayer(matchID, teamID, playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.InsertMatchPlayer(&schema.MatchPlayer{
				ID:          uuid.New().String(),
				MatchID:     matchID,
				TeamID:      teamID,
				PlayerID:    playerID,
				FromLineups: true,
			})
		}
		return nil, err
	}

	return mp, nil
}

func (s *DBStore) InsertMatchPlayer(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error) {
	return matchPlayer, matchPlayer.Insert(s.db, boil.Infer())
}

func (s *DBStore) UpdateMatchPlayer(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error) {
	if _, err := matchPlayer.Update(s.db, boil.Whitelist("position", "jersey_number", "from_lineups")); err != nil {
		return nil, err
	}
	return matchPlayer, nil
}

func (s *DBStore) GetMatchPlayers(matchID string) (schema.MatchPlayerSlice, error) {
	return schema.MatchPlayers(
		qm.Where("match_id = ?", matchID),
	).All(s.db)
}

func (s *DBStore) UpdateMatchPlayerScore(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error) {
	if _, err := matchPlayer.Update(s.db, boil.Whitelist("score")); err != nil {
		return nil, err
	}
	return matchPlayer, nil
}

func (s *DBStore) UpdateMatchPlayerPosition(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error) {
	if _, err := matchPlayer.Update(s.db, boil.Whitelist("position")); err != nil {
		return nil, err
	}
	return matchPlayer, nil
}

func (s *DBStore) UpdateMatchPlayerPlayedSeconds(matchPlayer *schema.MatchPlayer) (*schema.MatchPlayer, error) {
	if _, err := matchPlayer.Update(s.db, boil.Whitelist("played_seconds")); err != nil {
		return nil, err
	}
	return matchPlayer, nil
}

func (s *DBStore) IncMatchPlayerScoreFast(matchID string, teamID string, playerID string, score float64) error {
	query := `
update match_players
   set score = coalesce(score, 0) + $1
 where match_id = $2
   and player_id = $3
   and team_id = $4
`

	if _, err := s.db.Exec(query, score, matchID, playerID, teamID); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) GetMatchRewards(matchID string) (schema.MatchRewardSlice, error) {
	return schema.MatchRewards(
		qm.Where("match_id = ?", matchID),
		qm.OrderBy("position"),
	).All(s.db)
}

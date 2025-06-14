package dbstore

import (
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) UpdateMatchLeaderboard(matchID string) error {
	// generate leaderboard
	_, err := s.db.Exec(`select update_match_leaderboard($1)`, matchID)

	return err
}

func (s *DBStore) SetTransactionIDForMatchLeaderboard(userID, matchID, transactionID string) error {
	_, err := schema.MatchLeaderboards(
		qm.Where("user_id != ?", userID),
		qm.And("match_id = ?", matchID),
	).UpdateAll(s.db, schema.M{"transaction_id": transactionID})
	return err
}

func (s *DBStore) GetMatchLeaderboard(matchID string) (schema.MatchLeaderboardSlice, error) {
	return schema.MatchLeaderboards(
		qm.Where("match_id = ?", matchID),
		qm.Load("User"),
		qm.Load("Game"),
		qm.OrderBy("position"),
	).All(s.db)
}

func (s *DBStore) GetTopMatchLeaderboard(matchID string, maxPosition int) (schema.MatchLeaderboardSlice, error) {
	return schema.MatchLeaderboards(
		qm.Where("match_id = ?", matchID),
		qm.Where("position <= ?", maxPosition),
		qm.Where("position is not null"),
		qm.OrderBy("position"),
		qm.Load("User"),
		qm.Load("Game"),
	).All(s.db)
}

func (s *DBStore) GetFullMatchLeaderboard(matchID string) (schema.MatchLeaderboardSlice, error) {
	return schema.MatchLeaderboards(
		qm.Where("match_id = ?", matchID),
		qm.Where("position is not null"),
		qm.OrderBy("position"),
	).All(s.db)
}

func (s *DBStore) GetMatchLeaderboardAtTime(matchID string, t time.Time) (schema.MatchLeaderboardSlice, error) {
	query := `
select 0 id,
       score,
       rank() over(order by score desc) as position,
       game_id,
       match_id,
       user_id
  from (select g.id game_id,
               g.match_id,
               g.user_id,
               sum(ev.score) score
          from game_events ev,
               games g
         where ev.created_at < $1
           and ev.game_id = g.id
           and g.match_id = $2
         group by g.id,
                  g.match_id,
                  g.user_id
       ) t
`

	var results schema.MatchLeaderboardSlice
	if err := queries.Raw(query, t, matchID).Bind(nil, s.db, &results); err != nil {
		return nil, err
	}

	return results, nil
}

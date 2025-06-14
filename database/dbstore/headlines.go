package dbstore

import (
	"database/sql"
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/palantir/stacktrace"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetGameCount(matchID string) (int, error) {
	cnt, err := schema.Games(
		qm.Where("match_id = ?", matchID),
	).Count(s.db)
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}

func (s *DBStore) GetRandomUserAvatars(matchID string, count int) ([]string, error) {
	query := `
select usr.avatar_url avatar
  from users usr,
       games g
 where g.match_id = $1
   and g.user_id = usr.id
   and usr.avatar_url is not null
 order by random()
 limit $2`

	var records []struct {
		Avatar string
	}

	if err := queries.Raw(query, matchID, count).Bind(nil, s.db, &records); err != nil {
		return nil, err
	}
	var result []string
	for _, el := range records {
		result = append(result, el.Avatar)
	}
	return result, nil
}

func (s *DBStore) GetMostPickedPlayer(matchID string, isStar bool) (*schema.Player, error) {
	query := `
select p.*
  from players p,
       (select player_id
          from (select gp.player_id,
                       count(*) cnt
                  from match_players mp,
                       game_picks gp,
                       games g
                 where mp.match_id = $1
                   and mp.player_id = gp.player_id
                   and mp.is_star = $2
                   and gp.game_id = g.id
                   and g.match_id = $3
                 group by gp.player_id
               ) t
         order by t.cnt desc
         limit 1
       ) gp
 where p.id = gp.player_id
`

	var player schema.Player
	if err := queries.Raw(query, matchID, isStar, matchID).Bind(nil, s.db, &player); err != nil {
		return nil, err
	}
	return &player, nil
}

func (s *DBStore) GetHighestTeamPickPercent(matchID string) (*schema.Team, float64, error) {
	query := `
select mp.team_id,
       count(*) count
  from match_players mp,
       game_picks gp,
       games g
 where mp.match_id = $1
   and mp.player_id = gp.player_id
   and gp.game_id = g.id
   and g.match_id = $2
 group by mp.team_id
 order by count desc
`

	type TeamStat struct {
		TeamID string
		Count  int
	}

	var stats []*TeamStat
	if err := queries.Raw(query, matchID, matchID).Bind(nil, s.db, &stats); err != nil {
		return nil, 0, err
	}

	// first is max one
	if len(stats) == 0 {
		return nil, 0, nil
	}

	total := 0
	for _, st := range stats {
		total += st.Count
	}

	// find team by id
	team, err := s.GetTeamByID(stats[0].TeamID)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return team, 0, nil
	}

	return team, float64(stats[0].Count) / float64(total) * 100, nil
}

func (s *DBStore) GetNewPlayersCount(matchID string) (int, error) {
	query := `
select count(*) count
  from games g
 where g.match_id = $1
   and not exists (
        select null
          from games g2
         where g2.user_id = g.user_id
           and g2.created_at < g.created_at
       )
`

	var cntResult struct {
		Count int
	}
	if err := queries.Raw(query, matchID).Bind(nil, s.db, &cntResult); err != nil {
		return 0, err
	}
	return cntResult.Count, nil
}

func (s *DBStore) GetNewPlayerAvatars(matchID string, count int) ([]string, error) {
	query := `
select usr.avatar_url avatar
  from games g,
       users usr
 where g.match_id = $1
   and g.user_id = usr.id
   and usr.avatar_url is not null
   and not exists (
        select null
          from games g2
         where g2.user_id = g.user_id
           and g2.created_at < g.created_at
       )
 order by random()
 limit $2
`

	var records []struct {
		Avatar string
	}
	if err := queries.Raw(query, matchID, count).Bind(nil, s.db, &records); err != nil {
		return nil, err
	}

	var result []string
	for _, el := range records {
		result = append(result, el.Avatar)
	}
	return result, nil
}

func (s *DBStore) GetMatchesForHeadlines(t time.Time) (schema.MatchSlice, error) {
	return schema.Matches(
		qm.WhereIn("status in ?", []interface{}{database.MatchStatusLineups, database.MatchStatusGame}...),
		// few hours before match starts
		qm.Or("match_time between ? and ?", t.Add(-time.Hour), t.Add(3*time.Hour)),
	).All(s.db)
}

func (s *DBStore) GetCurrentMatchTime(matchID string) (int, int, error) {
	// pick last match event
	ev, err := schema.MatchEvents(
		qm.Where("match_id = ?", matchID),
		qm.OrderBy("match_event_id desc"),
		qm.Limit(1),
	).One(s.db)

	if err == sql.ErrNoRows {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	return ev.Minute, ev.Second, nil
}

func (s *DBStore) GetCurrentMatchMinute(matchID string) (int, error) {
	minute, _, err := s.GetCurrentMatchTime(matchID)
	return minute, err
}

func (s *DBStore) GetTopPlayerByAvg(matchID string) (*schema.Player, error) {
	return schema.Players(
		qm.InnerJoin("match_players on match_players.player_id = players.id"),
		qm.Where("match_players.match_id = ?", matchID),
		qm.OrderBy("players.avg_score desc nulls last"),
		qm.Limit(1),
	).One(s.db)
}

func (s *DBStore) GetAvgPlayerGoals(playerID string) (float64, error) {
	query := `
select coalesce(avg(goal_count),0) avg_goal_count
  from (
        with m as (
          select match_id
            from matches m,
                 match_players mp
           where m.match_time < current_timestamp
             and m.status = 'e'
             and mp.match_id = m.id
             and mp.player_id = $1
           order by m.match_time desc
           limit 5
        )
        select sum(case when ev.type = 45 then 1 else 0 end) goal_count
          from m,
               match_events ev
         where ev.match_id = m.match_id
           and ev.player_id = $1
           and ev.status = 1
         group by m.match_id
       ) t
`

	var scoreResult struct {
		AvgGoalCount float64
	}

	if err := queries.Raw(query, playerID).Bind(nil, s.db, &scoreResult); err != nil {
		return 0, err
	}

	return scoreResult.AvgGoalCount, nil
}

func (s *DBStore) GetTopPlayer(matchID string, fromMin int, toMin int) (*schema.Player, float64, error) {
	query := `
select ev.player_id,
       sum(coalesce(ev.points, 0)) points
  from match_events ev
 where ev.match_id = $1
   and ev.minute >= $2
   and ev.minute < $3
   and ev.player_id is not null
   and ev.points is not null
   and ev.status = 1
 group by ev.player_id
 order by points desc
 limit 1
`

	var topPlayer struct {
		PlayerID string
		Points   float64
	}

	if err := queries.Raw(query, matchID, fromMin, toMin).Bind(nil, s.db, &topPlayer); err != nil {
		return nil, 0, err
	}

	player, err := s.GetPlayerByID(topPlayer.PlayerID)
	if err != nil {
		return nil, 0, err
	}

	return player, topPlayer.Points, nil
}

func (s *DBStore) GetWorstPlayer(matchID string, fromMin int, toMin int) (*schema.Player, float64, error) {
	query := `
select ev.player_id,
       sum(coalesce(ev.points, 0)) points
  from match_events ev
 where ev.match_id = $1
   and ev.minute >= $2
   and ev.minute < $3
   and ev.player_id is not null
   and ev.points is not null
   and ev.status = 1
 group by ev.player_id
 order by points asc
 limit 1
`

	var worstPlayer struct {
		PlayerID string
		Points   float64
	}

	if err := queries.Raw(query, matchID, fromMin, toMin).Bind(nil, s.db, &worstPlayer); err != nil {
		return nil, 0, err
	}

	player, err := s.GetPlayerByID(worstPlayer.PlayerID)
	if err != nil {
		return nil, 0, err
	}

	return player, worstPlayer.Points, nil
}

func (s *DBStore) GetMostTraded(matchID string, fromMin, toMin int) (*schema.Player, int, error) {
	// TODO: remove picks made after game started
	query := `
select gp.player_id,
       count(*) count
  from game_picks gp,
       games g
 where gp.game_id = g.id
   and g.match_id = $1
   and gp.minute >= $2
   and gp.minute < $3
   and (gp.minute > 0 or gp.second > 0)
 group by gp.player_id
 order by count desc
 limit 1
`

	var playerCount struct {
		PlayerID string
		Count    int
	}

	if err := queries.Raw(query, matchID, fromMin, toMin).Bind(nil, s.db, &playerCount); err != nil {
		return nil, 0, err
	}

	player, err := s.GetPlayerByID(playerCount.PlayerID)
	if err != nil {
		return nil, 0, err
	}

	return player, playerCount.Count, nil
}

func (s *DBStore) GetPowerUpUsages(matchID string, fromMin, toMin int) (int, error) {
	cnt, err := schema.GamePowerups(
		qm.InnerJoin("games on games.id = game_powerups.game_id"),
		qm.Where("games.match_id = ?", matchID),
		qm.Where("game_powerups.minute >= ?", fromMin),
		qm.Where("game_powerups.minute < ?", toMin),
	).Count(s.db)
	return int(cnt), err
}

func (s *DBStore) GetMostUsedPowerUp(matchID string, fromMin, toMin int) (*schema.Powerup, int, error) {
	query := `
select p.id,
       count(distinct g.user_id) count
  from game_powerups gp,
       games g,
       powerups p
 where gp.game_id = g.id
   and g.match_id = $1
   and gp.minute >= $2
   and gp.minute < $3
   and gp.powerup_id = p.id
 group by p.id
 order by count desc
 limit 1
`
	var puCount struct {
		ID    string
		Count int
	}

	if err := queries.Raw(query, matchID, fromMin, toMin).Bind(nil, s.db, &puCount); err != nil {
		return nil, 0, err
	}

	powerUp, err := schema.Powerups(
		qm.Where("id = ?", puCount.ID),
	).One(s.db)
	if err != nil {
		return nil, 0, err
	}

	return powerUp, puCount.Count, nil
}

func (s *DBStore) GetBestUser(matchID string, fromMin, toMin int) (*schema.User, float64, error) {
	query := `
select g.user_id,
       sum(ev.score) total_score
  from game_events ev,
       games g
 where ev.game_id = g.id
   and g.match_id = $1
   and ev.minute >= $2
   and ev.minute < $3
 group by g.user_id
 order by total_score desc
 limit 1
`
	var bestUser struct {
		UserID     string
		TotalScore float64
	}

	if err := queries.Raw(query, matchID, fromMin, toMin).Bind(nil, s.db, &bestUser); err != nil {
		return nil, 0, err
	}

	// select user info
	user, err := schema.Users(
		qm.Where("id = ?", bestUser.UserID),
	).One(s.db)
	if err != nil {
		return nil, 0, err
	}

	return user, bestUser.TotalScore, nil
}

func (s *DBStore) GetTeamsWithScore(matchID string) ([]*database.TeamWithValue, error) {
	query := `
select ev.team_id,
       coalesce(sum(ev.points), 0) total_score
  from match_events ev
 where ev.match_id = $1
   and ev.team_id is not null
   and ev.status = 1
 group by ev.team_id
 order by total_score desc
`
	var teamInfo []struct {
		TeamID     string
		TotalScore float64
	}

	if err := queries.Raw(query, matchID).Bind(nil, s.db, &teamInfo); err != nil {
		return nil, err
	}

	var result []*database.TeamWithValue
	for _, el := range teamInfo {
		team, err := s.GetTeamByID(el.TeamID)
		if err != nil {
			return nil, err
		}

		result = append(result, &database.TeamWithValue{
			Team:  team,
			Value: el.TotalScore,
		})
	}

	return result, nil
}

func (s *DBStore) GetTeamsWithPossession(matchID string) ([]*database.TeamWithValue, error) {
	query := `
select ev.team_id,
       count(1) pass_count
  from match_events ev
 where ev.match_id = $1
   and ev.team_id is not null
   and ev.type in (2,3)
   and ev.status = 1
 group by ev.team_id
 order by pass_count desc
`
	var teamInfo []struct {
		TeamID    string
		PassCount float64
	}

	if err := queries.Raw(query, matchID).Bind(nil, s.db, &teamInfo); err != nil {
		return nil, err
	}

	var totalPassCnt float64
	if len(teamInfo) > 0 {
		totalPassCnt += teamInfo[0].PassCount
		if len(teamInfo) > 1 {
			totalPassCnt += teamInfo[1].PassCount
		}
	}

	// find team info
	var result []*database.TeamWithValue
	for _, el := range teamInfo {
		team, err := s.GetTeamByID(el.TeamID)
		if err != nil {
			return nil, err
		}

		var possession float64
		if totalPassCnt > 0 {
			possession = el.PassCount / totalPassCnt
		}
		result = append(result, &database.TeamWithValue{
			Team:  team,
			Value: possession,
		})
	}

	return result, nil
}

func (s *DBStore) GetTeamsWithShots(matchID string) ([]*database.TeamWithValue, error) {
	query := `
select ev.team_id,
       count(1) shot_count
  from match_events ev
 where ev.match_id = $1
   and ev.team_id is not null
   and ev.type in (43,44,45,53,63)
   and ev.status = 1
 group by ev.team_id
 order by shot_count desc
`
	var teamInfo []struct {
		TeamID    string
		ShotCount float64
	}

	if err := queries.Raw(query, matchID).Bind(nil, s.db, &teamInfo); err != nil {
		return nil, err
	}

	// find team info
	var result []*database.TeamWithValue
	for _, el := range teamInfo {
		team, err := s.GetTeamByID(el.TeamID)
		if err != nil {
			return nil, err
		}

		result = append(result, &database.TeamWithValue{
			Team:  team,
			Value: el.ShotCount,
		})
	}

	return result, nil
}

func (s *DBStore) GetMatchReward(matchID string) (float64, error) {
	rewards, err := schema.MatchRewards(
		qm.Where("match_id = ?", matchID),
	).All(s.db)
	if err != nil {
		return 0, err
	}

	var total float64 = 0
	for _, r := range rewards {
		total += r.Amount
	}

	return total, nil
}

func (s *DBStore) GetMatchWinner(matchID string) (*schema.User, float64, error) {
	lb, err := schema.MatchLeaderboards(
		qm.Where("match_id = ?", matchID),
		qm.Where("score is not null"),
		qm.Load("User"),
		qm.OrderBy("score desc"),
		qm.Limit(1),
	).One(s.db)

	if err != nil {
		return nil, 0, err
	}

	return lb.R.User, lb.Score.Float64, nil
}

func (s *DBStore) UpdateHeadlines(matchID string, headlines schema.MatchHeadlineSlice) error {
	return s.Transaction(func(store database.Store) error {
		_, err := schema.MatchHeadlines(
			qm.Where("match_id = ?", matchID),
		).DeleteAll(s.db)
		if err != nil {
			return stacktrace.Propagate(err, "cannot delete match headlines")
		}

		// insert new headlines
		for _, el := range headlines {
			if err = el.Insert(s.db, boil.Infer()); err != nil {
				return err
			}
		}

		return nil
	})
}

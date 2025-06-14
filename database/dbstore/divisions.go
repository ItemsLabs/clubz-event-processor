package dbstore

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/gameon-app-inc/fanclash-event-processor/types"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetGameWeeks(limit int) (schema.GameWeekSlice, error) {
	return schema.GameWeeks(qm.Limit(limit)).All(s.db)
}
func (s *DBStore) GetGameWeek(id string) (*schema.GameWeek, error) {
	return schema.GameWeeks(qm.Where("id = ?", id), qm.Load("Season")).One(s.db)
}

func (s *DBStore) GetCurrentGameWeek() (*schema.GameWeek, error) {
	return schema.GameWeeks(
		qm.Where("status = ?", "l"),
		qm.Load("Season"),
	).One(s.db)
}

func (s *DBStore) GetCurrentWeekLeaderboard(
	seasonStartAt, seasonEndAt, weekStartAt, weekEndAt time.Time,
	userID string,
) ([]*types.LeaderboardEntry, error) {
	//TODO create index:
	//	create index match_leaderboard_user_id_division_id_game_id_index
	//	 on match_leaderboard (user_id, division_id, game_id);
	query := `
with latest_user_divisions AS (
    select user_id,
           division_id
      from (
        select ud.user_id,
               ud.division_id,
               ROW_NUMBER() OVER (PARTITION BY ud.user_id ORDER BY ud.join_date DESC) AS rn
          from user_divisions ud
      ) t
     where t.rn = 1
),
all_week_matches as (
    select ml.user_id,
           ml.score,
           ml.position,
           row_number() over(partition by ml.user_id order by ml.score desc) as match_rank
      from match_leaderboard ml
     where ml.match_id in (
           select id 
             from matches
            where match_time >= $1 and match_time <= $2
           )
),
top_week_matches as (
    select *
      from all_week_matches
     where match_rank <= 5
),
season_users as (
    select u.id as user_id
      from users u
     where exists(
           select null
             from match_leaderboard ml,
                  matches m
            where ml.match_id = m.id
              and ml.user_id = u.id
              and m.match_time >= $3 and m.match_time <= $4
           )
),
ranked_leaderboard as (
    select season_users.user_id,
           divisions.id as division_id,
           coalesce(divisions.tier, 0) as division_tier,
           sum(coalesce(top_week_matches.score,0)) total_score
      from season_users
      left join top_week_matches on top_week_matches.user_id = season_users.user_id
      left join latest_user_divisions on latest_user_divisions.user_id = season_users.user_id
      left join divisions on divisions.id = latest_user_divisions.division_id
     group by season_users.user_id,
              divisions.id,
              divisions.tier
     order by divisions.tier nulls last,
              total_score desc
)
select ranked_leaderboard.user_id,
       users.name user_name,
       total_score,
       RANK() OVER (PARTITION BY division_tier ORDER BY total_score DESC nulls last) AS rank,
       (users.id = $5) current_user,
       division_tier,
       division_id,
       (select avg(t.position)
          from all_week_matches t
         where t.user_id = ranked_leaderboard.user_id
       ) as week_average_rank,
       (select count(*)
          from all_week_matches t
         where t.user_id = ranked_leaderboard.user_id
       ) as week_matches_played
  from ranked_leaderboard,
       users
 where ranked_leaderboard.user_id = users.id
 order by ranked_leaderboard.division_tier,
          rank`

	spew.Dump(weekStartAt, weekEndAt, seasonStartAt, seasonEndAt, userID)
	var leaderboard []*types.LeaderboardEntry
	err := queries.Raw(query, weekStartAt, weekEndAt, seasonStartAt, seasonEndAt, userID).Bind(nil, s.db, &leaderboard)
	if err != nil {
		return nil, err
	}
	return leaderboard, nil

}

func (s *DBStore) GetConcludedWeekLeaderboard(weekID, userID string) ([]*types.LeaderboardEntry, error) {
	query := `
	SELECT
		ugwh.user_id,
		u.name AS user_name,
		ugwh.week_points as total_score,
		ugwh.week_division_position as rank,
		(ugwh.user_id = $1) AS "current_user",
		ugwh.week_division_tier as division_tier,
		ugwh.week_division_id as division_id,
		ugwh.week_average_rank,
		ugwh.week_matches_played
	FROM user_game_week_histories ugwh
	JOIN 
	    users u ON ugwh.user_id = u.id
	WHERE ugwh.game_week_id = $2
	ORDER BY
     	division_tier, rank;`
	var leaderboard []*types.LeaderboardEntry
	err := queries.Raw(query, userID, weekID).Bind(nil, s.db, &leaderboard)
	if err != nil {
		return nil, err
	}
	return leaderboard, nil

}

func (s *DBStore) GetLatestUserGameWeekHistory(userID string) (*schema.UserGameWeekHistory, error) {
	return schema.UserGameWeekHistories(
		qm.Where("user_id = ?", userID),
		qm.OrderBy("created_at DESC"),
	).One(s.db)
}

func (s *DBStore) ListUserGameWeekHistories(userID string) (schema.UserGameWeekHistorySlice, error) {
	return schema.UserGameWeekHistories(
		qm.Where("user_id = ?", userID),
		qm.OrderBy("created_at DESC"),
	).All(s.db)
}

func (s *DBStore) GetUserGameWeekHistory(userID string, ID string) (*schema.UserGameWeekHistory, error) {
	return schema.UserGameWeekHistories(
		qm.Where("user_id = ?", userID),
		qm.Where("id = ?", ID),
		qm.OrderBy("created_at DESC"),
	).One(s.db)
}

func (s *DBStore) GetDivisionRewards(weekID string) (schema.DivisionRewardSlice, error) {
	return schema.DivisionRewards(
		qm.Where("week_id = ?", weekID),
		qm.Load(schema.DivisionRewardRels.Reward),
	).All(s.db)
}

func (s *DBStore) GetDivisionByID(divisionID string) (*schema.Division, error) {
	return schema.Divisions(qm.Where("id = ?", divisionID)).One(s.db)
}

func (s *DBStore) GetGameWeekDivision(weekID, divisionID string) (*schema.GameWeekDivision, error) {
	return schema.GameWeekDivisions(
		qm.Where("week_id = ?", weekID),
		qm.Where("division_id = ?", divisionID),
	).One(s.db)
}

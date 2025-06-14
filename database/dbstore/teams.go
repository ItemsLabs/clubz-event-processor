package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetTeamByID(teamID string) (*schema.Team, error) {
	return schema.Teams(
		qm.Where("id = ?", teamID),
	).One(s.db)
}

func (s *DBStore) GetAvgTeamScore(teamID string) (float64, error) {
	query := `
select coalesce(avg(match_score),0) avg_match_score
  from (with m as (
          select id
            from matches 
           where (home_team_id = $1 or away_team_id = $1)
             and match_time < current_timestamp
             and status = 'e'
           order by match_time desc
           limit 5
        )
        select sum(ev.points) match_score
          from m,
               match_events ev
         where ev.match_id = m.id
           and ev.team_id = $1
         group by m.id
       ) t
 where match_score is not null
`

	var scoreResult struct {
		AvgMatchScore float64
	}

	if err := queries.Raw(query, teamID).Bind(nil, s.db, &scoreResult); err != nil {
		return 0, err
	}

	return scoreResult.AvgMatchScore, nil
}

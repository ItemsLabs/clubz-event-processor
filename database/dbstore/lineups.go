package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

func (s *DBStore) GetLineupsPickInfo(matchID string) ([]*database.LineupPickInfo, error) {
	query := `
select g.id game_id,
       g.user_id user_id,
       (select array_agg(mp.player_id::text)
          from match_players mp,
               game_picks gp
         where mp.match_id = g.match_id
           and mp.player_id = gp.player_id
           and (mp.from_lineups = false or (mp.from_lineups = true and mp.position = 's'))
           and gp.game_id = g.id
           and gp.ended_at is null
       ) non_lineup_picks
  from games g
 where g.match_id = $1
`

	var result []*database.LineupPickInfo
	if err := queries.Raw(query, matchID).Bind(nil, s.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

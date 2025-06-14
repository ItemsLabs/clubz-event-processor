package handlers

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/palantir/stacktrace"
	"github.com/volatiletech/null/v8"
)

type SortableEventSlice schema.MatchEventSlice

func (s SortableEventSlice) Len() int {
	return len(s)
}
func (s SortableEventSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableEventSlice) Less(i, j int) bool {
	if s[i].Minute == s[j].Minute {
		return s[i].Second < s[j].Second
	}
	return s[i].Minute < s[j].Minute
}

func UpdatePlayedTime(store database.Store, matchID string) error {
	match, err := store.GetMatchByID(matchID)
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot get match players for match %s", matchID))
	}

	// first select match players for a match
	matchPlayers, err := store.GetMatchPlayers(matchID)
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot get match players for match %s", matchID))
	}

	events, err := store.GetMatchEvents(matchID)
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot get match events for match %s", matchID))
	}

	// sort events by minute
	sort.Sort(SortableEventSlice(events))

	var playerTime = make(map[string]int)
	var prevTime int
	var homeActivePlayers = make(map[string]struct{})
	var awayActivePlayers = make(map[string]struct{})
	var activePlayers = make(map[string]struct{})

	for _, ev := range events {
		if ev.Type == ActionLineUp {
			var payload lineupsPayload
			if err := json.Unmarshal([]byte(ev.Payload.String), &payload); err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot unmarshal lineups payload for event %d", ev.ID))
			}

			eventPlayers := make(map[string]struct{})
			for _, el := range payload.Players {
				if el.Position == database.PositionGoalkeeper ||
					el.Position == database.PositionDefender ||
					el.Position == database.PositionMidfielder ||
					el.Position == database.PositionForward {

					eventPlayers[el.ID] = struct{}{}
				}
			}

			if ev.TeamID.String == match.HomeTeamID {
				homeActivePlayers = eventPlayers
			} else if ev.TeamID.String == match.AwayTeamID {
				awayActivePlayers = eventPlayers
			}
			// update active players map
			activePlayers = make(map[string]struct{})
			for playerID := range homeActivePlayers {
				activePlayers[playerID] = struct{}{}
			}
			for playerID := range awayActivePlayers {
				activePlayers[playerID] = struct{}{}
			}
		} else if ev.Type == ActionSubstitution {
			var payload substitutionPayload
			if err := json.Unmarshal([]byte(ev.Payload.String), &payload); err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot unmarshal substitution payload for event %d", ev.ID))
			}

			delete(activePlayers, payload.OutPlayerID)
			activePlayers[payload.InPlayerID] = struct{}{}
		}

		currTime := ev.Minute*60 + ev.Second
		if currTime > 0 {
			timeSinceLastEvent := currTime - prevTime
			if timeSinceLastEvent > 0 {
				for playerID := range activePlayers {
					if _, ok := playerTime[playerID]; !ok {
						playerTime[playerID] = 0
					}
					playerTime[playerID] += timeSinceLastEvent
				}
			}

			prevTime = currTime
		}
	}

	return store.Transaction(func(store database.Store) error {
		for _, mp := range matchPlayers {
			if _, ok := playerTime[mp.PlayerID]; ok {
				mp.PlayedSeconds = null.IntFrom(playerTime[mp.PlayerID])
				if _, err := store.UpdateMatchPlayerPlayedSeconds(mp); err != nil {
					return stacktrace.Propagate(err,
						"cannot update match player %s played seconds %d", mp.ID, mp.PlayedSeconds.Int)
				}
			}
		}

		return nil
	})
}

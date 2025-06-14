package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/gameon-app-inc/fanclash-event-processor/processor"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
)

var (
	puCache     = map[int]map[int]struct{}{}
	puCacheInit = false
)

func NewPointsHandler(store database.Store) *BaseEventHandler {
	initPuCache(store)

	return NewBaseEventHandler(store, []localHandler{
		pointsHandler,
		leaderboardHandler,
	}, false)
}

func initPuCache(store database.Store) {
	if !puCacheInit {
		actions, err := store.GetPowerUpActions()
		if err != nil {
			logrus.WithError(err).Fatal("cannot get power-up actions from database")
		}

		// init power up cache
		for _, pu := range actions {
			if _, ok := puCache[pu.PowerupID]; !ok {
				puCache[pu.PowerupID] = map[int]struct{}{}
			}

			puCache[pu.PowerupID][pu.ActionID] = struct{}{}
		}

		puCacheInit = true
	}
}

func pointsHandler(store database.Store, event *processor.Event) error {
	if event.Points != nil && *event.Points != 0 {
		points := strconv.FormatFloat(*event.Points, 'f', 2, 64)
		eventType := event.Type
		if eventType == 10003 {
			eventType = 45
		} else if eventType == 10004 {
			eventType = 49
		}
		// preliminary exit
		if event.TeamID == nil || event.PlayerID == nil {
			logrus.WithField("id", event.ID).Warn("team_id or player_id is missing for points event")
			return nil
		}

		return store.Transaction(func(store database.Store) error {
			match, err := store.GetMatchByID(event.MatchID)
			if err != nil {
				return stacktrace.Propagate(err,
					fmt.Sprintf("cannot get match by id %s", event.MatchID))
			}

			// do not process points on ended match
			if match.Status == database.MatchStatusEnded {
				logrus.WithField("id", event.ID).Error("point events are not processed on ended matches")

				// set status as ignored
				if err := store.UpdateMatchEventStatus(&schema.MatchEvent{ID: event.ID, Status: database.MatchEventStatusIgnored}); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot set match event %d status to ignored", event.ID))
				}

				return nil
			}

			if event.PlayerID != nil {
				// update match player score
				if err := store.IncMatchPlayerScoreFast(event.MatchID, *event.TeamID, *event.PlayerID, *event.Points); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("GetMatchPlayer failed (match_id: %s, team_id: %s, player_id: %s, score: %f)",
							event.MatchID, *event.TeamID, *event.PlayerID, *event.Points))
				}

				matchTimeGaps := CalculateMatchTimeGaps(match)

				// find real match event id
				matchEvent, err := store.GetMatchEventByMatchEventID(event.MatchID, event.MatchEventID)
				if err != nil {
					return stacktrace.Propagate(err,
						"GetMatchEventByID failed (match_event_id: %d)", event.MatchEventID)
				}

				var picks schema.GamePickSlice
				if matchEvent.HasRealTimestamp {
					logrus.WithFields(logrus.Fields{
						"match_event_id": event.MatchEventID,
						"timestamp":      event.Timestamp,
					}).Info("real timestamp received for event")
					// select picks by minute and second
					picks, err = store.GetActivePicksAtTime(event.MatchID, *event.PlayerID, matchEvent.Timestamp)
					if err != nil {
						return stacktrace.Propagate(err,
							fmt.Sprintf("GetActivePicksAtTime failed (match_id: %s, player_id: %s)", event.MatchID, *event.PlayerID))
					}
				} else {
					logrus.WithFields(logrus.Fields{
						"match_event_id": event.MatchEventID,
						"timestamp":      event.Timestamp,
					}).Info("event don't have real timestamp")
					// select picks by minute and second
					picks, err = store.GetActivePicksAtMinSec(event.MatchID, *event.PlayerID, event.Minute, event.Second)
					if err != nil {
						return stacktrace.Propagate(err,
							fmt.Sprintf("GetActivePicksAtMinSec failed (match_id: %s, player_id: %s)", event.MatchID, *event.PlayerID))
					}
				}

				// no picks, no need for further execution
				if len(picks) == 0 {
					return nil
				}

				// select all powerups
				gamePu, err := store.GetActiveGamePowerUpsAtTime(event.MatchID, matchEvent.Timestamp)
				if err != nil {
					return stacktrace.Propagate(err,
						`GetGamePowerUps failed (GetActiveGamePowerUpsAtTime)`)
				}

				// distribute them by game + position
				gamePuMap := map[string]map[int]*schema.GamePowerup{}
				for _, el := range gamePu {
					if _, ok := gamePuMap[el.GameID]; !ok {
						gamePuMap[el.GameID] = map[int]*schema.GamePowerup{}
					}
					gamePuMap[el.GameID][el.Position] = el
				}

				// calculate power-up end for each power-up
				puEndMap := map[string]time.Time{}
				for _, el := range gamePu {
					end := CalculatePowerUpDuration(
						el.CreatedAt, time.Second*time.Duration(el.Duration), matchTimeGaps)
					puEndMap[el.ID] = end

					// end power-ups
					if event.Timestamp.After(end) {
						el.EndedAt = null.TimeFrom(end)
						if err := store.UpdateEndedAt(el); err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot set ended_at of powerup %s", el.ID))
						}

						// increase game version for each pu
						if err := store.IncGameVersion(el.GameID); err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot increase game %s version", el.GameID))
						}

						player, err := store.GetPlayerByID(*event.PlayerID)
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot get player %s", *event.PlayerID))
						}
						// notify game update
						EventIDString := strconv.Itoa(event.ID)
						var eventName string
						if eventType != 0 {
							eventName, err = store.GetActionNameByActionID(eventType)
							if err != nil {
								return stacktrace.Propagate(err, fmt.Sprintf("cannot get action name by action id %d", eventType))
							}
						}

						minutes := strconv.Itoa(event.Minute)
						seconds := strconv.Itoa(event.Second)
						if err := notifyGameUpdate(store, el.GameID, el.R.Game.UserID, points, points, player.ImageURL.String, EventIDString, minutes, seconds, event.MatchID, player.NormalizedName.String, eventName, 0, 0, ""); err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot notify game %s update", el.GameID))
						}
					}
				}
				// var powerUpName string
				// insert game events
				for _, p := range picks {
					ev := &schema.GameEvent{
						ID:              uuid.New().String(),
						Minute:          event.Minute,
						Second:          event.Second,
						Type:            eventType,
						InitialScore:    *event.Points,
						Score:           *event.Points,
						GameID:          p.GameID,
						PlayerID:        *event.PlayerID,
						TeamID:          *event.TeamID,
						GamePickID:      p.ID,
						MatchEventID:    null.IntFrom(event.ID),
						BoostMultiplier: 1.0,
						NFTImage:        null.StringFrom(p.R.Player.ImageURL.String),
						NFTMultiplier:   1.0,
					}

					// find powerup for that combination of game + position
					if _, ok := gamePuMap[p.GameID]; ok {
						if _, ok := gamePuMap[p.GameID][p.Position]; ok {
							powerup := gamePuMap[p.GameID][p.Position]

							// get calculated pu end
							end := puEndMap[powerup.ID]

							// event is within pu active time
							if !event.Timestamp.After(end) {
								// check whether powerup matches action
								if _, ok := puCache[powerup.PowerupID]; ok {
									if _, ok := puCache[powerup.PowerupID][eventType]; ok {
										// apply powerup to game event
										ev.PowerupID = null.StringFrom(powerup.ID)
										ev.Score = handlePowerupCalculation(event, powerup)
										ev.BoostMultiplier = powerup.Multiplier
									}
								}
							}
						}
					}
					var multiplier float64
					if p.AssignedPlayerID.Valid {
						playerInfo, err := store.GetAssignedPlayersWithNFTDetails([]string{p.AssignedPlayerID.String})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot get player %s", p.AssignedPlayerID.String))
						}
						if len(playerInfo) == 0 {
							return fmt.Errorf("no player found for assigned player ID %s", p.AssignedPlayerID.String)
						}

						player := playerInfo[0].Player
						nft := playerInfo[0].NFT

						var claiming, defence, distribution, dribbling, passing, shooting, stopping float64
						var nftImage string
						switch player.Rarity.String {
						case "common", "Common":
							claiming = nft.CommonClaiming.Float64
							defence = nft.CommonDefence.Float64
							distribution = nft.CommonDistribution.Float64
							dribbling = nft.CommonDribbling.Float64
							passing = nft.CommonPassing.Float64
							shooting = nft.CommonShooting.Float64
							stopping = nft.CommonStopping.Float64
							nftImage = nft.CommonImage.String
						case "uncommon", "Uncommon":
							claiming = nft.UncommonClaiming.Float64
							defence = nft.UncommonDefence.Float64
							distribution = nft.UncommonDistribution.Float64
							dribbling = nft.UncommonDribbling.Float64
							passing = nft.UncommonPassing.Float64
							shooting = nft.UncommonShooting.Float64
							stopping = nft.UncommonStopping.Float64
							nftImage = nft.UncommonImage.String
						case "rare", "Rare":
							claiming = nft.RareClaiming.Float64
							defence = nft.RareDefence.Float64
							distribution = nft.RareDistribution.Float64
							dribbling = nft.RareDribbling.Float64
							passing = nft.RarePassing.Float64
							shooting = nft.RareShooting.Float64
							stopping = nft.RareStopping.Float64
							nftImage = nft.RareImage.String
						case "ultra_rare", "UltraRare":
							claiming = nft.UltraRareClaiming.Float64
							defence = nft.UltraRareDefence.Float64
							distribution = nft.UltraRareDistribution.Float64
							dribbling = nft.UltraRareDribbling.Float64
							passing = nft.UltraRarePassing.Float64
							shooting = nft.UltraRareShooting.Float64
							stopping = nft.UltraRareStopping.Float64
							nftImage = nft.UltraRareImage.String
						case "legendary", "Legendary":
							claiming = nft.LegendaryClaiming.Float64
							defence = nft.LegendaryDefence.Float64
							distribution = nft.LegendaryDistribution.Float64
							dribbling = nft.LegendaryDribbling.Float64
							passing = nft.LegendaryPassing.Float64
							shooting = nft.LegendaryShooting.Float64
							stopping = nft.LegendaryStopping.Float64
							nftImage = nft.LegendaryImage.String
						default:
							return fmt.Errorf("unknown rarity %s for player %s", player.Rarity.String, p.AssignedPlayerID.String)
						}

						fmt.Println("claiming", claiming, "defence", defence, "distribution", distribution, "dribbling", dribbling, "passing", passing, "shooting", shooting, "stopping", stopping)

						eventComplete, err := store.GetActionsByTypeID(eventType)
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot get action %d", eventType))
						}
						if len(eventComplete) == 0 {
							return fmt.Errorf("no action found for type %d", eventType)
						}

						if *event.Points > 0 {
							for _, action := range eventComplete {
								switch action.NFTCategory.String {
								case "distribution":
									multiplier = distribution
								case "shooting":
									multiplier = shooting
								case "passing":
									multiplier = passing
								case "dribbling":
									multiplier = dribbling
								case "defence":
									multiplier = defence
								case "disciplinary":
									multiplier = 1.0
								case "stopping":
									multiplier = stopping
								case "claiming":
									multiplier = claiming
								default:
									multiplier = 1.0
								}
								ev.Score *= multiplier
								ev.NFTMultiplier = multiplier
								ev.NFTImage = null.StringFrom(nftImage)
							}
						}

					}

					err := store.InsertGameEvent(ev)
					if err != nil {
						return stacktrace.Propagate(err, "cannot insert game events")
					}
					boostMultiplier := ev.Score / ev.InitialScore
					// notify game update
					poweredUpScore := strconv.FormatFloat(ev.Score, 'f', 2, 64)
					initialScore := strconv.FormatFloat(ev.InitialScore, 'f', 2, 64)

					minutes := strconv.Itoa(event.Minute)
					seconds := strconv.Itoa(event.Second)
					actionName, err := store.GetActionNameByActionID(eventType)
					if err != nil {
						return stacktrace.Propagate(err, fmt.Sprintf("cannot get action %d name", eventType))
					}
					if err = notifyGameUpdate(store, ev.GameID, p.R.Game.UserID, poweredUpScore, initialScore, p.R.Player.ImageURL.String, ev.ID, minutes, seconds, event.MatchID, p.R.Player.NormalizedName.String, actionName, multiplier, boostMultiplier, ev.NFTImage.String); err != nil {
						return stacktrace.Propagate(err, fmt.Sprintf("cannot notify game %s update", p.R.Game.ID))
					}
				}
			}

			return nil
		})
	}

	return nil
}

const (
	shieldPowerupName  = "shield"  // Shield: Ignore negative points
	reversePowerupName = "reverse" // Reverse: All negative points are treated as positive
)

func leaderboardHandler(store database.Store, event *processor.Event) error {
	// recalculate leaderboard
	if event.Points != nil && *event.Points != 0 {
		// this shouldn't break processing
		if err := store.UpdateMatchLeaderboard(event.MatchID); err != nil {
			logrus.WithError(err).Error("cannot update match leaderboard")
			return nil
		}

		GetDebouncedSendMatchLeaderboard()(event.MatchID)
	}

	return nil
}

func handlePowerupCalculation(event *processor.Event, powerup *schema.GamePowerup) float64 {
	score := *event.Points
	if powerup.R.Powerup.Conditions.Valid && powerup.R.Powerup.Conditions.JSON != nil {
		var cond []conditionPayload
		_ = json.Unmarshal(powerup.R.Powerup.Conditions.JSON, &cond)
		// TODO: sort the conditions by order if field exists
		for _, c := range cond {
			// TODO: Use Entity, Field + Expression + Value to determine the value of score of the condition
			// for the time being We can use "Name" here since they are general conditions
			if score < 0 {
				switch c.Name {
				case shieldPowerupName:
					score = 0
				case reversePowerupName:
					score = score * -1
				}
			}
			// TODO: Add support for actual conditional calculations with expressions (see conditionPayload for layout)
		}
		// TODO: Use ".Expression" for more complex rules maybe, parse conditionals and nested ones.
	}
	if score > 0 {
		score *= powerup.Multiplier
	}
	return score
}

package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/gameon-app-inc/fanclash-event-processor/processor"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
)

const (
	// system ones
	ActionPeriodStart  = 10001
	ActionPeriodEnd    = 10002
	ActionGoal         = 10003
	ActionSelfGoal     = 10004
	ActionLineUp       = 10005
	ActionMatchEnd     = 10006
	ActionSubstitution = 10007
	ActionCancel       = 30000

	MaxPuPerGame = 4
)

func NewSystemHandler(store database.Store) *BaseEventHandler {
	return NewBaseEventHandler(store, []localHandler{
		handleMatchTimeChange,
		handlePeriodChange,
		handleGoal,
		handleLineUps,
		handleMatchEnd,
		handleSubstitution,
		handleCancel,
	}, false)
}

func handleMatchTimeChange(store database.Store, event *processor.Event) error {
	if event.Minute > 0 && event.Second > 0 {
		if err := store.UpdateMatchTime(event.MatchID, event.Minute, event.Second); err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot update match %s minute and second", event.MatchID))
		}
	}

	return nil
}

func handlePeriodChange(store database.Store, event *processor.Event) error {
	// period ids:
	// 1 - first half
	// 2 - second half
	// 3 - first period of extra time
	// 4 - second period of extra time
	// 5 - penalty shoot out
	// 14 - post-game
	// 15 - pre-game
	// 16 - pre-match

	// period start or period end
	if event.Type == ActionPeriodStart || event.Type == ActionPeriodEnd {
		// send headlines
		GetDebouncedSendMatchHeadlines()(event.MatchID)

		type eventPayload struct {
			PeriodID int `json:"period_id"`
		}

		match, err := store.GetMatchByIDWithTeams(event.MatchID)
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot match by id %s", event.MatchID))
		}

		// ignore status change on ended match
		if match.Status == database.MatchStatusEnded {
			// set status as ignored
			if err := store.UpdateMatchEventStatus(&schema.MatchEvent{ID: event.ID, Status: database.MatchEventStatusIgnored}); err != nil {
				return stacktrace.Propagate(err,
					fmt.Sprintf("cannot set match event %d status to ignored", event.ID))
			}

			return nil
		}

		prevStatus := match.Status
		sStartAlreadySet := match.SStart.Valid

		// unmarshal payload
		if event.Payload != nil {
			var payload eventPayload
			if err := json.Unmarshal([]byte(*event.Payload), &payload); err != nil {
				return stacktrace.Propagate(err, "cannot unmarshal payload")
			}

			if payload.PeriodID != 1 &&
				payload.PeriodID != 2 &&
				payload.PeriodID != 3 &&
				payload.PeriodID != 4 &&
				payload.PeriodID != 5 &&
				payload.PeriodID != 14 &&
				payload.PeriodID != 15 &&
				payload.PeriodID != 16 {

				// set status as ignored
				if err := store.UpdateMatchEventStatus(&schema.MatchEvent{ID: event.ID, Status: database.MatchEventStatusIgnored}); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot set match event %d status to ignored", event.ID))
				}

				//return fmt.Errorf("unknown period %d", payload.PeriodID)
				return nil
			}

			var shouldSendFirstTimeSummary bool
			var kickOffEvent bool
			var updateFields []string
			if event.Type == ActionPeriodStart {
				if payload.PeriodID == 1 {
					kickOffEvent = true
					match.FStart = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodFirstHalf
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "f_start", "period", "status")
				} else if payload.PeriodID == 2 {
					match.SStart = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodSecondHalf
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "s_start", "period", "status")
					// we should sent first time summary if this is first `Second Time Start` event in a game
					shouldSendFirstTimeSummary = sStartAlreadySet
				} else if payload.PeriodID == 3 {
					match.X1Start = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodFirstExt
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "x1_start", "period", "status")
				} else if payload.PeriodID == 4 {
					match.X2Start = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodSecondExt
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "x2_start", "period", "status")
				} else if payload.PeriodID == 5 {
					match.PStart = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodPenalties
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "p_start", "period", "status")
				} else if payload.PeriodID == 14 {
					match.Period = database.MatchPeriodPostGame
					updateFields = append(updateFields, "period")
				}
			} else if event.Type == ActionPeriodEnd {
				if payload.PeriodID == 1 {
					match.FEnd = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodHalfTime
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "f_end", "period", "status")
				} else if payload.PeriodID == 2 {
					match.SEnd = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodBreakX1
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "s_end", "period", "status")
				} else if payload.PeriodID == 3 {
					match.X1End = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodBreakX2
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "x1_end", "period", "status")
				} else if payload.PeriodID == 4 {
					match.X2End = null.TimeFrom(event.Timestamp)
					match.Period = database.MatchPeriodBreakP
					match.Status = database.MatchStatusGame
					updateFields = append(updateFields, "x2_end", "period", "status")
				} else if payload.PeriodID == 5 {
					match.PEnd = null.TimeFrom(event.Timestamp)
					match.Status = database.MatchStatusGame
					match.Period = database.MatchPeriodPostGame
					updateFields = append(updateFields, "p_end", "period", "status")
				} else if payload.PeriodID == 14 {
					match.Period = database.MatchPeriodPostGame
					updateFields = append(updateFields, "period")
				}
			}

			return store.Transaction(func(store database.Store) error {
				// pass match status to game
				if err = store.UpdateMatchGamesStatus(match.ID, gameStatusFromMatchStatus(match.Status)); err != nil {
					return stacktrace.Propagate(err, "cannot update game status")
				}

				// some periods do not make update
				if len(updateFields) > 0 {
					if _, err := store.UpdateMatch(match, updateFields); err != nil {
						return stacktrace.Propagate(err, "cannot update match periods")
					}

					if err := notifyMatchUpdate(store, match); err != nil {
						return stacktrace.Propagate(err, "cannot notify match update")
					}
				}

				// since we receive 2 period start events: 1 for each team
				// we should guard against double push for kick-off
				// so if match status changed to game and we receive kick-off event
				// then send notification
				if prevStatus != database.MatchStatusGame && kickOffEvent {
					title := "Match started ⚽"
					message := fmt.Sprintf("%s vs %s has kicked off!", GetTeamName(match.R.HomeTeam), GetTeamName(match.R.AwayTeam))
					payload := map[string]string{
						"match_id": match.ID,
						"type":     "kick_off",
					}
					if err = sendPushForAllUsers(store, title, message, match.ID, payload); err != nil {
						return stacktrace.Propagate(err, "cannot send kick-off push notification")
					}

					// on kick-off event game premium flag should be synced with user premium flag
					if err = store.SyncGamePremiumFlags(match.ID); err != nil {
						return stacktrace.Propagate(err, "cannot sync game premium flags")
					}
				}

				// check whether we should send first time summary
				if shouldSendFirstTimeSummary {
					// calculate powerup usage
					puUsage, err := store.GetPowerUpCountForGame(match.ID)
					if err != nil {
						return stacktrace.Propagate(err, "GetPowerUpCountForGame failed (match_id: %s)", match.ID)
					}

					// get whole leaderboard and send push for each user
					leaderboard, err := store.GetMatchLeaderboard(match.ID)
					if err != nil {
						return stacktrace.Propagate(err, "GetMatchLeaderboard failed (match_id: %s)", match.ID)
					}

					if len(leaderboard) > 0 {
						leader := leaderboard[0]

						title := "2nd half started ⚽"
						for _, entry := range leaderboard {
							var message string
							if entry.Position.Valid && entry.Score.Valid {
								tpl := "You're currently ranked %d out of %d with %d pts. Global leader is %s with %d pts. You have %d powerups left."

								usedPu := puUsage[entry.GameID]
								message = fmt.Sprintf(tpl, entry.Position.Int, len(leaderboard), int(entry.Score.Float64),
									leader.R.User.Name, int(leader.Score.Float64), MaxPuPerGame-usedPu)
							}

							err = sendPush(store, entry.UserID, match.ID, title, message, map[string]string{
								"match_id": match.ID,
								"type":     "half_time",
							})
							if err != nil {
								return stacktrace.Propagate(err, "cannot send match half-time push for user %s", entry.UserID)
							}
						}
					}

				}

				return nil
			})
		}
	}

	return nil
}

func handleGoal(store database.Store, event *processor.Event) error {
	// goal or self goal
	if event.Type == ActionGoal || event.Type == ActionSelfGoal {
		match, err := store.GetMatchByIDWithTeams(event.MatchID)
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot get match %s", event.MatchID))
		}

		if event.TeamID == nil {
			logrus.Error("team_id is missing for goal-related event")
			return nil
		}

		if *event.TeamID != match.HomeTeamID && *event.TeamID != match.AwayTeamID {
			return fmt.Errorf("unknown team_id %s", *event.TeamID)
		}

		return store.Transaction(func(store database.Store) error {
			if (event.Type == ActionGoal && match.HomeTeamID == *event.TeamID) ||
				(event.Type == ActionSelfGoal && match.AwayTeamID == *event.TeamID) { // home score

				if _, err = store.IncHomeScore(match); err != nil {
					return stacktrace.Propagate(err, "cannot increase home score")
				}
			} else if (event.Type == ActionGoal && match.AwayTeamID == *event.TeamID) ||
				(event.Type == ActionSelfGoal && match.HomeTeamID == *event.TeamID) { // away score

				if _, err = store.IncAwayScore(match); err != nil {
					return stacktrace.Propagate(err, "cannot increase away score")
				}
			} else {
				return errors.New("UNKNOWN")
			}

			if err := notifyMatchUpdate(store, match); err != nil {
				return stacktrace.Propagate(err, "cannot notify match update")
			}

			return nil
		})
	}

	return nil
}

func handleLineUps(store database.Store, event *processor.Event) error {
	if event.Type == ActionLineUp {
		return store.Transaction(func(store database.Store) error {
			if event.TeamID == nil {
				return errors.New("lineups action don't have team_id")
			}

			if event.Payload == nil {
				return errors.New("lineups actions has empty payload")
			}

			// unmarshal payload
			var payload lineupsPayload
			if err := json.Unmarshal([]byte(*event.Payload), &payload); err != nil {
				return stacktrace.Propagate(err, "cannot unmarshal payload")
			}

			// insert match players
			for _, player := range payload.Players {
				// get corresponding match player
				mp, err := store.GetMatchPlayer(event.MatchID, *event.TeamID, player.ID)
				if err == sql.ErrNoRows {
					// insert new match player
					if _, err := store.InsertMatchPlayer(&schema.MatchPlayer{
						ID:           uuid.New().String(),
						MatchID:      event.MatchID,
						TeamID:       *event.TeamID,
						PlayerID:     player.ID,
						Position:     null.StringFrom(player.Position),
						JerseyNumber: null.IntFrom(player.JerseyNumber),
						FromLineups:  true,
					}); err != nil {
						return stacktrace.Propagate(err, "cannot insert match player")
					}
				} else if err != nil {
					// some unknown error
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot get match player for (match_id: %s, team_id: %s, player_id: %s)",
							event.MatchID, *event.TeamID, player.ID))
				} else {
					// update position and jersey_number
					mp.Position = null.StringFrom(player.Position)
					mp.JerseyNumber = null.IntFrom(player.JerseyNumber)
					mp.FromLineups = true
					if _, err = store.UpdateMatchPlayer(mp); err != nil {
						return stacktrace.Propagate(err, "cannot update match player")
					}
				}
			}

			if cnt, err := calculatePrecedingLineUpActions(store, event.MatchID, event.MatchEventID); err != nil {
				return stacktrace.Propagate(err, "cannot calculate preceding line up actions")
			} else if cnt == 0 {
				logrus.Info("skipping lineups processing, cause it's first event")
				// We stop processing of lineups event if this is first event,
				// cause lineup events are created separately for each team.
				// And we want to start actual processing once we have lineups for both teams
				return nil
			} else {
				logrus.WithField("cnt", cnt).Info("actually processing lineups event")
			}

			// select all match players
			players, err := store.GetMatchPlayers(event.MatchID)
			if err != nil {
				return stacktrace.Propagate(err, "cannot get match players")
			}

			// check number of teams with match players
			var teamID = map[string]struct{}{}
			for _, p := range players {
				if p.FromLineups {
					teamID[p.TeamID] = struct{}{}
				}
			}

			// get match by id
			match, err := store.GetMatchByIDWithTeams(event.MatchID)
			if err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot get match by id %s", event.MatchID))
			}

			// ignore if match already has lineups
			if !match.HasLineups {
				match.HasLineups = true
				if match.Status == database.MatchStatusWaiting {
					match.Status = database.MatchStatusLineups
				}

				if _, err = store.UpdateMatch(match, []string{"has_lineups", "status"}); err != nil {
					return stacktrace.Propagate(err, "UpdateMatch failed")
				}

				// pass lineups status to game
				if err = store.UpdateMatchGamesStatus(match.ID, gameStatusFromMatchStatus(match.Status)); err != nil {
					return stacktrace.Propagate(err, "UpdateMatchGamesStatus failed")
				}

				// notify match update
				if err := notifyMatchUpdate(store, match); err != nil {
					return stacktrace.Propagate(err, "cannot notify match update")
				}

				if err := sendLineupNotifications(store, match); err != nil {
					return stacktrace.Propagate(err, "sendLineupNotifications failed")
				}
			}

			return nil
		})
	}

	return nil
}
func handleMatchEnd(store database.Store, event *processor.Event) error {
	if event.Type == ActionMatchEnd {
		if err := store.UpdateMatchLeaderboard(event.MatchID); err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot update match leaderboard (%s)", event.MatchID))
		}
		gameWeek, err := store.GetCurrentGameWeek()
		if err != nil {
			return stacktrace.Propagate(err, "GetCurrentGameWeek failed")
		}
		// get match by id
		match, err := store.GetMatchByIDWithTeams(event.MatchID)
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot get match by id %s", event.MatchID))
		}

		// already ended, stop execution
		if match.Status == database.MatchStatusEnded {
			// set status as ignored
			if err := store.UpdateMatchEventStatus(&schema.MatchEvent{ID: event.ID, Status: database.MatchEventStatusIgnored}); err != nil {
				return stacktrace.Propagate(err,
					fmt.Sprintf("cannot set match event %d status to ignored", event.ID))
			}

			return nil
		}

		match.Status = database.MatchStatusEnded
		match.MatchEnd = null.TimeFrom(event.Timestamp)

		if err := UpdatePlayedTime(store, event.MatchID); err != nil {
			logrus.WithField("match_id", event.MatchID).WithError(err).Error("cannot update played time for match")
		}

		return store.Transaction(func(store database.Store) error {
			if _, err = store.UpdateMatch(match, []string{"match_end", "status"}); err != nil {
				return stacktrace.Propagate(err, "cannot update match status to ended")
			}

			if err := notifyMatchUpdate(store, match); err != nil {
				return stacktrace.Propagate(err, "cannot notify match update")
			}

			// pass status to games
			if err = store.UpdateMatchGamesStatus(match.ID, gameStatusFromMatchStatus(match.Status)); err != nil {
				return stacktrace.Propagate(err, "cannot update game status")
			}

			// get whole leaderboard and send push for each user
			leaderboard, err := store.GetMatchLeaderboard(match.ID)
			if err != nil {
				return stacktrace.Propagate(err, "GetMatchLeaderboard failed (match_id: %s)", match.ID)
			}

			title := "Full-time 🏁"
			for _, entry := range leaderboard {
				var message string
				if entry.Position.Valid {
					message = fmt.Sprintf("You finished %d out of %d for %s vs %s",
						entry.Position.Int, len(leaderboard), GetTeamName(match.R.HomeTeam), GetTeamName(match.R.AwayTeam))
				} else {
					message = fmt.Sprintf("You finished %s vs %s",
						GetTeamName(match.R.HomeTeam), GetTeamName(match.R.AwayTeam))
				}

				err = sendPush(store, entry.UserID, match.ID, title, message, map[string]string{
					"match_id": match.ID,
					"type":     "full_time",
				})
				if err != nil {
					return stacktrace.Propagate(err, "cannot send match full-time push for user %d", entry.ID)
				}
			}

			if !match.Rewarded {
				// get match rewards
				rewards, err := store.GetMatchRewards(match.ID)
				if err != nil {
					return stacktrace.Propagate(err, "cannot get match rewards for")
				}

				// get match players
				players, err := store.GetMatchLeaderboard(match.ID)
				if err != nil {
					return stacktrace.Propagate(err, "cannot get top match leaderboard")
				}

				// select user info for each of top players and check which of them are premium users
				var premiumUsers = make(map[string]struct{})
				for _, player := range players {
					if player.R.Game.SubscriptionTier == database.SubscriptionTierPremium {
						premiumUsers[player.UserID] = struct{}{}
					}
				}

				// calculate rewards for each player
				userRewards := CalculateRewards(rewards, players)
				for _, r := range userRewards {
					var text = fmt.Sprintf("reward for %d position in match %s vs %s",
						r.Position, GetTeamName(match.R.HomeTeam), GetTeamName(match.R.AwayTeam))
					GameID, err := store.GetGameByUserIDMatchID(r.UserID, match.ID)
					if err != nil {
						return stacktrace.Propagate(err, fmt.Sprintf("cannot get game by user id: %s, and match id: %s", r.UserID, match.ID))
					}
					user, err := store.GetUserByID(r.UserID)
					if err != nil {
						return stacktrace.Propagate(err, fmt.Sprintf("cannot get user by id: %s", r.UserID))
					}
					// Handle credits
					if r.Reward.Credits > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Amount:     r.Reward.Credits,
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "v", // Virtual currency
							Delivered:  true,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create credit transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
					}

					// Handle game tokens
					if r.Reward.GameToken > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Amount:     r.Reward.GameToken,
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "g", // Game token
							Delivered:  false,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create game token transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
					}

					// Handle LAPT tokens
					if r.Reward.LaptToken > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Amount:     r.Reward.LaptToken,
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "l", // LAPT token
							Delivered:  false,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create LAPT token transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
					}

					// Handle event tickets
					if r.Reward.EventTickets > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Quantity:   int64(r.Reward.EventTickets),
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "e", // Event ticket
							Delivered:  false,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:           uuid.New().String(),
							EventTickets: r.Reward.EventTickets,
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Ticket(s)", strconv.Itoa(int(r.Reward.EventTickets)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.digitaloceanspaces.com/notification-icons/Tickets.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
						// Send a slack alert with the user and the prize that he has won
						slackMessage := fmt.Sprintf(
							"User ID: %s\nName: %s\nEmail: %s\nReal Name: %s\nHas won %d tickets",
							user.ID,
							user.Name,
							user.Email.String,    // Since Email is a null.String, use .String
							user.RealName.String, // Similarly handle RealName
							r.Reward.EventTickets,
						)
						// Send the alert to Slack
						SendSlackAlert(slackMessage)
					}

					// Handle balls
					if r.Reward.Balls > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Quantity:   int64(r.Reward.Balls),
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "b", // Balls
							Delivered:  false,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create ball transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:        uuid.New().String(),
							Ball:      r.Reward.Balls,
							CreatedAt: null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Ball", strconv.Itoa(int(r.Reward.Balls)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.digitaloceanspaces.com/notification-icons/Frame.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create Ball AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
						slackMessage := fmt.Sprintf(
							"User ID: %s\nName: %s\nEmail: %s\nReal Name: %s\nHas won %d balls",
							user.ID,
							user.Name,
							user.Email.String,    // Since Email is a null.String, use .String
							user.RealName.String, // Similarly handle RealName
							r.Reward.Balls,
						)
						// Send the alert to Slack
						SendSlackAlert(slackMessage)
					}

					// Handle signed balls
					if r.Reward.SignedBalls > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Quantity:   int64(r.Reward.SignedBalls),
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "a", // Signed Balls
							Delivered:  false,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed ball transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:         uuid.New().String(),
							SignedBall: r.Reward.SignedBalls,
							CreatedAt:  null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Signed Ball", strconv.Itoa(int(r.Reward.SignedBalls)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.digitaloceanspaces.com/notification-icons/Frame.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed ball AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
						slackMessage := fmt.Sprintf(
							"User ID: %s\nName: %s\nEmail: %s\nReal Name: %s\nHas won %d signed balls",
							user.ID,
							user.Name,
							user.Email.String,    // Since Email is a null.String, use .String
							user.RealName.String, // Similarly handle RealName
							r.Reward.SignedBalls,
						)
						// Send the alert to Slack
						SendSlackAlert(slackMessage)
					}

					// Handle shirts
					if r.Reward.Shirts > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Quantity:   int64(r.Reward.Shirts),
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "s", // Shirts
							Delivered:  false,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create shirt transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:        uuid.New().String(),
							Shirt:     r.Reward.Shirts,
							CreatedAt: null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Shirt", strconv.Itoa(int(r.Reward.Shirts)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.digitaloceanspaces.com/notification-icons/Shirt.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create Shirt AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
						slackMessage := fmt.Sprintf(
							"User ID: %s\nName: %s\nEmail: %s\nReal Name: %s\nHas won %d shirts",
							user.ID,
							user.Name,
							user.Email.String,    // Since Email is a null.String, use .String
							user.RealName.String, // Similarly handle RealName
							r.Reward.Shirts,
						)
						// Send the alert to Slack
						SendSlackAlert(slackMessage)
					}

					// Handle signed shirts
					if r.Reward.SignedShirts > 0 {
						_, err := store.InsertUserTransaction(&schema.Transaction{
							ID:         uuid.New().String(),
							UserID:     r.UserID,
							Quantity:   int64(r.Reward.SignedShirts),
							Text:       null.StringFrom(text),
							MatchID:    null.StringFrom(match.ID),
							ObjectType: "h", // Signed Shirts
							Delivered:  false,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt transaction for user: %s, and match: %s", r.UserID, match.ID))
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:          uuid.New().String(),
							SignedShirt: r.Reward.SignedShirts,
							CreatedAt:   null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Signed shirt", strconv.Itoa(int(r.Reward.SignedShirts)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.digitaloceanspaces.com/notification-icons/Shirt.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
						slackMessage := fmt.Sprintf(
							"User ID: %s\nName: %s\nEmail: %s\nReal Name: %s\nHas won %d signed shirts",
							user.ID,
							user.Name,
							user.Email.String,    // Since Email is a null.String, use .String
							user.RealName.String, // Similarly handle RealName
							r.Reward.SignedShirts,
						)
						// Send the alert to Slack
						SendSlackAlert(slackMessage)
					}

					// Handle card pack rewards
					if r.Reward.SeasonPack1 > 0 {
						// _, err := store.InsertAssignedCardPack(&schema.AssignedCardPack{
						// 	ID:                 uuid.New().String(),
						// 	UserID:             r.UserID,
						// 	CardPackTypeID:     "f9adef06-85fd-498b-83f1-86ac25a16367", // Pack ID for kickoff_pack_1
						// 	StoreTransactionID: null.StringFrom(uuid.New().String()),   // Use the transaction ID or generate a new one
						// 	Opened:             false,
						// })
						// if err != nil {
						// 	return stacktrace.Propagate(err, fmt.Sprintf("cannot create kickoff pack 1 reward for user: %s, and match: %s", r.UserID, match.ID))
						// }
						for i := 0; i < r.Reward.SeasonPack1; i++ {
							_, err = store.InsertUserTransaction(&schema.Transaction{
								ID:         uuid.New().String(),
								UserID:     r.UserID,
								Quantity:   int64(r.Reward.SeasonPack1),
								Text:       null.StringFrom(text),
								MatchID:    null.StringFrom(match.ID),
								ObjectType: "4",
								Delivered:  false,
							})
							if err != nil {
								return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt transaction for user: %s, and match: %s", r.UserID, match.ID))
							}
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:          uuid.New().String(),
							SeasonPack1: r.Reward.SeasonPack1,
							CreatedAt:   null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Card Pack", strconv.Itoa(int(r.Reward.SeasonPack1)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.cdn.digitaloceanspaces.com/notification-icons/Packs.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
					}

					if r.Reward.SeasonPack2 > 0 {
						// _, err := store.InsertAssignedCardPack(&schema.AssignedCardPack{
						// 	ID:                 uuid.New().String(),
						// 	UserID:             r.UserID,
						// 	CardPackTypeID:     "d15c5248-cdfe-45bc-aec6-83dbd2814965", // Pack ID for kickoff_pack_2
						// 	StoreTransactionID: null.StringFrom(uuid.New().String()),   // Use the transaction ID or generate a new one
						// 	Opened:             false,
						// })
						// if err != nil {
						// 	return stacktrace.Propagate(err, fmt.Sprintf("cannot create kickoff pack 2 reward for user: %s, and match: %s", r.UserID, match.ID))
						// }
						for i := 0; i < r.Reward.SeasonPack2; i++ {
							_, err = store.InsertUserTransaction(&schema.Transaction{
								ID:         uuid.New().String(),
								UserID:     r.UserID,
								Quantity:   int64(r.Reward.SeasonPack2),
								Text:       null.StringFrom(text),
								MatchID:    null.StringFrom(match.ID),
								ObjectType: "5",
								Delivered:  false,
							})
							if err != nil {
								return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt transaction for user: %s, and match: %s", r.UserID, match.ID))
							}
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:          uuid.New().String(),
							SeasonPack2: r.Reward.SeasonPack2,
							CreatedAt:   null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Card Pack", strconv.Itoa(int(r.Reward.SeasonPack2)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.cdn.digitaloceanspaces.com/notification-icons/Packs.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
					}

					if r.Reward.SeasonPack3 > 0 {
						// _, err := store.InsertAssignedCardPack(&schema.AssignedCardPack{
						// 	ID:                 uuid.New().String(),
						// 	UserID:             r.UserID,
						// 	CardPackTypeID:     "b4825f20-c611-4316-98ee-fc9433d9cd00", // Pack ID for kickoff_pack_3
						// 	StoreTransactionID: null.StringFrom(uuid.New().String()),   // Use the transaction ID or generate a new one
						// 	Opened:             false,
						// })
						// if err != nil {
						// 	return stacktrace.Propagate(err, fmt.Sprintf("cannot create kickoff pack 3 reward for user: %s, and match: %s", r.UserID, match.ID))
						// }
						for i := 0; i < r.Reward.SeasonPack3; i++ {
							_, err = store.InsertUserTransaction(&schema.Transaction{
								ID:         uuid.New().String(),
								UserID:     r.UserID,
								Quantity:   int64(r.Reward.SeasonPack3),
								Text:       null.StringFrom(text),
								MatchID:    null.StringFrom(match.ID),
								ObjectType: "6",
								Delivered:  false,
							})
							if err != nil {
								return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt transaction for user: %s, and match: %s", r.UserID, match.ID))
							}
						}
						Reward, err := store.CreateReward(&schema.Reward{
							ID:          uuid.New().String(),
							SeasonPack3: r.Reward.SeasonPack3,
							CreatedAt:   null.TimeFrom(event.Timestamp).Time,
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create event ticket Reward for user: %s, and match: %s", r.UserID, match.ID))
						}
						message := fmt.Sprintf("You won %s Card Pack", strconv.Itoa(int(r.Reward.SeasonPack3)))
						_, err = store.CreateAppInbox(&schema.AppInbox{
							ID:           uuid.New().String(),
							UserID:       null.StringFrom(r.UserID),
							Title:        message,
							Description:  message,
							Category:     "claim_prize",
							ImageURL:     null.StringFrom("https://laliga.ams3.cdn.digitaloceanspaces.com/notification-icons/Packs.png"),
							Priority:     "Medium",
							CreatedAt:    null.TimeFrom(event.Timestamp).Time,
							UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
							Read:         false,
							Claimed:      false,
							MatchIDID:    null.StringFrom(match.ID),
							RewardID:     null.StringFrom(Reward.ID),
							GameWeekIDID: null.StringFrom(gameWeek.ID),
							Clear:        false,
							GameID:       null.StringFrom(GameID.ID),
						})
						if err != nil {
							return stacktrace.Propagate(err, fmt.Sprintf("cannot create signed shirt AppInbox for user: %s, and match: %s", r.UserID, match.ID))
						}
					}

					message := fmt.Sprintf("You finished %d out of %d for %s vs %s",
						r.Position, len(leaderboard), GetTeamName(match.R.HomeTeam), GetTeamName(match.R.AwayTeam))

					Reward, err := store.CreateReward(&schema.Reward{
						ID:           uuid.New().String(),
						Credits:      r.Reward.Credits,
						GameToken:    r.Reward.GameToken,
						LaptToken:    r.Reward.LaptToken,
						EventTickets: r.Reward.EventTickets,
						Ball:         r.Reward.Balls,
						SignedBall:   r.Reward.SignedBalls,
						Shirt:        r.Reward.Shirts,
						SignedShirt:  r.Reward.SignedShirts,
						CreatedAt:    null.TimeFrom(event.Timestamp).Time,
					})

					if err != nil {
						return stacktrace.Propagate(err, fmt.Sprintf("cannot create reward for user: %s, and match: %s", r.UserID, match.ID))
					}

					_, err = store.CreateAppInbox(&schema.AppInbox{
						ID:           uuid.New().String(),
						UserID:       null.StringFrom(r.UserID),
						Title:        message,
						Description:  message,
						Category:     "match_end",
						ImageURL:     null.StringFrom("https://laliga.ams3.digitaloceanspaces.com/notification-icons/whistle.png"),
						Priority:     "Medium",
						CreatedAt:    null.TimeFrom(event.Timestamp).Time,
						UpdatedAt:    null.TimeFrom(event.Timestamp).Time,
						Read:         false,
						Claimed:      false,
						MatchIDID:    null.StringFrom(match.ID),
						GameWeekIDID: null.StringFrom(gameWeek.ID),
						Clear:        false,
						RewardID:     null.StringFrom(Reward.ID),
						GameID:       null.StringFrom(GameID.ID),
					})
					if err != nil {
						return stacktrace.Propagate(err, fmt.Sprintf("cannot create credit transaction for user: %s, and match: %s", r.UserID, match.ID))
					}
					// // Link the transaction ID back to the leaderboard
					// err = store.SetTransactionIDForMatchLeaderboard(r.UserID, match.ID, r.TransactionID)
					// if err != nil {
					// 	return stacktrace.Propagate(err, fmt.Sprintf("cannot update leaderboard for user: %s, and match: %s", r.UserID, match.ID))
					// }
				}

				// mark match as rewarded
				match.Rewarded = true
				if _, err = store.UpdateMatch(match, []string{"rewarded"}); err != nil {
					return stacktrace.Propagate(err, fmt.Sprintf("cannot set rewarded flag for match %s", match.ID))
				}
			}

			return nil
		})
	}

	return nil
}
func handleSubstitution(store database.Store, event *processor.Event) error {
	// turned off Substitution fanclash logic
	if true {
		return nil
	}

	if event.Type == ActionSubstitution {
		var payload substitutionPayload

		// unmarshal payload
		if err := json.Unmarshal([]byte(*event.Payload), &payload); err != nil {
			return stacktrace.Propagate(err, "cannot unmarshal payload")
		}

		// get in player info
		inPlayer, err := store.GetPlayerByID(payload.InPlayerID)
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot find in player %s", payload.InPlayerID))
		}

		// get out player info
		outPlayer, err := store.GetPlayerByID(payload.OutPlayerID)
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot find out player %s", payload.OutPlayerID))
		}

		return store.Transaction(func(store database.Store) error {
			match, err := store.GetMatchByID(event.MatchID)
			if err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot get match by id %s", event.MatchID))
			}

			// player who is going to leave this match
			outMp, err := store.GetOrCreateMatchPlayer(event.MatchID, *event.TeamID, payload.OutPlayerID)
			if err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot get match player for out player %s", payload.OutPlayerID))
			}

			// player who is coming into match
			inMp, err := store.GetOrCreateMatchPlayer(event.MatchID, *event.TeamID, payload.InPlayerID)
			if err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot get match player for in player %s", payload.InPlayerID))
			}

			// TODO: find more proper way of handling this
			// if out player already has "s" position it means that this is duplicate swap event
			// it should be ignored
			if outMp.Position.String == "s" {
				return nil
			}

			// by default new position is position of outgoing player
			newPosition := outMp.Position.String
			if payload.InPlayerPosition != "" {
				newPosition = payload.InPlayerPosition
			}
			// if out player already has "s" position - ignore this event

			// mark out player as substitution
			outMp.Position = null.StringFrom("s")
			if _, err = store.UpdateMatchPlayerPosition(outMp); err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot mark out match player as substitution %s", payload.OutPlayerID))
			}

			// set actual position for incoming player
			inMp.Position = null.StringFrom(newPosition)
			if _, err = store.UpdateMatchPlayerPosition(inMp); err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot change position for in match player %s", payload.InPlayerID))
			}

			// notify match update
			if err := notifyMatchUpdate(store, match); err != nil {
				return stacktrace.Propagate(err, "cannot notify match update")
			}

			// auto-replace player for this user
			// find users that has this player picks
			picks, err := store.GetActivePicksAtTime(event.MatchID, payload.OutPlayerID, event.Timestamp)
			if err != nil {
				return stacktrace.Propagate(err,
					fmt.Sprintf("GetActivePicksAtTime failed (match_id: %s, player_id: %s)", event.MatchID, payload.OutPlayerID))
			}

			for _, p := range picks {
				// closed pick - ignore
				if p.EndedAt.Valid {
					continue
				}

				// check whether player already have pick for in player
				gamePicks, err := store.GetGamePicks(p.GameID)
				if err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot get game picks for game %s", p.GameID))
				}

				haveThisPlayer := false
				for _, gamePick := range gamePicks {
					if !gamePick.EndedAt.Valid {
						if gamePick.PlayerID == payload.InPlayerID {
							haveThisPlayer = true
							break
						}
					}
				}
				if haveThisPlayer {
					continue
				}

				p.EndedAt = null.TimeFrom(event.Timestamp)
				if err = store.UpdatePickEndedAt(p); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("UpdatePickEndedAt failed for pick %s", p.ID))
				}
				// create new pick at same position pointing new player
				newPick := &schema.GamePick{
					ID:          uuid.New().String(),
					CreatedAt:   event.Timestamp,
					Position:    p.Position,
					Score:       0,
					GameID:      p.GameID,
					PlayerID:    payload.InPlayerID,
					Version:     0,
					Minute:      event.Minute,
					Second:      event.Second,
					UserSwapped: false,
				}
				if _, err = store.InsertPick(newPick); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot create replace pick for pick %s", p.ID))
				}
				player, err := store.GetPlayerByID(payload.InPlayerID)
				if err != nil {
					return stacktrace.Propagate(err, fmt.Sprintf("cannot get player %s", *event.PlayerID))
				}
				// emit game updated event
				minutes := strconv.Itoa(event.Minute)
				seconds := strconv.Itoa(event.Second)
				if err = notifyGameUpdate(store, p.GameID, p.R.Game.UserID, "0", "0", player.ImageURL.String, "", minutes, seconds, event.MatchID, p.R.Player.NormalizedName.String, "Substitution", 0, 0, ""); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot notify game %s update", p.R.Game.ID))
				}

				// non-star was replaced to star, notify user
				if false { //!outMp.IsStar && inMp.IsStar {
					err = sendPush(
						store,
						p.R.Game.UserID,
						event.MatchID,
						fmt.Sprintf("Change %s ⚠️", GetPlayerName(inPlayer)),
						fmt.Sprintf("%s has been subbed for a Star player! Come back to swap them 📲", GetPlayerName(outPlayer)),
						map[string]string{
							"new_pick_id":   newPick.ID,
							"in_player_id":  payload.InPlayerID,
							"out_player_id": payload.OutPlayerID,
							"match_id":      event.MatchID,
						},
					)
					if err != nil {
						return stacktrace.Propagate(err,
							fmt.Sprintf("cannot send underdog %s replaced to star %s push", payload.OutPlayerID, payload.InPlayerID))
					}
				}
			}

			return nil
		})
	}

	return nil
}

func handleCancel(store database.Store, event *processor.Event) error {
	if event.Type == ActionCancel {
		if event.Payload == nil {
			logrus.WithField("id", event.ID).Error("payload is empty for cancel action")
			return nil
		}

		var payload struct {
			ID int `json:"id"`
		}

		if err := json.Unmarshal([]byte(*event.Payload), &payload); err != nil {
			logrus.WithError(err).WithField("id", event.ID).Error("cannot unmarshal payload")
		}

		// get source match event
		matchEvent, err := store.GetMatchEventByID(payload.ID)
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot get match event by id %d", payload.ID))
		}

		match, err := store.GetMatchByIDWithTeams(event.MatchID)
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot get match by id %s", match.ID))
		}

		// ignore any cancel event on ended match
		if match.Status == database.MatchStatusEnded {
			logrus.WithField("id", event.ID).Error("cancel events are not processed on ended matches")

			// set status as ignored
			matchEvent.Status = database.MatchEventStatusIgnored
			if err := store.UpdateMatchEventStatus(matchEvent); err != nil {
				return stacktrace.Propagate(err,
					fmt.Sprintf("cannot set match event %d status to ignored", matchEvent.ID))
			}

			return nil
		}

		// we cancel only specific actions:
		// - point actions
		// - goals
		if IsPointAction(matchEvent.Type) {
			return store.Transaction(func(s database.Store) error {
				// mark match event as cancelled
				matchEvent.Status = database.MatchEventStatusCancelled
				if err := store.UpdateMatchEventStatus(matchEvent); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot set match event %d status to cancelled", matchEvent.ID))
				}

				// delete all game events that is part of match event
				if err := store.DeleteGameEventsByMatchEventID(matchEvent.ID); err != nil {
					return stacktrace.Propagate(err,
						fmt.Sprintf("cannot delete game events by match event id %d", matchEvent.ID))
				}

				// decrease match player score
				if matchEvent.Points.Valid && matchEvent.TeamID.Valid && matchEvent.PlayerID.Valid {
					mp, err := store.GetMatchPlayer(matchEvent.MatchID, matchEvent.TeamID.String, matchEvent.PlayerID.String)
					if err != nil {
						return stacktrace.Propagate(err,
							fmt.Sprintf("GetMatchPlayer failed (match_id: %s, team_id: %s, player_id: %s)",
								matchEvent.MatchID, matchEvent.TeamID.String, matchEvent.PlayerID.String))
					}

					mp.Score.Float64 -= matchEvent.Points.Float64

					if _, err := store.UpdateMatchPlayerScore(mp); err != nil {
						return stacktrace.Propagate(err,
							fmt.Sprintf("UpdateMatchPlayerScore failed (id: %s)", mp.ID))
					}
				}

				return nil
			})
		} else if matchEvent.Type == ActionGoal || matchEvent.Type == ActionSelfGoal {
			if !matchEvent.TeamID.Valid {
				return errors.New("team_id is missing for goal-related event")
			}

			if matchEvent.TeamID.String != match.HomeTeamID && matchEvent.TeamID.String != match.AwayTeamID {
				return fmt.Errorf("unknown team_id %s", matchEvent.TeamID.String)
			}

			return store.Transaction(func(store database.Store) error {
				if (matchEvent.Type == ActionGoal && match.HomeTeamID == matchEvent.TeamID.String) ||
					(matchEvent.Type == ActionSelfGoal && match.AwayTeamID == matchEvent.TeamID.String) { // home score

					if _, err = store.DecHomeScore(match); err != nil {
						return stacktrace.Propagate(err, "cannot decrease home score")
					}
				} else if (matchEvent.Type == ActionGoal && match.AwayTeamID == matchEvent.TeamID.String) ||
					(matchEvent.Type == ActionSelfGoal && match.HomeTeamID == matchEvent.TeamID.String) { // away score

					if _, err = store.DecAwayScore(match); err != nil {
						return stacktrace.Propagate(err, "cannot decrease away score")
					}
				} else {
					return errors.New("UNKNOWN")
				}

				if err := notifyMatchUpdate(store, match); err != nil {
					return stacktrace.Propagate(err, "cannot notify match update")
				}

				return nil
			})
		}
	}

	return nil
}

// handleOneHourToGoNotification sends a notification to all users who are not part of the match one hour before it starts.
func handleOneHourToGoNotification(store database.Store, matchID string) error {
	// Retrieve match details
	match, err := store.GetMatchByIDWithTeams(matchID)
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot get match by id %s", matchID))
	}

	// Retrieve all user IDs not part of the match
	userIDs, err := store.GetUserIDsNotInMatch(matchID)
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot get user ids not in match %s", matchID))
	}

	// Construct notification message
	title := fmt.Sprintf("One hour to go: %s vs. %s", GetTeamName(match.R.HomeTeam), GetTeamName(match.R.AwayTeam))
	description := "Pick your squad and get in the game now!"

	// Send notifications to each user
	for _, userID := range userIDs {
		err := sendPush(store, userID, matchID, title, description, map[string]string{
			"match_id": matchID,
			"type":     "one_hour_to_go",
		})
		if err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot send notification to user %s", userID))
		}
	}

	return nil
}

// sendLineupNotifications sends lineup notifications to users
func sendLineupNotifications(store database.Store, match *schema.Match) error {
	logrus.Info("sending lineup notifications")
	userIDs, err := store.GetUserIDsByMatchID(match.ID)
	if err != nil {
		return stacktrace.Propagate(err, "GetUserIDsByMatchID failed")
	}
	for _, userID := range userIDs {
		var playing []string
		var notPlaying []string
		gameInfo, err := store.GetGameByUserIDMatchID(userID, match.ID)
		if err != nil {
			return stacktrace.Propagate(err, "GetGameByUserIDMatchID failed")
		}

		for _, gamePick := range gameInfo.R.GamePicks {
			player, err := GetPlayerByIDCached(store, gamePick.PlayerID)
			if err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("GetPlayerByIDCached failed (id: %s)", gamePick.PlayerID))
			}
			playerName := GetPlayerName(player)
			matchPlayer, err := store.GetMatchPlayerWithoutTeam(match.ID, player.ID)
			if err != nil {
				continue
				//return stacktrace.Propagate(err, "GetMatchPlayer failed")
			}

			if matchPlayer.FromLineups && matchPlayer.Position.Valid && matchPlayer.Position.String != "s" {
				playing = append(playing, playerName)
			} else {
				notPlaying = append(notPlaying, playerName)
			}
		}

		if len(notPlaying) == 0 {
			title := "All 4 starting ✅"
			message := "Line-up's in and all 4 of your players are starting! Nice one 👊"
			payload := map[string]string{
				"match_id": match.ID,
				"type":     "lineups",
			}
			err = sendPush(store, userID, match.ID, title, message, payload)
		} else {
			playingStr := strings.Join(playing, ", ")
			notPlayingStr := strings.Join(notPlaying, ", ")
			title := fmt.Sprintf("Change needed for %s ⚠️", notPlayingStr)
			message := fmt.Sprintf("%s is/are not in the starting 11! Consider replacing them. Currently playing: %s", notPlayingStr, playingStr)
			payload := map[string]string{
				"match_id": match.ID,
				"type":     "lineups",
			}
			err = sendPush(store, userID, match.ID, title, message, payload)
		}
		if err != nil {
			return stacktrace.Propagate(err, "cannot send push notification")
		}
	}
	logrus.WithField("count", len(userIDs)).Info("processed lineups users")

	// lineup notifications for non-playing user
	var pushTitle string
	pushText := fmt.Sprintf("Teams in for %s vs %s - pick your 4 players to win 💪",
		GetTeamName(match.R.HomeTeam), GetTeamName(match.R.AwayTeam))
	err = sendPushForAllUsersNotInMatch(store, pushTitle, pushText, match.ID, map[string]string{
		"match_id": match.ID,
		"type":     "lineups",
	})

	return nil
}

func calculatePrecedingLineUpActions(store database.Store, matchID string, currMatchEventID int) (int, error) {
	events, err := store.GetMatchEvents(matchID)
	if err != nil {
		return 0, stacktrace.Propagate(err, "GetMatchEvents failed")
	}
	count := 0
	for _, event := range events {
		if event.MatchEventID < currMatchEventID && event.Type == ActionLineUp {
			count++
		}
	}
	return count, nil
}

var slackWebhookURL = "https://hooks.slack.com/services/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"

type SlackMessage struct {
	Text string `json:"text"`
}

func SendSlackAlert(message string) error {
	slackMessage := SlackMessage{
		Text: message,
	}

	jsonMessage, err := json.Marshal(slackMessage)
	if err != nil {
		return fmt.Errorf("could not marshal message: %v", err)
	}

	req, err := http.NewRequest("POST", slackWebhookURL, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-ok response status: %v", resp.Status)
	}

	return nil
}

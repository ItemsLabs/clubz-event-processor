package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/hashicorp/go-multierror"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
)

const (
	NotificationTypePowerUpReminder  = 1
	NotificationTypeWorstPlayerAlert = 2
	NotificationTypeRankChangeAlert  = 3
)

func SendMatchNotifications(store database.Store) error {
	// matches to check
	// TODO: create new method for getting matches
	matches, err := store.GetMatchesForHeadlines(time.Now())
	if err != nil {
		return err
	}

	var errs error
	for _, m := range matches {
		if err := SendPowerUpReminderNotification(store, m); err != nil {
			errs = multierror.Append(errs, stacktrace.Propagate(err,
				fmt.Sprintf("SendPowerUpReminderNotification failed for match %s", m.ID)))
		}

		if err := SendWorstPlayerAlert(store, m); err != nil {
			errs = multierror.Append(errs, stacktrace.Propagate(err,
				fmt.Sprintf("SendWorstPlayerAlert failed for match %s", m.ID)))
		}
	}

	return errs
}

func SendPowerUpReminderNotification(store database.Store, match *schema.Match) error {
	now := time.Now()
	// check more than 35 min passed after match start and match is active (defend against old matches)
	if match.Status == database.MatchStatusGame && now.After(match.MatchTime.Add(35*time.Minute)) {
		shouldNotify, err := shouldNotify(store, match.ID, "", NotificationTypePowerUpReminder)
		if err != nil {
			return err
		}
		if !shouldNotify {
			return nil
		}

		// get games without powerups
		games, err := store.GetGamesWithoutPowerUps(match.ID)
		if err != nil {
			return err
		}

		// insert notification for such users
		err = store.Transaction(func(store database.Store) error {
			for _, g := range games {
				err := sendPush(
					store, g.UserID, match.ID,
					"Boost Reminder 📣",
					"You have boosts available to use! Activate now to climb up the ranks 📈",
					map[string]string{
						"match_id": match.ID,
						"type":     "powerup_reminder",
					})
				if err != nil {
					return err
				}
			}

			return store.CreateMatchNotification(&schema.MatchNotification{
				MatchID: match.ID,
				Type:    NotificationTypePowerUpReminder,
			})
		})
		if err != nil {
			return err
		}

		logrus.WithField("match_id", match.ID).Info("power-up reminder successfully sent")
	}

	return nil
}

func SendWorstPlayerAlert(store database.Store, match *schema.Match) error {
	now := time.Now()
	// check more than 25 min passed after match start and match is active (defend against old matches)
	if match.Status == database.MatchStatusGame && now.After(match.MatchTime.Add(25*time.Minute)) {
		shouldNotify, err := shouldNotify(store, match.ID, "", NotificationTypeWorstPlayerAlert)
		if err != nil {
			return err
		}
		if !shouldNotify {
			return nil
		}

		player, points, err := store.GetWorstPlayer(match.ID, 0, 25)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logrus.WithField("match_id", match.ID).Info("No players have scored yet for this match.")
				return nil // Skip notification if no players have scored
			} else {
				return stacktrace.Propagate(err, "GetWorstPlayer failed")
			}
		}

		// get picks with shitty player
		picks, err := store.GetActivePicksAtTime(match.ID, player.ID, now)
		if err != nil {
			return stacktrace.Propagate(err, "GetActivePicksAtTime failed")
		}

		title := fmt.Sprintf("Time to sub %s❓", player.FullName.String)
		message := fmt.Sprintf("%s is your worst performing player with %0.0f pts", player.FullName.String, points)

		err = store.Transaction(func(store database.Store) error {
			for _, p := range picks {
				err := sendPush(
					store, p.R.Game.UserID, match.ID,
					title,
					message,
					map[string]string{
						"match_id": match.ID,
						"type":     "worst_player_alert",
					})
				if err != nil {
					return err
				}
			}

			return store.CreateMatchNotification(&schema.MatchNotification{
				MatchID: match.ID,
				Type:    NotificationTypeWorstPlayerAlert,
			})
		})
		if err != nil {
			return err
		}

		logrus.WithField("match_id", match.ID).Info("worst player reminder successfully sent")
	}

	return nil
}

func SendRankChangedNotification(store database.Store) error {
	now := time.Now()
	matches, err := store.GetMatchesForHeadlines(now)
	if err != nil {
		return err
	}

	type rankInfo struct {
		Position   int
		Percentile float64
	}

	for _, m := range matches {
		// match in game status and current minute between 10 and 85
		if m.Status == database.MatchStatusGame && m.Minute > 10 && m.Minute < 85 {
			// check current leaderboard
			current, err := store.GetMatchLeaderboard(m.ID)
			if err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot get match %s current leaderboard", m.ID))
			}

			prev, err := store.GetMatchLeaderboardAtTime(m.ID, now.Add(-time.Minute*5))
			if err != nil {
				return stacktrace.Propagate(err, fmt.Sprintf("cannot get match %s leaderboard at time", m.ID))
			}

			// build map of user info
			var currRankMap = map[string]*rankInfo{}
			var prevRankMap = map[string]*rankInfo{}

			// fill curr and prev map
			for _, el := range current {
				if el.Position.Valid {
					// calculate percentile rank
					currRankMap[el.UserID] = &rankInfo{
						Position:   el.Position.Int,
						Percentile: float64(el.Position.Int) / float64(len(current)),
					}
				}
			}

			if len(current) > 0 {
				for _, el := range prev {
					if el.Position.Valid {
						// percentile rank over current
						prevRankMap[el.UserID] = &rankInfo{
							Position:   el.Position.Int,
							Percentile: float64(el.Position.Int) / float64(len(current)),
						}
					}
				}
			}

			// start to send rank change notifications
			return store.Transaction(func(store database.Store) error {
				for userID, currInfo := range currRankMap {
					prevInfo, ok := prevRankMap[userID]
					if !ok {
						continue
					}

					// change diff
					diff := prevInfo.Percentile - currInfo.Percentile

					// percentile rank change over 25%
					if math.Abs(diff) > 0.25 {
						shouldNotify, err := shouldNotify(store, m.ID, userID, NotificationTypeRankChangeAlert)
						if err != nil {
							return err
						}

						// already notified
						if !shouldNotify {
							continue
						}

						// rise alert
						if diff > 0 {
							err = sendPush(store, userID, m.ID,
								"Your rank is rising 📈",
								fmt.Sprintf("Your global rank has risen from %d to %d 🔥", prevInfo.Position, currInfo.Position),
								map[string]string{
									"match_id": m.ID,
									"type":     "rank_rise",
								})

							if err != nil {
								return err
							}
						} else if diff < 0 {
							// drop alert
							err = sendPush(store, userID, m.ID,
								"Your rank is falling 📉",
								fmt.Sprintf("Your global rank has dropped from %d to %d - time to use a powerup❓",
									prevInfo.Position, currInfo.Position),
								map[string]string{
									"match_id": m.ID,
									"type":     "rank_drop",
								})

							if err != nil {
								return err
							}
						}

						// mark that we already sent notification for user
						err = store.CreateMatchNotification(&schema.MatchNotification{
							MatchID: m.ID,
							UserID:  null.StringFrom(userID),
							Type:    NotificationTypeRankChangeAlert,
						})
						if err != nil {
							return err
						}
					}
				}

				return nil
			})
		}
	}

	return nil
}

func shouldNotify(store database.Store, matchID string, userID string, typ int) (bool, error) {
	_, err := store.GetMatchNotification(matchID, userID, typ)
	// notification exists do nothing
	if err == nil {
		return false, nil
	} else if err != sql.ErrNoRows {
		return false, err
	}

	return true, nil
}

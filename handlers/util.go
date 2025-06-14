package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/config"
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/palantir/stacktrace"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
)

func sendPushForAllUsers(store database.Store, title, message string, matchID string, payload map[string]string) error {
	users, err := store.GetUserIDsByMatchID(matchID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to get user ids by match id")
	}

	sent := make(map[string]struct{}, len(users))

	for _, userID := range users {
		user, err := store.GetUserByID(userID)
		if err != nil {
			return stacktrace.Propagate(err, "failed to get user by id")
		}
		if ok, err := sendPushTrackingSent(
			store,
			userID,
			payload["match_id"],
			title,
			message,
			payload,
			sent,
		); err != nil {
			return stacktrace.Propagate(err, "failed to send push")
		} else if ok {
			globalMixPanelSender.Send(
				"Push Notification Sent to User", userID, map[string]interface{}{
					"event_name":  "Push Notification",
					"title":       title,
					"description": message,
					"distinct_id": userID,
					"matchID":     matchID,
					"type":        "sent",
					"username":    user.Name,
				},
			)
		}
	}
	return nil
}

func sendPushForAllUsersNotInMatch(
	store database.Store,
	title, message string,
	matchID string,
	payload map[string]string,
) error {
	allUsers, err := store.GetAllUserIDsFromTable()
	if err != nil {
		return stacktrace.Propagate(err, "failed to get all user ids")
	}

	matchUsers, err := store.GetUserIDsByMatchID(matchID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to get user ids by match id")
	}

	matchUserMap := make(map[string]struct{}, len(matchUsers))
	for _, userID := range matchUsers {
		matchUserMap[userID] = struct{}{}
	}

	sent := make(map[string]struct{}, len(allUsers))
	sentCount := 0
	for _, userID := range allUsers {
		user, err := store.GetUserByID(userID)
		if err != nil {
			return stacktrace.Propagate(err, "failed to get user by id")
		}
		if _, inMatch := matchUserMap[userID]; !inMatch {
			if ok, err := sendPushTrackingSent(
				store,
				userID,
				payload["match_id"],
				title,
				message,
				payload,
				sent,
			); err != nil {
				return stacktrace.Propagate(err, "failed to send push")
			} else if ok {
				sentCount++
				globalMixPanelSender.Send(
					"Push Notification Sent to User Not in Match", userID, map[string]interface{}{
						"event_name":  "Push Notification",
						"title":       title,
						"description": message,
						"distinct_id": userID,
						"matchID":     matchID,
						"type":        "sent",
						"username":    user.Name,
					},
				)
			}
		}
	}
	logrus.WithField("count", sentCount).Info("total sent pushes to non-match users")
	return nil
}

func sendPush(store database.Store, userID, matchID, title, message string, payload map[string]string) error {
	pushToken, err := getPushTokenByUser(store, userID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to get push token")
	}
	return sendPushWithPushToken(store, userID, pushToken, matchID, title, message, payload)
}

func sendPushTrackingSent(
	store database.Store,
	userID, matchID, title, message string,
	payload map[string]string,
	sent map[string]struct{},
) (bool, error) {
	pushToken, err := getPushTokenByUser(store, userID)
	if err != nil {
		return false, stacktrace.Propagate(err, "failed to get push token")
	}
	// track sent push to this token
	if _, ok := sent[pushToken]; ok {
		return false, nil
	}
	sent[pushToken] = struct{}{}

	return true, sendPushWithPushToken(store, userID, pushToken, matchID, title, message, payload)
}

func getPushTokenByUser(store database.Store, userID string) (string, error) {
	if userID != "" {
		user, err := store.GetUserByID(userID)
		if err != nil {
			return "", fmt.Errorf("cannot get user %s by id: %v", userID, err)
		}
		return user.FirebaseID.String, nil
	}
	return "", nil
}

func sendPushWithPushToken(
	store database.Store,
	userID, pushToken, matchID, title, message string,
	payload map[string]string,
) error {
	var url string
	if userID != "" {
		if matchID != "" {
			gameID, err := store.GetGameByUserIDMatchID(userID, matchID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return stacktrace.Propagate(err, "failed to get game by user ID %s and match ID %s", userID, matchID)
			}
			if gameID != nil {
				url = fmt.Sprintf("/game/%s", gameID.ID)
			}
		}
	}

	var data = map[string]interface{}{
		"user_id":    userID,
		"match_id":   matchID,
		"title":      title,
		"message":    message,
		"payload":    payload,
		"push_token": pushToken,
		"url":        url,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return stacktrace.Propagate(err, "cannot marshal data")
	}
	pushNotification := schema.PushNotification{
		UserID:  userID,
		MatchID: null.StringFrom(matchID),
		Title:   null.StringFrom(title),
		Message: null.StringFrom(message),
		SentAt:  time.Now(),
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return stacktrace.Propagate(err, "failed to marshal payload to JSON")
	}
	pushNotification.Payload = null.JSONFrom(jsonPayload)

	count, err := store.CountUserNotificationsLastHour(userID)
	if err != nil {
		fmt.Printf("Error counting notifications for user %s: %v\n", userID, err)
	} else {
		fmt.Printf("Number of notifications sent to user %s in the last hour: %d\n", userID, count)
	}
	if count < 3 {
		if err := store.InsertPushNotification(&pushNotification); err != nil {
			return stacktrace.Propagate(err, "failed to insert push notification")
		}
		_, err = store.InsertAMQPEvent(
			&schema.AmqpEvent{
				Exchange: config.RMQFCMExchange(),
				Type:     "push_notification",
				Data:     string(b),
			},
		)
	}
	if err != nil {
		return stacktrace.Propagate(err, "cannot insert amqp_event")
	}
	user, err := store.GetUserByID(userID)
	if err != nil {
		return stacktrace.Propagate(err, "failed to get user by id")
	}
	// Track the push notification event in Mixpanel
	globalMixPanelSender.Send(
		"Push Notification Sent", userID, map[string]interface{}{
			"event_name":  "Push Notification",
			"title":       title,
			"description": message,
			"distinct_id": userID,
			"matchID":     matchID,
			"type":        "sent",
			"username":    user.Name,
		},
	)

	return nil
}

func notifyMatchUpdate(store database.Store, match *schema.Match) error {
	var payload = map[string]interface{}{
		"match_id": match.ID,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return stacktrace.Propagate(err, "cannot marshal match updated payload")
	}

	_, err = store.InsertAMQPEvent(
		&schema.AmqpEvent{
			Exchange: config.RMQFCMExchange(),
			Type:     "match_updated",
			Data:     string(b),
		},
	)
	if err != nil {
		return stacktrace.Propagate(err, "cannot insert match updated amqp_event")
	}

	// Track match update in Mixpanel
	globalMixPanelSender.Send(
		"Match Updated", match.ID, map[string]interface{}{
			"matchID": match.ID,
			"status":  match.Status,
		},
	)

	return nil
}

func notifyGameUpdate(
	store database.Store,
	gameID, userID string,
	gameScore string,
	gameInitialScore string,
	playerImage string,
	gameEventID string,
	eventMinute string,
	eventSecond string,
	eventMatchId string,
	normalizedName string,
	gameEvent string,
	nftMultiplier float64,
	boostMultiplier float64,
	nftImage string,
) error {
	//mx.Lock()
	//defer mx.Unlock()
	//key := fmt.Sprintf("%s_%s", gameID, userID)ß
	//if val, ok := lastGameNotify[key]; ok {How can
	//	// less than 30 secs passed from last update
	//	// ignore this event
	//	if time.Now().Sub(val) < (time.Second * 30) {
	//		return nil
	//	}
	//} else {
	//	lastGameNotify[key] = time.Now()
	//}

	// insert fcm notification that game was changed
	var payload = map[string]interface{}{
		"user_id":            userID,
		"game_id":            gameID,
		"game_initial_score": gameInitialScore,
		"game_score":         gameScore,
		"player_image":       playerImage,
		"game_event_id":      gameEventID,
		"game_event":         gameEvent,
		"event_minute":       eventMinute,
		"event_second":       eventSecond,
		"normalized_name":    normalizedName,
		"match_id":           eventMatchId,
		"nft_multiplier":     nftMultiplier,
		"boost_multiplier":   boostMultiplier,
		"nft_image":          nftImage,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return stacktrace.Propagate(err, "cannot marshal game updated payload")
	}

	_, err = store.InsertAMQPEvent(
		&schema.AmqpEvent{
			Exchange: config.RMQGameUpdatesExchange(),
			Type:     "update",
			Data:     string(b),
		},
	)
	if err != nil {
		return stacktrace.Propagate(err, "cannot insert game updated amqp_event")
	}
	return nil
}

// player cache functionality
var playersCache = cache.New(5*time.Minute, 10*time.Minute)

func GetPlayerByIDCached(store database.Store, playerID string) (*schema.Player, error) {
	if player, ok := playersCache.Get(playerID); ok {
		return player.(*schema.Player), nil
	}

	// get from db
	player, err := store.GetPlayerByID(playerID)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("cannot get player %s by id", playerID))
	}

	// put into cache
	playersCache.SetDefault(playerID, player)

	return player, nil
}

func GetPlayerName(player *schema.Player) string {
	if player == nil {
		return "Player Name"
	}

	if player.FullName.Valid {
		return player.FullName.String
	}
	if player.NickName.Valid {
		return player.FullName.String
	}
	if player.FirstName.Valid || player.LastName.Valid {
		return player.FirstName.String + " " + player.LastName.String
	}

	return "Player Name"
}

func GetTeamName(team *schema.Team) string {
	return team.Name
}

func gameStatusFromMatchStatus(status string) string {
	switch status {
	case database.MatchStatusUnknown:
		return database.GameStatusWaiting
	case database.MatchStatusWaiting:
		return database.GameStatusWaiting
	case database.MatchStatusLineups:
		return database.GameStatusWaiting
	case database.MatchStatusGame:
		return database.GameStatusGameplay
	case database.MatchStatusEnded:
		return database.GameStatusFinished
	case database.MatchStatusCancelled:
		return database.GameStatusFinished
	default:
		return database.GameStatusWaiting
	}
}

func IsPointAction(actionType int) bool {
	return actionType < 100
}

func ConvertSubscriptionTier(val int) string {
	switch val {
	case 1:
		return "premium"
	case 2:
		return "lite"
	default:
		return "none"
	}
}

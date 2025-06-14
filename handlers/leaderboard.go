package handlers

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gameon-app-inc/fanclash-event-processor/config"
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type leaderboardEntry struct {
	ID               int      `json:"id"`
	Position         *int     `json:"position"`
	UserName         string   `json:"user_name"`
	UserID           string   `json:"user_id"`
	UserAvatarURL    *string  `json:"user_avatar_url"`
	Premium          bool     `json:"premium"`
	SubscriptionTier string   `json:"subscription_tier"`
	Influencer       bool     `json:"influencer"`
	Score            *float64 `json:"score"`
}

type leaderboard struct {
	MatchID string              `json:"match_id"`
	Entries []*leaderboardEntry `json:"entries"`
}

func SendLeaderboardToFCM(store database.Store, matchID string) {
	entries, err := store.GetMatchLeaderboard(matchID)
	if err != nil {
		logrus.WithError(err).Error("cannot get match leaderboard")
		return
	}

	// push leaderboard to fcm
	resultEntries := make([]*leaderboardEntry, 0, len(entries))

	for _, entry := range entries {
		resultEntries = append(resultEntries, &leaderboardEntry{
			ID:               entry.ID,
			Position:         entry.Position.Ptr(),
			UserName:         entry.R.User.Name,
			UserID:           entry.R.User.ID,
			UserAvatarURL:    entry.R.User.AvatarURL.Ptr(),
			Score:            entry.Score.Ptr(),
			Premium:          entry.R.Game.Premium,
			SubscriptionTier: ConvertSubscriptionTier(entry.R.Game.SubscriptionTier),
			Influencer:       entry.R.User.Influencer,
		})
	}

	// marshal to json
	body, err := json.Marshal(&leaderboard{matchID, resultEntries})
	if err != nil {
		logrus.WithError(err).Error("cannot marshal leaderboard entries")
		return
	}

	// send to fcm
	_, err = store.InsertAMQPEvent(&schema.AmqpEvent{
		Exchange: config.RMQGameUpdatesExchange(),
		Type:     "leaderboard_updated",
		Data:     string(body),
	})
	if err != nil {
		logrus.WithError(err).Error(err, "cannot insert leaderboard amqp_event")
		return
	}

	logrus.WithField("match_id", matchID).Info("successfully sent leaderboard")
}

func SendLeaderboardToRedis(store database.Store, matchID string) {
	println("SendLeaderboardToRedis starting")
	ctx := context.Background()
	client := config.NewRedisClient()
	defer client.Close()

	entries, err := store.GetMatchLeaderboard(matchID)
	if err != nil {
		logrus.WithError(err).Error("cannot get match leaderboard")
		return
	}

	// Prepare leaderboard data for AMQP
	resultEntries := make([]*leaderboardEntry, 0, len(entries))
	for _, entry := range entries {
		resultEntries = append(resultEntries, &leaderboardEntry{
			ID:               entry.ID,
			Position:         entry.Position.Ptr(),
			UserName:         entry.R.User.Name,
			UserID:           entry.R.User.ID,
			UserAvatarURL:    entry.R.User.AvatarURL.Ptr(),
			Score:            entry.Score.Ptr(),
			Premium:          entry.R.Game.Premium,
			SubscriptionTier: ConvertSubscriptionTier(entry.R.Game.SubscriptionTier),
			Influencer:       entry.R.User.Influencer,
		})
	}

	// Marshal to JSON
	body, err := json.Marshal(&leaderboard{matchID, resultEntries})
	if err != nil {
		logrus.WithError(err).Error("cannot marshal leaderboard entries")
		return
	}

	// Send to AMQP
	_, err = store.InsertAMQPEvent(&schema.AmqpEvent{
		Exchange: config.RMQGameUpdatesExchange(),
		Type:     "leaderboard_updated",
		Data:     string(body),
	})
	if err != nil {
		logrus.WithError(err).Error("cannot insert leaderboard amqp_event")
		return
	}

	pipe := client.Pipeline()
	defer pipe.Close()

	for _, el := range entries {
		entryID := "entry:" + strconv.Itoa(el.ID)
		key := "match:" + matchID + ":" + entryID
		position, _ := el.Position.Value()
		userAvatarUrl, _ := el.R.User.AvatarURL.Value()
		var score float64
		if el.Score.Valid {
			score = el.Score.Float64
		} else {
			score = 0
		}

		data := map[string]interface{}{
			"position":        position,
			"user_name":       el.R.User.Name,
			"user_id":         el.UserID,
			"user_avatar_url": userAvatarUrl,
			"score":           score,
		}
		if el.R.Game.Premium {
			data["premium"] = el.R.Game.Premium
		}
		if el.R.User.Influencer {
			data["influencer"] = el.R.User.Influencer
		}
		if el.R.Game.SubscriptionTier != 0 {
			data["subscription_tier"] = ConvertSubscriptionTier(el.R.Game.SubscriptionTier)
		}

		pipe.HMSet(ctx, key, data)

		if err := pipe.ZAdd(context.Background(), "match:"+matchID+":scores", &redis.Z{Score: score, Member: entryID}).Err(); err != nil {
			logrus.WithError(err).Error("failed to update fans_playing field in Redis")
			return
		}
	}

	if err = pipe.Set(ctx, "match:"+matchID+":fans_playing", len(entries), 0).Err(); err != nil {
		logrus.WithError(err).Error("failed to update fans_playing field in Redis")
		return
	}

	// Execute all pipeline operations
	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to batch store leaderboard entries and update fans_playing in Redis")
		return
	}

	logrus.WithField("match_id", matchID).Info("successfully sent leaderboard")
}

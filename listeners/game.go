package listeners

import (
	"context"
	"encoding/json"

	"github.com/gameon-app-inc/fanclash-event-processor/amqp"
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/handlers"
	"github.com/sirupsen/logrus"
)

type gameListener struct {
	store database.Store
}

func newGameListener(store database.Store) *gameListener {
	return &gameListener{
		store: store,
	}
}

func (l *gameListener) OnConnect() {

}

func (l *gameListener) OnReconnect() {

}

func (l *gameListener) OnDeliveryReceived(sess amqp.Session, d amqp.Delivery) {
	if d.Type == "new" || d.Type == "game_updated" {
		var payload struct {
			MatchID string `json:"match_id"`
			GameID  string `json:"game_id"`
		}

		// ack by default
		_ = d.Ack(false)

		if err := json.Unmarshal(d.Body, &payload); err != nil {
			logrus.WithError(err).Errorf("cannot unmarshal new game payload")
		}

		logrus.WithFields(logrus.Fields{
			"game_id":  payload.GameID,
			"match_id": payload.MatchID,
		}).Info("new game")

		handlers.GetDebouncedSendMatchLeaderboard()(payload.MatchID)
		handlers.GetDebouncedSendMatchHeadlines()(payload.MatchID)
	} else {
		logrus.WithField("type", d.Type).Error("unknown delivery with type")
		_ = d.Nack(false, false)
	}
}

func StartGameEventsListener(ctx context.Context, store database.Store, url, exchange, queue string) {
	amqp.NewSubscriber(
		ctx,
		url,
		newGameListener(store),
		amqp.DurableQueueDefiner(queue, exchange),
		amqp.NoPrefetchConsumerDefiner(queue),
	)
}

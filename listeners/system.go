package listeners

import (
	"context"
	"encoding/json"
	"github.com/gameon-app-inc/fanclash-event-processor/amqp"
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/handlers"
	"github.com/sirupsen/logrus"
)

type systemListener struct {
	store database.Store
}

func newSystemListener(store database.Store) *systemListener {
	return &systemListener{
		store: store,
	}
}

func (l *systemListener) OnConnect() {

}

func (l *systemListener) OnReconnect() {

}

func (l *systemListener) OnDeliveryReceived(sess amqp.Session, d amqp.Delivery) {
	if d.Type == "update_played_time" {
		var payload struct {
			MatchID string `json:"match_id"`
		}

		// ack by default
		_ = d.Ack(false)

		if err := json.Unmarshal(d.Body, &payload); err != nil {
			logrus.WithError(err).Errorf("cannot unmarshal update_played_time payload")
		}

		logrus.WithFields(logrus.Fields{
			"match_id": payload.MatchID,
		}).Info("update played time")

		if err := handlers.UpdatePlayedTime(l.store, payload.MatchID); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"match_id": payload.MatchID,
			}).Info("update played time error")
		} else {
			logrus.WithFields(logrus.Fields{
				"match_id": payload.MatchID,
			}).Info("update played time success")
		}
	} else {
		logrus.WithField("type", d.Type).Info("unknown type")

		// ack by default
		_ = d.Ack(false)
	}
}

func StartSystemEventsListener(ctx context.Context, store database.Store, url, exchange, queue string) {
	amqp.NewSubscriber(
		ctx,
		url,
		newSystemListener(store),
		amqp.DurableQueueDefiner(queue, exchange),
		amqp.NoPrefetchConsumerDefiner(queue),
	)
}

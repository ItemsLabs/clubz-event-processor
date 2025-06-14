package processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/amqp"
	"github.com/gameon-app-inc/fanclash-event-processor/config"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

type AMQPReader struct {
	subscriber *amqp.Subscriber
	processor  EventProcessor
}

func (r *AMQPReader) OnConnect() {
	logrus.Info("amqp reader connected")
}

func (r *AMQPReader) OnReconnect() {
	logrus.Info("amqp reader reconnected")
}

func (r *AMQPReader) OnDeliveryReceived(sess amqp.Session, d amqp.Delivery) {
	event := new(Event)

	// unmarshal data, if not success reject this message
	if err := json.Unmarshal(d.Body, event); err != nil {
		logrus.WithField("body", string(d.Body)).Error("invalid payload")
		logrus.WithError(err).Error("cannot unmarshal incoming delivery, reject")
		if err = d.Reject(false); err != nil {
			logrus.WithError(err).Error("cannot reject delivery")
		}
	}

	// received match event
	err := r.processor.NewEvent(event)
	if err == nil {
		if err := d.Ack(false); err != nil {
			logrus.WithError(err).Error("cannot ack delivery")
		}
	} else {
		originalErr := err
		err = stacktrace.RootCause(err)
		if err == ErrOutOfSequence || err == ErrUnprocessedEventOutOfSequence || err == ErrOutOfSequenceAfterRestore {
			logrus.WithError(err).Error("event is came out of sequence, drop it for proper processing")

			if err := d.Ack(false); err != nil {
				logrus.WithError(err).Error("cannot ack delivery")
			}
		} else {
			logrus.WithError(originalErr).Error("cannot process match event, nack")
			// some delay before nack
			time.Sleep(time.Millisecond * 500)
			if err = d.Nack(false, true); err != nil {
				logrus.WithError(err).Error("cannot nack delivery")
			}
		}
	}
}

func NewAMQPReader(ctx context.Context, url string, processor EventProcessor) *AMQPReader {
	r := &AMQPReader{
		processor: processor,
	}
	r.subscriber = amqp.NewSubscriber(
		ctx,
		url,
		r,
		amqp.DurableQueueDefiner(config.RMQProcessorQueue(), config.RMQMatchEventExchange()),
		amqp.NoPrefetchConsumerDefiner(config.RMQProcessorQueue()),
	)

	return r
}

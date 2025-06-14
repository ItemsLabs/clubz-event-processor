package amqp

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type QueueDefiner func(Session) error
type ConsumerDefiner func(Session) (<-chan amqp.Delivery, error)

type Delivery amqp.Delivery

func (d Delivery) Ack(multiple bool) error {
	return amqp.Delivery(d).Ack(multiple)
}
func (d Delivery) Reject(requeue bool) error {
	return amqp.Delivery(d).Reject(requeue)
}
func (d Delivery) Nack(multiple, requeue bool) error {
	return amqp.Delivery(d).Nack(multiple, requeue)
}

type SubscriberListener interface {
	OnConnect()
	OnReconnect()
	OnDeliveryReceived(Session, Delivery)
}

type Subscriber struct {
	ctx             context.Context
	url             string
	listener        SubscriberListener
	queueDefiner    QueueDefiner
	consumerDefiner ConsumerDefiner
}

func (s *Subscriber) run() {
	conn := NewConnector(s.ctx)
	connections := conn.Connect(s.url)

	var connectedOnce = false

	for conn := range connections {
		sess, ok := <-conn

		if !ok {
			logrus.Error("cannot get session")
			continue
		}

		// notify listener
		if s.listener != nil {
			if connectedOnce {
				s.listener.OnReconnect()
			} else {
				connectedOnce = true
				s.listener.OnConnect()
			}
		}

		if s.queueDefiner != nil {
			err := s.queueDefiner(sess)
			if err != nil {
				logrus.WithError(err).Error("cannot define queue")
				continue
			}
		}

		deliveries, err := s.consumerDefiner(sess)
		if err != nil {
			logrus.WithError(err).Error("cannot define consumer")
		}

		working := true
		for working {
			select {
			case msg := <-deliveries:
				// empty message, looks like connection was gone, need reconnect
				if msg.Type == "" && len(msg.Body) == 0 {
					working = false
					break
				}

				if s.listener != nil {
					s.listener.OnDeliveryReceived(sess, Delivery(msg))
				}
			case <-s.ctx.Done():
				logrus.Info("shut down subscriber")
				return
			}
		}
	}
}

func NewSubscriber(ctx context.Context, url string, listener SubscriberListener, queueDefiner QueueDefiner, consumerDefiner ConsumerDefiner) *Subscriber {
	sub := &Subscriber{
		ctx:             ctx,
		url:             url,
		listener:        listener,
		queueDefiner:    queueDefiner,
		consumerDefiner: consumerDefiner,
	}

	go sub.run()

	return sub
}

func DurableQueueDefiner(queue string, exchange string) QueueDefiner {
	return func(sub Session) error {
		if _, err := sub.QueueDeclare(queue, true, false, false, false, nil); err != nil {
			return errors.New(fmt.Sprintf("cannot consume from exclusive queue: %q, %v", queue, err))
		}

		if err := sub.QueueBind(queue, "", exchange, false, nil); err != nil {
			return errors.New(fmt.Sprintf("cannot consume without a binding to exchange: %q, %v", exchange, err))
		}

		return nil
	}
}

func NoPrefetchConsumerDefiner(queue string) ConsumerDefiner {
	return func(sub Session) (<-chan amqp.Delivery, error) {
		err := sub.Qos(1, 0, false)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot set Qos, %v", err))
		}

		deliveries, err := sub.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot consume from: %q, %v", queue, err))
		}

		return deliveries, err
	}
}

func PrefetchConsumerDefiner(queue string, prefetchCount int) ConsumerDefiner {
	return func(sub Session) (<-chan amqp.Delivery, error) {
		err := sub.Qos(prefetchCount, 0, false)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot set Qos, %v", err))
		}

		deliveries, err := sub.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("cannot consume from: %q, %v", queue, err))
		}

		return deliveries, err
	}
}

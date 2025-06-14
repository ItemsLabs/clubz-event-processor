package amqp

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	ReconnectTimeout = time.Second * 3
)

type Session struct {
	*amqp.Connection
	*amqp.Channel
}

type Connector struct {
	ctx context.Context
}

func (ac *Connector) Connect(url string) chan chan Session {
	sessions := make(chan chan Session)

	go func() {
		defer close(sessions)
		for {
			fnDone := make(chan bool)
			go func() {
				ac.connectCycle(sessions, url)
				close(fnDone)
			}()

			select {
			case <-fnDone:
			case <-ac.ctx.Done():
				logrus.Info("shutting down connector")
				return
			}

			time.Sleep(ReconnectTimeout)
		}
	}()

	return sessions
}

func (ac *Connector) connectCycle(sessions chan chan Session, url string) {
	sess := make(chan Session)

	// make connection to reader
	// "sessions" is not buffered channel, so we will wait until some reader appears
	select {
	case sessions <- sess:
	case <-ac.ctx.Done():
		logrus.Info("shutting down connection")
		close(sess)
		return
	}

	logrus.WithField("url", url).Info("dial rabbitmq")

	conn, err := amqp.Dial(url)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Error("cannot (re)dial")
		close(sess)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.WithError(err).Info("cannot create channel")
		close(sess)
		return
	}

	select {
	case sess <- Session{conn, ch}:
	case <-ac.ctx.Done():
		logrus.Info("shutting down new session")
		close(sess)
		return
	}

	return
}

func NewConnector(c context.Context) *Connector {
	return &Connector{
		ctx: c,
	}
}

package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gameon-app-inc/fanclash-event-processor/listeners"

	"github.com/gameon-app-inc/fanclash-event-processor/handlers"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"github.com/gameon-app-inc/fanclash-event-processor/database/dbstore"

	"github.com/gameon-app-inc/fanclash-event-processor/config"
	"github.com/gameon-app-inc/fanclash-event-processor/processor"
	_ "github.com/jackc/pgx/stdlib"
)

func openDB() (*sql.DB, error) {
	// Start a database connection.
	db, err := sql.Open("pgx", config.DatabaseURL())
	if err != nil {
		return nil, err
	}

	// Actually test the connection against the database, so we catch
	// problematic connections early.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := openDB()
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	defer db.Close()

	store := dbstore.New(db)
	stateStore := processor.NewDBStateStore(store)
	eventSource := processor.NewStoreEventSource(store)

	handlers.InitMatchDebouncedFunctions(store)

	factories := []processor.EventProcessorFactory{
		processor.NewBaseEventProcessorFactory(store, stateStore, eventSource, handlers.NewPointsHandler(store), 1),
		processor.NewBaseEventProcessorFactory(store, stateStore, eventSource, handlers.NewSystemHandler(store), 2),
		processor.NewBaseEventProcessorFactory(store, stateStore, eventSource, handlers.NewHeadlinesHandler(store), 3),
	}

	c := cron.New()
	_, _ = c.AddFunc("@every 5m", func() {
		if err := handlers.SendAllHeadlines(store); err != nil {
			logrus.WithError(err).Error("cannot send all headlines")
		}
	})
	_, _ = c.AddFunc("@every 1m", func() {
		if err := handlers.SendMatchNotifications(store); err != nil {
			logrus.WithError(err).Error("cannot send match notifications")
		}
	})
	_, _ = c.AddFunc("@every 5m", func() {
		if err := handlers.SendRankChangedNotification(store); err != nil {
			logrus.WithError(err).Error("cannot send rank change notification")
		}
	})
	c.Start()

	processor.NewAMQPReader(
		context.Background(),
		config.RMQConnectionURL(),
		processor.NewCombinedProcessor(factories),
	)
	logrus.Info("amqp reader started")

	// listener for game events
	listeners.StartGameEventsListener(
		context.Background(),
		store,
		config.RMQConnectionURL(),
		config.RMQGamesExchange(),
		config.RMQGamesListenerQueue(),
	)

	// listener for system events
	listeners.StartSystemEventsListener(
		context.Background(),
		store,
		config.RMQConnectionURL(),
		config.RMQSystemExchange(),
		config.RMQSystemListenerQueue(),
	)

	done := make(chan bool)
	<-done
}

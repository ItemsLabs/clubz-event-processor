package handlers

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/processor"
)

type localHandler func(store database.Store, event *processor.Event) error

type BaseEventHandler struct {
	store         database.Store
	handlers      []localHandler
	transactional bool
}

func NewBaseEventHandler(store database.Store, handlers []localHandler, globalTransaction bool) *BaseEventHandler {
	return &BaseEventHandler{
		store:         store,
		handlers:      handlers,
		transactional: globalTransaction,
	}
}

func (h *BaseEventHandler) Reset() error {
	return nil
}

func (h *BaseEventHandler) Handle(event *processor.Event, alreadyProcessed bool) error {
	// ignore for now
	if alreadyProcessed {
		return nil
	}

	// global transaction for whole processing
	if h.transactional {
		return h.store.Transaction(func(store database.Store) error {
			return h.handle(store, event)
		})
	} else {
		// each handler responsible for it's own transaction processing
		return h.handle(h.store, event)
	}
}

func (h *BaseEventHandler) handle(store database.Store, event *processor.Event) error {
	for _, handler := range h.handlers {
		if err := handler(store, event); err != nil {
			return err
		}
	}

	return nil
}

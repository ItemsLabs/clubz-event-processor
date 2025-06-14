package processor

import (
	"fmt"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/palantir/stacktrace"
)

type BaseEventProcessorFactory struct {
	store         database.Store
	stateStore    StateStore
	eventSource   EventSource
	handler       EventHandler
	processorType int
}

func NewBaseEventProcessorFactory(
	store database.Store,
	stateStore StateStore,
	eventSource EventSource,
	handler EventHandler,
	processorType int) *BaseEventProcessorFactory {

	return &BaseEventProcessorFactory{
		store:         store,
		stateStore:    stateStore,
		eventSource:   eventSource,
		handler:       handler,
		processorType: processorType,
	}
}

func (f *BaseEventProcessorFactory) NewEventProcessor(matchID string) (EventProcessor, error) {
	matchProc, err := f.store.GetOrCreateMatchEventProcessor(matchID, f.processorType)
	if err != nil {
		return nil, stacktrace.Propagate(err,
			fmt.Sprintf("GetOrCreateMatchEventProcessor failed (match_id: %s, type: %d)", matchID, f.processorType))
	}

	return NewBaseEventProcessor(&EventProcessorMeta{
		ID:      matchProc.ID,
		MatchID: matchProc.MatchID,
	}, f.eventSource, f.handler, f.stateStore), nil
}

package processor

import (
	"fmt"
	"sync"

	"github.com/palantir/stacktrace"
)

type CombinedProcessor struct {
	factories []EventProcessorFactory

	mx        sync.Mutex
	procCache map[string][]EventProcessor
}

func NewCombinedProcessor(factories []EventProcessorFactory) *CombinedProcessor {
	return &CombinedProcessor{
		factories: factories,
		procCache: make(map[string][]EventProcessor),
	}
}

func (r *CombinedProcessor) NewEvent(ev *Event) error {
	processors, err := r.getEventProcessorsForMatch(ev.MatchID)
	if err != nil {
		return err
	}

	for _, proc := range processors {
		if err = proc.NewEvent(ev); err != nil {
			return stacktrace.Propagate(err, fmt.Sprintf("cannot process event: %d", ev.ID))
		}
	}

	return nil
}

func (r *CombinedProcessor) getEventProcessorsForMatch(matchID string) ([]EventProcessor, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	// check for proc in cache
	processors, ok := r.procCache[matchID]
	if ok {
		return processors, nil
	}

	// if no processors in cache, create new ones
	processors = make([]EventProcessor, 0, len(r.factories))
	for _, f := range r.factories {
		proc, err := f.NewEventProcessor(matchID)
		if err != nil {
			return nil, stacktrace.Propagate(err, fmt.Sprintf("cannot create event processor for match %s", matchID))
		}
		processors = append(processors, proc)
	}

	// put in cache
	r.procCache[matchID] = processors

	return processors, nil
}

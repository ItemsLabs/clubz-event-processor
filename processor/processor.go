package processor

import (
	"errors"
	"fmt"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

var (
	ErrOutOfSequence                 = errors.New("event out of sequence")
	ErrUnprocessedEventOutOfSequence = errors.New("unprocessed event out of sequence, something wrong with db")
	ErrOutOfSequenceAfterRestore     = errors.New("event out of sequence after restore")
)

type EventProcessorMeta struct {
	ID      string
	MatchID string
}

type BaseEventProcessor struct {
	meta       *EventProcessorMeta
	source     EventSource
	handler    EventHandler
	stateStore StateStore

	cachedCurrentID *int

	logger *logrus.Entry
}

func NewBaseEventProcessor(meta *EventProcessorMeta, source EventSource, handler EventHandler, stateStore StateStore) *BaseEventProcessor {
	return &BaseEventProcessor{
		meta:       meta,
		source:     source,
		handler:    handler,
		stateStore: stateStore,
		logger:     logrus.WithField("processor_id", meta.ID),
	}
}

func (p *BaseEventProcessor) process(ev *Event) error {
	if ev.ID == 0 {
		return nil
	}

	start := time.Now()

	currentID, err := p.getCurrentID()
	if err != nil {
		return err
	}

	// event from past - ignore
	if ev.MatchEventID <= currentID {
		p.logger.WithFields(logrus.Fields{
			"current_id": currentID,
			"new_id":     ev.MatchEventID,
		}).Info("received event from past")
		return nil
	}

	if ev.MatchEventID != (currentID + 1) {
		return ErrOutOfSequence
	}

	if err = p.handler.Handle(ev, false); err != nil {
		return stacktrace.Propagate(err,
			fmt.Sprintf("cannot handle event (processor_id: %s, match_id: %s, event_id: %d)",
				p.meta.ID, p.meta.MatchID, ev.MatchEventID))
	}

	if err = p.setCurrentID(ev.MatchEventID); err != nil {
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"match_event_id": ev.MatchEventID,
		"elapsed":        time.Since(start),
	}).Info("successfully processed event")

	return nil
}

// restore state and return unprocessed events
func (p *BaseEventProcessor) restore() ([]*Event, error) {
	p.logger.Info("restoring processor")

	currentID, err := p.getCurrentID()
	if err != nil {
		return nil, err
	}

	// now we should get all events from event source
	events, err := p.source.GetEvents(p.meta.MatchID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "cannot get events from event source (processor_id: %s)", p.meta.ID)
	}

	// notify handler that we reset processing state
	if err = p.handler.Reset(); err != nil {
		return nil, stacktrace.Propagate(err, "cannot reset handler (processor_id: %s)", p.meta.ID)
	}

	// find max processed index
	var startUnprocessed int
	for idx, event := range events {
		if event.MatchEventID <= currentID {
			if err = p.handler.Handle(event, true); err != nil {
				return nil, stacktrace.Propagate(err, "cannot ensure event for handler (processor_id: %s)", p.meta.ID)
			}
			startUnprocessed = idx + 1
		} else {
			// all already processed events are passed, stop
			break
		}
	}

	return events[startUnprocessed:], nil
}

func (p *BaseEventProcessor) setCurrentID(currentID int) error {
	if err := p.stateStore.SetCurrentID(p.meta.ID, currentID); err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot set current id (processor: %s)", p.meta.ID))
	}
	p.cachedCurrentID = &currentID
	return nil
}

func (p *BaseEventProcessor) getCurrentID() (int, error) {
	// use cache to decrease db time
	if p.cachedCurrentID != nil {
		return *p.cachedCurrentID, nil
	}

	currentID, err := p.stateStore.GetCurrentID(p.meta.ID)
	if err != nil {
		return 0, stacktrace.Propagate(err, fmt.Sprintf("cannot get current id (processor: %s)", p.meta.ID))
	}

	p.cachedCurrentID = &currentID
	return currentID, nil
}

func (p *BaseEventProcessor) NewEvent(newEvent *Event) error {
	err := p.process(newEvent)
	if err == nil {
		return nil
	}

	// event out of sequence, restore
	if err == ErrOutOfSequence {
		unprocessed, err := p.restore()
		if err != nil {
			return err
		}

		// process unprocessed
		for _, ev := range unprocessed {
			if err = p.process(ev); err != nil {
				if err == ErrOutOfSequence {
					return ErrUnprocessedEventOutOfSequence
				}
				return err
			}
		}

		// check conditions after restore
		currentID, err := p.getCurrentID()
		if err != nil {
			return err
		}

		if currentID < newEvent.MatchEventID {
			p.logger.WithFields(logrus.Fields{
				"current_id":   currentID,
				"new_event_id": newEvent.MatchEventID,
			}).Error("event after restore event if out of sequence")
			return ErrOutOfSequenceAfterRestore
		}

		return nil
	} else {
		// some other error
		return err
	}
}

package processor

import (
	"fmt"
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/palantir/stacktrace"
)

func FullProcessMatch(store database.Store, matchID string, processor EventProcessor) error {
	ev, err := store.GetLatestMatchEvent(matchID)
	if err != nil {
		return stacktrace.Propagate(err, fmt.Sprintf("cannot get match by id %s", matchID))
	}

	if err := processor.NewEvent(FromMatchEvent(ev)); err != nil {
		return stacktrace.Propagate(err, "cannot process event")
	}

	return nil
}

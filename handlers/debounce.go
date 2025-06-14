package handlers

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
)

var (
	debouncedSendLeaderboard func(string)
	debouncedSendHeadlines   func(string)
)

type MatchFn func(s database.Store, matchID string)

func DebounceMatchFn(store database.Store, fn MatchFn, delay time.Duration) func(string) {
	type timeoutRec struct {
		Timer *time.Timer
		Count int
	}

	var mx sync.Mutex
	var incoming = make(chan string)
	var delayedCalls = make(map[string]*timeoutRec)

	go func() {
		for matchID := range incoming {
			mx.Lock()
			if _, ok := delayedCalls[matchID]; ok {
				// match func was called recently, increase delayed counter to make sure
				// function will be called after timeout
				delayedCalls[matchID].Count++
			} else {
				// execute func
				go fn(store, matchID)
				// create new timer for this match
				rec := &timeoutRec{Timer: time.NewTimer(delay)}
				delayedCalls[matchID] = rec

				// run clean up func
				go func(timer *time.Timer, matchID string) {
					// wait until timer is ends
					<-timer.C

					mx.Lock()
					// check whether time is still here
					if _, ok := delayedCalls[matchID]; ok {
						// cache delayed call count before delete from map
						delayedCount := delayedCalls[matchID].Count
						delete(delayedCalls, matchID)
						// put value into incoming if needed to repeat exec cycle
						if delayedCount > 0 {
							go func() {
								incoming <- matchID
							}()
						}
					}
					mx.Unlock()
				}(rec.Timer, matchID)
			}
			mx.Unlock()
		}
	}()

	return func(matchID string) {
		go func() { incoming <- matchID }()
	}
}

func InitMatchDebouncedFunctions(store database.Store) {
	debouncedSendLeaderboard = DebounceMatchFn(store, SendLeaderboardToRedis, time.Second*5)
	debouncedSendHeadlines = DebounceMatchFn(store, func(s database.Store, matchID string) {
		if err := SendHeadlinesForMatch(s, matchID); err != nil {
			logrus.WithError(err).WithField("match_id", matchID).Error("cannot send headlines for match")
		}
	}, time.Second*30)
}

func GetDebouncedSendMatchLeaderboard() func(string) {
	return debouncedSendLeaderboard
}

func GetDebouncedSendMatchHeadlines() func(string) {
	return debouncedSendHeadlines
}

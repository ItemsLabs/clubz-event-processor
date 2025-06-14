package handlers

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestDebounceMatchFn(t *testing.T) {
	var c1, c2 int
	fn := DebounceMatchFn(nil, func(s database.Store, matchID string) {
		if matchID == "test" {
			c1++
		} else if matchID == "test2" {
			c2++
		}
	}, time.Second)

	start := time.Now()
	for time.Since(start) < (time.Second / 2) {
		fn("test")
		fn("test2")
		time.Sleep(time.Millisecond)
	}

	time.Sleep(time.Second)
	fn("test")
	fn("test2")
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 2, c1)
	assert.Equal(t, 2, c2)
}

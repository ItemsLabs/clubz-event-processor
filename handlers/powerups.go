package handlers

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"time"
)

type TimeGap struct {
	Start time.Time
	End   time.Time
}

func CalculateMatchTimeGaps(match *schema.Match) []TimeGap {
	var times = []time.Time{
		match.FEnd.Time,
		match.SStart.Time,
		match.SEnd.Time,
		match.X1Start.Time,
		match.X1End.Time,
		match.X2Start.Time,
		match.X2End.Time,
		match.PStart.Time,
	}

	var result []TimeGap
	for i := 0; i < len(times); i += 2 {
		start := times[i]
		end := times[i+1]

		if start.IsZero() {
			continue
		}

		if end.IsZero() {
			result = append(result, TimeGap{Start: start})
		} else {
			result = append(result, TimeGap{start, end})
		}
	}

	return result
}

func CalculatePowerUpDuration(startAt time.Time, duration time.Duration, timeGaps []TimeGap) time.Time {
	// initial end, before applying gap time
	initialEnd := startAt.Add(duration)

	// make span
	var extraDuration time.Duration

	// check whether end is between one of gaps
	for _, gap := range timeGaps {
		if gap.End.IsZero() {
			// power up inside unclosed gap, it means that we cannot calculate power up end
			// use some super high constant, to show it will ended in super future
			if initialEnd.After(gap.Start) {
				return initialEnd.Add(time.Hour * 999)
			}
			continue
		}

		// overlaps
		if gap.Start.Before(initialEnd) && gap.End.After(startAt) {
			if gap.Start.After(startAt) || gap.Start == startAt {
				// add whole gap duration
				extraDuration += gap.End.Sub(gap.Start)
			} else {
				// add duration from start to gap end
				extraDuration += gap.End.Sub(startAt)
			}
		}
	}

	return initialEnd.Add(extraDuration)
}

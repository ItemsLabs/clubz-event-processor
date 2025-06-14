package handlers

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"testing"
	"time"
)

func TestCalculatePowerUpDuration_UnclosedGap(t *testing.T) {
	now := time.Now()
	match := &schema.Match{
		FEnd: null.TimeFrom(now.Add(time.Minute * 5)),
	}
	result := CalculatePowerUpDuration(now, time.Minute*10, CalculateMatchTimeGaps(match))

	assert.Equal(t, now.Add(time.Minute*10).Add(time.Hour*999), result)
}

func TestCalculatePowerUpDuration_NoGaps(t *testing.T) {
	now := time.Now()
	match := &schema.Match{
		FEnd:   null.TimeFrom(now.Add(-time.Minute * 10)),
		SStart: null.TimeFrom(now.Add(-time.Minute * 5)),
	}
	result := CalculatePowerUpDuration(now, time.Minute*10, CalculateMatchTimeGaps(match))

	assert.Equal(t, now.Add(time.Minute*10), result)
}

func TestCalculatePowerUpDuration_OneGap(t *testing.T) {
	now := time.Now()
	match := &schema.Match{
		FEnd:   null.TimeFrom(now.Add(time.Minute * 3)),
		SStart: null.TimeFrom(now.Add(time.Minute * 7)),
	}
	result := CalculatePowerUpDuration(now, time.Minute*10, CalculateMatchTimeGaps(match))

	assert.Equal(t, now.Add(time.Minute*10).Add(time.Minute*(7-3)), result)
}

func TestCalculatePowerUpDuration_MultipleGaps(t *testing.T) {
	now := time.Now()
	match := &schema.Match{
		FEnd:   null.TimeFrom(now.Add(time.Minute * 1)),
		SStart: null.TimeFrom(now.Add(time.Minute * 2)),

		SEnd:    null.TimeFrom(now.Add(time.Minute * 3)),
		X1Start: null.TimeFrom(now.Add(time.Minute * 4)),

		X1End:   null.TimeFrom(now.Add(time.Minute * 5)),
		X2Start: null.TimeFrom(now.Add(time.Minute * 6)),

		X2End:  null.TimeFrom(now.Add(time.Minute * 7)),
		PStart: null.TimeFrom(now.Add(time.Minute * 8)),
	}
	result := CalculatePowerUpDuration(now, time.Minute*20, CalculateMatchTimeGaps(match))

	assert.Equal(t, now.Add(time.Minute*20).Add(time.Minute*((2-1)+(4-3)+(6-5)+(8-7))), result)
}

func TestCalculatePowerUpDuration_PowerUpStartDuringTimeBreak(t *testing.T) {
	now := time.Now()
	match := &schema.Match{
		FEnd:   null.TimeFrom(now.Add(-time.Minute * 3)),
		SStart: null.TimeFrom(now.Add(time.Minute * 2)),
	}
	result := CalculatePowerUpDuration(now, time.Minute*10, CalculateMatchTimeGaps(match))

	assert.Equal(t, now.Add(time.Minute*10).Add(time.Minute*2), result)
}

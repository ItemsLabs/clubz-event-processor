package handlers

import "github.com/volatiletech/null/v8"

type lineupsPayload struct {
	Players []struct {
		ID           string `json:"id"`
		JerseyNumber int    `json:"jersey_number"`
		Position     string `json:"position"`
	} `json:"players"`
}

type substitutionPayload struct {
	InPlayerPosition string `json:"in_player_position"`
	InPlayerID       string `json:"player_in_id"`
	OutPlayerID      string `json:"player_out_id"`
}

type EntityType int

const (
	EventEntity EntityType = iota + 1
	PowerUpEntity
)

type conditionPayload struct {
	Name       string      `json:"name"`
	Entity     EntityType  `json:"entity,omitempty"`     // the entity type of the condition. (e.g. "event" or "powerup")
	Field      null.String `json:"field,omitempty"`      // the field of the object type referred
	Expression null.String `json:"expression,omitempty"` // Can be any valid boolean expression (e.g. "score < 0")
	Value      interface{} `json:"value,omitempty"`      // Should be correctly type inferred in the parser
	Order      null.Int    `json:"order,omitempty"`      // Order of condition evaluation
}

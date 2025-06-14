package types

import (
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/null/v8"
)

type ChatMessage struct {
	ID        string      `boil:"id" json:"id"`
	Message   string      `boil:"message" json:"message"`
	MatchID   string      `boil:"match_id" json:"match_id,omitempty"`
	RoomID    string      `boil:"room_id" json:"room_id,omitempty"`
	SenderID  string      `boil:"sender_id" json:"sender_id,omitempty"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at"`
	UpdatedAt time.Time   `boil:"updated_at" json:"updated_at"`
	AvatarURL null.String `boil:"avatar_url" json:"avatar_url"`
	UserName  string      `boil:"user_name" json:"user_name"`
}

type MatchPlayerWithPPG struct {
	*schema.MatchPlayer
	PPG float64 `json:"ppg"`
}

type MatchPlayerWithPPGSlice []*MatchPlayerWithPPG

type LeaderboardEntry struct {
	DivisionID        null.String  `boil:"division_id" json:"division_id" toml:"division_id" yaml:"division_id"`
	DivisionTier      null.Int64   `boil:"division_tier" json:"division_tier" toml:"division_tier" yaml:"division_tier"`
	UserID            string       `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	UserName          null.String  `boil:"user_name" json:"user_name" toml:"user_name" yaml:"user_name"`
	TotalScore        null.Float64 `boil:"total_score" json:"total_score" toml:"total_score" yaml:"total_score"`
	Rank              null.Int64   `boil:"rank" json:"rank" toml:"rank" yaml:"rank"`
	CurrentUser       bool         `boil:"current_user" json:"current_user" toml:"current_user" yaml:"current_user"`
	WeekAverageRank   null.Float64 `boil:"week_average_rank" json:"week_average_rank" toml:"week_average_rank" yaml:"week_average_rank"`
	WeekMatchesPlayed null.Int64   `boil:"week_matches_played" json:"week_matches_played" toml:"week_matches_played" yaml:"week_matches_played"`
}

type GlobalLeaderboardEntry struct {
	UserID     string       `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	UserName   null.String  `boil:"user_name" json:"user_name" toml:"user_name" yaml:"user_name"`
	TotalScore null.Float64 `boil:"total_score" json:"total_score" toml:"total_score" yaml:"total_score"`
	Rank       null.Int64   `boil:"rank" json:"rank" toml:"rank" yaml:"rank"`
}

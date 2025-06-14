package handlers

import (
	"sort"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
)

type UserReward struct {
	UserID        string `json:"user_id"`
	Position      int    `json:"position"`
	Reward        Reward `json:"reward"`
	TransactionID string `json:"transaction_id,omitempty"` // Include this if it's meant to be set later
}

type Reward struct {
	ID           int     `json:"id,omitempty"`
	Credits      float64 `json:"credits"`
	GameToken    float64 `json:"game_token"`
	LaptToken    float64 `json:"lapt_token"`
	EventTickets int     `json:"event_tickets"`
	Balls        int     `json:"balls"`
	SignedBalls  int     `json:"signed_balls"`
	Shirts       int     `json:"shirts"`
	SignedShirts int     `json:"signed_shirts"`
	KickoffPack1 int     `json:"kickoff_pack_1"`
	KickoffPack2 int     `json:"kickoff_pack_2"`
	KickoffPack3 int     `json:"kickoff_pack_3"`
	SeasonPack1  int     `json:"season_pack_1"`
	SeasonPack2  int     `json:"season_pack_2"`
	SeasonPack3  int     `json:"season_pack_3"`
}

func CalculateRewards(rewards schema.MatchRewardSlice, users schema.MatchLeaderboardSlice) []*UserReward {
	result := make([]*UserReward, len(users))

	sort.Slice(rewards, func(i, j int) bool {
		return rewards[i].MinPosition < rewards[j].MinPosition
	})

	binarySearch := func(ranges *schema.MatchRewardSlice, target int) int {
		left, right := 0, len(*ranges)-1
		for left <= right {
			mid := left + (right-left)/2
			if (*ranges)[mid].MinPosition <= target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return right // return high index
	}

	findReward := func(position int) Reward {
		idx := binarySearch(&rewards, position)
		if idx >= 0 {
			r := rewards[idx]
			if !r.MaxPosition.Valid || position <= r.MaxPosition.Int {
				return Reward{
					ID:           r.ID,
					Credits:      r.Amount,
					GameToken:    r.Game,
					LaptToken:    r.Lapt,
					EventTickets: r.Event,
					Balls:        r.Balls,
					SignedBalls:  r.SignedBalls,
					Shirts:       r.Shirts,
					SignedShirts: r.SignedShirts,
					KickoffPack1: r.KickoffPack1,
					KickoffPack2: r.KickoffPack2,
					KickoffPack3: r.KickoffPack3,
					SeasonPack1:  r.SeasonPack1,
					SeasonPack2:  r.SeasonPack2,
					SeasonPack3:  r.SeasonPack3,
				}
			}
		}
		// Return an empty Reward struct if no matching reward is found
		return Reward{}
	}

	for idx, usr := range users {
		reward := findReward(usr.Position.Int)
		result[idx] = &UserReward{
			UserID:   usr.UserID,
			Reward:   reward,
			Position: usr.Position.Int,
		}
	}
	return result
}

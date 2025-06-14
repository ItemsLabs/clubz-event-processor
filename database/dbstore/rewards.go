package dbstore

import (
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetRewardByID retrieves a reward by its ID.
func (s *DBStore) GetRewardByID(rewardID string) (*schema.Reward, error) {
	return schema.Rewards(
		qm.Where("id = ?", rewardID),
	).One(s.db)
}

// CreateReward creates a new reward.
func (s *DBStore) CreateReward(reward *schema.Reward) (*schema.Reward, error) {
	reward.CreatedAt = time.Now()
	reward.UpdatedAt = time.Now()
	return reward, reward.Insert(s.db, boil.Infer())
}

// UpdateReward updates a reward by its ID.
func (s *DBStore) UpdateReward(reward *schema.Reward) error {
	reward.UpdatedAt = time.Now()
	_, err := reward.Update(s.db, boil.Infer())
	return err
}

// DeleteRewardByID deletes a reward by its ID.
func (s *DBStore) DeleteRewardByID(rewardID string) error {
	reward, err := s.GetRewardByID(rewardID)
	if err != nil {
		return err
	}

	_, err = reward.Delete(s.db)
	return err
}

// GetRewardsByName retrieves all rewards with a specific name.
func (s *DBStore) GetRewardsByName(name string) (schema.RewardSlice, error) {
	return schema.Rewards(
		qm.Where("name = ?", name),
	).All(s.db)
}

// GetAllRewards retrieves all rewards, ordered by creation date.
func (s *DBStore) GetAllRewards() (schema.RewardSlice, error) {
	return schema.Rewards(
		qm.OrderBy("created_at desc"),
	).All(s.db)
}

// GetRewardsByCreditsRange retrieves rewards within a specific range of credits.
func (s *DBStore) GetRewardsByCreditsRange(minCredits, maxCredits float64) (schema.RewardSlice, error) {
	return schema.Rewards(
		qm.Where("credits >= ? AND credits <= ?", minCredits, maxCredits),
	).All(s.db)
}

// GetRecentRewards retrieves the most recent rewards, limited by the provided amount.
func (s *DBStore) GetRecentRewards(limit int) (schema.RewardSlice, error) {
	return schema.Rewards(
		qm.OrderBy("created_at desc"),
		qm.Limit(limit),
	).All(s.db)
}

// UpdateRewardCredits updates the credits of a reward by its ID.
func (s *DBStore) UpdateRewardCredits(rewardID string, credits float64) error {
	reward, err := s.GetRewardByID(rewardID)
	if err != nil {
		return err
	}

	reward.Credits = credits
	reward.UpdatedAt = time.Now()
	_, err = reward.Update(s.db, boil.Whitelist("credits", "updated_at"))
	return err
}

// UpdateRewardGameToken updates the game token of a reward by its ID.
func (s *DBStore) UpdateRewardGameToken(rewardID string, gameToken float64) error {
	reward, err := s.GetRewardByID(rewardID)
	if err != nil {
		return err
	}

	reward.GameToken = gameToken
	reward.UpdatedAt = time.Now()
	_, err = reward.Update(s.db, boil.Whitelist("game_token", "updated_at"))
	return err
}

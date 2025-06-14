package dbstore

import (
	"fmt"
	"math/rand"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// CreateNFTBucket creates a new NFT bucket in the database with detailed error handling.
func (s *DBStore) CreateNFTBucket(nb *schema.NFTBucket) (*schema.NFTBucket, error) {
	if err := nb.Insert(s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to insert NFT bucket: %w", err)
	}
	return nb, nil
}

// GetNFTBucketByID retrieves an NFT bucket by its ID with detailed error handling.
func (s *DBStore) GetNFTBucketByID(id string) (*schema.NFTBucket, error) {
	nb, err := schema.FindNFTBucket(s.db, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find NFT bucket by ID %s: %w", id, err)
	}
	return nb, nil
}

// GetAllNFTBuckets retrieves all NFT buckets from the database with detailed error handling.
func (s *DBStore) GetAllNFTBuckets() (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets().All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all NFT buckets: %w", err)
	}
	return nbs, nil
}

// UpdateNFTBucket updates an existing NFT bucket with detailed error handling.
func (s *DBStore) UpdateNFTBucket(nb *schema.NFTBucket) error {
	_, err := nb.Update(s.db, boil.Whitelist("name", "team_id", "age", "game_position", "position", "star_rating", "updated_at"))
	if err != nil {
		return fmt.Errorf("failed to update NFT bucket: %w", err)
	}
	return nil
}

// DeleteNFTBucket deletes an NFT bucket by its ID with detailed error handling.
func (s *DBStore) DeleteNFTBucket(id string) error {
	nb, err := schema.FindNFTBucket(s.db, id)
	if err != nil {
		return fmt.Errorf("failed to find NFT bucket for deletion: %w", err)
	}
	_, err = nb.Delete(s.db)
	if err != nil {
		return fmt.Errorf("failed to delete NFT bucket: %w", err)
	}
	return nil
}

// GetNFTBucketsByTeamID retrieves all NFT buckets for a given team ID with detailed error handling.
func (s *DBStore) GetNFTBucketsByTeamID(teamID string) (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets(
		qm.Where("team_id = ?", teamID),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by team ID %s: %w", teamID, err)
	}
	return nbs, nil
}

// CountNFTBuckets returns the total number of NFT buckets with detailed error handling.
func (s *DBStore) CountNFTBuckets() (int64, error) {
	count, err := schema.NFTBuckets().Count(s.db)
	if err != nil {
		return 0, fmt.Errorf("failed to count NFT buckets: %w", err)
	}
	return count, nil
}

// GetNFTsByRarity retrieves all NFTs with a specific rarity.
func (s *DBStore) GetNFTsByRarity(rarity string) (schema.NFTBucketSlice, error) {
	var mods []qm.QueryMod

	switch rarity {
	case "common":
		mods = append(mods, qm.Where("common_claiming IS NOT NULL"))
	case "uncommon":
		mods = append(mods, qm.Where("uncommon_claiming IS NOT NULL"))
	case "rare":
		mods = append(mods, qm.Where("rare_claiming IS NOT NULL"))
	case "ultra_rare":
		mods = append(mods, qm.Where("ultra_rare_claiming IS NOT NULL"))
	case "legendary":
		mods = append(mods, qm.Where("legendary_claiming IS NOT NULL"))
	default:
		return nil, fmt.Errorf("invalid rarity: %s", rarity)
	}

	nbs, err := schema.NFTBuckets(mods...).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by %s rarity: %w", rarity, err)
	}
	return nbs, nil
}

// GetNFTsByTeamAndRarity retrieves all NFTs from a specific team with a specific rarity.
func (s *DBStore) GetNFTsByTeamAndRarity(teamID, rarity string) (schema.NFTBucketSlice, error) {
	var mods []qm.QueryMod
	mods = append(mods, qm.Where("team_id = ?", teamID))
	switch rarity {
	case "common":
		mods = append(mods, qm.Where("common_claiming IS NOT NULL"))
	case "uncommon":
		mods = append(mods, qm.Where("uncommon_claiming IS NOT NULL"))
	case "rare":
		mods = append(mods, qm.Where("rare_claiming IS NOT NULL"))
	case "ultra_rare":
		mods = append(mods, qm.Where("ultra_rare_claiming IS NOT NULL"))
	case "legendary":
		mods = append(mods, qm.Where("legendary_claiming IS NOT NULL"))
	default:
		return nil, fmt.Errorf("invalid rarity: %s", rarity)
	}

	nbs, err := schema.NFTBuckets(mods...).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by team ID %s and %s rarity: %w", teamID, rarity, err)
	}
	return nbs, nil
}

// GetRandomNFTByTeamAndRarity retrieves a random NFT from a specific team with a specified rarity.
func (s *DBStore) GetRandomNFTByTeamAndRarity(teamID string, rarity string) (*schema.NFTBucket, error) {
	var mods []qm.QueryMod
	mods = append(mods, qm.Where("team_id = ?", teamID))

	switch rarity {
	case "common":
		mods = append(mods, qm.Where("common_claiming IS NOT NULL"))
	case "uncommon":
		mods = append(mods, qm.Where("uncommon_claiming IS NOT NULL"))
	case "rare":
		mods = append(mods, qm.Where("rare_claiming IS NOT NULL"))
	case "ultra_rare":
		mods = append(mods, qm.Where("ultra_rare_claiming IS NOT NULL"))
	case "legendary":
		mods = append(mods, qm.Where("legendary_claiming IS NOT NULL"))
	default:
		return nil, fmt.Errorf("invalid rarity: %s", rarity)
	}

	// Retrieve all matching NFTs
	nfts, err := schema.NFTBuckets(mods...).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFTs by team ID %s and %s rarity: %w", teamID, rarity, err)
	}

	if len(nfts) == 0 {
		return nil, fmt.Errorf("no NFTs found for team ID %s with %s rarity", teamID, rarity)
	}

	// Select a random NFT from the list
	randomIndex := rand.Intn(len(nfts))
	return nfts[randomIndex], nil
}

// GetRandomNFTBucketByTeamID retrieves a random NFT bucket for a given team ID.
func (s *DBStore) GetRandomNFTBucketByTeamID(teamID string) (*schema.NFTBucket, error) {
	// Retrieve all NFT buckets for the specified team ID
	nfts, err := schema.NFTBuckets(
		qm.Where("team_id = ?", teamID),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by team ID %s: %w", teamID, err)
	}

	if len(nfts) == 0 {
		return nil, fmt.Errorf("no NFTs found for team ID %s", teamID)
	}

	// Select a random NFT from the list
	randomIndex := rand.Intn(len(nfts))
	return nfts[randomIndex], nil
}

// GetNFTBucketsByPosition retrieves all NFT buckets with a specific position.
func (s *DBStore) GetNFTBucketsByPosition(position string) (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets(
		qm.Where("position = ?", position),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by position %s: %w", position, err)
	}
	return nbs, nil
}

// GetNFTBucketsByGamePosition retrieves all NFT buckets with a specific game position.
func (s *DBStore) GetNFTBucketsByGamePosition(gamePosition string) (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets(
		qm.Where("game_position = ?", gamePosition),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by game position %s: %w", gamePosition, err)
	}
	return nbs, nil
}

// GetNFTBucketsByPositionAndTeam retrieves all NFT buckets for a specific team and position.
func (s *DBStore) GetNFTBucketsByPositionAndTeam(teamID, position string) (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets(
		qm.Where("team_id = ? AND position = ?", teamID, position),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by team ID %s and position %s: %w", teamID, position, err)
	}
	return nbs, nil
}

// GetNFTBucketsByTeamAndNationality retrieves all NFT buckets for a specific team and nationality.
func (s *DBStore) GetNFTBucketsByTeamAndNationality(teamID, nationality string) (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets(
		qm.Where("team_id = ? AND nationality = ?", teamID, nationality),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by team ID %s and nationality %s: %w", teamID, nationality, err)
	}
	return nbs, nil
}

// GetNFTBucketsByNationality retrieves all NFT buckets with a specific nationality.
func (s *DBStore) GetNFTBucketsByNationality(nationality string) (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets(
		qm.Where("nationality = ?", nationality),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by nationality %s: %w", nationality, err)
	}
	return nbs, nil
}

// GetNFTBucketsByAgeRange retrieves all NFT buckets within a specific age range.
func (s *DBStore) GetNFTBucketsByAgeRange(minAge, maxAge int) (schema.NFTBucketSlice, error) {
	nbs, err := schema.NFTBuckets(
		qm.Where("age BETWEEN ? AND ?", minAge, maxAge),
	).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by age range %d-%d: %w", minAge, maxAge, err)
	}
	return nbs, nil
}

// GetNFTBucketsByRarityAndPosition retrieves all NFT buckets with a specific rarity and position.
func (s *DBStore) GetNFTBucketsByRarityAndPosition(rarity, position string) (schema.NFTBucketSlice, error) {
	var mods []qm.QueryMod

	switch rarity {
	case "common":
		mods = append(mods, qm.Where("common_claiming IS NOT NULL"))
	case "uncommon":
		mods = append(mods, qm.Where("uncommon_claiming IS NOT NULL"))
	case "rare":
		mods = append(mods, qm.Where("rare_claiming IS NOT NULL"))
	case "ultra_rare":
		mods = append(mods, qm.Where("ultra_rare_claiming IS NOT NULL"))
	case "legendary":
		mods = append(mods, qm.Where("legendary_claiming IS NOT NULL"))
	default:
		return nil, fmt.Errorf("invalid rarity: %s", rarity)
	}

	mods = append(mods, qm.Where("position = ?", position))

	nbs, err := schema.NFTBuckets(mods...).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFT buckets by rarity %s and position %s: %w", rarity, position, err)
	}
	return nbs, nil
}

// CountNFTBucketsByRarity returns the total number of NFT buckets for each rarity.
func (s *DBStore) CountNFTBucketsByRarity() (map[string]int64, error) {
	counts := make(map[string]int64)

	var err error
	counts["common"], err = schema.NFTBuckets(qm.Where("common_claiming IS NOT NULL")).Count(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count common NFT buckets: %w", err)
	}

	counts["uncommon"], err = schema.NFTBuckets(qm.Where("uncommon_claiming IS NOT NULL")).Count(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count uncommon NFT buckets: %w", err)
	}

	counts["rare"], err = schema.NFTBuckets(qm.Where("rare_claiming IS NOT NULL")).Count(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count rare NFT buckets: %w", err)
	}

	counts["ultra_rare"], err = schema.NFTBuckets(qm.Where("ultra_rare_claiming IS NOT NULL")).Count(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count ultra rare NFT buckets: %w", err)
	}

	counts["legendary"], err = schema.NFTBuckets(qm.Where("legendary_claiming IS NOT NULL")).Count(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to count legendary NFT buckets: %w", err)
	}

	return counts, nil
}

// AssignNFTToTeam updates or creates a new NFT bucket with a specific team assignment with detailed error handling.
func (s *DBStore) AssignNFTToTeam(nftID, teamID string) (*schema.NFTBucket, error) {
	nb := &schema.NFTBucket{
		ID:     nftID,
		TeamID: null.StringFrom(teamID),
	}
	if err := nb.Upsert(s.db, true, []string{"id"}, boil.Whitelist("team_id"), boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to upsert NFT bucket: %w", err)
	}
	return nb, nil
}

// GetRandomNFTBucket retrieves a random NFT bucket.
func (s *DBStore) GetRandomNFTBucket() (*schema.NFTBucket, error) {
	nbs, err := schema.NFTBuckets().All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all NFT buckets: %w", err)
	}

	if len(nbs) == 0 {
		return nil, fmt.Errorf("no NFT buckets found")
	}

	randomIndex := rand.Intn(len(nbs))
	return nbs[randomIndex], nil
}

// GetNFTBucketByName retrieves an NFT bucket by its name.
func (s *DBStore) GetNFTBucketByName(name string) (*schema.NFTBucket, error) {
	nb, err := schema.NFTBuckets(qm.Where("name = ?", name)).One(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to find NFT bucket by name %s: %w", name, err)
	}
	return nb, nil
}

// GetRandomNFTBucketByTeamAndRarityAndStarRating retrieves a random NFT bucket based on team ID, rarity, and star rating.
func (s *DBStore) GetRandomNFTBucketByTeamAndRarityAndStarRating(teamID string, rarity string, starRating int) (*schema.NFTBucket, error) {
	var mods []qm.QueryMod
	mods = append(mods, qm.Where("team_id = ?", teamID))
	mods = append(mods, qm.Where("star_rating = ?", starRating))

	nfts, err := schema.NFTBuckets(mods...).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve NFTs by team ID %s, rarity %s, and star rating %d: %w", teamID, rarity, starRating, err)
	}

	if len(nfts) == 0 {
		return nil, fmt.Errorf("no NFTs found for team ID %s with rarity %s and star rating %d", teamID, rarity, starRating)
	}

	randomIndex := rand.Intn(len(nfts))
	return nfts[randomIndex], nil
}

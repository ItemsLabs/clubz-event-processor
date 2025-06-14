package dbstore

import (
	"fmt"
	"strings"

	"github.com/gameon-app-inc/fanclash-event-processor/database"
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// CreateAssignedPlayer creates a new assigned player in the database.
func (s *DBStore) CreateAssignedPlayer(ap *schema.AssignedPlayer) (*schema.AssignedPlayer, error) {
	// Input validation (example: ensuring non-empty PlayerID and UserID)
	if ap.PlayerNFTID.String == "" {
		return nil, fmt.Errorf("create assigned player: playerID and userID must not be empty")
	}

	if err := ap.Insert(s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("create assigned player: failed to insert: %w", err)
	}
	return ap, nil
}

// GetAssignedPlayerByID retrieves an assigned player by its ID.
func (s *DBStore) GetAssignedPlayerByID(id string) (*schema.AssignedPlayer, error) {
	if id == "" {
		return nil, fmt.Errorf("get assigned player by ID: ID must not be empty")
	}

	ap, err := schema.FindAssignedPlayer(s.db, id)
	if err != nil {
		return nil, fmt.Errorf("get assigned player by ID: failed to find player: %w", err)
	}
	return ap, nil
}

// GetAllAssignedPlayers retrieves all assigned players from the database.
func (s *DBStore) GetAllAssignedPlayers() (schema.AssignedPlayerSlice, error) {
	return schema.AssignedPlayers().All(s.db)
}

// UpdateAssignedPlayer updates an existing assigned player.
func (s *DBStore) UpdateAssignedPlayer(ap *schema.AssignedPlayer) error {
	if ap.ID == "" {
		return fmt.Errorf("update assigned player: ID must not be empty")
	}

	_, err := ap.Update(s.db, boil.Whitelist("player_id", "user_id", "updated_at"))
	if err != nil {
		return fmt.Errorf("update assigned player: failed to update: %w", err)
	}
	return nil
}

// DeleteAssignedPlayer deletes an assigned player by its ID.
func (s *DBStore) DeleteAssignedPlayer(id string) error {
	if id == "" {
		return fmt.Errorf("delete assigned player: ID must not be empty")
	}

	ap, err := schema.FindAssignedPlayer(s.db, id)
	if err != nil {
		return fmt.Errorf("delete assigned player: find player: %w", err)
	}
	_, err = ap.Delete(s.db)
	if err != nil {
		return fmt.Errorf("delete assigned player: failed to delete: %w", err)
	}
	return nil
}

// GetAssignedPlayersByUserID retrieves all assigned players for a given user ID.
func (s *DBStore) GetAssignedPlayersByUserID(userID string) (schema.AssignedPlayerSlice, error) {
	return schema.AssignedPlayers(
		qm.Where("user_id = ?", userID),
	).All(s.db)
}

// GetAssignedPlayersByPlayerID retrieves all assigned players for a given player ID.
func (s *DBStore) GetAssignedPlayersByPlayerID(playerID string) (schema.AssignedPlayerSlice, error) {
	return schema.AssignedPlayers(
		qm.Where("player_id = ?", playerID),
	).All(s.db)
}

// PlayerPopularity returns a map of player IDs to their count of assignments.
func (s *DBStore) PlayerPopularity() (map[string]int, error) {
	var results []struct {
		PlayerID string `boil:"player_id"`
		Count    int    `boil:"count"`
	}

	query := `
SELECT player_id, COUNT(*) AS count
FROM assigned_players
GROUP BY player_id
ORDER BY count DESC
`
	if _, err := s.db.Query(query, &results); err != nil { // Adjusted to a direct query execution
		return nil, fmt.Errorf("player popularity: %w", err)
	}

	popularity := make(map[string]int)
	for _, result := range results {
		popularity[result.PlayerID] = result.Count
	}
	return popularity, nil
}

func (s *DBStore) UserEngagement() (map[string]float64, error) {
	var results []struct {
		UserID     string  `boil:"user_id"`
		AvgChanges float64 `boil:"avg_changes"`
	}
	query := `
SELECT user_id, AVG(changes) AS avg_changes
FROM assigned_player_history
GROUP BY user_id
`
	_, err := s.db.Query(query, &results)
	if err != nil {
		return nil, fmt.Errorf("user engagement: %w", err)
	}

	engagement := make(map[string]float64)
	for _, result := range results {
		engagement[result.UserID] = result.AvgChanges
	}
	return engagement, nil
}

func (s *DBStore) TeamBalanceAnalysis() (map[string]map[string]int, error) {
	var results []struct {
		UserID   string `boil:"user_id"`
		PlayerID string `boil:"player_id"`
		Count    int    `boil:"count"`
	}
	query := `
SELECT user_id, player_id, COUNT(*) AS count
FROM assigned_players
GROUP BY user_id, player_id
`
	_, err := s.db.Query(query, &results)
	if err != nil {
		return nil, fmt.Errorf("team balance analysis: %w", err)
	}

	balance := make(map[string]map[string]int)
	for _, result := range results {
		if _, ok := balance[result.UserID]; !ok {
			balance[result.UserID] = make(map[string]int)
		}
		balance[result.UserID][result.PlayerID] = result.Count
	}
	return balance, nil
}

// CountAllEntries counts all entries in the assigned_players table.
func (s *DBStore) CountAllEntries() (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM assigned_players"
	if err := s.db.QueryRow(query).Scan(&count); err != nil {
		return 0, fmt.Errorf("count all entries: %w", err)
	}
	return count, nil
}

// GetAssignedPlayersByUUIDs retrieves multiple assigned players by their UUIDs.
func (s *DBStore) GetAssignedPlayersByUUIDs(uuids []string) (schema.AssignedPlayerSlice, error) {
	if len(uuids) == 0 {
		return nil, fmt.Errorf("get assigned players by UUIDs: UUIDs array must not be empty")
	}

	// Build the query with a WHERE IN clause
	query := fmt.Sprintf("id IN (%s)", strings.Join(strings.Split(strings.Repeat("?", len(uuids)), ""), ","))
	args := make([]interface{}, len(uuids))
	for i, uuid := range uuids {
		args[i] = uuid
	}

	assignedPlayers, err := schema.AssignedPlayers(qm.Where(query, args...)).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("get assigned players by UUIDs: %w", err)
	}

	return assignedPlayers, nil
}

// GetAssignedPlayersWithNFTDetails retrieves assigned players by their UUIDs and fetches corresponding NFT details using nft_id.
func (s *DBStore) GetAssignedPlayersWithNFTDetails(uuids []string) ([]struct {
	Player *schema.AssignedPlayer
	NFT    *schema.NFTBucket
}, error) {
	// Get assigned players by UUIDs
	assignedPlayers, err := s.GetAssignedPlayersByUUIDs(uuids)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned players: %w", err)
	}

	// Fetch corresponding NFT details using nft_id
	var result []struct {
		Player *schema.AssignedPlayer
		NFT    *schema.NFTBucket
	}

	for _, player := range assignedPlayers {
		nft, err := s.GetNFTBucketByID(player.PlayerNFTID.String)
		if err != nil {
			return nil, fmt.Errorf("failed to get NFT bucket by ID %s: %w", player.PlayerNFTID.String, err)
		}

		result = append(result, struct {
			Player *schema.AssignedPlayer
			NFT    *schema.NFTBucket
		}{
			Player: player,
			NFT:    nft,
		})
	}

	return result, nil
}

type PlayerWithNFTDetails struct {
	Player *schema.AssignedPlayer
	NFT    *schema.NFTBucket
}

// GetAssignedPlayersByNFTIDs retrieves assigned players by an array of NFT IDs.
func (s *DBStore) GetAssignedPlayersByNFTIDs(nftIDs []string) ([]database.PlayerWithNFTDetails, error) {
	if len(nftIDs) == 0 {
		return nil, fmt.Errorf("get assigned players by NFT IDs: NFT IDs array must not be empty")
	}

	// Create a query with a WHERE IN clause
	query := fmt.Sprintf("nft_id IN (%s)", strings.Join(strings.Split(strings.Repeat("?", len(nftIDs)), ""), ","))
	args := make([]interface{}, len(nftIDs))
	for i, id := range nftIDs {
		args[i] = id
	}

	assignedPlayers, err := schema.AssignedPlayers(qm.Where(query, args...)).All(s.db)
	if err != nil {
		return nil, fmt.Errorf("get assigned players by NFT IDs: %w", err)
	}

	// Fetch corresponding NFT details using nft_id
	var result []database.PlayerWithNFTDetails

	for _, player := range assignedPlayers {
		nft, err := s.GetNFTBucketByID(player.NFTID.String)
		if err != nil {
			return nil, fmt.Errorf("failed to get NFT bucket by ID %s: %w", player.NFTID.String, err)
		}

		result = append(result, database.PlayerWithNFTDetails{
			Player: player,
			NFT:    nft,
		})
	}

	return result, nil
}

func (s *DBStore) InsertAssignedCardPack(acp *schema.AssignedCardPack) (*schema.AssignedCardPack, error) {
	// Input validation: Ensure non-empty UserID and CardPackTypeID
	if acp.UserID == "" || acp.CardPackTypeID == "" {
		return nil, fmt.Errorf("insert assigned card pack: UserID and CardPackTypeID must not be empty")
	}

	// Insert the assigned card pack into the database
	if err := acp.Insert(s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert assigned card pack: failed to insert: %w", err)
	}
	return acp, nil
}

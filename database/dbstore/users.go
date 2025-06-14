package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetUserByID(userID string) (*schema.User, error) {
	return schema.Users(qm.Where("id = ?", userID)).One(s.db)
}

func (s *DBStore) UpdateFinishedGamesByUserID(userID, gameID string) error {
	query := `
        UPDATE users
        SET finished_games = array_append(finished_games, $2)
        WHERE id = $1;
    `

	_, err := s.db.Exec(query, userID, gameID)
	return err
}

func (s *DBStore) GetAllUserIDsFromTable() ([]string, error) {
	var userIDs []string
	users, err := schema.Users(qm.OrderBy("updated_at desc")).All(s.db)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	return userIDs, nil
}

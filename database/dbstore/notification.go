package dbstore

import (
	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (s *DBStore) GetMatchNotification(matchID string, userID string, typ int) (*schema.MatchNotification, error) {
	var mods []qm.QueryMod
	mods = append(mods,
		qm.Where("match_id = ?", matchID),
		qm.Where("type = ?", typ),
	)

	if userID != "" {
		mods = append(mods, qm.Where("user_id = ?", userID))
	}

	return schema.MatchNotifications(mods...).One(s.db)
}

func (s *DBStore) CreateMatchNotification(notification *schema.MatchNotification) error {
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}
	return notification.Insert(s.db, boil.Infer())
}

func (s *DBStore) InsertPushNotification(pushNotification *schema.PushNotification) error {
	query := `
        INSERT INTO push_notifications (user_id, match_id, title, message, payload, sent_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := s.db.Exec(query,
		pushNotification.UserID,
		pushNotification.MatchID,
		pushNotification.Title,
		pushNotification.Message,
		pushNotification.Payload,
		pushNotification.SentAt,
	)
	return err
}

func (s *DBStore) CountUserNotificationsLastHour(userID string) (int, error) {
	// SQL query to count notifications sent to a specific user in the last hour
	query := `
        SELECT COUNT(*)
        FROM push_notifications
        WHERE user_id = $1 AND sent_at >= NOW() - INTERVAL '1 hour'
    `
	var count int
	err := s.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

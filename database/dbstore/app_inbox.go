package dbstore

import (
	"time"

	"github.com/gameon-app-inc/fanclash-event-processor/database/schema"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetAppInboxByUserID retrieves all app inbox items for a specific user, ordered by creation date.
func (s *DBStore) GetAppInboxByUserID(userID string) (schema.AppInboxSlice, error) {
	return schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.OrderBy("created_at desc"),
	).All(s.db)
}

// GetAppInboxByID retrieves a single app inbox item by its ID and user ID.
func (s *DBStore) GetAppInboxByID(id string, userID string) (*schema.AppInbox, error) {
	return schema.AppInboxes(
		qm.Where("id = ?", id),
		qm.Where("user_id = ?", userID),
	).One(s.db)
}

// CreateAppInbox creates a new app inbox item.
func (s *DBStore) CreateAppInbox(inbox *schema.AppInbox) (*schema.AppInbox, error) {
	return inbox, inbox.Insert(s.db, boil.Infer())
}

// MarkAppInboxAsRead marks an app inbox item as read.
func (s *DBStore) MarkAppInboxAsRead(id string, userID string) error {
	inbox, err := s.GetAppInboxByID(id, userID)
	if err != nil {
		return err
	}

	inbox.Read = true
	_, err = inbox.Update(s.db, boil.Whitelist("read", "updated_at"))
	return err
}

// DeleteAppInboxByID deletes a specific app inbox item by its ID and user ID.
func (s *DBStore) DeleteAppInboxByID(id string, userID string) error {
	inbox, err := s.GetAppInboxByID(id, userID)
	if err != nil {
		return err
	}

	_, err = inbox.Delete(s.db)
	return err
}

// GetUnreadAppInboxCountByUserID returns the number of unread inbox items for a specific user.
func (s *DBStore) GetUnreadAppInboxCountByUserID(userID string) (int64, error) {
	return schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.Where("read = false"),
	).Count(s.db)
}

// DeleteAllAppInboxesByUserID deletes all app inbox items for a specific user.
func (s *DBStore) DeleteAllAppInboxesByUserID(userID string) error {
	_, err := schema.AppInboxes(
		qm.Where("user_id = ?", userID),
	).DeleteAll(s.db)
	return err
}

// GetRecentAppInboxByUserID retrieves the most recent inbox items for a user, with a limit.
func (s *DBStore) GetRecentAppInboxByUserID(userID string, limit int) (schema.AppInboxSlice, error) {
	return schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.OrderBy("created_at desc"),
		qm.Limit(limit),
	).All(s.db)
}

// SetAppInboxesAsReadByUserID marks all inbox items as read for a specific user.
func (s *DBStore) SetAppInboxesAsReadByUserID(userID string) error {
	_, err := schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.Where("read = false"),
	).UpdateAll(s.db, map[string]interface{}{"read": true, "updated_at": time.Now()})
	return err
}

// MarkAppInboxAsClaimed marks an app inbox item as claimed.
func (s *DBStore) MarkAppInboxAsClaimed(id string, userID string) error {
	inbox, err := s.GetAppInboxByID(id, userID)
	if err != nil {
		return err
	}

	inbox.Claimed = true
	inbox.ClamedAt = null.TimeFrom(time.Now())
	_, err = inbox.Update(s.db, boil.Whitelist("claimed", "clamed_at", "updated_at"))
	return err
}

// MarkAppInboxAsUnclaimed marks an app inbox item as unclaimed.
func (s *DBStore) MarkAppInboxAsUnclaimed(id string, userID string) error {
	inbox, err := s.GetAppInboxByID(id, userID)
	if err != nil {
		return err
	}

	inbox.Claimed = false
	inbox.ClamedAt = null.Time{}
	_, err = inbox.Update(s.db, boil.Whitelist("claimed", "clamed_at", "updated_at"))
	return err
}

// GetClaimedAppInboxByUserID retrieves all claimed app inbox items for a specific user, ordered by the claimed date.
func (s *DBStore) GetClaimedAppInboxByUserID(userID string) (schema.AppInboxSlice, error) {
	return schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.Where("claimed = true"),
		qm.OrderBy("clamed_at desc"),
	).All(s.db)
}

// GetUnclaimedAppInboxByUserID retrieves all unclaimed app inbox items for a specific user, ordered by creation date.
func (s *DBStore) GetUnclaimedAppInboxByUserID(userID string) (schema.AppInboxSlice, error) {
	return schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.Where("claimed = false"),
		qm.OrderBy("created_at desc"),
	).All(s.db)
}

// SetAppInboxesAsClaimedByUserID marks all inbox items as claimed for a specific user.
func (s *DBStore) SetAppInboxesAsClaimedByUserID(userID string) error {
	_, err := schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.Where("claimed = false"),
	).UpdateAll(s.db, map[string]interface{}{"claimed": true, "clamed_at": time.Now(), "updated_at": time.Now()})
	return err
}

// DeleteAllClaimedAppInboxesByUserID deletes all claimed app inbox items for a specific user.
func (s *DBStore) DeleteAllClaimedAppInboxesByUserID(userID string) error {
	_, err := schema.AppInboxes(
		qm.Where("user_id = ?", userID),
		qm.Where("claimed = true"),
	).DeleteAll(s.db)
	return err
}

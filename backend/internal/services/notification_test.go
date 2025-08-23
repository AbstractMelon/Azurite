package services

import (
	"encoding/json"
	"testing"

	"github.com/azurite/backend/internal/models"
)

func TestNotificationService_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	testData := map[string]interface{}{
		"mod_id":   123,
		"mod_name": "Test Mod",
	}

	err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", testData)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify notification was created
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count notifications: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 notification, got %d", count)
	}

	// Verify notification data
	var title, message, data string
	err = db.QueryRow(`
		SELECT title, message, data FROM notifications
		WHERE user_id = ?
	`, userID).Scan(&title, &message, &data)
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}

	if title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %s", title)
	}
	if message != "Test message" {
		t.Errorf("Expected message 'Test message', got %s", message)
	}

	var parsedData map[string]interface{}
	err = json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		t.Fatalf("Failed to parse notification data: %v", err)
	}

	if parsedData["mod_name"] != "Test Mod" {
		t.Errorf("Expected mod_name 'Test Mod', got %v", parsedData["mod_name"])
	}
}

func TestNotificationService_CreateWithNilData(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
	if err != nil {
		t.Fatalf("Create() with nil data error = %v", err)
	}

	// Verify notification was created
	var data string
	err = db.QueryRow("SELECT data FROM notifications WHERE user_id = ?", userID).Scan(&data)
	if err != nil {
		t.Fatalf("Failed to get notification data: %v", err)
	}

	if data != "" {
		t.Errorf("Expected empty data for nil input, got %s", data)
	}
}

func TestNotificationService_GetByUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Create test notifications
	for i := 1; i <= 5; i++ {
		err := service.Create(userID, models.NotificationTypeModApproved,
			"Test Title", "Test message", nil)
		if err != nil {
			t.Fatalf("Failed to create notification %d: %v", i, err)
		}
	}

	notifications, total, err := service.GetByUser(userID, 1, 3)
	if err != nil {
		t.Fatalf("GetByUser() error = %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(notifications) != 3 {
		t.Errorf("Expected 3 notifications per page, got %d", len(notifications))
	}

	// Test second page
	notifications, total, err = service.GetByUser(userID, 2, 3)
	if err != nil {
		t.Fatalf("GetByUser() page 2 error = %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(notifications) != 2 {
		t.Errorf("Expected 2 notifications on page 2, got %d", len(notifications))
	}

	// All notifications should be unread initially
	for _, notification := range notifications {
		if notification.IsRead {
			t.Error("Expected notification to be unread initially")
		}
		if notification.UserID != userID {
			t.Errorf("Expected notification to belong to user %d, got %d", userID, notification.UserID)
		}
	}
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Create notification
	err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Get notification ID
	var notificationID int
	err = db.QueryRow("SELECT id FROM notifications WHERE user_id = ?", userID).Scan(&notificationID)
	if err != nil {
		t.Fatalf("Failed to get notification ID: %v", err)
	}

	// Mark as read
	err = service.MarkAsRead(notificationID, userID)
	if err != nil {
		t.Fatalf("MarkAsRead() error = %v", err)
	}

	// Verify notification is marked as read
	var isRead bool
	err = db.QueryRow("SELECT is_read FROM notifications WHERE id = ?", notificationID).Scan(&isRead)
	if err != nil {
		t.Fatalf("Failed to check read status: %v", err)
	}

	if !isRead {
		t.Error("Expected notification to be marked as read")
	}
}

func TestNotificationService_MarkAsReadUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")

	// Create notification for first user
	err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Get notification ID
	var notificationID int
	err = db.QueryRow("SELECT id FROM notifications WHERE user_id = ?", userID).Scan(&notificationID)
	if err != nil {
		t.Fatalf("Failed to get notification ID: %v", err)
	}

	// Try to mark as read with different user
	err = service.MarkAsRead(notificationID, otherUserID)
	if err == nil {
		t.Error("Expected error when marking notification as read by different user")
	}
}

func TestNotificationService_MarkAllAsRead(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Create multiple notifications
	for i := 1; i <= 3; i++ {
		err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
		if err != nil {
			t.Fatalf("Failed to create notification %d: %v", i, err)
		}
	}

	// Mark all as read
	err := service.MarkAllAsRead(userID)
	if err != nil {
		t.Fatalf("MarkAllAsRead() error = %v", err)
	}

	// Verify all notifications are marked as read
	var unreadCount int
	err = db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0", userID).Scan(&unreadCount)
	if err != nil {
		t.Fatalf("Failed to count unread notifications: %v", err)
	}

	if unreadCount != 0 {
		t.Errorf("Expected 0 unread notifications, got %d", unreadCount)
	}
}

func TestNotificationService_GetUnreadCount(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Initially no unread notifications
	count, err := service.GetUnreadCount(userID)
	if err != nil {
		t.Fatalf("GetUnreadCount() error = %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 unread notifications initially, got %d", count)
	}

	// Create notifications
	for i := 1; i <= 5; i++ {
		err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
		if err != nil {
			t.Fatalf("Failed to create notification %d: %v", i, err)
		}
	}

	count, err = service.GetUnreadCount(userID)
	if err != nil {
		t.Fatalf("GetUnreadCount() after creation error = %v", err)
	}

	if count != 5 {
		t.Errorf("Expected 5 unread notifications, got %d", count)
	}

	// Mark some as read
	notifications, _, err := service.GetByUser(userID, 1, 10)
	if err != nil {
		t.Fatalf("GetByUser() error = %v", err)
	}

	for i := 0; i < 2; i++ {
		err = service.MarkAsRead(notifications[i].ID, userID)
		if err != nil {
			t.Fatalf("MarkAsRead() error = %v", err)
		}
	}

	count, err = service.GetUnreadCount(userID)
	if err != nil {
		t.Fatalf("GetUnreadCount() after marking read error = %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 unread notifications after marking 2 as read, got %d", count)
	}
}

func TestNotificationService_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Create notification
	err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Get notification ID
	var notificationID int
	err = db.QueryRow("SELECT id FROM notifications WHERE user_id = ?", userID).Scan(&notificationID)
	if err != nil {
		t.Fatalf("Failed to get notification ID: %v", err)
	}

	// Delete notification
	err = service.Delete(notificationID, userID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify notification was deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notifications WHERE id = ?", notificationID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count notifications: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected notification to be deleted, but it still exists")
	}
}

func TestNotificationService_DeleteUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")

	// Create notification for first user
	err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Get notification ID
	var notificationID int
	err = db.QueryRow("SELECT id FROM notifications WHERE user_id = ?", userID).Scan(&notificationID)
	if err != nil {
		t.Fatalf("Failed to get notification ID: %v", err)
	}

	// Try to delete with different user
	err = service.Delete(notificationID, otherUserID)
	if err == nil {
		t.Error("Expected error when deleting notification by different user")
	}
}

func TestNotificationService_NotifyModRejected(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	service.NotifyModRejected(userID, "Test Mod", "Contains malware")

	// Verify notification was created
	var title, message string
	err := db.QueryRow(`
		SELECT title, message FROM notifications
		WHERE user_id = ? AND type = ?
	`, userID, models.NotificationTypeModRejected).Scan(&title, &message)
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}

	if title != "Mod Rejected" {
		t.Errorf("Expected title 'Mod Rejected', got %s", title)
	}
	if message != "Your mod 'Test Mod' has been rejected. Reason: Contains malware" {
		t.Errorf("Unexpected message: %s", message)
	}
}

func TestNotificationService_NotifyModApproved(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	service.NotifyModApproved(userID, "Test Mod")

	// Verify notification was created
	var title, message string
	err := db.QueryRow(`
		SELECT title, message FROM notifications
		WHERE user_id = ? AND type = ?
	`, userID, models.NotificationTypeModApproved).Scan(&title, &message)
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}

	if title != "Mod Approved" {
		t.Errorf("Expected title 'Mod Approved', got %s", title)
	}
	if message != "Your mod 'Test Mod' has been approved and is now live!" {
		t.Errorf("Unexpected message: %s", message)
	}
}

func TestNotificationService_NotifyNewComment(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	service.NotifyNewComment(userID, "Test Mod", "John Doe")

	// Verify notification was created
	var title, message string
	err := db.QueryRow(`
		SELECT title, message FROM notifications
		WHERE user_id = ? AND type = ?
	`, userID, models.NotificationTypeNewComment).Scan(&title, &message)
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}

	if title != "New Comment" {
		t.Errorf("Expected title 'New Comment', got %s", title)
	}
	if message != "John Doe commented on your mod 'Test Mod'" {
		t.Errorf("Unexpected message: %s", message)
	}
}

func TestNotificationService_NotifyModMilestone(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	service.NotifyModMilestone(userID, "Test Mod", "downloads", 1000)

	// Verify notification was created
	var title, message string
	err := db.QueryRow(`
		SELECT title, message FROM notifications
		WHERE user_id = ? AND type = ?
	`, userID, models.NotificationTypeModMilestone).Scan(&title, &message)
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}

	if title != "Mod Milestone Reached" {
		t.Errorf("Expected title 'Mod Milestone Reached', got %s", title)
	}
	if message != "Your mod 'Test Mod' has reached 1000 downloads!" {
		t.Errorf("Unexpected message: %s", message)
	}
}

func TestNotificationService_NotifyGameRequest(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	admin1ID := createTestUser(t, db, "admin1", "admin1@example.com")
	admin2ID := createTestUser(t, db, "admin2", "admin2@example.com")

	adminIDs := []int{admin1ID, admin2ID}
	service.NotifyGameRequest(adminIDs, "John Doe", "New Game")

	// Verify notifications were created for both admins
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM notifications
		WHERE type = ? AND user_id IN (?, ?)
	`, models.NotificationTypeGameRequest, admin1ID, admin2ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count notifications: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 notifications (one for each admin), got %d", count)
	}

	// Check one of the notifications
	var title, message string
	err = db.QueryRow(`
		SELECT title, message FROM notifications
		WHERE user_id = ? AND type = ?
	`, admin1ID, models.NotificationTypeGameRequest).Scan(&title, &message)
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}

	if title != "New Game Request" {
		t.Errorf("Expected title 'New Game Request', got %s", title)
	}
	if message != "John Doe has requested to add 'New Game' to the platform" {
		t.Errorf("Unexpected message: %s", message)
	}
}

func TestNotificationService_CleanupOldNotifications(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	service := NewNotificationService(db, config)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Create some notifications
	for i := 1; i <= 5; i++ {
		err := service.Create(userID, models.NotificationTypeModApproved, "Test Title", "Test message", nil)
		if err != nil {
			t.Fatalf("Failed to create notification %d: %v", i, err)
		}
	}

	// Manually update some notifications to be old (older than 30 days)
	_, err := db.Exec(`
		UPDATE notifications
		SET created_at = datetime('now', '-35 days')
		WHERE user_id = ?
		LIMIT 3
	`, userID)
	if err != nil {
		t.Fatalf("Failed to update notification dates: %v", err)
	}

	// Cleanup old notifications (older than 30 days)
	err = service.CleanupOldNotifications(30)
	if err != nil {
		t.Fatalf("CleanupOldNotifications() error = %v", err)
	}

	// Verify only recent notifications remain
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count remaining notifications: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 notifications to remain, got %d", count)
	}
}

package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/azurite/backend/internal/config"
	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/models"
	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	db     *database.DB
	config *config.Config
}

func NewNotificationService(db *database.DB, cfg *config.Config) *NotificationService {
	return &NotificationService{
		db:     db,
		config: cfg,
	}
}

func (s *NotificationService) Create(userID int, notificationType, title, message string, data interface{}) error {
	dataJSON := ""
	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		dataJSON = string(dataBytes)
	}

	_, err := s.db.Exec(`
		INSERT INTO notifications (user_id, type, title, message, data)
		VALUES (?, ?, ?, ?, ?)
	`, userID, notificationType, title, message, dataJSON)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	user, err := s.getUserForNotification(userID)
	if err == nil && user.NotifyEmail {
		go s.sendEmailNotification(user, title, message)
	}

	return nil
}

func (s *NotificationService) GetByUser(userID int, page, perPage int) ([]models.Notification, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT id, user_id, type, title, message, data, is_read, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&notification.Title, &notification.Message, &notification.Data,
			&notification.IsRead, &notification.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, total, nil
}

func (s *NotificationService) MarkAsRead(notificationID, userID int) error {
	result, err := s.db.Exec(`
		UPDATE notifications SET is_read = 1
		WHERE id = ? AND user_id = ?
	`, notificationID, userID)

	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("notification not found or unauthorized")
	}

	return nil
}

func (s *NotificationService) MarkAllAsRead(userID int) error {
	_, err := s.db.Exec(`
		UPDATE notifications SET is_read = 1
		WHERE user_id = ? AND is_read = 0
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

func (s *NotificationService) GetUnreadCount(userID int) (int, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM notifications
		WHERE user_id = ? AND is_read = 0
	`, userID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (s *NotificationService) Delete(notificationID, userID int) error {
	result, err := s.db.Exec(`
		DELETE FROM notifications
		WHERE id = ? AND user_id = ?
	`, notificationID, userID)

	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("notification not found or unauthorized")
	}

	return nil
}

func (s *NotificationService) NotifyModRejected(ownerID int, modName, reason string) {
	title := "Mod Rejected"
	message := fmt.Sprintf("Your mod '%s' has been rejected. Reason: %s", modName, reason)
	data := map[string]interface{}{
		"mod_name": modName,
		"reason":   reason,
	}
	s.Create(ownerID, models.NotificationTypeModRejected, title, message, data)
}

func (s *NotificationService) NotifyModApproved(ownerID int, modName string) {
	title := "Mod Approved"
	message := fmt.Sprintf("Your mod '%s' has been approved and is now live!", modName)
	data := map[string]interface{}{
		"mod_name": modName,
	}
	s.Create(ownerID, models.NotificationTypeModApproved, title, message, data)
}

func (s *NotificationService) NotifyNewComment(ownerID int, modName, commenterName string) {
	title := "New Comment"
	message := fmt.Sprintf("%s commented on your mod '%s'", commenterName, modName)
	data := map[string]interface{}{
		"mod_name":       modName,
		"commenter_name": commenterName,
	}
	s.Create(ownerID, models.NotificationTypeNewComment, title, message, data)
}

func (s *NotificationService) NotifyModMilestone(ownerID int, modName, milestone string, count int) {
	title := "Mod Milestone Reached"
	message := fmt.Sprintf("Your mod '%s' has reached %d %s!", modName, count, milestone)
	data := map[string]interface{}{
		"mod_name":  modName,
		"milestone": milestone,
		"count":     count,
	}
	s.Create(ownerID, models.NotificationTypeModMilestone, title, message, data)
}

func (s *NotificationService) NotifyGameRequest(adminIDs []int, requestorName, gameName string) {
	title := "New Game Request"
	message := fmt.Sprintf("%s has requested to add '%s' to the platform", requestorName, gameName)
	data := map[string]interface{}{
		"requestor_name": requestorName,
		"game_name":      gameName,
	}

	for _, adminID := range adminIDs {
		s.Create(adminID, models.NotificationTypeGameRequest, title, message, data)
	}
}

func (s *NotificationService) getUserForNotification(userID int) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, username, email, display_name, notify_email
		FROM users WHERE id = ?
	`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.DisplayName, &user.NotifyEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *NotificationService) sendEmailNotification(user *models.User, title, message string) {
	if s.config.Email.SMTPHost == "" || s.config.Email.SMTPUsername == "" {
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.Email.FromEmail)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "[Azurite] "+title)

	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Hello %s,</h2>
			<p>%s</p>
			<hr>
			<p>Best regards,<br>The Azurite Team</p>
			<p><small>You received this email because you have email notifications enabled. You can change this in your account settings.</small></p>
		</body>
		</html>
	`, user.DisplayName, message)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(
		s.config.Email.SMTPHost,
		s.config.Email.SMTPPort,
		s.config.Email.SMTPUsername,
		s.config.Email.SMTPPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
	}
}

func (s *NotificationService) sendSimpleEmail(to, subject, body string) error {
	if s.config.Email.SMTPHost == "" || s.config.Email.SMTPUsername == "" {
		fmt.Println("Email not configured")
		return errors.New("email not configured")
	}

	auth := smtp.PlainAuth("",
		s.config.Email.SMTPUsername,
		s.config.Email.SMTPPassword,
		s.config.Email.SMTPHost,
	)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(
		s.config.Email.SMTPHost+":"+strconv.Itoa(s.config.Email.SMTPPort),
		auth,
		s.config.Email.FromEmail,
		[]string{to},
		msg,
	)

	if err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", to, err)
	}

	fmt.Println("Email sent successfully")

	return err
}

func (s *NotificationService) SendPasswordResetEmail(email, token string) error {
	subject := "[Azurite] Password Reset Request"
	body := fmt.Sprintf(`
You have requested a password reset for your Azurite account.

Click the following link to reset your password:
http://localhost:3000/reset-password?token=%s

If you did not request this, please ignore this email.

The Azurite Team
	`, token)

	return s.sendSimpleEmail(email, subject, body)
}

func (s *NotificationService) SendWelcomeEmail(user *models.User) error {
	subject := "[Azurite] Welcome to Azurite!"
	body := fmt.Sprintf(`
Welcome to Azurite, %s!

Your account has been successfully created. You can now:
- Browse and download mods
- Upload your own mods
- Interact with the community

Visit Azurite: http://localhost:3000

Happy modding!
The Azurite Team
	`, user.DisplayName)

	return s.sendSimpleEmail(user.Email, subject, body)
}

func (s *NotificationService) CleanupOldNotifications(days int) error {
	_, err := s.db.Exec(`
		DELETE FROM notifications
		WHERE created_at < datetime('now', '-' || ? || ' days')
	`, days)

	if err != nil {
		return fmt.Errorf("failed to cleanup old notifications: %w", err)
	}

	return nil
}

package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/models"
)

type AdminService struct {
	db *database.DB
}

func NewAdminService(db *database.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

func (s *AdminService) BanUser(req *models.BanCreateRequest, bannedBy int) (*models.Ban, error) {
	if req.UserID == nil && req.IPAddress == "" {
		return nil, errors.New("either user ID or IP address must be provided")
	}

	var expiresAt *time.Time
	if req.Duration > 0 {
		expiry := time.Now().Add(time.Duration(req.Duration) * time.Hour * 24)
		expiresAt = &expiry
	}

	result, err := s.db.Exec(`
		INSERT INTO bans (user_id, ip_address, game_id, reason, banned_by, expires_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.UserID, req.IPAddress, req.GameID, req.Reason, bannedBy, expiresAt, true)

	if err != nil {
		return nil, fmt.Errorf("failed to create ban: %w", err)
	}

	banID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get ban ID: %w", err)
	}

	return s.GetBanByID(int(banID))
}

func (s *AdminService) GetBanByID(id int) (*models.Ban, error) {
	ban := &models.Ban{}
	var username, displayName, bannedByUsername, bannedByDisplayName sql.NullString
	var gameName sql.NullString

	err := s.db.QueryRow(`
		SELECT b.id, b.user_id, b.ip_address, b.game_id, b.reason, b.banned_by,
		       b.expires_at, b.is_active, b.created_at,
		       u.username, u.display_name, bu.username, bu.display_name, g.name
		FROM bans b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN users bu ON b.banned_by = bu.id
		LEFT JOIN games g ON b.game_id = g.id
		WHERE b.id = ?
	`, id).Scan(
		&ban.ID, &ban.UserID, &ban.IPAddress, &ban.GameID, &ban.Reason,
		&ban.BannedBy, &ban.ExpiresAt, &ban.IsActive, &ban.CreatedAt,
		&username, &displayName, &bannedByUsername, &bannedByDisplayName, &gameName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("ban not found")
		}
		return nil, fmt.Errorf("failed to get ban: %w", err)
	}

	if ban.UserID != nil {
		ban.User = &models.User{
			ID:          *ban.UserID,
			Username:    username.String,
			DisplayName: displayName.String,
		}
	}

	ban.BannedByUser = &models.User{
		ID:          ban.BannedBy,
		Username:    bannedByUsername.String,
		DisplayName: bannedByDisplayName.String,
	}

	if ban.GameID != nil {
		ban.Game = &models.Game{
			ID:   *ban.GameID,
			Name: gameName.String,
		}
	}

	return ban, nil
}

func (s *AdminService) ListBans(page, perPage int, active bool, gameID *int) ([]models.Ban, int64, error) {
	offset := (page - 1) * perPage

	whereClause := "WHERE b.is_active = ?"
	args := []interface{}{active}

	if gameID != nil {
		whereClause += " AND b.game_id = ?"
		args = append(args, *gameID)
	}

	var total int64
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM bans b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN users bu ON b.banned_by = bu.id
		LEFT JOIN games g ON b.game_id = g.id
		%s
	`, whereClause)

	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bans: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT b.id, b.user_id, b.ip_address, b.game_id, b.reason, b.banned_by,
		       b.expires_at, b.is_active, b.created_at,
		       u.username, u.display_name, bu.username, bu.display_name, g.name
		FROM bans b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN users bu ON b.banned_by = bu.id
		LEFT JOIN games g ON b.game_id = g.id
		%s
		ORDER BY b.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, perPage, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bans: %w", err)
	}
	defer rows.Close()

	var bans []models.Ban
	for rows.Next() {
		var ban models.Ban
		var username, displayName, bannedByUsername, bannedByDisplayName sql.NullString
		var gameName sql.NullString

		err := rows.Scan(
			&ban.ID, &ban.UserID, &ban.IPAddress, &ban.GameID, &ban.Reason,
			&ban.BannedBy, &ban.ExpiresAt, &ban.IsActive, &ban.CreatedAt,
			&username, &displayName, &bannedByUsername, &bannedByDisplayName, &gameName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan ban: %w", err)
		}

		if ban.UserID != nil {
			ban.User = &models.User{
				ID:          *ban.UserID,
				Username:    username.String,
				DisplayName: displayName.String,
			}
		}

		ban.BannedByUser = &models.User{
			ID:          ban.BannedBy,
			Username:    bannedByUsername.String,
			DisplayName: bannedByDisplayName.String,
		}

		if ban.GameID != nil {
			ban.Game = &models.Game{
				ID:   *ban.GameID,
				Name: gameName.String,
			}
		}

		bans = append(bans, ban)
	}

	return bans, total, nil
}

func (s *AdminService) UnbanUser(banID int) error {
	result, err := s.db.Exec(`
		UPDATE bans SET is_active = 0
		WHERE id = ? AND is_active = 1
	`, banID)

	if err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("ban not found or already inactive")
	}

	return nil
}

func (s *AdminService) IsUserBanned(userID int, gameID *int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM bans
		WHERE user_id = ? AND is_active = 1 AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	args := []interface{}{userID}

	if gameID != nil {
		query += " AND (game_id IS NULL OR game_id = ?)"
		args = append(args, *gameID)
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ban status: %w", err)
	}

	return count > 0, nil
}

func (s *AdminService) IsIPBanned(ipAddress string, gameID *int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM bans
		WHERE ip_address = ? AND is_active = 1 AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
	`
	args := []interface{}{ipAddress}

	if gameID != nil {
		query += " AND (game_id IS NULL OR game_id = ?)"
		args = append(args, *gameID)
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check IP ban status: %w", err)
	}

	return count > 0, nil
}

func (s *AdminService) GetUserStats() (map[string]int, error) {
	stats := make(map[string]int)

	var totalUsers int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_active = 1").Scan(&totalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}
	stats["total_users"] = totalUsers

	var newUsers int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM users
		WHERE is_active = 1 AND created_at >= datetime('now', '-30 days')
	`).Scan(&newUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get new users: %w", err)
	}
	stats["new_users_30_days"] = newUsers

	var bannedUsers int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM bans
		WHERE is_active = 1 AND user_id IS NOT NULL
	`).Scan(&bannedUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get banned users: %w", err)
	}
	stats["banned_users"] = bannedUsers

	var bannedIPs int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM bans
		WHERE is_active = 1 AND ip_address IS NOT NULL AND ip_address != ''
	`).Scan(&bannedIPs)
	if err != nil {
		return nil, fmt.Errorf("failed to get banned IPs: %w", err)
	}
	stats["banned_ips"] = bannedIPs

	return stats, nil
}

func (s *AdminService) GetModStats() (map[string]int, error) {
	stats := make(map[string]int)

	var totalMods int
	err := s.db.QueryRow("SELECT COUNT(*) FROM mods").Scan(&totalMods)
	if err != nil {
		return nil, fmt.Errorf("failed to get total mods: %w", err)
	}
	stats["total_mods"] = totalMods

	var activeMods int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM mods
		WHERE is_scanned = 1 AND scan_result = 'clean' AND is_rejected = 0
	`).Scan(&activeMods)
	if err != nil {
		return nil, fmt.Errorf("failed to get active mods: %w", err)
	}
	stats["active_mods"] = activeMods

	var rejectedMods int
	err = s.db.QueryRow("SELECT COUNT(*) FROM mods WHERE is_rejected = 1").Scan(&rejectedMods)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejected mods: %w", err)
	}
	stats["rejected_mods"] = rejectedMods

	var pendingMods int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM mods
		WHERE is_scanned = 0 OR scan_result = 'pending'
	`).Scan(&pendingMods)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending mods: %w", err)
	}
	stats["pending_mods"] = pendingMods

	var newMods int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM mods
		WHERE created_at >= datetime('now', '-7 days')
	`).Scan(&newMods)
	if err != nil {
		return nil, fmt.Errorf("failed to get new mods: %w", err)
	}
	stats["new_mods_7_days"] = newMods

	return stats, nil
}

func (s *AdminService) GetSystemStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalDownloads int
	err := s.db.QueryRow("SELECT SUM(downloads) FROM mods").Scan(&totalDownloads)
	if err != nil {
		return nil, fmt.Errorf("failed to get total downloads: %w", err)
	}
	stats["total_downloads"] = totalDownloads

	var totalLikes int
	err = s.db.QueryRow("SELECT SUM(likes) FROM mods").Scan(&totalLikes)
	if err != nil {
		return nil, fmt.Errorf("failed to get total likes: %w", err)
	}
	stats["total_likes"] = totalLikes

	var totalComments int
	err = s.db.QueryRow("SELECT COUNT(*) FROM comments WHERE is_active = 1").Scan(&totalComments)
	if err != nil {
		return nil, fmt.Errorf("failed to get total comments: %w", err)
	}
	stats["total_comments"] = totalComments

	var totalGames int
	err = s.db.QueryRow("SELECT COUNT(*) FROM games WHERE is_active = 1").Scan(&totalGames)
	if err != nil {
		return nil, fmt.Errorf("failed to get total games: %w", err)
	}
	stats["total_games"] = totalGames

	var pendingGameRequests int
	err = s.db.QueryRow("SELECT COUNT(*) FROM game_requests WHERE status = 'pending'").Scan(&pendingGameRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending game requests: %w", err)
	}
	stats["pending_game_requests"] = pendingGameRequests

	return stats, nil
}

func (s *AdminService) UpdateUserRole(userID int, role string) error {
	validRoles := []string{models.RoleUser, models.RoleAdmin, models.RoleCommunityModerator, models.RoleWikiMaintainer}
	isValid := false
	for _, validRole := range validRoles {
		if role == validRole {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("invalid role")
	}

	_, err := s.db.Exec("UPDATE users SET role = ? WHERE id = ?", role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

func (s *AdminService) DeactivateUser(userID int) error {
	_, err := s.db.Exec("UPDATE users SET is_active = 0 WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

func (s *AdminService) ActivateUser(userID int) error {
	_, err := s.db.Exec("UPDATE users SET is_active = 1 WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

func (s *AdminService) GetRecentActivity(limit int) ([]map[string]interface{}, error) {
	activities := make([]map[string]interface{}, 0)

	rows, err := s.db.Query(`
		SELECT 'mod_upload' as type, m.name, m.created_at, u.display_name
		FROM mods m
		JOIN users u ON m.owner_id = u.id
		WHERE m.created_at >= datetime('now', '-7 days')
		UNION ALL
		SELECT 'user' as type, u.display_name, u.created_at, u.display_name
		FROM users u
		WHERE u.created_at >= datetime('now', '-7 days')
		UNION ALL
		SELECT 'game_request' as type, gr.name, gr.created_at, u.display_name
		FROM game_requests gr
		JOIN users u ON gr.requested_by = u.id
		WHERE gr.created_at >= datetime('now', '-7 days')
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var activityType, name, createdAt string
		var userName sql.NullString

		err := rows.Scan(&activityType, &name, &createdAt, &userName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		activity := map[string]interface{}{
			"type":       activityType,
			"name":       name,
			"created_at": createdAt,
		}

		if userName.Valid {
			activity["user"] = userName.String
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// ListPendingMods returns a paginated list of mods pending review
func (s *AdminService) ListPendingMods(page, perPage int) ([]models.Mod, int64, error) {
	offset := (page - 1) * perPage

	// Count total pending mods
	var total int64
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM mods
		WHERE is_scanned = 0 AND is_rejected = 0
	`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending mods: %w", err)
	}

	// Get pending mods with related data
	rows, err := s.db.Query(`
		SELECT m.id, m.name, m.slug, m.description, m.short_description, m.icon,
		       m.version, m.game_version, m.game_id, m.owner_id, m.downloads, m.likes,
		       m.source_website, m.contact_info, m.is_rejected, m.rejection_reason,
		       m.is_scanned, m.scan_result, m.created_at, m.updated_at,
		       g.name as game_name, g.slug as game_slug,
		       u.username, u.display_name
		FROM mods m
		LEFT JOIN games g ON m.game_id = g.id
		LEFT JOIN users u ON m.owner_id = u.id
		WHERE m.is_scanned = 0 AND m.is_rejected = 0
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pending mods: %w", err)
	}
	defer rows.Close()

	var mods []models.Mod
	for rows.Next() {
		var mod models.Mod
		var gameName, gameSlug, username, displayName sql.NullString

		err := rows.Scan(
			&mod.ID, &mod.Name, &mod.Slug, &mod.Description, &mod.ShortDescription, &mod.Icon,
			&mod.Version, &mod.GameVersion, &mod.GameID, &mod.OwnerID, &mod.Downloads, &mod.Likes,
			&mod.SourceWebsite, &mod.ContactInfo, &mod.IsRejected, &mod.RejectionReason,
			&mod.IsScanned, &mod.ScanResult, &mod.CreatedAt, &mod.UpdatedAt,
			&gameName, &gameSlug, &username, &displayName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mod: %w", err)
		}

		// Add game info if available
		if gameName.Valid {
			mod.Game = &models.Game{
				ID:   mod.GameID,
				Name: gameName.String,
				Slug: gameSlug.String,
			}
		}

		// Add owner info if available
		if username.Valid {
			mod.Owner = &models.User{
				ID:          mod.OwnerID,
				Username:    username.String,
				DisplayName: displayName.String,
			}
		}

		mods = append(mods, mod)
	}

	return mods, total, nil
}

func (s *AdminService) CleanupExpiredBans() error {
	_, err := s.db.Exec(`
		UPDATE bans SET is_active = 0
		WHERE is_active = 1 AND expires_at IS NOT NULL AND expires_at <= CURRENT_TIMESTAMP
	`)

	if err != nil {
		return fmt.Errorf("failed to cleanup expired bans: %w", err)
	}

	return nil
}

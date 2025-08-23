package services

import (
	"testing"

	"github.com/azurite/backend/internal/models"
)

func TestAdminService_BanUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.BanCreateRequest{
		UserID:   &userID,
		GameID:   &gameID,
		Reason:   "Inappropriate behavior",
		Duration: 7, // 7 days
	}

	ban, err := service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("BanUser() error = %v", err)
	}

	if ban.UserID == nil || *ban.UserID != userID {
		t.Errorf("Expected user ID %d, got %v", userID, ban.UserID)
	}
	if ban.GameID == nil || *ban.GameID != gameID {
		t.Errorf("Expected game ID %d, got %v", gameID, ban.GameID)
	}
	if ban.Reason != req.Reason {
		t.Errorf("Expected reason %s, got %s", req.Reason, ban.Reason)
	}
	if ban.BannedBy != adminID {
		t.Errorf("Expected banned_by %d, got %d", adminID, ban.BannedBy)
	}
	if !ban.IsActive {
		t.Error("Expected ban to be active")
	}
	if ban.ExpiresAt == nil {
		t.Error("Expected expiration date to be set")
	}
}

func TestAdminService_BanUserPermanent(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")

	req := &models.BanCreateRequest{
		UserID:   &userID,
		Reason:   "Severe violation",
		Duration: 0, // Permanent ban
	}

	ban, err := service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("BanUser() permanent error = %v", err)
	}

	if ban.ExpiresAt != nil {
		t.Error("Expected permanent ban to have no expiration date")
	}
}

func TestAdminService_BanIP(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	adminID := createTestUser(t, db, "admin", "admin@example.com")

	req := &models.BanCreateRequest{
		IPAddress: "192.168.1.100",
		Reason:    "Spam from this IP",
		Duration:  1, // 1 day
	}

	ban, err := service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("BanUser() IP error = %v", err)
	}

	if ban.IPAddress != req.IPAddress {
		t.Errorf("Expected IP address %s, got %s", req.IPAddress, ban.IPAddress)
	}
	if ban.UserID != nil {
		t.Errorf("Expected no user ID for IP ban, got %v", ban.UserID)
	}
}

func TestAdminService_BanUserInvalidRequest(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	adminID := createTestUser(t, db, "admin", "admin@example.com")

	req := &models.BanCreateRequest{
		// No UserID or IPAddress
		Reason:   "Invalid request",
		Duration: 1,
	}

	_, err := service.BanUser(req, adminID)
	if err == nil {
		t.Error("Expected error for ban request without user ID or IP address")
	}
}

func TestAdminService_GetBanByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")

	// Create ban
	req := &models.BanCreateRequest{
		UserID:   &userID,
		Reason:   "Test ban",
		Duration: 1,
	}

	createdBan, err := service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("Failed to create ban: %v", err)
	}

	// Get ban by ID
	ban, err := service.GetBanByID(createdBan.ID)
	if err != nil {
		t.Fatalf("GetBanByID() error = %v", err)
	}

	if ban.ID != createdBan.ID {
		t.Errorf("Expected ID %d, got %d", createdBan.ID, ban.ID)
	}
	if ban.User == nil {
		t.Error("Expected user to be populated")
	}
	if ban.BannedByUser == nil {
		t.Error("Expected banned_by user to be populated")
	}
}

func TestAdminService_ListBans(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create test bans
	for i := 1; i <= 3; i++ {
		req := &models.BanCreateRequest{
			UserID:   &userID,
			GameID:   &gameID,
			Reason:   "Test ban",
			Duration: 1,
		}
		_, err := service.BanUser(req, adminID)
		if err != nil {
			t.Fatalf("Failed to create ban %d: %v", i, err)
		}
	}

	bans, total, err := service.ListBans(1, 10, true, nil)
	if err != nil {
		t.Fatalf("ListBans() error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(bans) != 3 {
		t.Errorf("Expected 3 bans, got %d", len(bans))
	}

	// Test filtering by game
	bans, total, err = service.ListBans(1, 10, true, &gameID)
	if err != nil {
		t.Fatalf("ListBans() with game filter error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3 for game filter, got %d", total)
	}
}

func TestAdminService_UnbanUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")

	// Create ban
	req := &models.BanCreateRequest{
		UserID:   &userID,
		Reason:   "Test ban",
		Duration: 7,
	}

	ban, err := service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("Failed to create ban: %v", err)
	}

	// Unban user
	err = service.UnbanUser(ban.ID)
	if err != nil {
		t.Fatalf("UnbanUser() error = %v", err)
	}

	// Verify ban is inactive
	updatedBan, err := service.GetBanByID(ban.ID)
	if err != nil {
		t.Fatalf("GetBanByID() after unban error = %v", err)
	}

	if updatedBan.IsActive {
		t.Error("Expected ban to be inactive after unbanning")
	}
}

func TestAdminService_UnbanNonExistentBan(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)

	err := service.UnbanUser(999)
	if err == nil {
		t.Error("Expected error when unbanning non-existent ban")
	}
}

func TestAdminService_IsUserBanned(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Initially not banned
	banned, err := service.IsUserBanned(userID, nil)
	if err != nil {
		t.Fatalf("IsUserBanned() error = %v", err)
	}
	if banned {
		t.Error("Expected user to not be banned initially")
	}

	// Create ban
	req := &models.BanCreateRequest{
		UserID:   &userID,
		GameID:   &gameID,
		Reason:   "Test ban",
		Duration: 7,
	}

	_, err = service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("Failed to create ban: %v", err)
	}

	// Now should be banned
	banned, err = service.IsUserBanned(userID, &gameID)
	if err != nil {
		t.Fatalf("IsUserBanned() after ban error = %v", err)
	}
	if !banned {
		t.Error("Expected user to be banned")
	}

	// Global ban check
	banned, err = service.IsUserBanned(userID, nil)
	if err != nil {
		t.Fatalf("IsUserBanned() global check error = %v", err)
	}
	if !banned {
		t.Error("Expected user to be globally banned")
	}
}

func TestAdminService_IsIPBanned(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	adminID := createTestUser(t, db, "admin", "admin@example.com")
	testIP := "192.168.1.100"

	// Initially not banned
	banned, err := service.IsIPBanned(testIP, nil)
	if err != nil {
		t.Fatalf("IsIPBanned() error = %v", err)
	}
	if banned {
		t.Error("Expected IP to not be banned initially")
	}

	// Create IP ban
	req := &models.BanCreateRequest{
		IPAddress: testIP,
		Reason:    "Spam IP",
		Duration:  1,
	}

	_, err = service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("Failed to create IP ban: %v", err)
	}

	// Now should be banned
	banned, err = service.IsIPBanned(testIP, nil)
	if err != nil {
		t.Fatalf("IsIPBanned() after ban error = %v", err)
	}
	if !banned {
		t.Error("Expected IP to be banned")
	}
}

func TestAdminService_GetUserStats(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)

	// Create test users
	user1ID := createTestUser(t, db, "user1", "user1@example.com")
	_ = createTestUser(t, db, "user2", "user2@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")

	// Ban one user
	req := &models.BanCreateRequest{
		UserID:   &user1ID,
		Reason:   "Test ban",
		Duration: 7,
	}
	_, err := service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("Failed to create ban: %v", err)
	}

	stats, err := service.GetUserStats()
	if err != nil {
		t.Fatalf("GetUserStats() error = %v", err)
	}

	if stats["total_users"] != 3 {
		t.Errorf("Expected 3 total users, got %d", stats["total_users"])
	}
	if stats["banned_users"] != 1 {
		t.Errorf("Expected 1 banned user, got %d", stats["banned_users"])
	}

	// Check that required keys exist
	requiredKeys := []string{"total_users", "new_users_30_days", "banned_users", "banned_ips"}
	for _, key := range requiredKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected stats to contain key %s", key)
		}
	}
}

func TestAdminService_GetModStats(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create test mods
	modID1 := createTestMod(t, db, "Mod 1", "mod-1", gameID, userID)
	modID2 := createTestMod(t, db, "Mod 2", "mod-2", gameID, userID)

	// Mark one mod as scanned and clean
	_, err := db.Exec("UPDATE mods SET is_scanned = 1, scan_result = 'clean' WHERE id = ?", modID1)
	if err != nil {
		t.Fatalf("Failed to update mod scan status: %v", err)
	}

	// Mark one mod as rejected
	_, err = db.Exec("UPDATE mods SET is_rejected = 1 WHERE id = ?", modID2)
	if err != nil {
		t.Fatalf("Failed to update mod rejection status: %v", err)
	}

	stats, err := service.GetModStats()
	if err != nil {
		t.Fatalf("GetModStats() error = %v", err)
	}

	if stats["total_mods"] != 2 {
		t.Errorf("Expected 2 total mods, got %d", stats["total_mods"])
	}
	if stats["active_mods"] != 1 {
		t.Errorf("Expected 1 active mod, got %d", stats["active_mods"])
	}
	if stats["rejected_mods"] != 1 {
		t.Errorf("Expected 1 rejected mod, got %d", stats["rejected_mods"])
	}

	// Check that required keys exist
	requiredKeys := []string{"total_mods", "active_mods", "rejected_mods", "pending_mods", "new_mods_7_days"}
	for _, key := range requiredKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected stats to contain key %s", key)
		}
	}
}

func TestAdminService_GetSystemStats(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Update mod with some stats
	_, err := db.Exec("UPDATE mods SET downloads = 100, likes = 50 WHERE id = ?", modID)
	if err != nil {
		t.Fatalf("Failed to update mod stats: %v", err)
	}

	// Create a comment
	_, err = db.Exec(`
		INSERT INTO comments (mod_id, user_id, content, is_active)
		VALUES (?, ?, ?, ?)
	`, modID, userID, "Test comment", true)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	stats, err := service.GetSystemStats()
	if err != nil {
		t.Fatalf("GetSystemStats() error = %v", err)
	}

	if stats["total_downloads"] != 100 {
		t.Errorf("Expected 100 total downloads, got %v", stats["total_downloads"])
	}
	if stats["total_likes"] != 50 {
		t.Errorf("Expected 50 total likes, got %v", stats["total_likes"])
	}
	if stats["total_comments"] != 1 {
		t.Errorf("Expected 1 total comment, got %v", stats["total_comments"])
	}
	if stats["total_games"] != 1 {
		t.Errorf("Expected 1 total game, got %v", stats["total_games"])
	}

	// Check that required keys exist
	requiredKeys := []string{"total_downloads", "total_likes", "total_comments", "total_games", "pending_game_requests"}
	for _, key := range requiredKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected stats to contain key %s", key)
		}
	}
}

func TestAdminService_UpdateUserRole(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Update user role
	err := service.UpdateUserRole(userID, models.RoleAdmin)
	if err != nil {
		t.Fatalf("UpdateUserRole() error = %v", err)
	}

	// Verify role was updated
	var role string
	err = db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		t.Fatalf("Failed to get user role: %v", err)
	}

	if role != models.RoleAdmin {
		t.Errorf("Expected role %s, got %s", models.RoleAdmin, role)
	}
}

func TestAdminService_UpdateUserRoleInvalid(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	err := service.UpdateUserRole(userID, "invalid_role")
	if err == nil {
		t.Error("Expected error for invalid role")
	}
}

func TestAdminService_DeactivateUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	err := service.DeactivateUser(userID)
	if err != nil {
		t.Fatalf("DeactivateUser() error = %v", err)
	}

	// Verify user is deactivated
	var isActive bool
	err = db.QueryRow("SELECT is_active FROM users WHERE id = ?", userID).Scan(&isActive)
	if err != nil {
		t.Fatalf("Failed to get user active status: %v", err)
	}

	if isActive {
		t.Error("Expected user to be deactivated")
	}
}

func TestAdminService_ActivateUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// First deactivate
	err := service.DeactivateUser(userID)
	if err != nil {
		t.Fatalf("DeactivateUser() error = %v", err)
	}

	// Then reactivate
	err = service.ActivateUser(userID)
	if err != nil {
		t.Fatalf("ActivateUser() error = %v", err)
	}

	// Verify user is activated
	var isActive bool
	err = db.QueryRow("SELECT is_active FROM users WHERE id = ?", userID).Scan(&isActive)
	if err != nil {
		t.Fatalf("Failed to get user active status: %v", err)
	}

	if !isActive {
		t.Error("Expected user to be activated")
	}
}

func TestAdminService_GetRecentActivity(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	_ = createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create a game request
	_, err := db.Exec(`
		INSERT INTO game_requests (name, description, requested_by, status)
		VALUES (?, ?, ?, ?)
	`, "New Game", "Please add this game", userID, "pending")
	if err != nil {
		t.Fatalf("Failed to create game request: %v", err)
	}

	activities, err := service.GetRecentActivity(10)
	if err != nil {
		t.Fatalf("GetRecentActivity() error = %v", err)
	}

	// Should have at least some activities (mod, user, game request)
	if len(activities) == 0 {
		t.Error("Expected some recent activities")
	}

	// Check that activities have required fields
	for _, activity := range activities {
		if activity["type"] == nil {
			t.Error("Expected activity to have type")
		}
		if activity["name"] == nil {
			t.Error("Expected activity to have name")
		}
		if activity["created_at"] == nil {
			t.Error("Expected activity to have created_at")
		}
	}
}

func TestAdminService_ListPendingMods(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create pending mods (not scanned)
	for i := 1; i <= 3; i++ {
		modID := createTestMod(t, db, "Pending Mod", "pending-mod", gameID, userID)
		// Ensure mods are not scanned (they should be by default)
		_, err := db.Exec("UPDATE mods SET is_scanned = 0 WHERE id = ?", modID)
		if err != nil {
			t.Fatalf("Failed to update mod scan status: %v", err)
		}
	}

	// Create a processed mod
	processedModID := createTestMod(t, db, "Processed Mod", "processed-mod", gameID, userID)
	_, err := db.Exec("UPDATE mods SET is_scanned = 1, scan_result = 'clean' WHERE id = ?", processedModID)
	if err != nil {
		t.Fatalf("Failed to update processed mod: %v", err)
	}

	mods, total, err := service.ListPendingMods(1, 10)
	if err != nil {
		t.Fatalf("ListPendingMods() error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected 3 pending mods, got %d", total)
	}
	if len(mods) != 3 {
		t.Errorf("Expected 3 mods in list, got %d", len(mods))
	}

	// All returned mods should be pending (not scanned)
	for _, mod := range mods {
		if mod.IsScanned {
			t.Error("Expected pending mods to not be scanned")
		}
		if mod.IsRejected {
			t.Error("Expected pending mods to not be rejected")
		}
		if mod.Game == nil {
			t.Error("Expected mod to have game info populated")
		}
		if mod.Owner == nil {
			t.Error("Expected mod to have owner info populated")
		}
	}
}

func TestAdminService_CleanupExpiredBans(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewAdminService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")

	// Create ban that will expire
	req := &models.BanCreateRequest{
		UserID:   &userID,
		Reason:   "Test ban",
		Duration: 1,
	}

	ban, err := service.BanUser(req, adminID)
	if err != nil {
		t.Fatalf("Failed to create ban: %v", err)
	}

	// Manually set ban to expired
	_, err = db.Exec(`
		UPDATE bans SET expires_at = datetime('now', '-1 day')
		WHERE id = ?
	`, ban.ID)
	if err != nil {
		t.Fatalf("Failed to set ban as expired: %v", err)
	}

	// Cleanup expired bans
	err = service.CleanupExpiredBans()
	if err != nil {
		t.Fatalf("CleanupExpiredBans() error = %v", err)
	}

	// Verify ban is now inactive
	updatedBan, err := service.GetBanByID(ban.ID)
	if err != nil {
		t.Fatalf("GetBanByID() after cleanup error = %v", err)
	}

	if updatedBan.IsActive {
		t.Error("Expected expired ban to be inactive after cleanup")
	}
}

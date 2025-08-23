package services

import (
	"testing"

	"github.com/azurite/backend/internal/auth"
	"github.com/azurite/backend/internal/models"
)

func TestUserService_Register(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	req := &models.RegisterRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}

	user, err := service.Register(req)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, user.Username)
	}
	if user.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, user.Email)
	}
	if user.DisplayName != req.DisplayName {
		t.Errorf("Expected display name %s, got %s", req.DisplayName, user.DisplayName)
	}
	if user.Role != models.RoleUser {
		t.Errorf("Expected role %s, got %s", models.RoleUser, user.Role)
	}
	if !user.IsActive {
		t.Error("Expected user to be active")
	}
	if user.PasswordHash == "" {
		t.Error("Expected password hash to be set")
	}
	if user.PasswordHash == req.Password {
		t.Error("Password hash should not be the same as password")
	}
}

func TestUserService_RegisterDuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// Create first user
	req1 := &models.RegisterRequest{
		Username:    "testuser",
		Email:       "test1@example.com",
		Password:    "password123",
		DisplayName: "Test User 1",
	}

	_, err := service.Register(req1)
	if err != nil {
		t.Fatalf("First Register() error = %v", err)
	}

	// Try to create user with same username
	req2 := &models.RegisterRequest{
		Username:    "testuser", // Same username
		Email:       "test2@example.com",
		Password:    "password123",
		DisplayName: "Test User 2",
	}

	_, err = service.Register(req2)
	if err == nil {
		t.Error("Expected error for duplicate username")
	}
}

func TestUserService_RegisterDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// Create first user
	req1 := &models.RegisterRequest{
		Username:    "testuser1",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User 1",
	}

	_, err := service.Register(req1)
	if err != nil {
		t.Fatalf("First Register() error = %v", err)
	}

	// Try to create user with same email
	req2 := &models.RegisterRequest{
		Username:    "testuser2",
		Email:       "test@example.com", // Same email
		Password:    "password123",
		DisplayName: "Test User 2",
	}

	_, err = service.Register(req2)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
}

func TestUserService_Login(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// Register a user first
	registerReq := &models.RegisterRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}

	user, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Login with correct credentials
	loginReq := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	authResponse, err := service.Login(loginReq)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if authResponse.Token == "" {
		t.Error("Expected non-empty token")
	}
	if authResponse.User.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, authResponse.User.ID)
	}
}

func TestUserService_LoginInvalidCredentials(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// Register a user first
	registerReq := &models.RegisterRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}

	_, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"Wrong password", "test@example.com", "wrongpassword"},
		{"Wrong email", "wrong@example.com", "password123"},
		{"Both wrong", "wrong@example.com", "wrongpassword"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginReq := &models.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			_, err := service.Login(loginReq)
			if err == nil {
				t.Error("Expected error for invalid credentials")
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	userID := createTestUser(t, db, "testuser", "test@example.com")

	user, err := service.GetByID(userID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if user.ID != userID {
		t.Errorf("Expected ID %d, got %d", userID, user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
}

func TestUserService_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	_, err := service.GetByID(999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	userID := createTestUser(t, db, "testuser", "test@example.com")

	user, err := service.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() error = %v", err)
	}

	if user.ID != userID {
		t.Errorf("Expected ID %d, got %d", userID, user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", user.Email)
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	userID := createTestUser(t, db, "testuser", "test@example.com")

	user, err := service.GetByUsername("testuser")
	if err != nil {
		t.Fatalf("GetByUsername() error = %v", err)
	}

	if user.ID != userID {
		t.Errorf("Expected ID %d, got %d", userID, user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
}

func TestUserService_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	userID := createTestUser(t, db, "testuser", "test@example.com")

	updateReq := &models.UserUpdateRequest{
		DisplayName:  "Updated Display Name",
		Bio:          "Updated bio",
		NotifyEmail:  false,
		NotifyInSite: false,
	}

	updatedUser, err := service.Update(userID, updateReq)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updatedUser.DisplayName != updateReq.DisplayName {
		t.Errorf("Expected display name %s, got %s", updateReq.DisplayName, updatedUser.DisplayName)
	}
	if updatedUser.Bio != updateReq.Bio {
		t.Errorf("Expected bio %s, got %s", updateReq.Bio, updatedUser.Bio)
	}
	if updatedUser.NotifyEmail != updateReq.NotifyEmail {
		t.Errorf("Expected notify_email %v, got %v", updateReq.NotifyEmail, updatedUser.NotifyEmail)
	}
	if updatedUser.NotifyInSite != updateReq.NotifyInSite {
		t.Errorf("Expected notify_in_site %v, got %v", updateReq.NotifyInSite, updatedUser.NotifyInSite)
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// Register user with known password
	registerReq := &models.RegisterRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "oldpassword123",
		DisplayName: "Test User",
	}

	user, err := service.Register(registerReq)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Update password
	err = service.UpdatePassword(user.ID, "oldpassword123", "newpassword123")
	if err != nil {
		t.Fatalf("UpdatePassword() error = %v", err)
	}

	// Verify old password no longer works
	loginReq1 := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "oldpassword123",
	}

	_, err = service.Login(loginReq1)
	if err == nil {
		t.Error("Expected error when logging in with old password")
	}

	// Verify new password works
	loginReq2 := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "newpassword123",
	}

	_, err = service.Login(loginReq2)
	if err != nil {
		t.Errorf("Expected successful login with new password, got error: %v", err)
	}
}

func TestUserService_UpdatePasswordWrongCurrent(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	userID := createTestUser(t, db, "testuser", "test@example.com")

	err := service.UpdatePassword(userID, "wrongcurrentpassword", "newpassword123")
	if err == nil {
		t.Error("Expected error when providing wrong current password")
	}
}

func TestUserService_RequestPasswordReset(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	userID := createTestUser(t, db, "testuser", "test@example.com")

	err := service.RequestPasswordReset("test@example.com")
	if err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}

	// Verify reset token was created
	var tokenCount int
	err = db.QueryRow("SELECT COUNT(*) FROM password_reset_tokens WHERE user_id = ?", userID).Scan(&tokenCount)
	if err != nil {
		t.Fatalf("Failed to check token count: %v", err)
	}

	if tokenCount != 1 {
		t.Errorf("Expected 1 reset token, got %d", tokenCount)
	}
}

func TestUserService_RequestPasswordResetNonExistentEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// This should not return an error (for security reasons)
	err := service.RequestPasswordReset("nonexistent@example.com")
	if err != nil {
		t.Errorf("RequestPasswordReset() with non-existent email should not error, got: %v", err)
	}
}

func TestUserService_ResetPassword(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// Create user and request password reset
	userID := createTestUser(t, db, "testuser", "test@example.com")
	err := service.RequestPasswordReset("test@example.com")
	if err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}

	// Get the reset token
	var token string
	err = db.QueryRow("SELECT token FROM password_reset_tokens WHERE user_id = ?", userID).Scan(&token)
	if err != nil {
		t.Fatalf("Failed to get reset token: %v", err)
	}

	// Reset password
	newPassword := "newpassword123"
	err = service.ResetPassword(token, newPassword)
	if err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}

	// Verify token was marked as used
	var used bool
	err = db.QueryRow("SELECT used FROM password_reset_tokens WHERE token = ?", token).Scan(&used)
	if err != nil {
		t.Fatalf("Failed to check token used status: %v", err)
	}

	if !used {
		t.Error("Expected reset token to be marked as used")
	}

	// Verify new password works
	loginReq := &models.LoginRequest{
		Email:    "test@example.com",
		Password: newPassword,
	}

	_, err = service.Login(loginReq)
	if err != nil {
		t.Errorf("Expected successful login with new password, got error: %v", err)
	}
}

func TestUserService_ResetPasswordInvalidToken(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	err := service.ResetPassword("invalid_token", "newpassword123")
	if err == nil {
		t.Error("Expected error for invalid reset token")
	}
}

func TestUserService_GetUserRoles(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Assign a game-specific role
	_, err := db.Exec(`
		INSERT INTO user_roles (user_id, game_id, role)
		VALUES (?, ?, ?)
	`, userID, gameID, models.RoleCommunityModerator)
	if err != nil {
		t.Fatalf("Failed to insert user role: %v", err)
	}

	roles, err := service.GetUserRoles(userID)
	if err != nil {
		t.Fatalf("GetUserRoles() error = %v", err)
	}

	if len(roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(roles))
	}

	role := roles[0]
	if role.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, role.UserID)
	}
	if role.GameID == nil || *role.GameID != gameID {
		t.Errorf("Expected game ID %d, got %v", gameID, role.GameID)
	}
	if role.Role != models.RoleCommunityModerator {
		t.Errorf("Expected role %s, got %s", models.RoleCommunityModerator, role.Role)
	}
}

func TestUserService_List(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	// Create test users
	createTestUser(t, db, "user1", "user1@example.com")
	createTestUser(t, db, "user2", "user2@example.com")
	createTestUser(t, db, "user3", "user3@example.com")

	users, total, err := service.List(1, 2)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users per page, got %d", len(users))
	}

	// Test second page
	users, total, err = service.List(2, 2)
	if err != nil {
		t.Fatalf("List() page 2 error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user on page 2, got %d", len(users))
	}
}

func TestUserService_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	config := setupTestConfig()
	authService := auth.NewService(config)
	notificationService := NewNotificationService(db, config)
	service := NewUserService(db, authService, notificationService)

	user := &models.User{
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Role:        models.RoleAdmin,
	}
	_, err := service.Create(user)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", user.Email)
	}
	if user.DisplayName != "Test User" {
		t.Errorf("Expected display name 'Test User', got %s", user.DisplayName)
	}
	if user.Role != models.RoleAdmin {
		t.Errorf("Expected role %s, got %s", models.RoleAdmin, user.Role)
	}
	if !user.IsActive {
		t.Error("Expected user to be active")
	}

	// Verify password can be checked (random password should have been set)
	if user.PasswordHash == "" {
		t.Error("Expected password hash to be set")
	}
}

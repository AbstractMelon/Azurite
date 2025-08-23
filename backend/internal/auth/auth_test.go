package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/azurite/backend/internal/config"
	"github.com/azurite/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func TestNewService(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:     "test-secret",
			GitHubID:      "github-id",
			GitHubSecret:  "github-secret",
			GoogleID:      "google-id",
			GoogleSecret:  "google-secret",
			DiscordID:     "discord-id",
			DiscordSecret: "discord-secret",
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	service := NewService(cfg)

	if service.config != cfg {
		t.Error("Expected config to be set")
	}

	if service.githubConfig == nil {
		t.Error("Expected GitHub config to be initialized")
	}

	if service.googleConfig == nil {
		t.Error("Expected Google config to be initialized")
	}

	if service.discordConfig == nil {
		t.Error("Expected Discord config to be initialized")
	}
}

func TestNewServiceWithMissingCredentials(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
			// Missing OAuth credentials
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	service := NewService(cfg)

	if service.githubConfig != nil {
		t.Error("Expected GitHub config to be nil when credentials missing")
	}

	if service.googleConfig != nil {
		t.Error("Expected Google config to be nil when credentials missing")
	}

	if service.discordConfig != nil {
		t.Error("Expected Discord config to be nil when credentials missing")
	}
}

func TestGenerateToken(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret-key-for-jwt",
		},
	}
	service := NewService(cfg)

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleUser,
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	// Verify token can be parsed
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		t.Fatalf("Expected token to be valid, got error: %v", err)
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		t.Fatal("Expected valid claims")
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}

	if claims.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, claims.Username)
	}

	if claims.Role != user.Role {
		t.Errorf("Expected Role %s, got %s", user.Role, claims.Role)
	}

	if claims.Issuer != "azurite" {
		t.Errorf("Expected Issuer 'azurite', got %s", claims.Issuer)
	}
}

func TestValidateToken(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret-key-for-jwt",
		},
	}
	service := NewService(cfg)

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     models.RoleUser,
	}

	// Generate a valid token
	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Expected token to be valid, got error: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}

	if claims.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, claims.Username)
	}

	if claims.Role != user.Role {
		t.Errorf("Expected Role %s, got %s", user.Role, claims.Role)
	}
}

func TestValidateTokenWithInvalidToken(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret-key-for-jwt",
		},
	}
	service := NewService(cfg)

	tests := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Invalid token", "invalid.token.here"},
		{"Malformed token", "not.a.jwt"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := service.ValidateToken(test.token)
			if err == nil {
				t.Error("Expected error for invalid token")
			}
		})
	}
}

func TestValidateTokenWithExpiredToken(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret-key-for-jwt",
		},
	}

	// Create an expired token manually
	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		Role:     models.RoleUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "azurite",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.Auth.JWTSecret))

	service := NewService(cfg)
	_, err := service.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestHasPermission(t *testing.T) {
	service := &Service{}

	userRoles := []models.UserRole{
		{UserID: 1, GameID: &[]int{1}[0], Role: models.RoleCommunityModerator},
		{UserID: 1, GameID: &[]int{2}[0], Role: models.RoleWikiMaintainer},
	}

	tests := []struct {
		name          string
		userID        int
		gameID        *int
		requiredRoles []string
		userMainRole  string
		expected      bool
	}{
		{
			name:          "Admin has all permissions",
			userID:        1,
			gameID:        &[]int{1}[0],
			requiredRoles: []string{models.RoleCommunityModerator},
			userMainRole:  models.RoleAdmin,
			expected:      true,
		},
		{
			name:          "User has specific game role",
			userID:        1,
			gameID:        &[]int{1}[0],
			requiredRoles: []string{models.RoleCommunityModerator},
			userMainRole:  models.RoleUser,
			expected:      true,
		},
		{
			name:          "User lacks required role",
			userID:        1,
			gameID:        &[]int{3}[0],
			requiredRoles: []string{models.RoleCommunityModerator},
			userMainRole:  models.RoleUser,
			expected:      false,
		},
		{
			name:          "User has main role",
			userID:        1,
			gameID:        nil,
			requiredRoles: []string{models.RoleUser},
			userMainRole:  models.RoleUser,
			expected:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := service.HasPermission(test.userID, test.gameID, test.requiredRoles, userRoles, test.userMainRole)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestCanModerateGame(t *testing.T) {
	service := &Service{}

	userRoles := []models.UserRole{
		{UserID: 1, GameID: &[]int{1}[0], Role: models.RoleCommunityModerator},
		{UserID: 1, GameID: &[]int{2}[0], Role: models.RoleWikiMaintainer},
	}

	tests := []struct {
		name         string
		userID       int
		gameID       int
		userMainRole string
		expected     bool
	}{
		{
			name:         "Admin can moderate any game",
			userID:       1,
			gameID:       1,
			userMainRole: models.RoleAdmin,
			expected:     true,
		},
		{
			name:         "Community moderator can moderate assigned game",
			userID:       1,
			gameID:       1,
			userMainRole: models.RoleUser,
			expected:     true,
		},
		{
			name:         "Wiki maintainer can moderate assigned game",
			userID:       1,
			gameID:       2,
			userMainRole: models.RoleUser,
			expected:     true,
		},
		{
			name:         "User cannot moderate unassigned game",
			userID:       1,
			gameID:       3,
			userMainRole: models.RoleUser,
			expected:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := service.CanModerateGame(test.userID, test.gameID, userRoles, test.userMainRole)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestCanManageWiki(t *testing.T) {
	service := &Service{}

	userRoles := []models.UserRole{
		{UserID: 1, GameID: &[]int{1}[0], Role: models.RoleWikiMaintainer},
		{UserID: 1, GameID: &[]int{2}[0], Role: models.RoleCommunityModerator},
	}

	tests := []struct {
		name         string
		userID       int
		gameID       int
		userMainRole string
		expected     bool
	}{
		{
			name:         "Admin can manage any wiki",
			userID:       1,
			gameID:       1,
			userMainRole: models.RoleAdmin,
			expected:     true,
		},
		{
			name:         "Wiki maintainer can manage assigned game wiki",
			userID:       1,
			gameID:       1,
			userMainRole: models.RoleUser,
			expected:     true,
		},
		{
			name:         "Community moderator cannot manage wiki",
			userID:       1,
			gameID:       2,
			userMainRole: models.RoleUser,
			expected:     false,
		},
		{
			name:         "User cannot manage unassigned wiki",
			userID:       1,
			gameID:       3,
			userMainRole: models.RoleUser,
			expected:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := service.CanManageWiki(test.userID, test.gameID, userRoles, test.userMainRole)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestParseUserID(t *testing.T) {
	service := &Service{}

	tests := []struct {
		input       string
		expected    int
		expectError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"-1", -1, false},
		{"abc", 0, true},
		{"", 0, true},
		{"123.45", 0, true},
	}

	for _, test := range tests {
		result, err := service.ParseUserID(test.input)
		if test.expectError {
			if err == nil {
				t.Errorf("ParseUserID(%q): expected error, got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseUserID(%q): expected no error, got %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("ParseUserID(%q): expected %d, got %d", test.input, test.expected, result)
			}
		}
	}
}

func TestIsValidRole(t *testing.T) {
	service := &Service{}

	tests := []struct {
		role     string
		expected bool
	}{
		{models.RoleUser, true},
		{models.RoleAdmin, true},
		{models.RoleCommunityModerator, true},
		{models.RoleWikiMaintainer, true},
		{"invalid_role", false},
		{"", false},
		{"ADMIN", false}, // Case sensitive
	}

	for _, test := range tests {
		result := service.IsValidRole(test.role)
		if result != test.expected {
			t.Errorf("IsValidRole(%q): expected %v, got %v", test.role, test.expected, result)
		}
	}
}

func TestGeneratePasswordResetToken(t *testing.T) {
	service := &Service{}

	token1 := service.GeneratePasswordResetToken()
	token2 := service.GeneratePasswordResetToken()

	if token1 == "" {
		t.Error("Expected non-empty token")
	}

	if len(token1) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("Expected token length 64, got %d", len(token1))
	}

	if token1 == token2 {
		t.Error("Expected different tokens to be generated")
	}
}

func TestGenerateState(t *testing.T) {
	state1 := GenerateState()
	state2 := GenerateState()

	if state1 == "" {
		t.Error("Expected non-empty state")
	}

	if len(state1) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("Expected state length 64, got %d", len(state1))
	}

	if state1 == state2 {
		t.Error("Expected different states to be generated")
	}
}

func TestGetAuthURLs(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:     "test-secret",
			GitHubID:      "github-id",
			GitHubSecret:  "github-secret",
			GoogleID:      "google-id",
			GoogleSecret:  "google-secret",
			DiscordID:     "discord-id",
			DiscordSecret: "discord-secret",
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	service := NewService(cfg)
	state := "test-state"

	githubURL := service.GetGitHubAuthURL(state)
	if githubURL == "" {
		t.Error("Expected non-empty GitHub auth URL")
	}
	if !strings.Contains(githubURL, state) {
		t.Error("Expected GitHub URL to contain state parameter")
	}

	googleURL := service.GetGoogleAuthURL(state)
	if googleURL == "" {
		t.Error("Expected non-empty Google auth URL")
	}
	if !strings.Contains(googleURL, state) {
		t.Error("Expected Google URL to contain state parameter")
	}

	discordURL := service.GetDiscordAuthURL(state)
	if discordURL == "" {
		t.Error("Expected non-empty Discord auth URL")
	}
	if !strings.Contains(discordURL, state) {
		t.Error("Expected Discord URL to contain state parameter")
	}
}

func TestGetAuthURLsWithoutCredentials(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
			// No OAuth credentials
		},
	}

	service := NewService(cfg)
	state := "test-state"

	if service.GetGitHubAuthURL(state) != "" {
		t.Error("Expected empty GitHub URL when credentials missing")
	}

	if service.GetGoogleAuthURL(state) != "" {
		t.Error("Expected empty Google URL when credentials missing")
	}

	if service.GetDiscordAuthURL(state) != "" {
		t.Error("Expected empty Discord URL when credentials missing")
	}
}

func TestExchangeCodeMethods(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:     "test-secret",
			GitHubID:      "github-id",
			GitHubSecret:  "github-secret",
			GoogleID:      "google-id",
			GoogleSecret:  "google-secret",
			DiscordID:     "discord-id",
			DiscordSecret: "discord-secret",
		},
	}

	service := NewService(cfg)

	// These will fail because we don't have real OAuth setup, but we test the error handling
	_, err := service.ExchangeGitHubCode("test-code")
	if err == nil {
		t.Error("Expected error when exchanging invalid GitHub code")
	}

	_, err = service.ExchangeGoogleCode("test-code")
	if err == nil {
		t.Error("Expected error when exchanging invalid Google code")
	}

	_, err = service.ExchangeDiscordCode("test-code")
	if err == nil {
		t.Error("Expected error when exchanging invalid Discord code")
	}
}

func TestExchangeCodeWithoutCredentials(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
			// No OAuth credentials
		},
	}

	service := NewService(cfg)

	_, err := service.ExchangeGitHubCode("test-code")
	if err == nil || err.Error() != "GitHub OAuth not configured" {
		t.Error("Expected GitHub OAuth not configured error")
	}

	_, err = service.ExchangeGoogleCode("test-code")
	if err == nil || err.Error() != "Google OAuth not configured" {
		t.Error("Expected Google OAuth not configured error")
	}

	_, err = service.ExchangeDiscordCode("test-code")
	if err == nil || err.Error() != "Discord OAuth not configured" {
		t.Error("Expected Discord OAuth not configured error")
	}
}

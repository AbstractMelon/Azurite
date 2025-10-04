package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/azurite/backend/internal/auth"
	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/models"
	"github.com/azurite/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type UserService struct {
	db           *database.DB
	auth         *auth.Service
	notification *NotificationService
}

func NewUserService(db *database.DB, authService *auth.Service, notificationService *NotificationService) *UserService {
	return &UserService{
		db:           db,
		auth:         authService,
		notification: notificationService,
	}
}

func (s *UserService) Register(req *models.RegisterRequest) (*models.User, error) {
	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	if !utils.ValidateUsername(req.Username) {
		return nil, errors.New("invalid username format")
	}

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? OR username = ?", req.Email, req.Username).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if count > 0 {
		return nil, errors.New("email or username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := s.db.Exec(`
		INSERT INTO users (username, email, password_hash, display_name, role, is_active, avatar, bio)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Username, req.Email, hashedPassword, req.DisplayName, models.RoleUser, true, "/static/placeholders/avatar.jpg", "User has not provided a bio")

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	return s.GetByID(int(userID))
}

func (s *UserService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user has a password (OAuth users don't have passwords)
	if user.PasswordHash == "" {
		return nil, errors.New("this account uses OAuth authentication. Please sign in with your OAuth provider")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	s.db.Exec("UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?", user.ID)

	token, err := s.auth.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, username, email, password_hash, display_name, avatar, bio, role,
		       is_active, email_verified, notify_email, notify_in_site,
		       github_id, discord_id, google_id, created_at, updated_at, last_login_at
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Avatar, &user.Bio, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.NotifyEmail, &user.NotifyInSite, &user.GitHubID, &user.DiscordID,
		&user.GoogleID, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, username, email, password_hash, display_name, avatar, bio, role,
		       is_active, email_verified, notify_email, notify_in_site,
		       github_id, discord_id, google_id, created_at, updated_at, last_login_at
		FROM users WHERE email = ?
	`, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Avatar, &user.Bio, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.NotifyEmail, &user.NotifyInSite, &user.GitHubID, &user.DiscordID,
		&user.GoogleID, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow(`
		SELECT id, username, email, password_hash, display_name, avatar, bio, role,
		       is_active, email_verified, notify_email, notify_in_site,
		       github_id, discord_id, google_id, created_at, updated_at, last_login_at
		FROM users WHERE username = ?
	`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Avatar, &user.Bio, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.NotifyEmail, &user.NotifyInSite, &user.GitHubID, &user.DiscordID,
		&user.GoogleID, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(userID int, req *models.UserUpdateRequest) (*models.User, error) {
	_, err := s.db.Exec(`
		UPDATE users SET display_name = ?, bio = ?, notify_email = ?, notify_in_site = ?
		WHERE id = ?
	`, req.DisplayName, req.Bio, req.NotifyEmail, req.NotifyInSite, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.GetByID(userID)
}

func (s *UserService) UpdatePassword(userID int, currentPassword, newPassword string) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if user has a password (OAuth users don't have passwords)
	if user.PasswordHash == "" {
		return errors.New("this account uses OAuth authentication and does not have a password")
	}

	if !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = s.db.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *UserService) HandleOAuthCallback(provider, code, state string) (*models.AuthResponse, error) {
	var token *oauth2.Token
	var err error

	switch provider {
	case "github":
		token, err = s.auth.ExchangeGitHubCode(code)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange GitHub code: %w", err)
		}
		return s.handleGitHubUser(token)

	case "google":
		token, err = s.auth.ExchangeGoogleCode(code)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange Google code: %w", err)
		}
		return s.handleGoogleUser(token)

	case "discord":
		token, err = s.auth.ExchangeDiscordCode(code)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange Discord code: %w", err)
		}
		return s.handleDiscordUser(token)

	default:
		return nil, errors.New("unsupported OAuth provider")
	}
}

func (s *UserService) handleGitHubUser(token *oauth2.Token) (*models.AuthResponse, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Bio       string `json:"bio"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub user: %w", err)
	}

	if githubUser.Email == "" {
		req, _ = http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		req.Header.Set("Authorization", "token "+token.AccessToken)

		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get GitHub emails: %w", err)
		}
		defer resp.Body.Close()

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&emails); err == nil {
			for _, email := range emails {
				if email.Primary {
					githubUser.Email = email.Email
					break
				}
			}
		}
	}

	return s.findOrCreateOAuthUser("github", fmt.Sprintf("%d", githubUser.ID), githubUser.Email, githubUser.Login, githubUser.Name, githubUser.AvatarURL, githubUser.Bio)
}

func (s *UserService) handleGoogleUser(token *oauth2.Token) (*models.AuthResponse, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user: %w", err)
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode Google user: %w", err)
	}

	username := utils.GenerateSlug(googleUser.Name)
	return s.findOrCreateOAuthUser("google", googleUser.ID, googleUser.Email, username, googleUser.Name, googleUser.Picture, "")
}

func (s *UserService) handleDiscordUser(token *oauth2.Token) (*models.AuthResponse, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Discord user: %w", err)
	}
	defer resp.Body.Close()

	var discordUser struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Email         string `json:"email"`
		Avatar        string `json:"avatar"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		return nil, fmt.Errorf("failed to decode Discord user: %w", err)
	}

	username := discordUser.Username + discordUser.Discriminator
	avatarURL := ""
	if discordUser.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordUser.ID, discordUser.Avatar)
	}

	return s.findOrCreateOAuthUser("discord", discordUser.ID, discordUser.Email, username, discordUser.Username, avatarURL, "")
}

func (s *UserService) findOrCreateOAuthUser(provider, providerID, email, username, displayName, avatar, bio string) (*models.AuthResponse, error) {
	var user *models.User
	var err error

	switch provider {
	case "github":
		user, err = s.getByProviderID("github_id", providerID)
	case "google":
		user, err = s.getByProviderID("google_id", providerID)
	case "discord":
		user, err = s.getByProviderID("discord_id", providerID)
	}

	if err == nil && user != nil {
		s.db.Exec("UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?", user.ID)
		token, err := s.auth.GenerateToken(user)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token: %w", err)
		}
		return &models.AuthResponse{Token: token, User: *user}, nil
	}

	if email != "" {
		user, err = s.GetByEmail(email)
		if err == nil {
			return s.linkOAuthAccount(user.ID, provider, providerID)
		}
	}

	return s.createOAuthUser(provider, providerID, email, username, displayName, avatar, bio)
}

func (s *UserService) getByProviderID(column, providerID string) (*models.User, error) {
	user := &models.User{}
	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, display_name, avatar, bio, role,
		       is_active, email_verified, notify_email, notify_in_site,
		       github_id, discord_id, google_id, created_at, updated_at, last_login_at
		FROM users WHERE %s = ?
	`, column)

	err := s.db.QueryRow(query, providerID).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Avatar, &user.Bio, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.NotifyEmail, &user.NotifyInSite, &user.GitHubID, &user.DiscordID,
		&user.GoogleID, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) linkOAuthAccount(userID int, provider, providerID string) (*models.AuthResponse, error) {
	var column string
	switch provider {
	case "github":
		column = "github_id"
	case "google":
		column = "google_id"
	case "discord":
		column = "discord_id"
	}

	query := fmt.Sprintf("UPDATE users SET %s = ?, last_login_at = CURRENT_TIMESTAMP WHERE id = ?", column)
	_, err := s.db.Exec(query, providerID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to link OAuth account: %w", err)
	}

	user, err := s.GetByID(userID)
	if err != nil {
		return nil, err
	}

	token, err := s.auth.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *UserService) createOAuthUser(provider, providerID, email, username, displayName, avatar, bio string) (*models.AuthResponse, error) {
	baseUsername := utils.GenerateSlug(username)
	finalUsername := baseUsername
	counter := 1

	for {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", finalUsername).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to check username availability: %w", err)
		}

		if count == 0 {
			break
		}

		finalUsername = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
	}

	var githubID, googleID, discordID sql.NullString

	switch provider {
	case "github":
		githubID = sql.NullString{String: providerID, Valid: true}
	case "google":
		googleID = sql.NullString{String: providerID, Valid: true}
	case "discord":
		discordID = sql.NullString{String: providerID, Valid: true}
	}

	// Set default values for OAuth users
	if avatar == "" {
		avatar = "/static/placeholders/avatar.jpg"
	}
	if bio == "" {
		bio = "User has not provided a bio"
	}

	// Insert with NULL password_hash for OAuth users
	result, err := s.db.Exec(`
		INSERT INTO users (username, email, password_hash, display_name, avatar, bio, role, is_active, email_verified,
		                  github_id, google_id, discord_id, last_login_at)
		VALUES (?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, finalUsername, email, displayName, avatar, bio, models.RoleUser, true, true,
		githubID, googleID, discordID)

	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	user, err := s.GetByID(int(userID))
	if err != nil {
		return nil, err
	}

	token, err := s.auth.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *UserService) RequestPasswordReset(email string) error {
	user, err := s.GetByEmail(email)
	if err != nil {
		// Always return nil to prevent email enumeration
		return nil
	}

	token := s.auth.GeneratePasswordResetToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	_, err = s.db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES (?, ?, ?)
	`, user.ID, token, expiresAt)

	if err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send the password reset email
	if s.notification != nil {
		go s.notification.SendPasswordResetEmail(email, token)
	}

	return nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	var userID int
	var expiresAt time.Time
	var used bool

	err := s.db.QueryRow(`
		SELECT user_id, expires_at, used FROM password_reset_tokens
		WHERE token = ?
	`, token).Scan(&userID, &expiresAt, &used)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("invalid reset token")
		}
		return fmt.Errorf("failed to get reset token: %w", err)
	}

	if used {
		return errors.New("reset token already used")
	}

	if time.Now().After(expiresAt) {
		return errors.New("reset token expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	_, err = tx.Exec("UPDATE password_reset_tokens SET used = 1 WHERE token = ?", token)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return tx.Commit()
}

func (s *UserService) GetUserRoles(userID int) ([]models.UserRole, error) {
	rows, err := s.db.Query(`
		SELECT id, user_id, game_id, role FROM user_roles WHERE user_id = ?
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []models.UserRole
	for rows.Next() {
		var role models.UserRole
		err := rows.Scan(&role.ID, &role.UserID, &role.GameID, &role.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *UserService) List(page, perPage int) ([]models.User, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_active = 1").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT id, username, email, display_name, avatar, bio, role,
		       is_active, email_verified, created_at, updated_at, last_login_at
		FROM users WHERE is_active = 1
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.DisplayName,
			&user.Avatar, &user.Bio, &user.Role, &user.IsActive,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

// Create creates a new user (for testing purposes)
func (s *UserService) Create(user *models.User) (*models.User, error) {
	if user.PasswordHash == "" {
		// Generate a default password hash for testing
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword123"), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	query := `
		INSERT INTO users (username, email, password_hash, display_name, avatar, bio, role,
		                  is_active, email_verified, notify_email, notify_in_site, github_id,
		                  discord_id, google_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`

	result, err := s.db.Exec(query,
		user.Username, user.Email, user.PasswordHash, user.DisplayName,
		user.Avatar, user.Bio, user.Role, user.IsActive, user.EmailVerified,
		user.NotifyEmail, user.NotifyInSite, user.GitHubID, user.DiscordID,
		user.GoogleID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = int(id)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Clear password hash from returned user
	user.PasswordHash = ""

	return user, nil
}

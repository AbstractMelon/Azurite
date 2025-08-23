package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/azurite/backend/internal/config"
	"github.com/azurite/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type Service struct {
	config        *config.Config
	githubConfig  *oauth2.Config
	googleConfig  *oauth2.Config
	discordConfig *oauth2.Config
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewService(cfg *config.Config) *Service {
	service := &Service{
		config: cfg,
	}

	if cfg.Auth.GitHubID != "" && cfg.Auth.GitHubSecret != "" {
		service.githubConfig = &oauth2.Config{
			ClientID:     cfg.Auth.GitHubID,
			ClientSecret: cfg.Auth.GitHubSecret,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
			RedirectURL:  fmt.Sprintf("http://%s:%s/api/auth/callback/github", cfg.Server.Host, cfg.Server.Port),
		}
	}

	if cfg.Auth.GoogleID != "" && cfg.Auth.GoogleSecret != "" {
		service.googleConfig = &oauth2.Config{
			ClientID:     cfg.Auth.GoogleID,
			ClientSecret: cfg.Auth.GoogleSecret,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
			RedirectURL:  fmt.Sprintf("http://%s:%s/api/auth/callback/google", cfg.Server.Host, cfg.Server.Port),
		}
	}

	if cfg.Auth.DiscordID != "" && cfg.Auth.DiscordSecret != "" {
		service.discordConfig = &oauth2.Config{
			ClientID:     cfg.Auth.DiscordID,
			ClientSecret: cfg.Auth.DiscordSecret,
			Scopes:       []string{"identify", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
			RedirectURL: fmt.Sprintf("http://%s:%s/api/auth/callback/discord", cfg.Server.Host, cfg.Server.Port),
		}
	}

	return service
}

func (s *Service) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "azurite",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Auth.JWTSecret))
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *Service) ExtractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func (s *Service) GetCurrentUser(c *gin.Context) (*models.User, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	return user, nil
}

func (s *Service) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := s.ExtractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := s.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		user := &models.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
		}

		c.Set("user", user)
		c.Next()
	}
}

func (s *Service) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.GetCurrentUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

func (s *Service) GetGitHubAuthURL(state string) string {
	if s.githubConfig == nil {
		return ""
	}
	return s.githubConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Service) GetGoogleAuthURL(state string) string {
	if s.googleConfig == nil {
		return ""
	}
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Service) GetDiscordAuthURL(state string) string {
	if s.discordConfig == nil {
		return ""
	}
	return s.discordConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Service) ExchangeGitHubCode(code string) (*oauth2.Token, error) {
	if s.githubConfig == nil {
		return nil, errors.New("GitHub OAuth not configured")
	}
	return s.githubConfig.Exchange(oauth2.NoContext, code)
}

func (s *Service) ExchangeGoogleCode(code string) (*oauth2.Token, error) {
	if s.googleConfig == nil {
		return nil, errors.New("Google OAuth not configured")
	}
	return s.googleConfig.Exchange(oauth2.NoContext, code)
}

func (s *Service) ExchangeDiscordCode(code string) (*oauth2.Token, error) {
	if s.discordConfig == nil {
		return nil, errors.New("Discord OAuth not configured")
	}
	return s.discordConfig.Exchange(oauth2.NoContext, code)
}

func GenerateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *Service) HasPermission(userID int, gameID *int, requiredRoles []string, userRoles []models.UserRole, userMainRole string) bool {
	for _, required := range requiredRoles {
		if userMainRole == required {
			return true
		}

		if gameID != nil {
			for _, role := range userRoles {
				if role.UserID == userID && role.GameID != nil && *role.GameID == *gameID && role.Role == required {
					return true
				}
			}
		}
	}

	return false
}

func (s *Service) CanModerateGame(userID int, gameID int, userRoles []models.UserRole, userMainRole string) bool {
	if userMainRole == models.RoleAdmin {
		return true
	}

	for _, role := range userRoles {
		if role.UserID == userID && role.GameID != nil && *role.GameID == gameID &&
			(role.Role == models.RoleCommunityModerator || role.Role == models.RoleWikiMaintainer) {
			return true
		}
	}

	return false
}

func (s *Service) CanManageWiki(userID int, gameID int, userRoles []models.UserRole, userMainRole string) bool {
	if userMainRole == models.RoleAdmin {
		return true
	}

	for _, role := range userRoles {
		if role.UserID == userID && role.GameID != nil && *role.GameID == gameID &&
			role.Role == models.RoleWikiMaintainer {
			return true
		}
	}

	return false
}

func (s *Service) ParseUserID(userIDStr string) (int, error) {
	return strconv.Atoi(userIDStr)
}

func (s *Service) IsValidRole(role string) bool {
	validRoles := []string{
		models.RoleUser,
		models.RoleAdmin,
		models.RoleCommunityModerator,
		models.RoleWikiMaintainer,
	}

	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}

	return false
}

func (s *Service) GeneratePasswordResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

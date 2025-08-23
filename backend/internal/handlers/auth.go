package handlers

import (
	"net/http"
	"strconv"

	"github.com/azurite/backend/internal/auth"
	"github.com/azurite/backend/internal/models"
	"github.com/azurite/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService         *services.UserService
	notificationService *services.NotificationService
	authService         *auth.Service
}

func NewAuthHandler(userService *services.UserService, notificationService *services.NotificationService, authService *auth.Service) *AuthHandler {
	return &AuthHandler{
		userService:         userService,
		notificationService: notificationService,
		authService:         authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	go h.notificationService.SendWelcomeEmail(user)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data: models.AuthResponse{
			Token: token,
			User:  *user,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	authResponse, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    authResponse,
	})
}

func (h *AuthHandler) GitHubAuth(c *gin.Context) {
	state := auth.GenerateState()
	url := h.authService.GetGitHubAuthURL(state)

	if url == "" {
		c.JSON(http.StatusNotImplemented, models.APIResponse{
			Success: false,
			Error:   "GitHub OAuth not configured",
		})
		return
	}

	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"url": url},
	})
}

func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	state := auth.GenerateState()
	url := h.authService.GetGoogleAuthURL(state)

	if url == "" {
		c.JSON(http.StatusNotImplemented, models.APIResponse{
			Success: false,
			Error:   "Google OAuth not configured",
		})
		return
	}

	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"url": url},
	})
}

func (h *AuthHandler) DiscordAuth(c *gin.Context) {
	state := auth.GenerateState()
	url := h.authService.GetDiscordAuthURL(state)

	if url == "" {
		c.JSON(http.StatusNotImplemented, models.APIResponse{
			Success: false,
			Error:   "Discord OAuth not configured",
		})
		return
	}

	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"url": url},
	})
}

func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	h.handleOAuthCallback(c, "github")
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	h.handleOAuthCallback(c, "google")
}

func (h *AuthHandler) DiscordCallback(c *gin.Context) {
	h.handleOAuthCallback(c, "discord")
}

func (h *AuthHandler) handleOAuthCallback(c *gin.Context, provider string) {
	code := c.Query("code")
	state := c.Query("state")

	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid state parameter",
		})
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	if code == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Authorization code not provided",
		})
		return
	}

	authResponse, err := h.userService.HandleOAuthCallback(provider, code, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    authResponse,
	})
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req models.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err := h.userService.RequestPasswordReset(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "If the email exists, a password reset link has been sent",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "If the email exists, a password reset link has been sent",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.PasswordResetConfirm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err := h.userService.ResetPassword(req.Token, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password reset successfully",
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	fullUser, err := h.userService.GetByID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user profile",
		})
		return
	}

	fullUser.PasswordHash = ""

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    fullUser,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	updatedUser, err := h.userService.Update(user.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	updatedUser.PasswordHash = ""

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    updatedUser,
	})
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err = h.userService.UpdatePassword(user.ID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password updated successfully",
	})
}

func (h *AuthHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	user.PasswordHash = ""
	user.Email = ""

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}

func (h *AuthHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userService.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	user.PasswordHash = ""
	user.Email = ""

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}

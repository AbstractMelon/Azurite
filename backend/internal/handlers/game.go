package handlers

import (
	"net/http"
	"strconv"

	"github.com/azurite/backend/internal/models"
	"github.com/azurite/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *services.GameService
	modService  *services.ModService
}

func NewGameHandler(gameService *services.GameService, modService *services.ModService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		modService:  modService,
	}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=200"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	game, err := h.gameService.Create(req.Name, req.Description, req.Icon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    game,
	})
}

func (h *GameHandler) GetGame(c *gin.Context) {
	slug := c.Param("slug")

	game, err := h.gameService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Game not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    game,
	})
}

func (h *GameHandler) ListGames(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	games, total, err := h.gameService.List(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.PaginatedResponse{
			Data:       games,
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (h *GameHandler) UpdateGame(c *gin.Context) {
	gameIDStr := c.Param("id")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid game ID",
		})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required,min=1,max=200"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	game, err := h.gameService.Update(gameID, req.Name, req.Description, req.Icon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    game,
	})
}

func (h *GameHandler) DeleteGame(c *gin.Context) {
	gameIDStr := c.Param("id")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid game ID",
		})
		return
	}

	err = h.gameService.Delete(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Game deleted successfully",
	})
}

func (h *GameHandler) GetGameMods(c *gin.Context) {
	slug := c.Param("slug")

	game, err := h.gameService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Game not found",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	sortBy := c.DefaultQuery("sort", "created")
	order := c.DefaultQuery("order", "desc")
	search := c.Query("search")
	tags := c.QueryArray("tags")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	mods, total, err := h.modService.ListByGame(game.ID, page, perPage, sortBy, order, tags, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.PaginatedResponse{
			Data:       mods,
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (h *GameHandler) CreateGameRequest(c *gin.Context) {
	var req models.GameRequestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Authentication required",
		})
		return
	}

	userModel := user.(*models.User)
	gameRequest, err := h.gameService.CreateRequest(&req, userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    gameRequest,
	})
}

func (h *GameHandler) ListGameRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	requests, total, err := h.gameService.ListRequests(page, perPage, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.PaginatedResponse{
			Data:       requests,
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (h *GameHandler) ApproveGameRequest(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request ID",
		})
		return
	}

	var req struct {
		AdminNotes string `json:"admin_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	game, err := h.gameService.ApproveRequest(requestID, req.AdminNotes)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    game,
		Message: "Game request approved successfully",
	})
}

func (h *GameHandler) DenyGameRequest(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request ID",
		})
		return
	}

	var req struct {
		AdminNotes string `json:"admin_notes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err = h.gameService.DenyRequest(requestID, req.AdminNotes)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Game request denied successfully",
	})
}

func (h *GameHandler) GetGameTags(c *gin.Context) {
	slug := c.Param("slug")

	game, err := h.gameService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Game not found",
		})
		return
	}

	tags, err := h.gameService.GetTags(game.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    tags,
	})
}

func (h *GameHandler) CreateGameTag(c *gin.Context) {
	slug := c.Param("slug")

	game, err := h.gameService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Game not found",
		})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	tag, err := h.gameService.CreateTag(game.ID, req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    tag,
	})
}

func (h *GameHandler) DeleteGameTag(c *gin.Context) {
	tagIDStr := c.Param("tagId")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid tag ID",
		})
		return
	}

	err = h.gameService.DeleteTag(tagID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tag deleted successfully",
	})
}

func (h *GameHandler) AssignModerator(c *gin.Context) {
	slug := c.Param("slug")

	game, err := h.gameService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Game not found",
		})
		return
	}

	var req struct {
		UserID int `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	err = h.gameService.AssignModerator(game.ID, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Moderator assigned successfully",
	})
}

func (h *GameHandler) RemoveModerator(c *gin.Context) {
	slug := c.Param("slug")
	userIDStr := c.Param("userId")

	game, err := h.gameService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Game not found",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	err = h.gameService.RemoveModerator(game.ID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Moderator removed successfully",
	})
}

func (h *GameHandler) GetGameModerators(c *gin.Context) {
	slug := c.Param("slug")

	game, err := h.gameService.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Game not found",
		})
		return
	}

	moderators, err := h.gameService.GetModerators(game.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    moderators,
	})
}

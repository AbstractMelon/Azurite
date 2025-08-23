package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/azurite/backend/internal/auth"
	"github.com/azurite/backend/internal/config"
	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/handlers"
	"github.com/azurite/backend/internal/middleware"
	"github.com/azurite/backend/internal/models"
	"github.com/azurite/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system env vars")
	}

	// Log loaded .env file status
	fmt.Printf("Loaded .env file: %v\n", err == nil)

	// Load configuration
	cfg := config.Load()

	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	if err := os.MkdirAll(cfg.Storage.ModsPath, 0755); err != nil {
		log.Fatal("Failed to create mods directory:", err)
	}

	if err := os.MkdirAll(cfg.Storage.ImagesPath, 0755); err != nil {
		log.Fatal("Failed to create images directory:", err)
	}

	authService := auth.NewService(cfg)
	notificationService := services.NewNotificationService(db, cfg)
	userService := services.NewUserService(db, authService, notificationService)
	gameService := services.NewGameService(db)
	modService := services.NewModService(db, cfg.Storage.ModsPath, cfg.Storage.ImagesPath)
	commentService := services.NewCommentService(db)
	documentationService := services.NewDocumentationService(db)
	adminService := services.NewAdminService(db)
	imageService := services.NewImageService(cfg.Storage.ImagesPath)

	authHandler := handlers.NewAuthHandler(userService, notificationService, authService)
	gameHandler := handlers.NewGameHandler(gameService, modService)

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CheckBan(db))
	router.Use(gin.Recovery())

	setupRoutes(router, authHandler, gameHandler, authService, userService, gameService, modService, commentService, notificationService, documentationService, adminService, imageService, cfg)

	// Capture shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down database...")

		// WAL checkpoint
		_, err := db.Exec("PRAGMA wal_checkpoint(FULL);")
		if err != nil {
			log.Println("Failed to checkpoint WAL:", err)
		}

		if err := db.Close(); err != nil {
			log.Println("Failed to close database:", err)
		}

		log.Println("Database closed cleanly")
		os.Exit(0)
	}()

	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	gameHandler *handlers.GameHandler,
	authService *auth.Service,
	userService *services.UserService,
	gameService *services.GameService,
	modService *services.ModService,
	commentService *services.CommentService,
	notificationService *services.NotificationService,
	documentationService *services.DocumentationService,
	adminService *services.AdminService,
	imageService *services.ImageService,
	cfg *config.Config,
) {
	api := router.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.RequestPasswordReset)
		auth.POST("/reset-password", authHandler.ResetPassword)

		// OAuth routes
		auth.GET("/github", authHandler.GitHubAuth)
		auth.GET("/google", authHandler.GoogleAuth)
		auth.GET("/discord", authHandler.DiscordAuth)
		auth.GET("/callback/github", authHandler.GitHubCallback)
		auth.GET("/callback/google", authHandler.GoogleCallback)
		auth.GET("/callback/discord", authHandler.DiscordCallback)

		// Protected auth routes
		authProtected := auth.Use(authService.RequireAuth())
		{
			authProtected.GET("/profile", authHandler.GetProfile)
			authProtected.PUT("/profile", authHandler.UpdateProfile)
			authProtected.PUT("/password", authHandler.UpdatePassword)

			// User's mods
			authProtected.GET("/mods", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)

				page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
				perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

				mods, total, err := modService.ListByOwner(user.ID, page, perPage)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Add like status for each mod
				for i := range mods {
					mods[i].IsLiked = modService.IsLiked(mods[i].ID, user.ID)
				}

				totalPages := int((total + int64(perPage) - 1) / int64(perPage))

				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data": gin.H{
						"data":        mods,
						"page":        page,
						"per_page":    perPage,
						"total":       total,
						"total_pages": totalPages,
					},
				})
			})
		}
	}

	// Search routes
	search := api.Group("/search")
	{
		search.GET("/mods", func(c *gin.Context) {
			query := c.Query("q")
			if query == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
				return
			}

			gameID, _ := strconv.Atoi(c.DefaultQuery("game_id", "0"))
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

			mods, total, err := modService.Search(query, gameID, page, perPage)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Add like status if user is authenticated
			if user, exists := c.Get("user"); exists {
				userModel := user.(*models.User)
				for i := range mods {
					mods[i].IsLiked = modService.IsLiked(mods[i].ID, userModel.ID)
				}
			}

			totalPages := int((total + int64(perPage) - 1) / int64(perPage))

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"data":        mods,
					"page":        page,
					"per_page":    perPage,
					"total":       total,
					"total_pages": totalPages,
					"query":       query,
				},
			})
		})
	}

	// User routes
	users := api.Group("/users")
	{
		users.GET("/:id", authHandler.GetUserByID)
		users.GET("/username/:username", authHandler.GetUserByUsername)
	}

	// Game routes
	games := api.Group("/games")
	{
		games.GET("", gameHandler.ListGames)
		games.GET("/:slug", gameHandler.GetGame)
		games.GET("/:slug/mods", gameHandler.GetGameMods)
		games.GET("/:slug/tags", gameHandler.GetGameTags)

		// Game requests
		gameRequests := games.Group("/requests")
		{
			gameRequests.POST("", authService.RequireAuth(), gameHandler.CreateGameRequest)
			gameRequests.GET("", authService.RequireAuth(), authService.RequireRole("admin"), gameHandler.ListGameRequests)
			gameRequests.POST("/:id/approve", authService.RequireAuth(), authService.RequireRole("admin"), gameHandler.ApproveGameRequest)
			gameRequests.POST("/:id/deny", authService.RequireAuth(), authService.RequireRole("admin"), gameHandler.DenyGameRequest)
		}

		games.GET("/:slug/moderators", gameHandler.GetGameModerators)

		// Protected game routes
		gamesProtected := games.Use(authService.RequireAuth(), authService.RequireRole("admin"))
		{
			gamesProtected.POST("", gameHandler.CreateGame)
			gamesProtected.POST("/:slug/tags", gameHandler.CreateGameTag)
			gamesProtected.DELETE("/:slug/tags/:tagId", gameHandler.DeleteGameTag)
			gamesProtected.POST("/:slug/moderators", gameHandler.AssignModerator)
			gamesProtected.DELETE("/:slug/moderators/:userId", gameHandler.RemoveModerator)
		}

		// Game management by ID (separate group to avoid conflicts)
		gamesByID := api.Group("/games/manage")
		gamesByID.Use(authService.RequireAuth(), authService.RequireRole("admin"))
		{
			gamesByID.PUT("/:id", gameHandler.UpdateGame)
			gamesByID.DELETE("/:id", gameHandler.DeleteGame)
		}
	}

	// Mod routes
	mods := api.Group("/mods")
	{
		mods.GET("/:gameSlug/:modSlug", func(c *gin.Context) {
			gameSlug := c.Param("gameSlug")
			modSlug := c.Param("modSlug")

			mod, err := modService.GetBySlug(gameSlug, modSlug)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Mod not found"})
				return
			}

			// Check if user has liked this mod
			if user, exists := c.Get("user"); exists {
				userModel := user.(*models.User)
				mod.IsLiked = modService.IsLiked(mod.ID, userModel.ID)
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": mod})
		})

		// Get mod by ID (for editing)
		mods.GET("/id/:id", func(c *gin.Context) {
			modID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mod ID"})
				return
			}

			mod, err := modService.GetByID(modID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Mod not found"})
				return
			}

			// Add like status if user is authenticated
			if user, exists := c.Get("user"); exists {
				userModel := user.(*models.User)
				mod.IsLiked = modService.IsLiked(mod.ID, userModel.ID)
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": mod})
		})

		// Protected mod routes
		modsProtected := mods.Use(authService.RequireAuth())
		{
			modsProtected.POST("", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)

				var req models.ModCreateRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				mod, err := modService.Create(&req, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{"success": true, "data": mod})
			})

			modsProtected.PUT("/:id", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				modID, _ := strconv.Atoi(c.Param("id"))

				var req models.ModUpdateRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				mod, err := modService.Update(modID, &req, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "data": mod})
			})

			modsProtected.DELETE("/:id", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				modID, _ := strconv.Atoi(c.Param("id"))

				err := modService.Delete(modID, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mod deleted"})
			})

			modsProtected.POST("/:id/like", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				modID, _ := strconv.Atoi(c.Param("id"))

				err := modService.Like(modID, user.ID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mod liked"})
			})

			modsProtected.DELETE("/:id/like", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				modID, _ := strconv.Atoi(c.Param("id"))

				err := modService.Unlike(modID, user.ID)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mod unliked"})
			})

			modsProtected.POST("/:id/files", middleware.FileUploadLimit(cfg.Storage.MaxFileSize), func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				modID, _ := strconv.Atoi(c.Param("id"))

				// Verify mod ownership
				mod, err := modService.GetByID(modID)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Mod not found"})
					return
				}
				if mod.OwnerID != user.ID {
					c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
					return
				}

				file, header, err := c.Request.FormFile("file")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
					return
				}
				defer file.Close()

				isMain := c.PostForm("is_main") == "true"

				modFile, err := modService.UploadFile(modID, file, header, isMain)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{"success": true, "data": modFile})
			})

			modsProtected.POST("/:id/reject", authService.RequireRole("admin", "community_moderator"), func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				modID, _ := strconv.Atoi(c.Param("id"))

				var req struct {
					Reason string `json:"reason" binding:"required"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				err := modService.Reject(modID, req.Reason, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mod rejected"})
			})

			modsProtected.POST("/:id/approve", authService.RequireRole("admin", "community_moderator"), func(c *gin.Context) {
				modID, _ := strconv.Atoi(c.Param("id"))

				err := modService.Approve(modID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mod approved"})
			})
		}
	}

	// Comment routes
	comments := api.Group("/comments")
	{
		comments.GET("/mod/:modId", func(c *gin.Context) {
			modID, _ := strconv.Atoi(c.Param("modId"))
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

			comments, total, err := commentService.ListByMod(modID, page, perPage)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			totalPages := int((total + int64(perPage) - 1) / int64(perPage))
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": models.PaginatedResponse{
					Data:       comments,
					Page:       page,
					PerPage:    perPage,
					Total:      total,
					TotalPages: totalPages,
				},
			})
		})

		commentsProtected := comments.Use(authService.RequireAuth())
		{
			commentsProtected.POST("", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				modID, _ := strconv.Atoi(c.Query("mod_id"))

				var req models.CommentCreateRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				comment, err := commentService.Create(&req, modID, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{"success": true, "data": comment})
			})

			commentsProtected.PUT("/:id", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				commentID, _ := strconv.Atoi(c.Param("id"))

				var req struct {
					Content string `json:"content" binding:"required"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				comment, err := commentService.Update(commentID, req.Content, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "data": comment})
			})

			commentsProtected.DELETE("/:id", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				commentID, _ := strconv.Atoi(c.Param("id"))

				isAdmin := user.Role == "admin"
				err := commentService.Delete(commentID, user.ID, isAdmin)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment deleted"})
			})
		}
	}

	// Notification routes
	notifications := api.Group("/notifications")
	notifications.Use(authService.RequireAuth())
	{
		notifications.GET("", func(c *gin.Context) {
			user, _ := authService.GetCurrentUser(c)
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

			notifications, total, err := notificationService.GetByUser(user.ID, page, perPage)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			totalPages := int((total + int64(perPage) - 1) / int64(perPage))
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": models.PaginatedResponse{
					Data:       notifications,
					Page:       page,
					PerPage:    perPage,
					Total:      total,
					TotalPages: totalPages,
				},
			})
		})

		notifications.GET("/unread-count", func(c *gin.Context) {
			user, _ := authService.GetCurrentUser(c)

			count, err := notificationService.GetUnreadCount(user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"count": count}})
		})

		notifications.PUT("/:id/read", func(c *gin.Context) {
			user, _ := authService.GetCurrentUser(c)
			notificationID, _ := strconv.Atoi(c.Param("id"))

			err := notificationService.MarkAsRead(notificationID, user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notification marked as read"})
		})

		notifications.PUT("/read-all", func(c *gin.Context) {
			user, _ := authService.GetCurrentUser(c)

			err := notificationService.MarkAllAsRead(user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "message": "All notifications marked as read"})
		})

		notifications.DELETE("/:id", func(c *gin.Context) {
			user, _ := authService.GetCurrentUser(c)
			notificationID, _ := strconv.Atoi(c.Param("id"))

			err := notificationService.Delete(notificationID, user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notification deleted"})
		})
	}

	// Documentation routes
	docs := api.Group("/docs")
	{
		docs.GET("/:gameSlug", func(c *gin.Context) {
			gameSlug := c.Param("gameSlug")
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

			game, err := gameService.GetBySlug(gameSlug)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
				return
			}

			docs, total, err := documentationService.ListByGame(game.ID, page, perPage)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			totalPages := int((total + int64(perPage) - 1) / int64(perPage))
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": models.PaginatedResponse{
					Data:       docs,
					Page:       page,
					PerPage:    perPage,
					Total:      total,
					TotalPages: totalPages,
				},
			})
		})

		docs.GET("/:gameSlug/:docSlug", func(c *gin.Context) {
			gameSlug := c.Param("gameSlug")
			docSlug := c.Param("docSlug")

			doc, err := documentationService.GetBySlug(gameSlug, docSlug)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Documentation not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": doc})
		})

		docsProtected := docs.Use(authService.RequireAuth())
		{
			docsProtected.POST("/:gameSlug", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				gameSlug := c.Param("gameSlug")

				game, err := gameService.GetBySlug(gameSlug)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
					return
				}

				userRoles, _ := userService.GetUserRoles(user.ID)
				canEdit := documentationService.CanUserEdit(user.ID, game.ID, userRoles, user.Role)
				if !canEdit {
					c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
					return
				}

				var req models.DocumentationCreateRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				doc, err := documentationService.Create(&req, game.ID, user.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{"success": true, "data": doc})
			})

			docsProtected.PUT("/:gameSlug/:docSlug", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				gameSlug := c.Param("gameSlug")
				docSlug := c.Param("docSlug")

				doc, err := documentationService.GetBySlug(gameSlug, docSlug)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Documentation not found"})
					return
				}

				userRoles, _ := userService.GetUserRoles(user.ID)
				canEdit := documentationService.CanUserEdit(user.ID, doc.GameID, userRoles, user.Role)

				var req models.DocumentationCreateRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
					return
				}

				updatedDoc, err := documentationService.Update(doc.ID, &req, user.ID, canEdit)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedDoc})
			})

			docsProtected.DELETE("/:gameSlug/:docSlug", func(c *gin.Context) {
				user, _ := authService.GetCurrentUser(c)
				gameSlug := c.Param("gameSlug")
				docSlug := c.Param("docSlug")

				doc, err := documentationService.GetBySlug(gameSlug, docSlug)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Documentation not found"})
					return
				}

				userRoles, _ := userService.GetUserRoles(user.ID)
				canDelete := documentationService.CanUserEdit(user.ID, doc.GameID, userRoles, user.Role)

				err = documentationService.Delete(doc.ID, user.ID, canDelete)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Documentation deleted"})
			})
		}
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(authService.RequireAuth(), authService.RequireRole("admin"))
	{
		admin.GET("/stats", func(c *gin.Context) {
			userStats, _ := adminService.GetUserStats()
			modStats, _ := adminService.GetModStats()
			systemStats, _ := adminService.GetSystemStats()

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"users":  userStats,
					"mods":   modStats,
					"system": systemStats,
				},
			})
		})

		admin.GET("/activity", func(c *gin.Context) {
			limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
			activity, err := adminService.GetRecentActivity(limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "data": activity})
		})

		admin.GET("/mods/pending", func(c *gin.Context) {
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

			mods, total, err := adminService.ListPendingMods(page, perPage)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			totalPages := int((total + int64(perPage) - 1) / int64(perPage))

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"data":        mods,
					"page":        page,
					"per_page":    perPage,
					"total":       total,
					"total_pages": totalPages,
				},
			})
		})

		admin.POST("/bans", func(c *gin.Context) {
			user, _ := authService.GetCurrentUser(c)

			var req models.BanCreateRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}

			ban, err := adminService.BanUser(&req, user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusCreated, gin.H{"success": true, "data": ban})
		})

		admin.GET("/bans", func(c *gin.Context) {
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
			active := c.DefaultQuery("active", "true") == "true"

			bans, total, err := adminService.ListBans(page, perPage, active, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			totalPages := int((total + int64(perPage) - 1) / int64(perPage))
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": models.PaginatedResponse{
					Data:       bans,
					Page:       page,
					PerPage:    perPage,
					Total:      total,
					TotalPages: totalPages,
				},
			})
		})

		admin.POST("/bans/:id/unban", func(c *gin.Context) {
			banID, _ := strconv.Atoi(c.Param("id"))

			err := adminService.UnbanUser(banID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User unbanned"})
		})
	}

	// File serving and CDN
	router.Static("/files/mods", cfg.Storage.ModsPath)
	router.Static("/files/images", cfg.Storage.ImagesPath)

	router.GET("/download/:gameSlug/:modSlug", func(c *gin.Context) {
		gameSlug := c.Param("gameSlug")
		modSlug := c.Param("modSlug")

		mod, err := modService.GetBySlug(gameSlug, modSlug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mod not found"})
			return
		}

		if len(mod.Files) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No files available"})
			return
		}

		var mainFile *models.ModFile
		for _, file := range mod.Files {
			if file.IsMain {
				mainFile = &file
				break
			}
		}

		if mainFile == nil {
			mainFile = &mod.Files[0]
		}

		modService.IncrementDownloadCount(mod.ID)

		c.Header("Content-Disposition", "attachment; filename="+mainFile.Filename)
		c.File(mainFile.FilePath)
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "azurite-api",
		})
	})
}

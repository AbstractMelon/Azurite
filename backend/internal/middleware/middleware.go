package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"https://azurite.dev",
		}

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := utils.GetClientIP(c.Request.RemoteAddr, c.GetHeader("X-Forwarded-For"), c.GetHeader("X-Real-IP"))
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n",
			timestamp.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func CheckBan(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := utils.GetClientIP(c.Request.RemoteAddr, c.GetHeader("X-Forwarded-For"), c.GetHeader("X-Real-IP"))

		var count int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM bans
			WHERE ip_address = ? AND is_active = 1 AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
		`, clientIP).Scan(&count)

		if err == nil && count > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Your IP address has been banned from this service",
			})
			c.Abort()
			return
		}

		if userInterface, exists := c.Get("user"); exists {
			if user, ok := userInterface.(map[string]interface{}); ok {
				if userID, ok := user["id"].(int); ok {
					err := db.QueryRow(`
						SELECT COUNT(*) FROM bans
						WHERE user_id = ? AND is_active = 1 AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
					`, userID).Scan(&count)

					if err == nil && count > 0 {
						c.JSON(http.StatusForbidden, gin.H{
							"error": "Your account has been banned",
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

func FileUploadLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("File too large. Maximum size is %s", utils.FormatFileSize(maxSize)),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidateContentType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")

		for _, allowed := range allowedTypes {
			if strings.HasPrefix(contentType, allowed) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error": "Unsupported content type",
		})
		c.Abort()
	}
}

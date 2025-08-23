package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSlug(text string) string {
	text = strings.ToLower(text)
	reg := regexp.MustCompile(`[^a-z0-9\s\-_]`)
	text = reg.ReplaceAllString(text, "")
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, "-")
	text = regexp.MustCompile(`-+`).ReplaceAllString(text, "-")
	text = strings.Trim(text, "-")

	if text == "" {
		text = fmt.Sprintf("item-%d", time.Now().Unix())
	}

	return text
}

func GenerateRandomToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GenerateFileHash(file multipart.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func SaveUploadedFile(file multipart.File, filename, directory string) error {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(directory, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	file.Seek(0, 0)
	_, err = io.Copy(dst, file)
	return err
}

func IsAllowedFileType(filename string, allowedTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range allowedTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}

func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeTypes := map[string]string{
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".7z":   "application/x-7z-compressed",
		".tar":  "application/x-tar",
		".gz":   "application/gzip",
		".jar":  "application/java-archive",
		".dll":  "application/x-msdownload",
		".exe":  "application/x-msdownload",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".webp": "image/webp",
		".json": "application/json",
		".txt":  "text/plain",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}

func ContainsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return usernameRegex.MatchString(username)
}

func SanitizeFilename(filename string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = reg.ReplaceAllString(filename, "_")
	filename = regexp.MustCompile(`_{2,}`).ReplaceAllString(filename, "_")
	return strings.Trim(filename, "_")
}

func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func StringSliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func RemoveStringFromSlice(slice []string, str string) []string {
	result := make([]string, 0)
	for _, s := range slice {
		if s != str {
			result = append(result, s)
		}
	}
	return result
}

func GetClientIP(remoteAddr, xForwardedFor, xRealIP string) string {
	if xRealIP != "" {
		return xRealIP
	}

	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if remoteAddr != "" {
		if ip := strings.Split(remoteAddr, ":"); len(ip) > 0 {
			return ip[0]
		}
	}

	return "unknown"
}

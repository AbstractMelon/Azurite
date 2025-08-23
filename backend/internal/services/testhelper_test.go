package services

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/azurite/backend/internal/config"
	"github.com/azurite/backend/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) *database.DB {
	// Create temporary directory
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Create database
	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	if err := createTestSchema(db.DB); err != nil {
		t.Fatalf("Failed to create test schema: %v", err)
	}

	return db
}

// createTestSchema creates the test database schema
func createTestSchema(db *sql.DB) error {
	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT,
		display_name VARCHAR(100) NOT NULL,
		avatar TEXT,
		bio TEXT,
		role VARCHAR(20) DEFAULT 'user',
		is_active BOOLEAN DEFAULT true,
		email_verified BOOLEAN DEFAULT false,
		notify_email BOOLEAN DEFAULT true,
		notify_in_site BOOLEAN DEFAULT true,
		github_id TEXT UNIQUE,
		discord_id TEXT UNIQUE,
		google_id TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login_at DATETIME
	);

	-- Games table
	CREATE TABLE IF NOT EXISTS games (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(200) NOT NULL,
		slug VARCHAR(200) UNIQUE NOT NULL,
		description TEXT,
		icon TEXT,
		is_active BOOLEAN DEFAULT true,
		mod_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Tags table
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(100) NOT NULL,
		slug VARCHAR(100) NOT NULL,
		game_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
		UNIQUE(slug, game_id)
	);

	-- Mods table
	CREATE TABLE IF NOT EXISTS mods (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(200) NOT NULL,
		slug VARCHAR(200) NOT NULL,
		description TEXT NOT NULL,
		short_description VARCHAR(500) NOT NULL,
		icon TEXT,
		version VARCHAR(50) NOT NULL,
		game_version VARCHAR(50) NOT NULL,
		game_id INTEGER NOT NULL,
		owner_id INTEGER NOT NULL,
		downloads INTEGER DEFAULT 0,
		likes INTEGER DEFAULT 0,
		source_website TEXT,
		contact_info TEXT,
		is_rejected BOOLEAN DEFAULT false,
		rejection_reason TEXT,
		is_scanned BOOLEAN DEFAULT false,
		scan_result VARCHAR(20) DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(slug, game_id)
	);

	-- Mod files table
	CREATE TABLE IF NOT EXISTS mod_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mod_id INTEGER NOT NULL,
		filename VARCHAR(255) NOT NULL,
		file_path TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		mime_type VARCHAR(100),
		hash TEXT NOT NULL,
		is_main BOOLEAN DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE
	);

	-- Mod tags junction table
	CREATE TABLE IF NOT EXISTS mod_tags (
		mod_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (mod_id, tag_id),
		FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);

	-- Mod dependencies junction table
	CREATE TABLE IF NOT EXISTS mod_dependencies (
		mod_id INTEGER NOT NULL,
		dependency_id INTEGER NOT NULL,
		PRIMARY KEY (mod_id, dependency_id),
		FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE,
		FOREIGN KEY (dependency_id) REFERENCES mods(id) ON DELETE CASCADE
	);

	-- Mod likes table
	CREATE TABLE IF NOT EXISTS mod_likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mod_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(mod_id, user_id)
	);

	-- Comments table
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mod_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		parent_id INTEGER,
		is_active BOOLEAN DEFAULT true,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (mod_id) REFERENCES mods(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
	);

	-- Game requests table
	CREATE TABLE IF NOT EXISTS game_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(200) NOT NULL,
		description TEXT NOT NULL,
		icon TEXT,
		requested_by INTEGER NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		admin_notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (requested_by) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Documentation table
	CREATE TABLE IF NOT EXISTS documentation (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		game_id INTEGER NOT NULL,
		title VARCHAR(200) NOT NULL,
		slug VARCHAR(200) NOT NULL,
		content TEXT NOT NULL,
		author_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(slug, game_id)
	);

	-- Notifications table
	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		type VARCHAR(50) NOT NULL,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		data TEXT,
		is_read BOOLEAN DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- User roles table
	CREATE TABLE IF NOT EXISTS user_roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		game_id INTEGER,
		role VARCHAR(50) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
		UNIQUE(user_id, game_id, role)
	);

	-- Bans table
	CREATE TABLE IF NOT EXISTS bans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		ip_address TEXT,
		game_id INTEGER,
		reason TEXT NOT NULL,
		banned_by INTEGER NOT NULL,
		expires_at DATETIME,
		is_active BOOLEAN DEFAULT true,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
		FOREIGN KEY (banned_by) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Password reset tokens table
	CREATE TABLE IF NOT EXISTS password_reset_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT UNIQUE NOT NULL,
		expires_at DATETIME NOT NULL,
		used BOOLEAN DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err := db.Exec(schema)
	return err
}

// setupTestConfig creates a test configuration
func setupTestConfig() *config.Config {
	return &config.Config{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret-key-for-testing-only",
		},
		Email: config.EmailConfig{
			SMTPHost:     "",
			SMTPPort:     587,
			SMTPUsername: "",
			SMTPPassword: "",
			FromEmail:    "test@azurite.test",
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}
}

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *database.DB, username, email string) int {
	result, err := db.Exec(`
		INSERT INTO users (username, email, display_name, password_hash, role, avatar, bio)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, username, email, username+" Display", "hashed_password", "user", "", "")

	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	return int(userID)
}

// createTestGame creates a test game in the database
func createTestGame(t *testing.T, db *database.DB, name, slug string) int {
	result, err := db.Exec(`
		INSERT INTO games (name, slug, description, is_active)
		VALUES (?, ?, ?, ?)
	`, name, slug, "Test game description", true)

	if err != nil {
		t.Fatalf("Failed to create test game: %v", err)
	}

	gameID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get game ID: %v", err)
	}

	return int(gameID)
}

// createTestMod creates a test mod in the database
func createTestMod(t *testing.T, db *database.DB, name, slug string, gameID, ownerID int) int {
	result, err := db.Exec(`
		INSERT INTO mods (name, slug, description, short_description, version, game_version, game_id, owner_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, name, slug, "Test mod description", "Short description", "1.0.0", "1.0", gameID, ownerID)

	if err != nil {
		t.Fatalf("Failed to create test mod: %v", err)
	}

	modID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get mod ID: %v", err)
	}

	return int(modID)
}

// cleanupTestDB closes the test database
func cleanupTestDB(db *database.DB) {
	if db != nil {
		db.Close()
	}
}

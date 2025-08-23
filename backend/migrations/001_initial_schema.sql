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

-- User roles table (for community moderators and wiki maintainers)
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
);

-- Sessions table for JWT token blacklisting
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token_hash TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_mods_game_id ON mods(game_id);
CREATE INDEX IF NOT EXISTS idx_mods_owner_id ON mods(owner_id);
CREATE INDEX IF NOT EXISTS idx_mods_slug ON mods(slug);
CREATE INDEX IF NOT EXISTS idx_mods_created_at ON mods(created_at);
CREATE INDEX IF NOT EXISTS idx_mods_downloads ON mods(downloads);
CREATE INDEX IF NOT EXISTS idx_mods_likes ON mods(likes);
CREATE INDEX IF NOT EXISTS idx_comments_mod_id ON comments(mod_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_bans_user_id ON bans(user_id);
CREATE INDEX IF NOT EXISTS idx_bans_ip_address ON bans(ip_address);
CREATE INDEX IF NOT EXISTS idx_bans_game_id ON bans(game_id);
CREATE INDEX IF NOT EXISTS idx_bans_is_active ON bans(is_active);

DROP TRIGGER IF EXISTS update_game_mod_count_insert;
DROP TRIGGER IF EXISTS update_game_mod_count_update;
DROP TRIGGER IF EXISTS update_game_mod_count_delete;

DROP TRIGGER IF EXISTS update_users_updated_at;
DROP TRIGGER IF EXISTS update_mods_updated_at;
DROP TRIGGER IF EXISTS update_comments_updated_at;
DROP TRIGGER IF EXISTS update_games_updated_at;
DROP TRIGGER IF EXISTS update_documentation_updated_at;
DROP TRIGGER IF EXISTS update_game_requests_updated_at;


-- Triggers to update mod count in games table
CREATE TRIGGER update_game_mod_count_insert
AFTER INSERT ON mods
WHEN NEW.is_rejected = 0 AND NEW.is_scanned = 1 AND NEW.scan_result = 'clean'
BEGIN
    UPDATE games SET mod_count = mod_count + 1 WHERE id = NEW.game_id;
END;

CREATE TRIGGER update_game_mod_count_update
AFTER UPDATE ON mods
WHEN (OLD.is_rejected = 1 OR OLD.is_scanned = 0 OR OLD.scan_result != 'clean') AND
     (NEW.is_rejected = 0 AND NEW.is_scanned = 1 AND NEW.scan_result = 'clean')
BEGIN
    UPDATE games SET mod_count = mod_count + 1 WHERE id = NEW.game_id;
END;

CREATE TRIGGER update_game_mod_count_delete
AFTER UPDATE ON mods
WHEN (OLD.is_rejected = 0 AND OLD.is_scanned = 1 AND OLD.scan_result = 'clean') AND
     (NEW.is_rejected = 1 OR NEW.is_scanned = 0 OR NEW.scan_result != 'clean')
BEGIN
    UPDATE games SET mod_count = mod_count - 1 WHERE id = NEW.game_id;
END;

-- Trigger to update updated_at timestamps
CREATE TRIGGER update_users_updated_at
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_mods_updated_at
AFTER UPDATE ON mods
FOR EACH ROW
BEGIN
    UPDATE mods SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_comments_updated_at
AFTER UPDATE ON comments
FOR EACH ROW
BEGIN
    UPDATE comments SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_games_updated_at
AFTER UPDATE ON games
FOR EACH ROW
BEGIN
    UPDATE games SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_documentation_updated_at
AFTER UPDATE ON documentation
FOR EACH ROW
BEGIN
    UPDATE documentation SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_game_requests_updated_at
AFTER UPDATE ON game_requests
FOR EACH ROW
BEGIN
    UPDATE game_requests SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

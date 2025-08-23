package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID            int            `json:"id" db:"id"`
	Username      string         `json:"username" db:"username"`
	Email         string         `json:"email" db:"email"`
	PasswordHash  string         `json:"-" db:"password_hash"`
	DisplayName   string         `json:"display_name" db:"display_name"`
	Avatar        string         `json:"avatar" db:"avatar"`
	Bio           string         `json:"bio" db:"bio"`
	Role          string         `json:"role" db:"role"`
	IsActive      bool           `json:"is_active" db:"is_active"`
	EmailVerified bool           `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	LastLoginAt   *time.Time     `json:"last_login_at" db:"last_login_at"`
	NotifyEmail   bool           `json:"notify_email" db:"notify_email"`
	NotifyInSite  bool           `json:"notify_in_site" db:"notify_in_site"`
	GitHubID      sql.NullString `json:"github_id" db:"github_id"`
	DiscordID     sql.NullString `json:"discord_id" db:"discord_id"`
	GoogleID      sql.NullString `json:"google_id" db:"google_id"`
}

type Game struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	ModCount    int       `json:"mod_count" db:"mod_count"`
}

type Mod struct {
	ID               int       `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Slug             string    `json:"slug" db:"slug"`
	Description      string    `json:"description" db:"description"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Icon             string    `json:"icon" db:"icon"`
	Version          string    `json:"version" db:"version"`
	GameVersion      string    `json:"game_version" db:"game_version"`
	GameID           int       `json:"game_id" db:"game_id"`
	OwnerID          int       `json:"owner_id" db:"owner_id"`
	Downloads        int       `json:"downloads" db:"downloads"`
	Likes            int       `json:"likes" db:"likes"`
	SourceWebsite    string    `json:"source_website" db:"source_website"`
	ContactInfo      string    `json:"contact_info" db:"contact_info"`
	IsRejected       bool      `json:"is_rejected" db:"is_rejected"`
	RejectionReason  string    `json:"rejection_reason" db:"rejection_reason"`
	IsScanned        bool      `json:"is_scanned" db:"is_scanned"`
	ScanResult       string    `json:"scan_result" db:"scan_result"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	Game             *Game     `json:"game,omitempty"`
	Owner            *User     `json:"owner,omitempty"`
	Tags             []Tag     `json:"tags,omitempty"`
	Dependencies     []Mod     `json:"dependencies,omitempty"`
	Files            []ModFile `json:"files,omitempty"`
	IsLiked          bool      `json:"is_liked,omitempty"`
}

type ModFile struct {
	ID        int       `json:"id" db:"id"`
	ModID     int       `json:"mod_id" db:"mod_id"`
	Filename  string    `json:"filename" db:"filename"`
	FilePath  string    `json:"file_path" db:"file_path"`
	FileSize  int64     `json:"file_size" db:"file_size"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	Hash      string    `json:"hash" db:"hash"`
	IsMain    bool      `json:"is_main" db:"is_main"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Tag struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	GameID    int       `json:"game_id" db:"game_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ModTag struct {
	ModID int `json:"mod_id" db:"mod_id"`
	TagID int `json:"tag_id" db:"tag_id"`
}

type ModDependency struct {
	ModID        int `json:"mod_id" db:"mod_id"`
	DependencyID int `json:"dependency_id" db:"dependency_id"`
}

type ModLike struct {
	ID        int       `json:"id" db:"id"`
	ModID     int       `json:"mod_id" db:"mod_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Comment struct {
	ID        int       `json:"id" db:"id"`
	ModID     int       `json:"mod_id" db:"mod_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	ParentID  *int      `json:"parent_id" db:"parent_id"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	User      *User     `json:"user,omitempty"`
	Replies   []Comment `json:"replies,omitempty"`
}

type GameRequest struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	RequestedBy int       `json:"requested_by" db:"requested_by"`
	Status      string    `json:"status" db:"status"`
	AdminNotes  string    `json:"admin_notes" db:"admin_notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	User        *User     `json:"user,omitempty"`
}

type Documentation struct {
	ID        int       `json:"id" db:"id"`
	GameID    int       `json:"game_id" db:"game_id"`
	Title     string    `json:"title" db:"title"`
	Slug      string    `json:"slug" db:"slug"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Author    *User     `json:"author,omitempty"`
	Game      *Game     `json:"game,omitempty"`
}

type Notification struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"`
	Title     string    `json:"title" db:"title"`
	Message   string    `json:"message" db:"message"`
	Data      string    `json:"data" db:"data"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserRole struct {
	ID     int    `json:"id" db:"id"`
	UserID int    `json:"user_id" db:"user_id"`
	GameID *int   `json:"game_id" db:"game_id"`
	Role   string `json:"role" db:"role"`
}

type Ban struct {
	ID           int        `json:"id" db:"id"`
	UserID       *int       `json:"user_id" db:"user_id"`
	IPAddress    string     `json:"ip_address" db:"ip_address"`
	GameID       *int       `json:"game_id" db:"game_id"`
	Reason       string     `json:"reason" db:"reason"`
	BannedBy     int        `json:"banned_by" db:"banned_by"`
	ExpiresAt    *time.Time `json:"expires_at" db:"expires_at"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	User         *User      `json:"user,omitempty"`
	BannedByUser *User      `json:"banned_by_user,omitempty"`
	Game         *Game      `json:"game,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetConfirm struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type ModCreateRequest struct {
	Name             string   `json:"name" binding:"required,min=1,max=200"`
	Description      string   `json:"description" binding:"required"`
	ShortDescription string   `json:"short_description" binding:"required,max=500"`
	Version          string   `json:"version" binding:"required"`
	GameVersion      string   `json:"game_version" binding:"required"`
	GameID           int      `json:"game_id" binding:"required"`
	SourceWebsite    string   `json:"source_website"`
	ContactInfo      string   `json:"contact_info"`
	Tags             []string `json:"tags"`
	Dependencies     []int    `json:"dependencies"`
}

type ModUpdateRequest struct {
	Name             string   `json:"name" binding:"required,min=1,max=200"`
	Description      string   `json:"description" binding:"required"`
	ShortDescription string   `json:"short_description" binding:"required,max=500"`
	Version          string   `json:"version" binding:"required"`
	GameVersion      string   `json:"game_version" binding:"required"`
	SourceWebsite    string   `json:"source_website"`
	ContactInfo      string   `json:"contact_info"`
	Tags             []string `json:"tags"`
	Dependencies     []int    `json:"dependencies"`
}

type CommentCreateRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	ParentID *int   `json:"parent_id"`
}

type GameRequestCreate struct {
	Name        string `json:"name" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"required"`
}

type DocumentationCreateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required"`
}

type UserUpdateRequest struct {
	DisplayName  string `json:"display_name" binding:"required,min=1,max=100"`
	Bio          string `json:"bio"`
	NotifyEmail  bool   `json:"notify_email"`
	NotifyInSite bool   `json:"notify_in_site"`
}

type BanCreateRequest struct {
	UserID    *int   `json:"user_id"`
	IPAddress string `json:"ip_address"`
	GameID    *int   `json:"game_id"`
	Reason    string `json:"reason" binding:"required"`
	Duration  int    `json:"duration"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

const (
	RoleUser               = "user"
	RoleAdmin              = "admin"
	RoleCommunityModerator = "community_moderator"
	RoleWikiMaintainer     = "wiki_maintainer"
)

const (
	GameRequestStatusPending  = "pending"
	GameRequestStatusApproved = "approved"
	GameRequestStatusDenied   = "denied"
)

const (
	NotificationTypeModRejected  = "mod_rejected"
	NotificationTypeNewComment   = "new_comment"
	NotificationTypeModMilestone = "mod_milestone"
	NotificationTypeGameRequest  = "game_request"
	NotificationTypeModApproved  = "mod_approved"
)

const (
	ScanResultPending = "pending"
	ScanResultClean   = "clean"
	ScanResultThreat  = "threat"
)

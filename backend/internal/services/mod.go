package services

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/models"
	"github.com/azurite/backend/internal/utils"
)

type ModService struct {
	db          *database.DB
	storagePath string
	imagesPath  string
}

func NewModService(db *database.DB, storagePath, imagesPath string) *ModService {
	return &ModService{
		db:          db,
		storagePath: storagePath,
		imagesPath:  imagesPath,
	}
}

func (s *ModService) Create(req *models.ModCreateRequest, ownerID int) (*models.Mod, error) {
	slug := utils.GenerateSlug(req.Name)
	originalSlug := slug
	counter := 1

	for {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM mods WHERE slug = ? AND game_id = ?", slug, req.GameID).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug availability: %w", err)
		}

		if count == 0 {
			break
		}

		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}

	result, err := s.db.Exec(`
		INSERT INTO mods (name, slug, description, short_description, version, game_version,
		                 game_id, owner_id, source_website, contact_info, is_scanned, scan_result)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.Name, slug, req.Description, req.ShortDescription, req.Version,
		req.GameVersion, req.GameID, ownerID, req.SourceWebsite, req.ContactInfo,
		false, models.ScanResultPending)

	if err != nil {
		return nil, fmt.Errorf("failed to create mod: %w", err)
	}

	modID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get mod ID: %w", err)
	}

	mod, err := s.GetByID(int(modID))
	if err != nil {
		return nil, err
	}

	if len(req.Tags) > 0 {
		err = s.updateModTags(int(modID), req.Tags, req.GameID)
		if err != nil {
			return nil, fmt.Errorf("failed to update tags: %w", err)
		}
	}

	if len(req.Dependencies) > 0 {
		err = s.updateModDependencies(int(modID), req.Dependencies)
		if err != nil {
			return nil, fmt.Errorf("failed to update dependencies: %w", err)
		}
	}

	go s.simulateMalwareScan(int(modID))

	return mod, nil
}

func (s *ModService) Update(modID int, req *models.ModUpdateRequest, ownerID int) (*models.Mod, error) {
	mod, err := s.GetByID(modID)
	if err != nil {
		return nil, err
	}

	if mod.OwnerID != ownerID {
		return nil, errors.New("unauthorized to update this mod")
	}

	_, err = s.db.Exec(`
		UPDATE mods SET name = ?, description = ?, short_description = ?, version = ?,
		               game_version = ?, source_website = ?, contact_info = ?
		WHERE id = ?
	`, req.Name, req.Description, req.ShortDescription, req.Version,
		req.GameVersion, req.SourceWebsite, req.ContactInfo, modID)

	if err != nil {
		return nil, fmt.Errorf("failed to update mod: %w", err)
	}

	if len(req.Tags) > 0 || len(req.Dependencies) > 0 {
		tx, err := s.db.Begin()
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback()

		if len(req.Tags) > 0 {
			_, err = tx.Exec("DELETE FROM mod_tags WHERE mod_id = ?", modID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete old tags: %w", err)
			}

			err = s.updateModTagsInTx(tx, modID, req.Tags, mod.GameID)
			if err != nil {
				return nil, fmt.Errorf("failed to update tags: %w", err)
			}
		}

		if len(req.Dependencies) > 0 {
			_, err = tx.Exec("DELETE FROM mod_dependencies WHERE mod_id = ?", modID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete old dependencies: %w", err)
			}

			err = s.updateModDependenciesInTx(tx, modID, req.Dependencies)
			if err != nil {
				return nil, fmt.Errorf("failed to update dependencies: %w", err)
			}
		}

		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return s.GetByID(modID)
}

func (s *ModService) GetByID(id int) (*models.Mod, error) {
	mod := &models.Mod{
		Game:  &models.Game{},
		Owner: &models.User{},
	}

	// Use sql.NullString for nullable fields
	var icon, rejectionReason sql.NullString

	err := s.db.QueryRow(`
        SELECT m.id, m.name, m.slug, m.description, m.short_description, m.icon,
               m.version, m.game_version, m.game_id, m.owner_id, m.downloads, m.likes,
               m.source_website, m.contact_info, m.is_rejected, m.rejection_reason,
               m.is_scanned, m.scan_result, m.created_at, m.updated_at,
               g.name, g.slug, u.username, u.display_name
        FROM mods m
        JOIN games g ON m.game_id = g.id
        JOIN users u ON m.owner_id = u.id
        WHERE m.id = ?
    `, id).Scan(
		&mod.ID, &mod.Name, &mod.Slug, &mod.Description, &mod.ShortDescription,
		&icon, &mod.Version, &mod.GameVersion, &mod.GameID, &mod.OwnerID,
		&mod.Downloads, &mod.Likes, &mod.SourceWebsite, &mod.ContactInfo,
		&mod.IsRejected, &rejectionReason, &mod.IsScanned, &mod.ScanResult,
		&mod.CreatedAt, &mod.UpdatedAt, &mod.Game.Name, &mod.Game.Slug,
		&mod.Owner.Username, &mod.Owner.DisplayName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("mod not found")
		}
		return nil, fmt.Errorf("failed to get mod: %w", err)
	}

	// Convert nullable fields
	if icon.Valid {
		mod.Icon = icon.String
	}
	if rejectionReason.Valid {
		mod.RejectionReason = rejectionReason.String
	}

	tags, err := s.getModTags(id)
	if err == nil {
		mod.Tags = tags
	}

	dependencies, err := s.getModDependencies(id)
	if err == nil {
		mod.Dependencies = dependencies
	}

	files, err := s.getModFiles(id)
	if err == nil {
		mod.Files = files
	}

	return mod, nil
}

func (s *ModService) GetBySlug(gameSlug, modSlug string) (*models.Mod, error) {
    mod := &models.Mod{
        Game:  &models.Game{},
        Owner: &models.User{},
    }

	// Use sql.NullString for nullable fields
    var icon, rejectionReason sql.NullString

    err := s.db.QueryRow(`
        SELECT m.id, m.name, m.slug, m.description, m.short_description, m.icon,
               m.version, m.game_version, m.game_id, m.owner_id, m.downloads, m.likes,
               m.source_website, m.contact_info, m.is_rejected, m.rejection_reason,
               m.is_scanned, m.scan_result, m.created_at, m.updated_at,
               g.name, g.slug, u.username, u.display_name
        FROM mods m
        JOIN games g ON m.game_id = g.id
        JOIN users u ON m.owner_id = u.id
        WHERE m.slug = ? AND g.slug = ?
    `, modSlug, gameSlug).Scan(
        &mod.ID, &mod.Name, &mod.Slug, &mod.Description, &mod.ShortDescription,
        &icon, &mod.Version, &mod.GameVersion, &mod.GameID, &mod.OwnerID,
        &mod.Downloads, &mod.Likes, &mod.SourceWebsite, &mod.ContactInfo,
        &mod.IsRejected, &rejectionReason, &mod.IsScanned, &mod.ScanResult,
        &mod.CreatedAt, &mod.UpdatedAt, &mod.Game.Name, &mod.Game.Slug,
        &mod.Owner.Username, &mod.Owner.DisplayName,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("mod not found")
        }
        return nil, fmt.Errorf("failed to get mod: %w", err)
    }

	// Convert nullable fields
    if icon.Valid {
        mod.Icon = icon.String
    }
    if rejectionReason.Valid {
        mod.RejectionReason = rejectionReason.String
    }

    // set nested IDs on already initialized structs
    mod.Game.ID = mod.GameID
    mod.Owner.ID = mod.OwnerID

	tags, err := s.getModTags(mod.ID)
	if err == nil {
		mod.Tags = tags
	}

	dependencies, err := s.getModDependencies(mod.ID)
	if err == nil {
		mod.Dependencies = dependencies
	}

	files, err := s.getModFiles(mod.ID)
	if err == nil {
		mod.Files = files
	}

	return mod, nil
}

func (s *ModService) ListByGame(gameID int, page, perPage int, sortBy, order string, tags []string, search string) ([]models.Mod, int64, error) {
	offset := (page - 1) * perPage

	whereClause := "WHERE m.game_id = ? AND m.is_rejected = 0 AND m.is_scanned = 1 AND m.scan_result = 'clean'"
	args := []interface{}{gameID}

	if search != "" {
		whereClause += " AND (m.name LIKE ? OR m.description LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if len(tags) > 0 {
		placeholders := strings.Repeat("?,", len(tags))
		placeholders = placeholders[:len(placeholders)-1]
		whereClause += fmt.Sprintf(" AND m.id IN (SELECT mt.mod_id FROM mod_tags mt JOIN tags t ON mt.tag_id = t.id WHERE t.slug IN (%s))", placeholders)
		for _, tag := range tags {
			args = append(args, tag)
		}
	}

	orderClause := "ORDER BY m.created_at DESC"
	switch sortBy {
	case "name":
		if order == "desc" {
			orderClause = "ORDER BY m.name DESC"
		} else {
			orderClause = "ORDER BY m.name ASC"
		}
	case "downloads":
		if order == "asc" {
			orderClause = "ORDER BY m.downloads ASC"
		} else {
			orderClause = "ORDER BY m.downloads DESC"
		}
	case "likes":
		if order == "asc" {
			orderClause = "ORDER BY m.likes ASC"
		} else {
			orderClause = "ORDER BY m.likes DESC"
		}
	case "updated":
		if order == "asc" {
			orderClause = "ORDER BY m.updated_at ASC"
		} else {
			orderClause = "ORDER BY m.updated_at DESC"
		}
	}

	var total int64
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM mods m
		JOIN games g ON m.game_id = g.id
		JOIN users u ON m.owner_id = u.id
		%s
	`, whereClause)

	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mods: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT m.id, m.name, m.slug, m.description, m.short_description, m.icon,
		       m.version, m.game_version, m.game_id, m.owner_id, m.downloads, m.likes,
		       m.source_website, m.contact_info, m.is_rejected, m.rejection_reason,
		       m.is_scanned, m.scan_result, m.created_at, m.updated_at,
		       g.name, g.slug, u.username, u.display_name
		FROM mods m
		JOIN games g ON m.game_id = g.id
		JOIN users u ON m.owner_id = u.id
		%s %s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	args = append(args, perPage, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get mods: %w", err)
	}
	defer rows.Close()

	var mods []models.Mod
	for rows.Next() {
		var mod models.Mod
		mod.Game = &models.Game{}
		mod.Owner = &models.User{}

		// Use sql.NullString for nullable fields
		var icon, rejectionReason sql.NullString

		err := rows.Scan(
			&mod.ID, &mod.Name, &mod.Slug, &mod.Description, &mod.ShortDescription,
			&icon, &mod.Version, &mod.GameVersion, &mod.GameID, &mod.OwnerID,
			&mod.Downloads, &mod.Likes, &mod.SourceWebsite, &mod.ContactInfo,
			&mod.IsRejected, &rejectionReason, &mod.IsScanned, &mod.ScanResult,
			&mod.CreatedAt, &mod.UpdatedAt, &mod.Game.Name, &mod.Game.Slug,
			&mod.Owner.Username, &mod.Owner.DisplayName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mod: %w", err)
		}

		// Convert nullable fields
		if icon.Valid {
			mod.Icon = icon.String
		}
		if rejectionReason.Valid {
			mod.RejectionReason = rejectionReason.String
		}

		tags, err := s.getModTags(mod.ID)
		if err == nil {
			mod.Tags = tags
		}

		mods = append(mods, mod)
	}

	return mods, total, nil
}

func (s *ModService) ListByOwner(ownerID int, page, perPage int) ([]models.Mod, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM mods WHERE owner_id = ?", ownerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mods: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT m.id, m.name, m.slug, m.description, m.short_description, m.icon,
		       m.version, m.game_version, m.game_id, m.owner_id, m.downloads, m.likes,
		       m.source_website, m.contact_info, m.is_rejected, m.rejection_reason,
		       m.is_scanned, m.scan_result, m.created_at, m.updated_at,
		       g.name, g.slug
		FROM mods m
		JOIN games g ON m.game_id = g.id
		WHERE m.owner_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`, ownerID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get mods: %w", err)
	}
	defer rows.Close()

	var mods []models.Mod
	for rows.Next() {
		var mod models.Mod
		mod.Game = &models.Game{}
		mod.Owner = &models.User{ID: ownerID}

		// Use sql.NullString for nullable fields
		var icon, rejectionReason sql.NullString

		err := rows.Scan(
			&mod.ID, &mod.Name, &mod.Slug, &mod.Description, &mod.ShortDescription,
			&icon, &mod.Version, &mod.GameVersion, &mod.GameID, &mod.OwnerID,
			&mod.Downloads, &mod.Likes, &mod.SourceWebsite, &mod.ContactInfo,
			&mod.IsRejected, &rejectionReason, &mod.IsScanned, &mod.ScanResult,
			&mod.CreatedAt, &mod.UpdatedAt, &mod.Game.Name, &mod.Game.Slug,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mod: %w", err)
		}

		// Convert nullable fields
		if icon.Valid {
			mod.Icon = icon.String
		}
		if rejectionReason.Valid {
			mod.RejectionReason = rejectionReason.String
		}

		mods = append(mods, mod)
	}

	return mods, total, nil
}

func (s *ModService) UploadFile(modID int, file multipart.File, header *multipart.FileHeader, isMain bool) (*models.ModFile, error) {
	allowedTypes := []string{".zip", ".rar", ".7z", ".jar", ".dll", ".exe", ".json", ".txt"}
	if !utils.IsAllowedFileType(header.Filename, allowedTypes) {
		return nil, errors.New("file type not allowed")
	}

	sanitizedFilename := utils.SanitizeFilename(header.Filename)
	hash, err := utils.GenerateFileHash(file)
	if err != nil {
		return nil, fmt.Errorf("failed to generate file hash: %w", err)
	}

	modDir := filepath.Join(s.storagePath, strconv.Itoa(modID))
	if err := os.MkdirAll(modDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create mod directory: %w", err)
	}

	filePath := filepath.Join(modDir, sanitizedFilename)
	if err := utils.SaveUploadedFile(file, sanitizedFilename, modDir); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	if isMain {
		_, err = s.db.Exec("UPDATE mod_files SET is_main = 0 WHERE mod_id = ?", modID)
		if err != nil {
			return nil, fmt.Errorf("failed to update main file status: %w", err)
		}
	}

	result, err := s.db.Exec(`
		INSERT INTO mod_files (mod_id, filename, file_path, file_size, mime_type, hash, is_main)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, modID, header.Filename, filePath, header.Size, utils.GetMimeType(header.Filename), hash, isMain)

	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	fileID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get file ID: %w", err)
	}

	go s.simulateFileScan(modID, int(fileID))

	return &models.ModFile{
		ID:       int(fileID),
		ModID:    modID,
		Filename: header.Filename,
		FilePath: filePath,
		FileSize: header.Size,
		MimeType: utils.GetMimeType(header.Filename),
		Hash:     hash,
		IsMain:   isMain,
	}, nil
}

func (s *ModService) Like(modID, userID int) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM mod_likes WHERE mod_id = ? AND user_id = ?", modID, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing like: %w", err)
	}

	if count > 0 {
		return errors.New("mod already liked")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO mod_likes (mod_id, user_id) VALUES (?, ?)", modID, userID)
	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}

	_, err = tx.Exec("UPDATE mods SET likes = likes + 1 WHERE id = ?", modID)
	if err != nil {
		return fmt.Errorf("failed to update like count: %w", err)
	}

	return tx.Commit()
}

func (s *ModService) Unlike(modID, userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec("DELETE FROM mod_likes WHERE mod_id = ? AND user_id = ?", modID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("like not found")
	}

	_, err = tx.Exec("UPDATE mods SET likes = likes - 1 WHERE id = ?", modID)
	if err != nil {
		return fmt.Errorf("failed to update like count: %w", err)
	}

	return tx.Commit()
}

func (s *ModService) IsLiked(modID, userID int) bool {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM mod_likes WHERE mod_id = ? AND user_id = ?", modID, userID).Scan(&count)
	return err == nil && count > 0
}

func (s *ModService) IncrementDownloadCount(modID int) error {
	_, err := s.db.Exec("UPDATE mods SET downloads = downloads + 1 WHERE id = ?", modID)
	return err
}

func (s *ModService) Reject(modID int, reason string, rejectedBy int) error {
	_, err := s.db.Exec(`
		UPDATE mods SET is_rejected = 1, rejection_reason = ?
		WHERE id = ?
	`, reason, modID)

	if err != nil {
		return fmt.Errorf("failed to reject mod: %w", err)
	}

	return nil
}

func (s *ModService) Approve(modID int) error {
	_, err := s.db.Exec(`
		UPDATE mods SET is_rejected = 0, rejection_reason = NULL
		WHERE id = ?
	`, modID)

	if err != nil {
		return fmt.Errorf("failed to approve mod: %w", err)
	}

	return nil
}

func (s *ModService) Search(query string, gameID int, page, perPage int) ([]models.Mod, int64, error) {
	offset := (page - 1) * perPage
	searchTerm := "%" + query + "%"

	whereClause := "WHERE m.is_rejected = 0 AND m.is_scanned = 1 AND m.scan_result = 'clean' AND (m.name LIKE ? OR m.description LIKE ? OR m.short_description LIKE ?)"
	args := []interface{}{searchTerm, searchTerm, searchTerm}

	if gameID > 0 {
		whereClause += " AND m.game_id = ?"
		args = append(args, gameID)
	}

	var total int64
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM mods m %s
	`, whereClause)

	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	searchQuery := fmt.Sprintf(`
		SELECT
			m.id, m.name, m.slug, m.description, m.short_description,
			m.icon, m.version, m.game_version, m.game_id, m.owner_id,
			m.downloads, m.likes, m.source_website, m.contact_info,
			m.is_rejected, m.rejection_reason, m.is_scanned, m.scan_result,
			m.created_at, m.updated_at,
			g.name, g.slug, g.icon,
			u.username, u.display_name, u.avatar
		FROM mods m
		LEFT JOIN games g ON m.game_id = g.id
		LEFT JOIN users u ON m.owner_id = u.id
		%s
		ORDER BY m.downloads DESC, m.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, perPage, offset)
	rows, err := s.db.Query(searchQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search mods: %w", err)
	}
	defer rows.Close()

	var mods []models.Mod
	for rows.Next() {
		var mod models.Mod
		var game models.Game
		var owner models.User

		// Use sql.NullString for nullable fields
		var modIcon, gameIcon, userAvatar, rejectionReason sql.NullString

		err := rows.Scan(
			&mod.ID, &mod.Name, &mod.Slug, &mod.Description, &mod.ShortDescription,
			&modIcon, &mod.Version, &mod.GameVersion, &mod.GameID, &mod.OwnerID,
			&mod.Downloads, &mod.Likes, &mod.SourceWebsite, &mod.ContactInfo,
			&mod.IsRejected, &rejectionReason, &mod.IsScanned, &mod.ScanResult,
			&mod.CreatedAt, &mod.UpdatedAt,
			&game.Name, &game.Slug, &gameIcon,
			&owner.Username, &owner.DisplayName, &userAvatar,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mod: %w", err)
		}

		// Convert nullable fields
		if modIcon.Valid {
			mod.Icon = modIcon.String
		}
		if gameIcon.Valid {
			game.Icon = gameIcon.String
		}
		if userAvatar.Valid {
			owner.Avatar = userAvatar.String
		}
		if rejectionReason.Valid {
			mod.RejectionReason = rejectionReason.String
		}

		game.ID = mod.GameID
		owner.ID = mod.OwnerID

		mod.Game = &game
		mod.Owner = &owner

		mods = append(mods, mod)
	}

	return mods, total, nil
}

func (s *ModService) Delete(modID int, ownerID int) error {
	mod, err := s.GetByID(modID)
	if err != nil {
		return err
	}

	if mod.OwnerID != ownerID {
		return errors.New("unauthorized to delete this mod")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	files, err := s.getModFiles(modID)
	if err == nil {
		for _, file := range files {
			os.Remove(file.FilePath)
		}
	}

	_, err = tx.Exec("DELETE FROM mod_files WHERE mod_id = ?", modID)
	if err != nil {
		return fmt.Errorf("failed to delete mod files: %w", err)
	}

	_, err = tx.Exec("DELETE FROM mod_tags WHERE mod_id = ?", modID)
	if err != nil {
		return fmt.Errorf("failed to delete mod tags: %w", err)
	}

	_, err = tx.Exec("DELETE FROM mod_dependencies WHERE mod_id = ? OR dependency_id = ?", modID, modID)
	if err != nil {
		return fmt.Errorf("failed to delete mod dependencies: %w", err)
	}

	_, err = tx.Exec("DELETE FROM mod_likes WHERE mod_id = ?", modID)
	if err != nil {
		return fmt.Errorf("failed to delete mod likes: %w", err)
	}

	_, err = tx.Exec("DELETE FROM comments WHERE mod_id = ?", modID)
	if err != nil {
		return fmt.Errorf("failed to delete mod comments: %w", err)
	}

	_, err = tx.Exec("DELETE FROM mods WHERE id = ?", modID)
	if err != nil {
		return fmt.Errorf("failed to delete mod: %w", err)
	}

	return tx.Commit()
}

func (s *ModService) getModTags(modID int) ([]models.Tag, error) {
	rows, err := s.db.Query(`
		SELECT t.id, t.name, t.slug, t.game_id, t.created_at
		FROM tags t
		JOIN mod_tags mt ON t.id = mt.tag_id
		WHERE mt.mod_id = ?
	`, modID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.GameID, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *ModService) getModDependencies(modID int) ([]models.Mod, error) {
	rows, err := s.db.Query(`
		SELECT m.id, m.name, m.slug, m.version, m.game_id
		FROM mods m
		JOIN mod_dependencies md ON m.id = md.dependency_id
		WHERE md.mod_id = ?
	`, modID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dependencies []models.Mod
	for rows.Next() {
		var dep models.Mod
		err := rows.Scan(&dep.ID, &dep.Name, &dep.Slug, &dep.Version, &dep.GameID)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

func (s *ModService) getModFiles(modID int) ([]models.ModFile, error) {
	rows, err := s.db.Query(`
		SELECT id, mod_id, filename, file_path, file_size, mime_type, hash, is_main, created_at
		FROM mod_files
		WHERE mod_id = ?
		ORDER BY is_main DESC, created_at ASC
	`, modID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.ModFile
	for rows.Next() {
		var file models.ModFile
		err := rows.Scan(&file.ID, &file.ModID, &file.Filename, &file.FilePath,
			&file.FileSize, &file.MimeType, &file.Hash, &file.IsMain, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (s *ModService) updateModTags(modID int, tagNames []string, gameID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	return s.updateModTagsInTx(tx, modID, tagNames, gameID)
}

func (s *ModService) updateModTagsInTx(tx *sql.Tx, modID int, tagNames []string, gameID int) error {
	for _, tagName := range tagNames {
		tagSlug := utils.GenerateSlug(tagName)

		var tagID int
		err := tx.QueryRow("SELECT id FROM tags WHERE slug = ? AND game_id = ?", tagSlug, gameID).Scan(&tagID)
		if err == sql.ErrNoRows {
			result, err := tx.Exec("INSERT INTO tags (name, slug, game_id) VALUES (?, ?, ?)", tagName, tagSlug, gameID)
			if err != nil {
				return err
			}
			id, err := result.LastInsertId()
			if err != nil {
				return err
			}
			tagID = int(id)
		} else if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT OR IGNORE INTO mod_tags (mod_id, tag_id) VALUES (?, ?)", modID, tagID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *ModService) updateModDependencies(modID int, dependencyIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	return s.updateModDependenciesInTx(tx, modID, dependencyIDs)
}

func (s *ModService) updateModDependenciesInTx(tx *sql.Tx, modID int, dependencyIDs []int) error {
	for _, depID := range dependencyIDs {
		if depID == modID {
			continue
		}

		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM mods WHERE id = ?", depID).Scan(&count)
		if err != nil || count == 0 {
			continue
		}

		_, err = tx.Exec("INSERT OR IGNORE INTO mod_dependencies (mod_id, dependency_id) VALUES (?, ?)", modID, depID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *ModService) simulateMalwareScan(modID int) {
	time.Sleep(5 * time.Second)

	s.db.Exec(`
		UPDATE mods SET is_scanned = 1, scan_result = 'clean'
		WHERE id = ?
	`, modID)
}

func (s *ModService) simulateFileScan(modID, fileID int) {
	time.Sleep(3 * time.Second)

	// Simulate occasional threat detection for demo purposes
	if fileID%13 == 0 {
		s.db.Exec(`
			UPDATE mods SET is_scanned = 1, scan_result = 'threat', is_rejected = 1,
			              rejection_reason = 'Malware detected in uploaded file'
			WHERE id = ?
		`, modID)
	} else {
		s.db.Exec(`
			UPDATE mods SET is_scanned = 1, scan_result = 'clean'
			WHERE id = ?
		`, modID)
	}
}

package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/models"
	"github.com/azurite/backend/internal/utils"
)

type DocumentationService struct {
	db *database.DB
}

func NewDocumentationService(db *database.DB) *DocumentationService {
	return &DocumentationService{
		db: db,
	}
}

func (s *DocumentationService) Create(req *models.DocumentationCreateRequest, gameID, authorID int) (*models.Documentation, error) {
	slug := utils.GenerateSlug(req.Title)
	originalSlug := slug
	counter := 1

	for {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM documentation WHERE slug = ? AND game_id = ?", slug, gameID).Scan(&count)
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
		INSERT INTO documentation (game_id, title, slug, content, author_id)
		VALUES (?, ?, ?, ?, ?)
	`, gameID, req.Title, slug, req.Content, authorID)

	if err != nil {
		return nil, fmt.Errorf("failed to create documentation: %w", err)
	}

	docID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get documentation ID: %w", err)
	}

	return s.GetByID(int(docID))
}

func (s *DocumentationService) GetByID(id int) (*models.Documentation, error) {
	doc := &models.Documentation{}
	var authorUsername, authorDisplayName, gameName, gameSlug string

	err := s.db.QueryRow(`
		SELECT d.id, d.game_id, d.title, d.slug, d.content, d.author_id,
		       d.created_at, d.updated_at, u.username, u.display_name,
		       g.name, g.slug
		FROM documentation d
		JOIN users u ON d.author_id = u.id
		JOIN games g ON d.game_id = g.id
		WHERE d.id = ?
	`, id).Scan(
		&doc.ID, &doc.GameID, &doc.Title, &doc.Slug, &doc.Content,
		&doc.AuthorID, &doc.CreatedAt, &doc.UpdatedAt,
		&authorUsername, &authorDisplayName, &gameName, &gameSlug,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("documentation not found")
		}
		return nil, fmt.Errorf("failed to get documentation: %w", err)
	}

	doc.Author = &models.User{
		ID:          doc.AuthorID,
		Username:    authorUsername,
		DisplayName: authorDisplayName,
	}

	doc.Game = &models.Game{
		ID:   doc.GameID,
		Name: gameName,
		Slug: gameSlug,
	}

	return doc, nil
}

func (s *DocumentationService) GetBySlug(gameSlug, docSlug string) (*models.Documentation, error) {
	doc := &models.Documentation{}
	var authorUsername, authorDisplayName, gameName string

	err := s.db.QueryRow(`
		SELECT d.id, d.game_id, d.title, d.slug, d.content, d.author_id,
		       d.created_at, d.updated_at, u.username, u.display_name, g.name
		FROM documentation d
		JOIN users u ON d.author_id = u.id
		JOIN games g ON d.game_id = g.id
		WHERE d.slug = ? AND g.slug = ?
	`, docSlug, gameSlug).Scan(
		&doc.ID, &doc.GameID, &doc.Title, &doc.Slug, &doc.Content,
		&doc.AuthorID, &doc.CreatedAt, &doc.UpdatedAt,
		&authorUsername, &authorDisplayName, &gameName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("documentation not found")
		}
		return nil, fmt.Errorf("failed to get documentation: %w", err)
	}

	doc.Author = &models.User{
		ID:          doc.AuthorID,
		Username:    authorUsername,
		DisplayName: authorDisplayName,
	}

	doc.Game = &models.Game{
		ID:   doc.GameID,
		Name: gameName,
		Slug: gameSlug,
	}

	return doc, nil
}

func (s *DocumentationService) ListByGame(gameID int, page, perPage int) ([]models.Documentation, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM documentation WHERE game_id = ?", gameID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documentation: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT d.id, d.game_id, d.title, d.slug, d.content, d.author_id,
		       d.created_at, d.updated_at, u.username, u.display_name
		FROM documentation d
		JOIN users u ON d.author_id = u.id
		WHERE d.game_id = ?
		ORDER BY d.created_at DESC
		LIMIT ? OFFSET ?
	`, gameID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get documentation: %w", err)
	}
	defer rows.Close()

	var docs []models.Documentation
	for rows.Next() {
		var doc models.Documentation
		var authorUsername, authorDisplayName string

		err := rows.Scan(
			&doc.ID, &doc.GameID, &doc.Title, &doc.Slug, &doc.Content,
			&doc.AuthorID, &doc.CreatedAt, &doc.UpdatedAt,
			&authorUsername, &authorDisplayName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan documentation: %w", err)
		}

		doc.Author = &models.User{
			ID:          doc.AuthorID,
			Username:    authorUsername,
			DisplayName: authorDisplayName,
		}

		docs = append(docs, doc)
	}

	return docs, total, nil
}

func (s *DocumentationService) Update(docID int, req *models.DocumentationCreateRequest, authorID int, canEdit bool) (*models.Documentation, error) {
	doc, err := s.GetByID(docID)
	if err != nil {
		return nil, err
	}

	if !canEdit && doc.AuthorID != authorID {
		return nil, errors.New("unauthorized to update this documentation")
	}

	_, err = s.db.Exec(`
		UPDATE documentation SET title = ?, content = ?
		WHERE id = ?
	`, req.Title, req.Content, docID)

	if err != nil {
		return nil, fmt.Errorf("failed to update documentation: %w", err)
	}

	return s.GetByID(docID)
}

func (s *DocumentationService) Delete(docID int, authorID int, canDelete bool) error {
	doc, err := s.GetByID(docID)
	if err != nil {
		return err
	}

	if !canDelete && doc.AuthorID != authorID {
		return errors.New("unauthorized to delete this documentation")
	}

	_, err = s.db.Exec("DELETE FROM documentation WHERE id = ?", docID)
	if err != nil {
		return fmt.Errorf("failed to delete documentation: %w", err)
	}

	return nil
}

func (s *DocumentationService) GetByAuthor(authorID int, page, perPage int) ([]models.Documentation, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM documentation WHERE author_id = ?", authorID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documentation: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT d.id, d.game_id, d.title, d.slug, d.content, d.author_id,
		       d.created_at, d.updated_at, g.name, g.slug
		FROM documentation d
		JOIN games g ON d.game_id = g.id
		WHERE d.author_id = ?
		ORDER BY d.created_at DESC
		LIMIT ? OFFSET ?
	`, authorID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get documentation: %w", err)
	}
	defer rows.Close()

	var docs []models.Documentation
	for rows.Next() {
		var doc models.Documentation
		var gameName, gameSlug string

		err := rows.Scan(
			&doc.ID, &doc.GameID, &doc.Title, &doc.Slug, &doc.Content,
			&doc.AuthorID, &doc.CreatedAt, &doc.UpdatedAt,
			&gameName, &gameSlug,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan documentation: %w", err)
		}

		doc.Game = &models.Game{
			ID:   doc.GameID,
			Name: gameName,
			Slug: gameSlug,
		}

		docs = append(docs, doc)
	}

	return docs, total, nil
}

func (s *DocumentationService) Search(gameID int, query string, page, perPage int) ([]models.Documentation, int64, error) {
	offset := (page - 1) * perPage
	searchTerm := "%" + query + "%"

	var total int64
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM documentation
		WHERE game_id = ? AND (title LIKE ? OR content LIKE ?)
	`, gameID, searchTerm, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT d.id, d.game_id, d.title, d.slug, d.content, d.author_id,
		       d.created_at, d.updated_at, u.username, u.display_name
		FROM documentation d
		JOIN users u ON d.author_id = u.id
		WHERE d.game_id = ? AND (d.title LIKE ? OR d.content LIKE ?)
		ORDER BY d.updated_at DESC
		LIMIT ? OFFSET ?
	`, gameID, searchTerm, searchTerm, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search documentation: %w", err)
	}
	defer rows.Close()

	var docs []models.Documentation
	for rows.Next() {
		var doc models.Documentation
		var authorUsername, authorDisplayName string

		err := rows.Scan(
			&doc.ID, &doc.GameID, &doc.Title, &doc.Slug, &doc.Content,
			&doc.AuthorID, &doc.CreatedAt, &doc.UpdatedAt,
			&authorUsername, &authorDisplayName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan search result: %w", err)
		}

		doc.Author = &models.User{
			ID:          doc.AuthorID,
			Username:    authorUsername,
			DisplayName: authorDisplayName,
		}

		docs = append(docs, doc)
	}

	return docs, total, nil
}

func (s *DocumentationService) CanUserEdit(userID, gameID int, userRoles []models.UserRole, userMainRole string) bool {
	if userMainRole == models.RoleAdmin {
		return true
	}

	for _, role := range userRoles {
		if role.UserID == userID && role.GameID != nil && *role.GameID == gameID &&
			(role.Role == models.RoleWikiMaintainer || role.Role == models.RoleCommunityModerator) {
			return true
		}
	}

	return false
}

func (s *DocumentationService) GetStats(gameID int) (map[string]int, error) {
	stats := make(map[string]int)

	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM documentation WHERE game_id = ?", gameID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get documentation stats: %w", err)
	}
	stats["total"] = total

	var contributors int
	err = s.db.QueryRow(`
		SELECT COUNT(DISTINCT author_id) FROM documentation WHERE game_id = ?
	`, gameID).Scan(&contributors)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributors count: %w", err)
	}
	stats["contributors"] = contributors

	return stats, nil
}

package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/models"
	"github.com/azurite/backend/internal/utils"
)

type GameService struct {
	db *database.DB
}

func NewGameService(db *database.DB) *GameService {
	return &GameService{
		db: db,
	}
}

func (s *GameService) Create(name, description, icon string) (*models.Game, error) {
	slug := utils.GenerateSlug(name)
	originalSlug := slug
	counter := 1

	for {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM games WHERE slug = ?", slug).Scan(&count)
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
		INSERT INTO games (name, slug, description, icon, is_active)
		VALUES (?, ?, ?, ?, ?)
	`, name, slug, description, icon, true)

	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	gameID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get game ID: %w", err)
	}

	return s.GetByID(int(gameID))
}

func (s *GameService) GetByID(id int) (*models.Game, error) {
	game := &models.Game{}
	err := s.db.QueryRow(`
		SELECT id, name, slug, description, icon, is_active, mod_count, created_at, updated_at
		FROM games WHERE id = ?
	`, id).Scan(
		&game.ID, &game.Name, &game.Slug, &game.Description, &game.Icon,
		&game.IsActive, &game.ModCount, &game.CreatedAt, &game.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("game not found")
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

func (s *GameService) GetBySlug(slug string) (*models.Game, error) {
	game := &models.Game{}
	err := s.db.QueryRow(`
		SELECT id, name, slug, description, icon, is_active, mod_count, created_at, updated_at
		FROM games WHERE slug = ?
	`, slug).Scan(
		&game.ID, &game.Name, &game.Slug, &game.Description, &game.Icon,
		&game.IsActive, &game.ModCount, &game.CreatedAt, &game.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("game not found")
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

func (s *GameService) List(page, perPage int) ([]models.Game, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM games WHERE is_active = 1").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count games: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT id, name, slug, description, icon, is_active, mod_count, created_at, updated_at
		FROM games WHERE is_active = 1
		ORDER BY name ASC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get games: %w", err)
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var game models.Game
		err := rows.Scan(
			&game.ID, &game.Name, &game.Slug, &game.Description, &game.Icon,
			&game.IsActive, &game.ModCount, &game.CreatedAt, &game.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan game: %w", err)
		}
		games = append(games, game)
	}

	return games, total, nil
}

func (s *GameService) Update(gameID int, name, description, icon string) (*models.Game, error) {
	_, err := s.db.Exec(`
		UPDATE games SET name = ?, description = ?, icon = ?
		WHERE id = ?
	`, name, description, icon, gameID)

	if err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return s.GetByID(gameID)
}

func (s *GameService) Delete(gameID int) error {
	var modCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM mods WHERE game_id = ?", gameID).Scan(&modCount)
	if err != nil {
		return fmt.Errorf("failed to check mod count: %w", err)
	}

	if modCount > 0 {
		return errors.New("cannot delete game with existing mods")
	}

	_, err = s.db.Exec("UPDATE games SET is_active = 0 WHERE id = ?", gameID)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	return nil
}

func (s *GameService) CreateRequest(req *models.GameRequestCreate, userID int) (*models.GameRequest, error) {
	result, err := s.db.Exec(`
		INSERT INTO game_requests (name, description, requested_by, status, icon, admin_notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.Name, req.Description, userID, models.GameRequestStatusPending, "", "No admin notes at this time")

	if err != nil {
		return nil, fmt.Errorf("failed to create game request: %w", err)
	}

	requestID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get request ID: %w", err)
	}

	return s.GetRequestByID(int(requestID))
}

func (s *GameService) GetRequestByID(id int) (*models.GameRequest, error) {
	request := &models.GameRequest{}
	var username, displayName string

	err := s.db.QueryRow(`
		SELECT gr.id, gr.name, gr.description, gr.icon, gr.requested_by, gr.status,
		       gr.admin_notes, gr.created_at, gr.updated_at,
		       u.username, u.display_name
		FROM game_requests gr
		JOIN users u ON gr.requested_by = u.id
		WHERE gr.id = ?
	`, id).Scan(
		&request.ID, &request.Name, &request.Description, &request.Icon,
		&request.RequestedBy, &request.Status, &request.AdminNotes,
		&request.CreatedAt, &request.UpdatedAt, &username, &displayName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("game request not found")
		}
		return nil, fmt.Errorf("failed to get game request: %w", err)
	}

	request.User = &models.User{
		ID:          request.RequestedBy,
		Username:    username,
		DisplayName: displayName,
	}

	return request, nil
}

func (s *GameService) ListRequests(page, perPage int, status string) ([]models.GameRequest, int64, error) {
	offset := (page - 1) * perPage

	whereClause := ""
	args := []interface{}{}

	if status != "" {
		whereClause = "WHERE gr.status = ?"
		args = append(args, status)
	}

	var total int64
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM game_requests gr
		JOIN users u ON gr.requested_by = u.id
		%s
	`, whereClause)

	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count game requests: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT gr.id, gr.name, gr.description, gr.icon, gr.requested_by, gr.status,
		       gr.admin_notes, gr.created_at, gr.updated_at,
		       u.username, u.display_name
		FROM game_requests gr
		JOIN users u ON gr.requested_by = u.id
		%s
		ORDER BY gr.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, perPage, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get game requests: %w", err)
	}
	defer rows.Close()

	var requests []models.GameRequest
	for rows.Next() {
		var request models.GameRequest
		var username, displayName string

		err := rows.Scan(
			&request.ID, &request.Name, &request.Description, &request.Icon,
			&request.RequestedBy, &request.Status, &request.AdminNotes,
			&request.CreatedAt, &request.UpdatedAt, &username, &displayName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan game request: %w", err)
		}

		request.User = &models.User{
			ID:          request.RequestedBy,
			Username:    username,
			DisplayName: displayName,
		}

		requests = append(requests, request)
	}

	return requests, total, nil
}

func (s *GameService) ApproveRequest(requestID int, adminNotes string) (*models.Game, error) {
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return nil, err
	}

	if request.Status != models.GameRequestStatusPending {
		return nil, errors.New("request is not pending")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	game, err := s.Create(request.Name, request.Description, request.Icon)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE game_requests SET status = ?, admin_notes = ?
		WHERE id = ?
	`, models.GameRequestStatusApproved, adminNotes, requestID)

	if err != nil {
		return nil, fmt.Errorf("failed to update request status: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return game, nil
}

func (s *GameService) DenyRequest(requestID int, adminNotes string) error {
	request, err := s.GetRequestByID(requestID)
	if err != nil {
		return err
	}

	if request.Status != models.GameRequestStatusPending {
		return errors.New("request is not pending")
	}

	_, err = s.db.Exec(`
		UPDATE game_requests SET status = ?, admin_notes = ?
		WHERE id = ?
	`, models.GameRequestStatusDenied, adminNotes, requestID)

	if err != nil {
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}

func (s *GameService) GetTags(gameID int) ([]models.Tag, error) {
	rows, err := s.db.Query(`
		SELECT id, name, slug, game_id, created_at
		FROM tags WHERE game_id = ?
		ORDER BY name ASC
	`, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.GameID, &tag.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *GameService) CreateTag(gameID int, name string) (*models.Tag, error) {
	slug := utils.GenerateSlug(name)

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM tags WHERE slug = ? AND game_id = ?", slug, gameID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check tag existence: %w", err)
	}

	if count > 0 {
		return nil, errors.New("tag already exists")
	}

	result, err := s.db.Exec(`
		INSERT INTO tags (name, slug, game_id)
		VALUES (?, ?, ?)
	`, name, slug, gameID)

	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	tagID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get tag ID: %w", err)
	}

	return &models.Tag{
		ID:     int(tagID),
		Name:   name,
		Slug:   slug,
		GameID: gameID,
	}, nil
}

func (s *GameService) DeleteTag(tagID int) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM mod_tags WHERE tag_id = ?", tagID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check tag usage: %w", err)
	}

	if count > 0 {
		return errors.New("cannot delete tag that is in use")
	}

	_, err = s.db.Exec("DELETE FROM tags WHERE id = ?", tagID)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (s *GameService) AssignModerator(gameID, userID int) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM user_roles WHERE user_id = ? AND game_id = ? AND role = ?",
		userID, gameID, models.RoleCommunityModerator).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	if count > 0 {
		return errors.New("user is already a moderator for this game")
	}

	_, err = s.db.Exec(`
		INSERT INTO user_roles (user_id, game_id, role)
		VALUES (?, ?, ?)
	`, userID, gameID, models.RoleCommunityModerator)

	if err != nil {
		return fmt.Errorf("failed to assign moderator: %w", err)
	}

	return nil
}

func (s *GameService) RemoveModerator(gameID, userID int) error {
	result, err := s.db.Exec(`
		DELETE FROM user_roles
		WHERE user_id = ? AND game_id = ? AND role = ?
	`, userID, gameID, models.RoleCommunityModerator)

	if err != nil {
		return fmt.Errorf("failed to remove moderator: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user is not a moderator for this game")
	}

	return nil
}

func (s *GameService) GetModerators(gameID int) ([]models.User, error) {
	rows, err := s.db.Query(`
		SELECT u.id, u.username, u.display_name, u.email, u.avatar, u.created_at
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.game_id = ? AND ur.role = ?
		ORDER BY u.display_name ASC
	`, gameID, models.RoleCommunityModerator)
	if err != nil {
		return nil, fmt.Errorf("failed to get moderators: %w", err)
	}
	defer rows.Close()

	var moderators []models.User
	for rows.Next() {
		var mod models.User
		err := rows.Scan(&mod.ID, &mod.Username, &mod.DisplayName, &mod.Email, &mod.Avatar, &mod.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan moderator: %w", err)
		}
		moderators = append(moderators, mod)
	}

	return moderators, nil
}

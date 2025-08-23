package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/azurite/backend/internal/database"
	"github.com/azurite/backend/internal/models"
)

type CommentService struct {
	db *database.DB
}

func NewCommentService(db *database.DB) *CommentService {
	return &CommentService{
		db: db,
	}
}

func (s *CommentService) Create(req *models.CommentCreateRequest, modID, userID int) (*models.Comment, error) {
	if req.ParentID != nil {
		var parentModID int
		err := s.db.QueryRow("SELECT mod_id FROM comments WHERE id = ?", *req.ParentID).Scan(&parentModID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("parent comment not found")
			}
			return nil, fmt.Errorf("failed to verify parent comment: %w", err)
		}

		if parentModID != modID {
			return nil, errors.New("parent comment belongs to different mod")
		}
	}

	result, err := s.db.Exec(`
		INSERT INTO comments (mod_id, user_id, content, parent_id)
		VALUES (?, ?, ?, ?)
	`, modID, userID, req.Content, req.ParentID)

	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment ID: %w", err)
	}

	return s.GetByID(int(commentID))
}

func (s *CommentService) GetByID(id int) (*models.Comment, error) {
	comment := &models.Comment{}
	var username, displayName, avatar string

	err := s.db.QueryRow(`
		SELECT c.id, c.mod_id, c.user_id, c.content, c.parent_id, c.is_active,
		       c.created_at, c.updated_at, u.username, u.display_name, u.avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`, id).Scan(
		&comment.ID, &comment.ModID, &comment.UserID, &comment.Content,
		&comment.ParentID, &comment.IsActive, &comment.CreatedAt, &comment.UpdatedAt,
		&username, &displayName, &avatar,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("comment not found")
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	comment.User = &models.User{
		ID:          comment.UserID,
		Username:    username,
		DisplayName: displayName,
		Avatar:      avatar,
	}

	return comment, nil
}

func (s *CommentService) ListByMod(modID int, page, perPage int) ([]models.Comment, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM comments
		WHERE mod_id = ? AND is_active = 1 AND parent_id IS NULL
	`, modID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT c.id, c.mod_id, c.user_id, c.content, c.parent_id, c.is_active,
		       c.created_at, c.updated_at, u.username, u.display_name, u.avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.mod_id = ? AND c.is_active = 1 AND c.parent_id IS NULL
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?
	`, modID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var username, displayName, avatar string

		err := rows.Scan(
			&comment.ID, &comment.ModID, &comment.UserID, &comment.Content,
			&comment.ParentID, &comment.IsActive, &comment.CreatedAt, &comment.UpdatedAt,
			&username, &displayName, &avatar,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment: %w", err)
		}

		comment.User = &models.User{
			ID:          comment.UserID,
			Username:    username,
			DisplayName: displayName,
			Avatar:      avatar,
		}

		replies, err := s.getReplies(comment.ID)
		if err == nil {
			comment.Replies = replies
		}

		comments = append(comments, comment)
	}

	return comments, total, nil
}

func (s *CommentService) Update(commentID int, content string, userID int) (*models.Comment, error) {
	comment, err := s.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	if comment.UserID != userID {
		return nil, errors.New("unauthorized to update this comment")
	}

	_, err = s.db.Exec(`
		UPDATE comments SET content = ?
		WHERE id = ?
	`, content, commentID)

	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return s.GetByID(commentID)
}

func (s *CommentService) Delete(commentID int, userID int, isAdmin bool) error {
	comment, err := s.GetByID(commentID)
	if err != nil {
		return err
	}

	if !isAdmin && comment.UserID != userID {
		return errors.New("unauthorized to delete this comment")
	}

	_, err = s.db.Exec(`
		UPDATE comments SET is_active = 0
		WHERE id = ?
	`, commentID)

	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (s *CommentService) Moderate(commentID int, isActive bool) error {
	_, err := s.db.Exec(`
		UPDATE comments SET is_active = ?
		WHERE id = ?
	`, isActive, commentID)

	if err != nil {
		return fmt.Errorf("failed to moderate comment: %w", err)
	}

	return nil
}

func (s *CommentService) GetByUser(userID int, page, perPage int) ([]models.Comment, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = ? AND is_active = 1", userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT c.id, c.mod_id, c.user_id, c.content, c.parent_id, c.is_active,
		       c.created_at, c.updated_at, m.name, m.slug, g.slug
		FROM comments c
		JOIN mods m ON c.mod_id = m.id
		JOIN games g ON m.game_id = g.id
		WHERE c.user_id = ? AND c.is_active = 1
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?
	`, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var modName, modSlug, gameSlug string

		err := rows.Scan(
			&comment.ID, &comment.ModID, &comment.UserID, &comment.Content,
			&comment.ParentID, &comment.IsActive, &comment.CreatedAt, &comment.UpdatedAt,
			&modName, &modSlug, &gameSlug,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment: %w", err)
		}

		comments = append(comments, comment)
	}

	return comments, total, nil
}

func (s *CommentService) getReplies(parentID int) ([]models.Comment, error) {
	rows, err := s.db.Query(`
		SELECT c.id, c.mod_id, c.user_id, c.content, c.parent_id, c.is_active,
		       c.created_at, c.updated_at, u.username, u.display_name, u.avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.parent_id = ? AND c.is_active = 1
		ORDER BY c.created_at ASC
	`, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}
	defer rows.Close()

	var replies []models.Comment
	for rows.Next() {
		var reply models.Comment
		var username, displayName, avatar string

		err := rows.Scan(
			&reply.ID, &reply.ModID, &reply.UserID, &reply.Content,
			&reply.ParentID, &reply.IsActive, &reply.CreatedAt, &reply.UpdatedAt,
			&username, &displayName, &avatar,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reply: %w", err)
		}

		reply.User = &models.User{
			ID:          reply.UserID,
			Username:    username,
			DisplayName: displayName,
			Avatar:      avatar,
		}

		replies = append(replies, reply)
	}

	return replies, nil
}

func (s *CommentService) GetStats(modID int) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM comments WHERE mod_id = ? AND is_active = 1", modID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get comment count: %w", err)
	}

	return count, nil
}

func (s *CommentService) ListForModeration(page, perPage int) ([]models.Comment, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	err := s.db.QueryRow("SELECT COUNT(*) FROM comments WHERE is_active = 0").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT c.id, c.mod_id, c.user_id, c.content, c.parent_id, c.is_active,
		       c.created_at, c.updated_at, u.username, u.display_name, u.avatar,
		       m.name, m.slug, g.slug
		FROM comments c
		JOIN users u ON c.user_id = u.id
		JOIN mods m ON c.mod_id = m.id
		JOIN games g ON m.game_id = g.id
		WHERE c.is_active = 0
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?
	`, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var username, displayName, avatar, modName, modSlug, gameSlug string

		err := rows.Scan(
			&comment.ID, &comment.ModID, &comment.UserID, &comment.Content,
			&comment.ParentID, &comment.IsActive, &comment.CreatedAt, &comment.UpdatedAt,
			&username, &displayName, &avatar, &modName, &modSlug, &gameSlug,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment: %w", err)
		}

		comment.User = &models.User{
			ID:          comment.UserID,
			Username:    username,
			DisplayName: displayName,
			Avatar:      avatar,
		}

		comments = append(comments, comment)
	}

	return comments, total, nil
}

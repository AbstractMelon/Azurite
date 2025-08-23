package services

import (
	"fmt"
	"testing"

	"github.com/azurite/backend/internal/models"
)

func TestCommentService_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	req := &models.CommentCreateRequest{
		Content:  "This is a test comment",
		ParentID: nil,
	}

	comment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if comment.Content != req.Content {
		t.Errorf("Expected content %s, got %s", req.Content, comment.Content)
	}
	if comment.ModID != modID {
		t.Errorf("Expected mod ID %d, got %d", modID, comment.ModID)
	}
	if comment.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, comment.UserID)
	}
	if comment.ParentID != nil {
		t.Errorf("Expected nil parent ID, got %v", comment.ParentID)
	}
	if !comment.IsActive {
		t.Error("Expected comment to be active")
	}
	if comment.User == nil {
		t.Error("Expected user to be populated")
	}
}

func TestCommentService_CreateReply(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create parent comment
	parentReq := &models.CommentCreateRequest{
		Content:  "Parent comment",
		ParentID: nil,
	}

	parentComment, err := service.Create(parentReq, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create parent comment: %v", err)
	}

	// Create reply
	replyReq := &models.CommentCreateRequest{
		Content:  "This is a reply",
		ParentID: &parentComment.ID,
	}

	replyComment, err := service.Create(replyReq, modID, userID)
	if err != nil {
		t.Fatalf("Create() reply error = %v", err)
	}

	if replyComment.Content != replyReq.Content {
		t.Errorf("Expected content %s, got %s", replyReq.Content, replyComment.Content)
	}
	if replyComment.ParentID == nil {
		t.Error("Expected parent ID to be set")
	}
	if *replyComment.ParentID != parentComment.ID {
		t.Errorf("Expected parent ID %d, got %d", parentComment.ID, *replyComment.ParentID)
	}
}

func TestCommentService_CreateReplyInvalidParent(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	req := &models.CommentCreateRequest{
		Content:  "This is a reply to non-existent comment",
		ParentID: &[]int{999}[0], // Non-existent parent
	}

	_, err := service.Create(req, modID, userID)
	if err == nil {
		t.Error("Expected error for invalid parent comment")
	}
}

func TestCommentService_CreateReplyDifferentMod(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID1 := createTestMod(t, db, "Test Mod 1", "test-mod-1", gameID, userID)
	modID2 := createTestMod(t, db, "Test Mod 2", "test-mod-2", gameID, userID)

	// Create comment on first mod
	parentReq := &models.CommentCreateRequest{
		Content:  "Comment on mod 1",
		ParentID: nil,
	}

	parentComment, err := service.Create(parentReq, modID1, userID)
	if err != nil {
		t.Fatalf("Failed to create parent comment: %v", err)
	}

	// Try to create reply on different mod
	replyReq := &models.CommentCreateRequest{
		Content:  "Reply on different mod",
		ParentID: &parentComment.ID,
	}

	_, err = service.Create(replyReq, modID2, userID)
	if err == nil {
		t.Error("Expected error when replying to comment from different mod")
	}
}

func TestCommentService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comment
	req := &models.CommentCreateRequest{
		Content:  "Test comment",
		ParentID: nil,
	}

	createdComment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Get comment by ID
	comment, err := service.GetByID(createdComment.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if comment.ID != createdComment.ID {
		t.Errorf("Expected ID %d, got %d", createdComment.ID, comment.ID)
	}
	if comment.Content != req.Content {
		t.Errorf("Expected content %s, got %s", req.Content, comment.Content)
	}
	if comment.User == nil {
		t.Error("Expected user to be populated")
	}
}

func TestCommentService_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)

	_, err := service.GetByID(999)
	if err == nil {
		t.Error("Expected error for non-existent comment")
	}
}

func TestCommentService_ListByMod(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create test comments
	for i := 1; i <= 3; i++ {
		req := &models.CommentCreateRequest{
			Content:  fmt.Sprintf("Test comment %d", i),
			ParentID: nil,
		}
		_, err := service.Create(req, modID, userID)
		if err != nil {
			t.Fatalf("Failed to create comment %d: %v", i, err)
		}
	}

	comments, total, err := service.ListByMod(modID, 1, 10)
	if err != nil {
		t.Fatalf("ListByMod() error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(comments) != 3 {
		t.Errorf("Expected 3 comments, got %d", len(comments))
	}

	// Test pagination
	comments, total, err = service.ListByMod(modID, 1, 2)
	if err != nil {
		t.Fatalf("ListByMod() pagination error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(comments) != 2 {
		t.Errorf("Expected 2 comments per page, got %d", len(comments))
	}
}

func TestCommentService_ListByModWithReplies(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create parent comment
	parentReq := &models.CommentCreateRequest{
		Content:  "Parent comment",
		ParentID: nil,
	}

	parentComment, err := service.Create(parentReq, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create parent comment: %v", err)
	}

	// Create replies
	for i := 1; i <= 2; i++ {
		replyReq := &models.CommentCreateRequest{
			Content:  fmt.Sprintf("Reply %d", i),
			ParentID: &parentComment.ID,
		}
		_, err := service.Create(replyReq, modID, userID)
		if err != nil {
			t.Fatalf("Failed to create reply %d: %v", i, err)
		}
	}

	comments, total, err := service.ListByMod(modID, 1, 10)
	if err != nil {
		t.Fatalf("ListByMod() error = %v", err)
	}

	if total != 1 {
		t.Errorf("Expected total 1 top-level comment, got %d", total)
	}
	if len(comments) != 1 {
		t.Errorf("Expected 1 top-level comment, got %d", len(comments))
	}

	comment := comments[0]
	if len(comment.Replies) != 2 {
		t.Errorf("Expected 2 replies, got %d", len(comment.Replies))
	}
}

func TestCommentService_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comment
	req := &models.CommentCreateRequest{
		Content:  "Original content",
		ParentID: nil,
	}

	comment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Update comment
	updatedContent := "Updated content"
	updatedComment, err := service.Update(comment.ID, updatedContent, userID)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updatedComment.Content != updatedContent {
		t.Errorf("Expected content %s, got %s", updatedContent, updatedComment.Content)
	}
}

func TestCommentService_UpdateUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comment with first user
	req := &models.CommentCreateRequest{
		Content:  "Original content",
		ParentID: nil,
	}

	comment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Try to update with different user
	_, err = service.Update(comment.ID, "Updated content", otherUserID)
	if err == nil {
		t.Error("Expected error when updating comment by different user")
	}
}

func TestCommentService_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comment
	req := &models.CommentCreateRequest{
		Content:  "Test comment",
		ParentID: nil,
	}

	comment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Delete comment
	err = service.Delete(comment.ID, userID, false)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify comment is inactive
	deletedComment, err := service.GetByID(comment.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if deletedComment.IsActive {
		t.Error("Expected comment to be inactive after deletion")
	}
}

func TestCommentService_DeleteAsAdmin(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comment with regular user
	req := &models.CommentCreateRequest{
		Content:  "Test comment",
		ParentID: nil,
	}

	comment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Delete comment as admin
	err = service.Delete(comment.ID, adminID, true) // isAdmin = true
	if err != nil {
		t.Fatalf("Delete() as admin error = %v", err)
	}

	// Verify comment is inactive
	deletedComment, err := service.GetByID(comment.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if deletedComment.IsActive {
		t.Error("Expected comment to be inactive after admin deletion")
	}
}

func TestCommentService_DeleteUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comment with first user
	req := &models.CommentCreateRequest{
		Content:  "Test comment",
		ParentID: nil,
	}

	comment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Try to delete with different user (not admin)
	err = service.Delete(comment.ID, otherUserID, false)
	if err == nil {
		t.Error("Expected error when deleting comment by different user")
	}
}

func TestCommentService_Moderate(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comment
	req := &models.CommentCreateRequest{
		Content:  "Test comment",
		ParentID: nil,
	}

	comment, err := service.Create(req, modID, userID)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	// Moderate comment (deactivate)
	err = service.Moderate(comment.ID, false)
	if err != nil {
		t.Fatalf("Moderate() error = %v", err)
	}

	// Verify comment is inactive
	moderatedComment, err := service.GetByID(comment.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if moderatedComment.IsActive {
		t.Error("Expected comment to be inactive after moderation")
	}

	// Moderate comment (reactivate)
	err = service.Moderate(comment.ID, true)
	if err != nil {
		t.Fatalf("Moderate() reactivate error = %v", err)
	}

	// Verify comment is active again
	reactivatedComment, err := service.GetByID(comment.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if !reactivatedComment.IsActive {
		t.Error("Expected comment to be active after reactivation")
	}
}

func TestCommentService_GetByUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comments for user
	for i := 1; i <= 3; i++ {
		req := &models.CommentCreateRequest{
			Content:  fmt.Sprintf("Comment %d", i),
			ParentID: nil,
		}
		_, err := service.Create(req, modID, userID)
		if err != nil {
			t.Fatalf("Failed to create comment %d: %v", i, err)
		}
	}

	comments, total, err := service.GetByUser(userID, 1, 10)
	if err != nil {
		t.Fatalf("GetByUser() error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(comments) != 3 {
		t.Errorf("Expected 3 comments, got %d", len(comments))
	}

	// All comments should belong to the user
	for _, comment := range comments {
		if comment.UserID != userID {
			t.Errorf("Expected comment to belong to user %d, got %d", userID, comment.UserID)
		}
	}
}

func TestCommentService_GetStats(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Initially no comments
	count, err := service.GetStats(modID)
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 comments initially, got %d", count)
	}

	// Create some comments
	for i := 1; i <= 5; i++ {
		req := &models.CommentCreateRequest{
			Content:  fmt.Sprintf("Comment %d", i),
			ParentID: nil,
		}
		_, err := service.Create(req, modID, userID)
		if err != nil {
			t.Fatalf("Failed to create comment %d: %v", i, err)
		}
	}

	count, err = service.GetStats(modID)
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if count != 5 {
		t.Errorf("Expected 5 comments, got %d", count)
	}

	// Delete one comment and verify count updates
	comments, _, err := service.ListByMod(modID, 1, 10)
	if err != nil {
		t.Fatalf("ListByMod() error = %v", err)
	}

	err = service.Delete(comments[0].ID, userID, false)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	count, err = service.GetStats(modID)
	if err != nil {
		t.Fatalf("GetStats() after delete error = %v", err)
	}

	if count != 4 {
		t.Errorf("Expected 4 comments after deletion, got %d", count)
	}
}

func TestCommentService_ListForModeration(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewCommentService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create comments
	for i := 1; i <= 3; i++ {
		req := &models.CommentCreateRequest{
			Content:  fmt.Sprintf("Comment %d", i),
			ParentID: nil,
		}
		_, err := service.Create(req, modID, userID)
		if err != nil {
			t.Fatalf("Failed to create comment %d: %v", i, err)
		}
	}

	// Initially no moderated comments
	comments, total, err := service.ListForModeration(1, 10)
	if err != nil {
		t.Fatalf("ListForModeration() error = %v", err)
	}

	if total != 0 {
		t.Errorf("Expected 0 moderated comments initially, got %d", total)
	}

	// Moderate some comments
	allComments, _, err := service.ListByMod(modID, 1, 10)
	if err != nil {
		t.Fatalf("ListByMod() error = %v", err)
	}

	for i := 0; i < 2; i++ {
		err = service.Moderate(allComments[i].ID, false)
		if err != nil {
			t.Fatalf("Moderate() error = %v", err)
		}
	}

	// Now should have moderated comments
	comments, total, err = service.ListForModeration(1, 10)
	if err != nil {
		t.Fatalf("ListForModeration() after moderation error = %v", err)
	}

	if total != 2 {
		t.Errorf("Expected 2 moderated comments, got %d", total)
	}
	if len(comments) != 2 {
		t.Errorf("Expected 2 comments in list, got %d", len(comments))
	}

	// All returned comments should be inactive
	for _, comment := range comments {
		if comment.IsActive {
			t.Error("Expected moderated comments to be inactive")
		}
	}
}

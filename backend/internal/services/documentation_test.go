package services

import (
	"testing"

	"github.com/azurite/backend/internal/models"
)

func TestDocumentationService_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Getting Started Guide",
		Content: "This is a comprehensive getting started guide...",
	}

	doc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if doc.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, doc.Title)
	}
	if doc.Content != req.Content {
		t.Errorf("Expected content %s, got %s", req.Content, doc.Content)
	}
	if doc.GameID != gameID {
		t.Errorf("Expected game ID %d, got %d", gameID, doc.GameID)
	}
	if doc.AuthorID != userID {
		t.Errorf("Expected author ID %d, got %d", userID, doc.AuthorID)
	}
	if doc.Slug == "" {
		t.Error("Expected non-empty slug")
	}
	if doc.Author == nil {
		t.Error("Expected author to be populated")
	}
	if doc.Game == nil {
		t.Error("Expected game to be populated")
	}
}

func TestDocumentationService_CreateDuplicateTitle(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Duplicate Title",
		Content: "First document",
	}

	// Create first document
	_, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("First Create() error = %v", err)
	}

	// Create second document with same title (should get different slug)
	req.Content = "Second document"
	doc2, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Second Create() error = %v", err)
	}

	// Should have different slug even with same title
	if doc2.Slug == "duplicate-title" {
		t.Error("Expected second document to have modified slug to avoid duplicates")
	}
}

func TestDocumentationService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Test Documentation",
		Content: "Test content",
	}

	createdDoc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	doc, err := service.GetByID(createdDoc.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if doc.ID != createdDoc.ID {
		t.Errorf("Expected ID %d, got %d", createdDoc.ID, doc.ID)
	}
	if doc.Title != req.Title {
		t.Errorf("Expected title %s, got %s", req.Title, doc.Title)
	}
	if doc.Author == nil {
		t.Error("Expected author to be populated")
	}
	if doc.Game == nil {
		t.Error("Expected game to be populated")
	}
}

func TestDocumentationService_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)

	_, err := service.GetByID(999)
	if err == nil {
		t.Error("Expected error for non-existent documentation")
	}
}

func TestDocumentationService_GetBySlug(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Test Documentation",
		Content: "Test content",
	}

	createdDoc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	doc, err := service.GetBySlug("test-game", createdDoc.Slug)
	if err != nil {
		t.Fatalf("GetBySlug() error = %v", err)
	}

	if doc.ID != createdDoc.ID {
		t.Errorf("Expected ID %d, got %d", createdDoc.ID, doc.ID)
	}
	if doc.Slug != createdDoc.Slug {
		t.Errorf("Expected slug %s, got %s", createdDoc.Slug, doc.Slug)
	}
}

func TestDocumentationService_GetBySlugNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)

	_, err := service.GetBySlug("non-existent-game", "non-existent-doc")
	if err == nil {
		t.Error("Expected error for non-existent documentation")
	}
}

func TestDocumentationService_ListByGame(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create test documentation
	for i := 1; i <= 5; i++ {
		req := &models.DocumentationCreateRequest{
			Title:   "Doc " + string(rune('0'+i)),
			Content: "Content for doc " + string(rune('0'+i)),
		}
		_, err := service.Create(req, gameID, userID)
		if err != nil {
			t.Fatalf("Failed to create doc %d: %v", i, err)
		}
	}

	docs, total, err := service.ListByGame(gameID, 1, 3)
	if err != nil {
		t.Fatalf("ListByGame() error = %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(docs) != 3 {
		t.Errorf("Expected 3 docs per page, got %d", len(docs))
	}

	// Test second page
	docs, total, err = service.ListByGame(gameID, 2, 3)
	if err != nil {
		t.Fatalf("ListByGame() page 2 error = %v", err)
	}

	if total != 5 {
		t.Errorf("Expected total 5, got %d", total)
	}
	if len(docs) != 2 {
		t.Errorf("Expected 2 docs on page 2, got %d", len(docs))
	}

	// All docs should have author populated
	for _, doc := range docs {
		if doc.Author == nil {
			t.Error("Expected author to be populated")
		}
	}
}

func TestDocumentationService_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Original Title",
		Content: "Original content",
	}

	doc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	updateReq := &models.DocumentationCreateRequest{
		Title:   "Updated Title",
		Content: "Updated content",
	}

	updatedDoc, err := service.Update(doc.ID, updateReq, userID, false)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updatedDoc.Title != updateReq.Title {
		t.Errorf("Expected title %s, got %s", updateReq.Title, updatedDoc.Title)
	}
	if updatedDoc.Content != updateReq.Content {
		t.Errorf("Expected content %s, got %s", updateReq.Content, updatedDoc.Content)
	}
}

func TestDocumentationService_UpdateUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Original Title",
		Content: "Original content",
	}

	doc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	updateReq := &models.DocumentationCreateRequest{
		Title:   "Updated Title",
		Content: "Updated content",
	}

	// Try to update with different user
	_, err = service.Update(doc.ID, updateReq, otherUserID, false)
	if err == nil {
		t.Error("Expected error when updating documentation by different user")
	}
}

func TestDocumentationService_UpdateAsEditor(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	editorID := createTestUser(t, db, "editor", "editor@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Original Title",
		Content: "Original content",
	}

	doc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	updateReq := &models.DocumentationCreateRequest{
		Title:   "Updated Title",
		Content: "Updated content",
	}

	// Update with editor privileges (canEdit = true)
	updatedDoc, err := service.Update(doc.ID, updateReq, editorID, true)
	if err != nil {
		t.Fatalf("Update() as editor error = %v", err)
	}

	if updatedDoc.Title != updateReq.Title {
		t.Errorf("Expected title %s, got %s", updateReq.Title, updatedDoc.Title)
	}
}

func TestDocumentationService_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Test Documentation",
		Content: "Test content",
	}

	doc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	err = service.Delete(doc.ID, userID, false)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify documentation was deleted
	_, err = service.GetByID(doc.ID)
	if err == nil {
		t.Error("Expected error when getting deleted documentation")
	}
}

func TestDocumentationService_DeleteUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Test Documentation",
		Content: "Test content",
	}

	doc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	// Try to delete with different user
	err = service.Delete(doc.ID, otherUserID, false)
	if err == nil {
		t.Error("Expected error when deleting documentation by different user")
	}
}

func TestDocumentationService_DeleteWithPermission(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	moderatorID := createTestUser(t, db, "moderator", "mod@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	req := &models.DocumentationCreateRequest{
		Title:   "Test Documentation",
		Content: "Test content",
	}

	doc, err := service.Create(req, gameID, userID)
	if err != nil {
		t.Fatalf("Failed to create documentation: %v", err)
	}

	// Delete with permission (canDelete = true)
	err = service.Delete(doc.ID, moderatorID, true)
	if err != nil {
		t.Fatalf("Delete() with permission error = %v", err)
	}

	// Verify documentation was deleted
	_, err = service.GetByID(doc.ID)
	if err == nil {
		t.Error("Expected error when getting deleted documentation")
	}
}

func TestDocumentationService_GetByAuthor(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create documentation for user
	for i := 1; i <= 3; i++ {
		req := &models.DocumentationCreateRequest{
			Title:   "User Doc " + string(rune('0'+i)),
			Content: "Content " + string(rune('0'+i)),
		}
		_, err := service.Create(req, gameID, userID)
		if err != nil {
			t.Fatalf("Failed to create doc %d: %v", i, err)
		}
	}

	docs, total, err := service.GetByAuthor(userID, 1, 10)
	if err != nil {
		t.Fatalf("GetByAuthor() error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(docs) != 3 {
		t.Errorf("Expected 3 docs, got %d", len(docs))
	}

	// All docs should belong to the user and have game info
	for _, doc := range docs {
		if doc.AuthorID != userID {
			t.Errorf("Expected doc to belong to user %d, got %d", userID, doc.AuthorID)
		}
		if doc.Game == nil {
			t.Error("Expected game to be populated")
		}
	}
}

func TestDocumentationService_Search(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create test documentation
	docs := []struct {
		title   string
		content string
	}{
		{"Installation Guide", "How to install mods for this game"},
		{"Configuration Manual", "Configure your game settings"},
		{"Troubleshooting", "Fix common installation problems"},
		{"Advanced Tips", "Tips for advanced users"},
	}

	for _, doc := range docs {
		req := &models.DocumentationCreateRequest{
			Title:   doc.title,
			Content: doc.content,
		}
		_, err := service.Create(req, gameID, userID)
		if err != nil {
			t.Fatalf("Failed to create doc: %v", err)
		}
	}

	// Search for "installation"
	results, total, err := service.Search(gameID, "installation", 1, 10)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if total != 2 {
		t.Errorf("Expected 2 results for 'installation', got %d", total)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 docs in results, got %d", len(results))
	}

	// Search for "configuration"
	results, total, err = service.Search(gameID, "configuration", 1, 10)
	if err != nil {
		t.Fatalf("Search() for configuration error = %v", err)
	}

	if total != 1 {
		t.Errorf("Expected 1 result for 'configuration', got %d", total)
	}

	// Search for non-existent term
	results, total, err = service.Search(gameID, "nonexistent", 1, 10)
	if err != nil {
		t.Fatalf("Search() for nonexistent error = %v", err)
	}

	if total != 0 {
		t.Errorf("Expected 0 results for 'nonexistent', got %d", total)
	}
}

func TestDocumentationService_CanUserEdit(t *testing.T) {
	service := NewDocumentationService(nil) // Don't need DB for this test

	userRoles := []models.UserRole{
		{UserID: 1, GameID: &[]int{1}[0], Role: models.RoleWikiMaintainer},
		{UserID: 1, GameID: &[]int{2}[0], Role: models.RoleCommunityModerator},
	}

	tests := []struct {
		name         string
		userID       int
		gameID       int
		userMainRole string
		expected     bool
	}{
		{
			name:         "Admin can edit any game",
			userID:       1,
			gameID:       1,
			userMainRole: models.RoleAdmin,
			expected:     true,
		},
		{
			name:         "Wiki maintainer can edit assigned game",
			userID:       1,
			gameID:       1,
			userMainRole: models.RoleUser,
			expected:     true,
		},
		{
			name:         "Community moderator can edit assigned game",
			userID:       1,
			gameID:       2,
			userMainRole: models.RoleUser,
			expected:     true,
		},
		{
			name:         "User cannot edit unassigned game",
			userID:       1,
			gameID:       3,
			userMainRole: models.RoleUser,
			expected:     false,
		},
		{
			name:         "Regular user cannot edit",
			userID:       2,
			gameID:       1,
			userMainRole: models.RoleUser,
			expected:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := service.CanUserEdit(test.userID, test.gameID, userRoles, test.userMainRole)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestDocumentationService_GetStats(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewDocumentationService(db)
	userID1 := createTestUser(t, db, "user1", "user1@example.com")
	userID2 := createTestUser(t, db, "user2", "user2@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Initially no documentation
	stats, err := service.GetStats(gameID)
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if stats["total"] != 0 {
		t.Errorf("Expected 0 total docs initially, got %d", stats["total"])
	}
	if stats["contributors"] != 0 {
		t.Errorf("Expected 0 contributors initially, got %d", stats["contributors"])
	}

	// Create documentation by multiple users
	for i := 1; i <= 3; i++ {
		req := &models.DocumentationCreateRequest{
			Title:   "Doc by User 1 - " + string(rune('0'+i)),
			Content: "Content",
		}
		_, err := service.Create(req, gameID, userID1)
		if err != nil {
			t.Fatalf("Failed to create doc %d: %v", i, err)
		}
	}

	// Create documentation by second user
	req := &models.DocumentationCreateRequest{
		Title:   "Doc by User 2",
		Content: "Content",
	}
	_, err = service.Create(req, gameID, userID2)
	if err != nil {
		t.Fatalf("Failed to create doc by user 2: %v", err)
	}

	stats, err = service.GetStats(gameID)
	if err != nil {
		t.Fatalf("GetStats() after creation error = %v", err)
	}

	if stats["total"] != 4 {
		t.Errorf("Expected 4 total docs, got %d", stats["total"])
	}
	if stats["contributors"] != 2 {
		t.Errorf("Expected 2 contributors, got %d", stats["contributors"])
	}
}

package services

import (
	"testing"

	"github.com/azurite/backend/internal/models"
)

func TestGameService_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)

	tests := []struct {
		name        string
		gameName    string
		description string
		icon        string
		wantErr     bool
	}{
		{
			name:        "Valid game creation",
			gameName:    "Test Game",
			description: "A test game",
			icon:        "icon.png",
			wantErr:     false,
		},
		{
			name:        "Game with empty description",
			gameName:    "Another Game",
			description: "",
			icon:        "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := service.Create(tt.gameName, tt.description, tt.icon)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if game.Name != tt.gameName {
					t.Errorf("Expected game name %s, got %s", tt.gameName, game.Name)
				}
				if game.Description != tt.description {
					t.Errorf("Expected description %s, got %s", tt.description, game.Description)
				}
				if game.Icon != tt.icon {
					t.Errorf("Expected icon %s, got %s", tt.icon, game.Icon)
				}
				if !game.IsActive {
					t.Error("Expected game to be active")
				}
				if game.Slug == "" {
					t.Error("Expected non-empty slug")
				}
			}
		})
	}
}

func TestGameService_CreateDuplicateSlug(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)

	// Create first game
	game1, err := service.Create("Test Game", "Description", "")
	if err != nil {
		t.Fatalf("Failed to create first game: %v", err)
	}

	// Create second game with same name (should get different slug)
	game2, err := service.Create("Test Game", "Description", "")
	if err != nil {
		t.Fatalf("Failed to create second game: %v", err)
	}

	if game1.Slug == game2.Slug {
		t.Error("Expected different slugs for games with same name")
	}
}

func TestGameService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")

	game, err := service.GetByID(gameID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if game.ID != gameID {
		t.Errorf("Expected ID %d, got %d", gameID, game.ID)
	}
	if game.Name != "Test Game" {
		t.Errorf("Expected name 'Test Game', got %s", game.Name)
	}
}

func TestGameService_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)

	_, err := service.GetByID(999)
	if err == nil {
		t.Error("Expected error for non-existent game")
	}
}

func TestGameService_GetBySlug(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")

	game, err := service.GetBySlug("test-game")
	if err != nil {
		t.Fatalf("GetBySlug() error = %v", err)
	}

	if game.ID != gameID {
		t.Errorf("Expected ID %d, got %d", gameID, game.ID)
	}
	if game.Slug != "test-game" {
		t.Errorf("Expected slug 'test-game', got %s", game.Slug)
	}
}

func TestGameService_List(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)

	// Create test games
	createTestGame(t, db, "Game 1", "game-1")
	createTestGame(t, db, "Game 2", "game-2")
	createTestGame(t, db, "Game 3", "game-3")

	games, total, err := service.List(1, 2)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(games) != 2 {
		t.Errorf("Expected 2 games per page, got %d", len(games))
	}

	// Test second page
	games, total, err = service.List(2, 2)
	if err != nil {
		t.Fatalf("List() page 2 error = %v", err)
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if len(games) != 1 {
		t.Errorf("Expected 1 game on page 2, got %d", len(games))
	}
}

func TestGameService_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Original Game", "original-game")

	updatedGame, err := service.Update(gameID, "Updated Game", "Updated description", "new-icon.png")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updatedGame.Name != "Updated Game" {
		t.Errorf("Expected name 'Updated Game', got %s", updatedGame.Name)
	}
	if updatedGame.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got %s", updatedGame.Description)
	}
	if updatedGame.Icon != "new-icon.png" {
		t.Errorf("Expected icon 'new-icon.png', got %s", updatedGame.Icon)
	}
}

func TestGameService_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")

	err := service.Delete(gameID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Game should be marked as inactive
	game, err := service.GetByID(gameID)
	if err != nil {
		t.Fatalf("GetByID() after delete error = %v", err)
	}
	if game.IsActive {
		t.Error("Expected game to be inactive after delete")
	}
}

func TestGameService_DeleteWithMods(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")
	userID := createTestUser(t, db, "testuser", "test@example.com")
	createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	err := service.Delete(gameID)
	if err == nil {
		t.Error("Expected error when deleting game with mods")
	}
}

func TestGameService_CreateRequest(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	req := &models.GameRequestCreate{
		Name:        "New Game Request",
		Description: "Please add this game",
	}

	gameRequest, err := service.CreateRequest(req, userID)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	if gameRequest.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, gameRequest.Name)
	}
	if gameRequest.Description != req.Description {
		t.Errorf("Expected description %s, got %s", req.Description, gameRequest.Description)
	}
	if gameRequest.RequestedBy != userID {
		t.Errorf("Expected requested_by %d, got %d", userID, gameRequest.RequestedBy)
	}
	if gameRequest.Status != models.GameRequestStatusPending {
		t.Errorf("Expected status %s, got %s", models.GameRequestStatusPending, gameRequest.Status)
	}
}

func TestGameService_ListRequests(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	// Create test requests
	req1 := &models.GameRequestCreate{Name: "Game 1", Description: "Description 1"}
	req2 := &models.GameRequestCreate{Name: "Game 2", Description: "Description 2"}

	service.CreateRequest(req1, userID)
	service.CreateRequest(req2, userID)

	requests, total, err := service.ListRequests(1, 10, "")
	if err != nil {
		t.Fatalf("ListRequests() error = %v", err)
	}

	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
	if len(requests) != 2 {
		t.Errorf("Expected 2 requests, got %d", len(requests))
	}

	// Test filtering by status
	requests, total, err = service.ListRequests(1, 10, models.GameRequestStatusPending)
	if err != nil {
		t.Fatalf("ListRequests() with status filter error = %v", err)
	}

	if total != 2 {
		t.Errorf("Expected total 2 pending requests, got %d", total)
	}
}

func TestGameService_ApproveRequest(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	req := &models.GameRequestCreate{
		Name:        "New Game",
		Description: "Please add this game",
	}

	gameRequest, err := service.CreateRequest(req, userID)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	adminNotes := "Approved by admin"
	game, err := service.ApproveRequest(gameRequest.ID, adminNotes)
	if err != nil {
		t.Fatalf("ApproveRequest() error = %v", err)
	}

	if game.Name != req.Name {
		t.Errorf("Expected game name %s, got %s", req.Name, game.Name)
	}

	// Verify request status was updated
	updatedRequest, err := service.GetRequestByID(gameRequest.ID)
	if err != nil {
		t.Fatalf("GetRequestByID() error = %v", err)
	}

	if updatedRequest.Status != models.GameRequestStatusApproved {
		t.Errorf("Expected status %s, got %s", models.GameRequestStatusApproved, updatedRequest.Status)
	}
	if updatedRequest.AdminNotes != adminNotes {
		t.Errorf("Expected admin notes %s, got %s", adminNotes, updatedRequest.AdminNotes)
	}
}

func TestGameService_DenyRequest(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	userID := createTestUser(t, db, "testuser", "test@example.com")

	req := &models.GameRequestCreate{
		Name:        "New Game",
		Description: "Please add this game",
	}

	gameRequest, err := service.CreateRequest(req, userID)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	adminNotes := "Denied - already exists"
	err = service.DenyRequest(gameRequest.ID, adminNotes)
	if err != nil {
		t.Fatalf("DenyRequest() error = %v", err)
	}

	// Verify request status was updated
	updatedRequest, err := service.GetRequestByID(gameRequest.ID)
	if err != nil {
		t.Fatalf("GetRequestByID() error = %v", err)
	}

	if updatedRequest.Status != models.GameRequestStatusDenied {
		t.Errorf("Expected status %s, got %s", models.GameRequestStatusDenied, updatedRequest.Status)
	}
	if updatedRequest.AdminNotes != adminNotes {
		t.Errorf("Expected admin notes %s, got %s", adminNotes, updatedRequest.AdminNotes)
	}
}

func TestGameService_CreateTag(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")

	tag, err := service.CreateTag(gameID, "Action")
	if err != nil {
		t.Fatalf("CreateTag() error = %v", err)
	}

	if tag.Name != "Action" {
		t.Errorf("Expected tag name 'Action', got %s", tag.Name)
	}
	if tag.GameID != gameID {
		t.Errorf("Expected game ID %d, got %d", gameID, tag.GameID)
	}
	if tag.Slug == "" {
		t.Error("Expected non-empty slug")
	}
}

func TestGameService_CreateDuplicateTag(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create first tag
	_, err := service.CreateTag(gameID, "Action")
	if err != nil {
		t.Fatalf("Failed to create first tag: %v", err)
	}

	// Try to create duplicate
	_, err = service.CreateTag(gameID, "Action")
	if err == nil {
		t.Error("Expected error when creating duplicate tag")
	}
}

func TestGameService_GetTags(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create test tags
	service.CreateTag(gameID, "Action")
	service.CreateTag(gameID, "Adventure")
	service.CreateTag(gameID, "RPG")

	tags, err := service.GetTags(gameID)
	if err != nil {
		t.Fatalf("GetTags() error = %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}

	// Tags should be sorted by name
	expectedNames := []string{"Action", "Adventure", "RPG"}
	for i, tag := range tags {
		if tag.Name != expectedNames[i] {
			t.Errorf("Expected tag %d to be %s, got %s", i, expectedNames[i], tag.Name)
		}
	}
}

func TestGameService_DeleteTag(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")

	tag, err := service.CreateTag(gameID, "Action")
	if err != nil {
		t.Fatalf("CreateTag() error = %v", err)
	}

	err = service.DeleteTag(tag.ID)
	if err != nil {
		t.Fatalf("DeleteTag() error = %v", err)
	}

	// Verify tag was deleted
	tags, err := service.GetTags(gameID)
	if err != nil {
		t.Fatalf("GetTags() error = %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected 0 tags after deletion, got %d", len(tags))
	}
}

func TestGameService_AssignModerator(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")
	userID := createTestUser(t, db, "moderator", "mod@example.com")

	err := service.AssignModerator(gameID, userID)
	if err != nil {
		t.Fatalf("AssignModerator() error = %v", err)
	}

	// Verify moderator was assigned
	moderators, err := service.GetModerators(gameID)
	if err != nil {
		t.Fatalf("GetModerators() error = %v", err)
	}

	if len(moderators) != 1 {
		t.Errorf("Expected 1 moderator, got %d", len(moderators))
	}
	if moderators[0].ID != userID {
		t.Errorf("Expected moderator ID %d, got %d", userID, moderators[0].ID)
	}
}

func TestGameService_AssignDuplicateModerator(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")
	userID := createTestUser(t, db, "moderator", "mod@example.com")

	// Assign moderator first time
	err := service.AssignModerator(gameID, userID)
	if err != nil {
		t.Fatalf("First AssignModerator() error = %v", err)
	}

	// Try to assign same moderator again
	err = service.AssignModerator(gameID, userID)
	if err == nil {
		t.Error("Expected error when assigning duplicate moderator")
	}
}

func TestGameService_RemoveModerator(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")
	userID := createTestUser(t, db, "moderator", "mod@example.com")

	// Assign moderator first
	err := service.AssignModerator(gameID, userID)
	if err != nil {
		t.Fatalf("AssignModerator() error = %v", err)
	}

	// Remove moderator
	err = service.RemoveModerator(gameID, userID)
	if err != nil {
		t.Fatalf("RemoveModerator() error = %v", err)
	}

	// Verify moderator was removed
	moderators, err := service.GetModerators(gameID)
	if err != nil {
		t.Fatalf("GetModerators() error = %v", err)
	}

	if len(moderators) != 0 {
		t.Errorf("Expected 0 moderators after removal, got %d", len(moderators))
	}
}

func TestGameService_RemoveNonExistentModerator(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewGameService(db)
	gameID := createTestGame(t, db, "Test Game", "test-game")
	userID := createTestUser(t, db, "user", "user@example.com")

	err := service.RemoveModerator(gameID, userID)
	if err == nil {
		t.Error("Expected error when removing non-existent moderator")
	}
}

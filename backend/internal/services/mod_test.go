package services

import (
	"mime/multipart"
	"os"
	"strings"
	"testing"

	"github.com/azurite/backend/internal/models"
)

// modMockMultipartFile implements multipart.File interface for testing
type modMockMultipartFile struct {
	*strings.Reader
}

func (m *modMockMultipartFile) Close() error {
	return nil
}

func newmodMockMultipartFile(content string) multipart.File {
	return &modMockMultipartFile{
		Reader: strings.NewReader(content),
	}
}

func TestNewModService(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	storagePath := "/tmp/storage"
	imagesPath := "/tmp/images"

	service := NewModService(db, storagePath, imagesPath)

	if service == nil {
		t.Fatal("NewModService returned nil")
	}

	if service.db != db {
		t.Error("Database not set correctly")
	}

	if service.storagePath != storagePath {
		t.Error("Storage path not set correctly")
	}

	if service.imagesPath != imagesPath {
		t.Error("Images path not set correctly")
	}
}

func TestModService_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	// Create test data
	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	t.Run("successful creation", func(t *testing.T) {
		req := &models.ModCreateRequest{
			Name:             "Test Mod",
			Description:      "Test description",
			ShortDescription: "Short desc",
			Version:          "1.0.0",
			GameVersion:      "1.0",
			GameID:           gameID,
			SourceWebsite:    "https://example.com",
			ContactInfo:      "test@example.com",
		}

		mod, err := service.Create(req, userID)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if mod == nil {
			t.Fatal("Create returned nil mod")
		}

		if mod.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, mod.Name)
		}

		if mod.Slug != "test-mod" {
			t.Errorf("Expected slug 'test-mod', got %s", mod.Slug)
		}

		if mod.OwnerID != userID {
			t.Errorf("Expected owner ID %d, got %d", userID, mod.OwnerID)
		}

		if mod.IsScanned {
			t.Error("Expected mod to not be scanned initially")
		}

		if mod.ScanResult != models.ScanResultPending {
			t.Errorf("Expected scan result to be pending, got %s", mod.ScanResult)
		}
	})

	t.Run("creation with tags", func(t *testing.T) {
		req := &models.ModCreateRequest{
			Name:             "Tagged Mod",
			Description:      "Test description",
			ShortDescription: "Short desc",
			Version:          "1.0.0",
			GameVersion:      "1.0",
			GameID:           gameID,
			Tags:             []string{"adventure", "magic"},
		}

		mod, err := service.Create(req, userID)
		if err != nil {
			t.Fatalf("Create with tags failed: %v", err)
		}

		if len(mod.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(mod.Tags))
		}
	})

	t.Run("creation with dependencies", func(t *testing.T) {
		// Create a dependency mod first
		depMod := createTestMod(t, db, "Dependency Mod", "dependency-mod", gameID, userID)

		req := &models.ModCreateRequest{
			Name:             "Dependent Mod",
			Description:      "Test description",
			ShortDescription: "Short desc",
			Version:          "1.0.0",
			GameVersion:      "1.0",
			GameID:           gameID,
			Dependencies:     []int{depMod},
		}

		mod, err := service.Create(req, userID)
		if err != nil {
			t.Fatalf("Create with dependencies failed: %v", err)
		}

		if len(mod.Dependencies) != 1 {
			t.Errorf("Expected 1 dependency, got %d", len(mod.Dependencies))
		}
	})

	t.Run("slug collision handling", func(t *testing.T) {
		// Create first mod
		req1 := &models.ModCreateRequest{
			Name:             "Duplicate Name",
			Description:      "Test description",
			ShortDescription: "Short desc",
			Version:          "1.0.0",
			GameVersion:      "1.0",
			GameID:           gameID,
		}

		mod1, err := service.Create(req1, userID)
		if err != nil {
			t.Fatalf("First mod creation failed: %v", err)
		}

		// Create second mod with same name
		req2 := &models.ModCreateRequest{
			Name:             "Duplicate Name",
			Description:      "Test description 2",
			ShortDescription: "Short desc 2",
			Version:          "1.0.1",
			GameVersion:      "1.0",
			GameID:           gameID,
		}

		mod2, err := service.Create(req2, userID)
		if err != nil {
			t.Fatalf("Second mod creation failed: %v", err)
		}

		if mod1.Slug == mod2.Slug {
			t.Error("Expected different slugs for duplicate names")
		}

		if mod2.Slug != "duplicate-name-1" {
			t.Errorf("Expected slug 'duplicate-name-1', got %s", mod2.Slug)
		}
	})
}

func TestModService_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("successful update", func(t *testing.T) {
		req := &models.ModUpdateRequest{
			Name:             "Updated Mod",
			Description:      "Updated description",
			ShortDescription: "Updated short desc",
			Version:          "2.0.0",
			GameVersion:      "2.0",
			SourceWebsite:    "https://updated.com",
			ContactInfo:      "updated@example.com",
		}

		mod, err := service.Update(modID, req, userID)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if mod.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, mod.Name)
		}

		if mod.Version != req.Version {
			t.Errorf("Expected version %s, got %s", req.Version, mod.Version)
		}
	})

	t.Run("unauthorized update", func(t *testing.T) {
		req := &models.ModUpdateRequest{
			Name: "Unauthorized Update",
		}

		_, err := service.Update(modID, req, otherUserID)
		if err == nil {
			t.Error("Expected unauthorized error")
		}

		if !strings.Contains(err.Error(), "unauthorized") {
			t.Errorf("Expected unauthorized error, got: %v", err)
		}
	})

	t.Run("update with tags and dependencies", func(t *testing.T) {
		depMod := createTestMod(t, db, "Dependency Mod", "dep-mod", gameID, userID)

		req := &models.ModUpdateRequest{
			Name:             "Updated with Relations",
			Description:      "Test description",
			ShortDescription: "Short desc",
			Version:          "1.0.0",
			GameVersion:      "1.0",
			Tags:             []string{"updated", "test"},
			Dependencies:     []int{depMod},
		}

		mod, err := service.Update(modID, req, userID)
		if err != nil {
			t.Fatalf("Update with relations failed: %v", err)
		}

		if len(mod.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(mod.Tags))
		}

		if len(mod.Dependencies) != 1 {
			t.Errorf("Expected 1 dependency, got %d", len(mod.Dependencies))
		}
	})

	t.Run("update non-existent mod", func(t *testing.T) {
		req := &models.ModUpdateRequest{
			Name: "Non-existent",
		}

		_, err := service.Update(99999, req, userID)
		if err == nil {
			t.Error("Expected error for non-existent mod")
		}
	})
}

func TestModService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("successful get", func(t *testing.T) {
		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if mod.ID != modID {
			t.Errorf("Expected ID %d, got %d", modID, mod.ID)
		}

		if mod.Name != "Test Mod" {
			t.Errorf("Expected name 'Test Mod', got %s", mod.Name)
		}

		if mod.Game == nil {
			t.Error("Game should not be nil")
		}

		if mod.Owner == nil {
			t.Error("Owner should not be nil")
		}
	})

	t.Run("get non-existent mod", func(t *testing.T) {
		_, err := service.GetByID(99999)
		if err == nil {
			t.Error("Expected error for non-existent mod")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
	})

	t.Run("get mod with nullable fields", func(t *testing.T) {
		// Create mod with nullable fields set
		_, err := db.Exec(`
			UPDATE mods SET icon = ?, rejection_reason = ? WHERE id = ?
		`, "test-icon.png", "Test rejection reason", modID)
		if err != nil {
			t.Fatalf("Failed to update mod with nullable fields: %v", err)
		}

		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if mod.Icon != "test-icon.png" {
			t.Errorf("Expected icon 'test-icon.png', got %s", mod.Icon)
		}

		if mod.RejectionReason != "Test rejection reason" {
			t.Errorf("Expected rejection reason 'Test rejection reason', got %s", mod.RejectionReason)
		}
	})
}

func TestModService_GetBySlug(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("successful get by slug", func(t *testing.T) {
		mod, err := service.GetBySlug("test-game", "test-mod")
		if err != nil {
			t.Fatalf("GetBySlug failed: %v", err)
		}

		if mod.ID != modID {
			t.Errorf("Expected ID %d, got %d", modID, mod.ID)
		}

		if mod.Game == nil || mod.Game.ID != gameID {
			t.Error("Game not set correctly")
		}

		if mod.Owner == nil || mod.Owner.ID != userID {
			t.Error("Owner not set correctly")
		}
	})

	t.Run("get non-existent mod by slug", func(t *testing.T) {
		_, err := service.GetBySlug("test-game", "non-existent")
		if err == nil {
			t.Error("Expected error for non-existent mod")
		}
	})

	t.Run("get mod with non-existent game", func(t *testing.T) {
		_, err := service.GetBySlug("non-existent-game", "test-mod")
		if err == nil {
			t.Error("Expected error for non-existent game")
		}
	})
}

func TestModService_ListByGame(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create multiple mods
	mod1ID := createTestMod(t, db, "First Mod", "first-mod", gameID, userID)
	mod2ID := createTestMod(t, db, "Second Mod", "second-mod", gameID, userID)
	mod3ID := createTestMod(t, db, "Third Mod", "third-mod", gameID, userID)

	// Set mods as scanned and clean for listing
	for _, modID := range []int{mod1ID, mod2ID, mod3ID} {
		_, err := db.Exec("UPDATE mods SET is_scanned = 1, scan_result = 'clean' WHERE id = ?", modID)
		if err != nil {
			t.Fatalf("Failed to update mod scan status: %v", err)
		}
	}

	t.Run("basic listing", func(t *testing.T) {
		mods, total, err := service.ListByGame(gameID, 1, 10, "created", "desc", nil, "")
		if err != nil {
			t.Fatalf("ListByGame failed: %v", err)
		}

		if total != 3 {
			t.Errorf("Expected total 3, got %d", total)
		}

		if len(mods) != 3 {
			t.Errorf("Expected 3 mods, got %d", len(mods))
		}
	})

	t.Run("pagination", func(t *testing.T) {
		mods, total, err := service.ListByGame(gameID, 1, 2, "created", "desc", nil, "")
		if err != nil {
			t.Fatalf("ListByGame with pagination failed: %v", err)
		}

		if total != 3 {
			t.Errorf("Expected total 3, got %d", total)
		}

		if len(mods) != 2 {
			t.Errorf("Expected 2 mods, got %d", len(mods))
		}

		// Get second page
		mods2, _, err := service.ListByGame(gameID, 2, 2, "created", "desc", nil, "")
		if err != nil {
			t.Fatalf("ListByGame page 2 failed: %v", err)
		}

		if len(mods2) != 1 {
			t.Errorf("Expected 1 mod on page 2, got %d", len(mods2))
		}
	})

	t.Run("search functionality", func(t *testing.T) {
		mods, total, err := service.ListByGame(gameID, 1, 10, "created", "desc", nil, "First")
		if err != nil {
			t.Fatalf("ListByGame with search failed: %v", err)
		}

		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}

		if len(mods) != 1 || mods[0].Name != "First Mod" {
			t.Error("Search didn't return correct mod")
		}
	})

	t.Run("sorting by name", func(t *testing.T) {
		mods, _, err := service.ListByGame(gameID, 1, 10, "name", "asc", nil, "")
		if err != nil {
			t.Fatalf("ListByGame with name sort failed: %v", err)
		}

		if len(mods) < 2 {
			t.Fatal("Need at least 2 mods for sorting test")
		}

		if mods[0].Name > mods[1].Name {
			t.Error("Mods not sorted by name ascending")
		}
	})

	t.Run("sorting by downloads desc", func(t *testing.T) {
		// Set different download counts
		_, err := db.Exec("UPDATE mods SET downloads = ? WHERE id = ?", 100, mod1ID)
		if err != nil {
			t.Fatal("Failed to update downloads")
		}
		_, err = db.Exec("UPDATE mods SET downloads = ? WHERE id = ?", 200, mod2ID)
		if err != nil {
			t.Fatal("Failed to update downloads")
		}

		mods, _, err := service.ListByGame(gameID, 1, 10, "downloads", "desc", nil, "")
		if err != nil {
			t.Fatalf("ListByGame with downloads sort failed: %v", err)
		}

		if len(mods) < 2 {
			t.Fatal("Need at least 2 mods for sorting test")
		}

		if mods[0].Downloads < mods[1].Downloads {
			t.Error("Mods not sorted by downloads descending")
		}
	})

	t.Run("with tags filter", func(t *testing.T) {
		// Create a tag and associate it with a mod
		result, err := db.Exec("INSERT INTO tags (name, slug, game_id) VALUES (?, ?, ?)", "Adventure", "adventure", gameID)
		if err != nil {
			t.Fatal("Failed to create tag")
		}
		tagID, _ := result.LastInsertId()

		_, err = db.Exec("INSERT INTO mod_tags (mod_id, tag_id) VALUES (?, ?)", mod1ID, int(tagID))
		if err != nil {
			t.Fatal("Failed to associate tag with mod")
		}

		mods, total, err := service.ListByGame(gameID, 1, 10, "created", "desc", []string{"adventure"}, "")
		if err != nil {
			t.Fatalf("ListByGame with tags failed: %v", err)
		}

		if total != 1 {
			t.Errorf("Expected total 1 with tag filter, got %d", total)
		}

		if len(mods) != 1 || mods[0].ID != mod1ID {
			t.Error("Tag filter didn't work correctly")
		}
	})
}

func TestModService_ListByOwner(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")

	// Create mods for different users
	createTestMod(t, db, "User Mod 1", "user-mod-1", gameID, userID)
	createTestMod(t, db, "User Mod 2", "user-mod-2", gameID, userID)
	createTestMod(t, db, "Other User Mod", "other-user-mod", gameID, otherUserID)

	t.Run("list by owner", func(t *testing.T) {
		mods, total, err := service.ListByOwner(userID, 1, 10)
		if err != nil {
			t.Fatalf("ListByOwner failed: %v", err)
		}

		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}

		if len(mods) != 2 {
			t.Errorf("Expected 2 mods, got %d", len(mods))
		}

		for _, mod := range mods {
			if mod.OwnerID != userID {
				t.Error("Returned mod with wrong owner ID")
			}
		}
	})

	t.Run("list by owner with pagination", func(t *testing.T) {
		mods, total, err := service.ListByOwner(userID, 1, 1)
		if err != nil {
			t.Fatalf("ListByOwner with pagination failed: %v", err)
		}

		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}

		if len(mods) != 1 {
			t.Errorf("Expected 1 mod per page, got %d", len(mods))
		}
	})
}

func TestModService_UploadFile(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	tempDir := t.TempDir()
	service := NewModService(db, tempDir, tempDir)

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("successful file upload", func(t *testing.T) {
		// Create a test file
		content := "test file content"
		file := newmodMockMultipartFile(content)

		header := &multipart.FileHeader{
			Filename: "test.zip",
			Size:     int64(len(content)),
		}

		modFile, err := service.UploadFile(modID, file, header, true)
		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}

		if modFile.ModID != modID {
			t.Errorf("Expected mod ID %d, got %d", modID, modFile.ModID)
		}

		if modFile.Filename != "test.zip" {
			t.Errorf("Expected filename 'test.zip', got %s", modFile.Filename)
		}

		if !modFile.IsMain {
			t.Error("Expected file to be marked as main")
		}

		// Check if file was created on disk
		if _, err := os.Stat(modFile.FilePath); os.IsNotExist(err) {
			t.Error("File was not created on disk")
		}
	})

	t.Run("upload with invalid file type", func(t *testing.T) {
		content := "test file content"
		file := newmodMockMultipartFile(content)

		header := &multipart.FileHeader{
			Filename: "test.invalid",
			Size:     int64(len(content)),
		}

		_, err := service.UploadFile(modID, file, header, false)
		if err == nil {
			t.Error("Expected error for invalid file type")
		}

		if !strings.Contains(err.Error(), "file type not allowed") {
			t.Errorf("Expected file type error, got: %v", err)
		}
	})

	t.Run("upload multiple files with main flag", func(t *testing.T) {
		// Upload first file as main
		content1 := "first file content"
		file1 := newmodMockMultipartFile(content1)
		header1 := &multipart.FileHeader{
			Filename: "first.zip",
			Size:     int64(len(content1)),
		}

		_, err := service.UploadFile(modID, file1, header1, true)
		if err != nil {
			t.Fatalf("First upload failed: %v", err)
		}

		// Upload second file as main (should unset first as main)
		content2 := "second file content"
		file2 := newmodMockMultipartFile(content2)
		header2 := &multipart.FileHeader{
			Filename: "second.zip",
			Size:     int64(len(content2)),
		}

		_, err = service.UploadFile(modID, file2, header2, true)
		if err != nil {
			t.Fatalf("Second upload failed: %v", err)
		}

		// Check that only the second file is marked as main
		files, err := service.getModFiles(modID)
		if err != nil {
			t.Fatalf("Failed to get mod files: %v", err)
		}

		mainCount := 0
		for _, file := range files {
			if file.IsMain {
				mainCount++
				if file.Filename != "second.zip" {
					t.Error("Wrong file marked as main")
				}
			}
		}

		if mainCount != 1 {
			t.Errorf("Expected exactly 1 main file, got %d", mainCount)
		}
	})
}

func TestModService_LikeUnlike(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("like mod", func(t *testing.T) {
		err := service.Like(modID, userID)
		if err != nil {
			t.Fatalf("Like failed: %v", err)
		}

		// Check if like was recorded
		if !service.IsLiked(modID, userID) {
			t.Error("Mod should be liked")
		}

		// Check if like count was incremented
		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("Failed to get mod: %v", err)
		}

		if mod.Likes != 1 {
			t.Errorf("Expected 1 like, got %d", mod.Likes)
		}
	})

	t.Run("like already liked mod", func(t *testing.T) {
		err := service.Like(modID, userID)
		if err == nil {
			t.Error("Expected error when liking already liked mod")
		}

		if !strings.Contains(err.Error(), "already liked") {
			t.Errorf("Expected 'already liked' error, got: %v", err)
		}
	})

	t.Run("unlike mod", func(t *testing.T) {
		err := service.Unlike(modID, userID)
		if err != nil {
			t.Fatalf("Unlike failed: %v", err)
		}

		// Check if like was removed
		if service.IsLiked(modID, userID) {
			t.Error("Mod should not be liked")
		}

		// Check if like count was decremented
		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("Failed to get mod: %v", err)
		}

		if mod.Likes != 0 {
			t.Errorf("Expected 0 likes, got %d", mod.Likes)
		}
	})

	t.Run("unlike not liked mod", func(t *testing.T) {
		err := service.Unlike(modID, userID)
		if err == nil {
			t.Error("Expected error when unliking not liked mod")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
	})
}

func TestModService_IncrementDownloadCount(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Get initial download count
	mod, err := service.GetByID(modID)
	if err != nil {
		t.Fatalf("Failed to get mod: %v", err)
	}
	initialDownloads := mod.Downloads

	// Increment download count
	err = service.IncrementDownloadCount(modID)
	if err != nil {
		t.Fatalf("IncrementDownloadCount failed: %v", err)
	}

	// Check if download count was incremented
	mod, err = service.GetByID(modID)
	if err != nil {
		t.Fatalf("Failed to get mod: %v", err)
	}

	expectedDownloads := initialDownloads + 1
	if mod.Downloads != expectedDownloads {
		t.Errorf("Expected %d downloads, got %d", expectedDownloads, mod.Downloads)
	}
}

func TestModService_RejectApprove(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	adminID := createTestUser(t, db, "admin", "admin@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("reject mod", func(t *testing.T) {
		reason := "Contains inappropriate content"
		err := service.Reject(modID, reason, adminID)
		if err != nil {
			t.Fatalf("Reject failed: %v", err)
		}

		// Check if mod was rejected
		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("Failed to get mod: %v", err)
		}

		if !mod.IsRejected {
			t.Error("Mod should be rejected")
		}

		if mod.RejectionReason != reason {
			t.Errorf("Expected rejection reason '%s', got '%s'", reason, mod.RejectionReason)
		}
	})

	t.Run("approve mod", func(t *testing.T) {
		err := service.Approve(modID)
		if err != nil {
			t.Fatalf("Approve failed: %v", err)
		}

		// Check if mod was approved
		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("Failed to get mod: %v", err)
		}

		if mod.IsRejected {
			t.Error("Mod should not be rejected")
		}

		if mod.RejectionReason != "" {
			t.Errorf("Expected empty rejection reason, got '%s'", mod.RejectionReason)
		}
	})
}

func TestModService_Search(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID1 := createTestGame(t, db, "Game One", "game-one")
	gameID2 := createTestGame(t, db, "Game Two", "game-two")

	// Create test mods
	mod1ID := createTestMod(t, db, "Adventure Mod", "adventure-mod", gameID1, userID)
	mod2ID := createTestMod(t, db, "Magic Quest", "magic-quest", gameID1, userID)
	mod3ID := createTestMod(t, db, "Racing Game", "racing-game", gameID2, userID)

	// Set all mods as scanned and clean
	for _, modID := range []int{mod1ID, mod2ID, mod3ID} {
		_, err := db.Exec("UPDATE mods SET is_scanned = 1, scan_result = 'clean' WHERE id = ?", modID)
		if err != nil {
			t.Fatalf("Failed to update mod scan status: %v", err)
		}
	}

	t.Run("search all games", func(t *testing.T) {
		mods, total, err := service.Search("Adventure", 0, 1, 10)
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}

		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}

		if len(mods) != 1 || mods[0].Name != "Adventure Mod" {
			t.Error("Search didn't return correct mod")
		}
	})

	t.Run("search specific game", func(t *testing.T) {
		mods, total, err := service.Search("Mod", gameID1, 1, 10)
		if err != nil {
			t.Fatalf("Search with game filter failed: %v", err)
		}

		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}

		if len(mods) != 1 || mods[0].GameID != gameID1 {
			t.Error("Game filter didn't work correctly")
		}
	})

	t.Run("search with pagination", func(t *testing.T) {
		mods, total, err := service.Search("", 0, 1, 2)
		if err != nil {
			t.Fatalf("Search with pagination failed: %v", err)
		}

		if total != 3 {
			t.Errorf("Expected total 3, got %d", total)
		}

		if len(mods) != 2 {
			t.Errorf("Expected 2 mods per page, got %d", len(mods))
		}
	})

	t.Run("search no results", func(t *testing.T) {
		mods, total, err := service.Search("NonExistentMod", 0, 1, 10)
		if err != nil {
			t.Fatalf("Search with no results failed: %v", err)
		}

		if total != 0 {
			t.Errorf("Expected total 0, got %d", total)
		}

		if len(mods) != 0 {
			t.Errorf("Expected 0 mods, got %d", len(mods))
		}
	})

	t.Run("search in description", func(t *testing.T) {
		// Update a mod's description
		_, err := db.Exec("UPDATE mods SET description = ? WHERE id = ?", "This is a unique description for testing", mod1ID)
		if err != nil {
			t.Fatal("Failed to update mod description")
		}

		mods, total, err := service.Search("unique description", 0, 1, 10)
		if err != nil {
			t.Fatalf("Search in description failed: %v", err)
		}

		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}

		if len(mods) != 1 || mods[0].ID != mod1ID {
			t.Error("Description search didn't work correctly")
		}
	})
}

func TestModService_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	tempDir := t.TempDir()
	service := NewModService(db, tempDir, tempDir)

	userID := createTestUser(t, db, "testuser", "test@example.com")
	otherUserID := createTestUser(t, db, "otheruser", "other@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create related data
	depModID := createTestMod(t, db, "Dependency Mod", "dep-mod", gameID, userID)

	// Add tag
	result, err := db.Exec("INSERT INTO tags (name, slug, game_id) VALUES (?, ?, ?)", "Test Tag", "test-tag", gameID)
	if err != nil {
		t.Fatal("Failed to create tag")
	}
	tagID, _ := result.LastInsertId()

	_, err = db.Exec("INSERT INTO mod_tags (mod_id, tag_id) VALUES (?, ?)", modID, int(tagID))
	if err != nil {
		t.Fatal("Failed to add mod tag")
	}

	// Add dependency
	_, err = db.Exec("INSERT INTO mod_dependencies (mod_id, dependency_id) VALUES (?, ?)", modID, depModID)
	if err != nil {
		t.Fatal("Failed to add dependency")
	}

	// Add like
	_, err = db.Exec("INSERT INTO mod_likes (mod_id, user_id) VALUES (?, ?)", modID, userID)
	if err != nil {
		t.Fatal("Failed to add like")
	}

	// Add comment
	_, err = db.Exec("INSERT INTO comments (mod_id, user_id, content) VALUES (?, ?, ?)", modID, userID, "Test comment")
	if err != nil {
		t.Fatal("Failed to add comment")
	}

	// Add file
	content := "test file content"
	file := newmodMockMultipartFile(content)
	header := &multipart.FileHeader{
		Filename: "test.zip",
		Size:     int64(len(content)),
	}

	modFile, err := service.UploadFile(modID, file, header, true)
	if err != nil {
		t.Fatal("Failed to upload file")
	}

	t.Run("successful deletion", func(t *testing.T) {
		err := service.Delete(modID, userID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Check if mod was deleted
		_, err = service.GetByID(modID)
		if err == nil {
			t.Error("Mod should be deleted")
		}

		// Check if file was deleted from disk
		if _, err := os.Stat(modFile.FilePath); !os.IsNotExist(err) {
			t.Error("File should be deleted from disk")
		}

		// Check if related data was deleted
		var count int

		// Check mod_tags
		err = db.QueryRow("SELECT COUNT(*) FROM mod_tags WHERE mod_id = ?", modID).Scan(&count)
		if err != nil || count != 0 {
			t.Error("Mod tags should be deleted")
		}

		// Check mod_dependencies
		err = db.QueryRow("SELECT COUNT(*) FROM mod_dependencies WHERE mod_id = ? OR dependency_id = ?", modID, modID).Scan(&count)
		if err != nil || count != 0 {
			t.Error("Mod dependencies should be deleted")
		}

		// Check mod_likes
		err = db.QueryRow("SELECT COUNT(*) FROM mod_likes WHERE mod_id = ?", modID).Scan(&count)
		if err != nil || count != 0 {
			t.Error("Mod likes should be deleted")
		}

		// Check comments
		err = db.QueryRow("SELECT COUNT(*) FROM comments WHERE mod_id = ?", modID).Scan(&count)
		if err != nil || count != 0 {
			t.Error("Mod comments should be deleted")
		}

		// Check mod_files
		err = db.QueryRow("SELECT COUNT(*) FROM mod_files WHERE mod_id = ?", modID).Scan(&count)
		if err != nil || count != 0 {
			t.Error("Mod files should be deleted")
		}
	})

	t.Run("unauthorized deletion", func(t *testing.T) {
		// Create new mod for this test
		newModID := createTestMod(t, db, "Another Mod", "another-mod", gameID, userID)

		err := service.Delete(newModID, otherUserID)
		if err == nil {
			t.Error("Expected unauthorized error")
		}

		if !strings.Contains(err.Error(), "unauthorized") {
			t.Errorf("Expected unauthorized error, got: %v", err)
		}
	})

	t.Run("delete non-existent mod", func(t *testing.T) {
		err := service.Delete(99999, userID)
		if err == nil {
			t.Error("Expected error for non-existent mod")
		}
	})
}

func TestModService_GetModTags(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Create tags
	result1, err := db.Exec("INSERT INTO tags (name, slug, game_id) VALUES (?, ?, ?)", "Adventure", "adventure", gameID)
	if err != nil {
		t.Fatal("Failed to create tag 1")
	}
	tagID1, _ := result1.LastInsertId()

	result2, err := db.Exec("INSERT INTO tags (name, slug, game_id) VALUES (?, ?, ?)", "Magic", "magic", gameID)
	if err != nil {
		t.Fatal("Failed to create tag 2")
	}
	tagID2, _ := result2.LastInsertId()

	// Associate tags with mod
	_, err = db.Exec("INSERT INTO mod_tags (mod_id, tag_id) VALUES (?, ?)", modID, int(tagID1))
	if err != nil {
		t.Fatal("Failed to associate tag 1")
	}
	_, err = db.Exec("INSERT INTO mod_tags (mod_id, tag_id) VALUES (?, ?)", modID, int(tagID2))
	if err != nil {
		t.Fatal("Failed to associate tag 2")
	}

	tags, err := service.getModTags(modID)
	if err != nil {
		t.Fatalf("getModTags failed: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}

	tagNames := make(map[string]bool)
	for _, tag := range tags {
		tagNames[tag.Name] = true
		if tag.GameID != gameID {
			t.Errorf("Tag has wrong game ID: expected %d, got %d", gameID, tag.GameID)
		}
	}

	if !tagNames["Adventure"] || !tagNames["Magic"] {
		t.Error("Expected Adventure and Magic tags")
	}
}

func TestModService_GetModDependencies(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)
	dep1ID := createTestMod(t, db, "Dependency 1", "dep-1", gameID, userID)
	dep2ID := createTestMod(t, db, "Dependency 2", "dep-2", gameID, userID)

	// Add dependencies
	_, err := db.Exec("INSERT INTO mod_dependencies (mod_id, dependency_id) VALUES (?, ?)", modID, dep1ID)
	if err != nil {
		t.Fatal("Failed to add dependency 1")
	}
	_, err = db.Exec("INSERT INTO mod_dependencies (mod_id, dependency_id) VALUES (?, ?)", modID, dep2ID)
	if err != nil {
		t.Fatal("Failed to add dependency 2")
	}

	deps, err := service.getModDependencies(modID)
	if err != nil {
		t.Fatalf("getModDependencies failed: %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	depIDs := make(map[int]bool)
	for _, dep := range deps {
		depIDs[dep.ID] = true
		if dep.GameID != gameID {
			t.Errorf("Dependency has wrong game ID: expected %d, got %d", gameID, dep.GameID)
		}
	}

	if !depIDs[dep1ID] || !depIDs[dep2ID] {
		t.Error("Expected both dependency IDs")
	}
}

func TestModService_GetModFiles(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Upload files
	content1 := "file 1 content"
	file1 := newmodMockMultipartFile(content1)
	header1 := &multipart.FileHeader{
		Filename: "file1.zip",
		Size:     int64(len(content1)),
	}

	content2 := "file 2 content"
	file2 := newmodMockMultipartFile(content2)
	header2 := &multipart.FileHeader{
		Filename: "file2.zip",
		Size:     int64(len(content2)),
	}

	_, err := service.UploadFile(modID, file1, header1, true)
	if err != nil {
		t.Fatal("Failed to upload file 1")
	}

	_, err = service.UploadFile(modID, file2, header2, false)
	if err != nil {
		t.Fatal("Failed to upload file 2")
	}

	files, err := service.getModFiles(modID)
	if err != nil {
		t.Fatalf("getModFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	// Check that main file comes first
	if !files[0].IsMain {
		t.Error("Main file should be first")
	}

	fileNames := make(map[string]bool)
	for _, file := range files {
		fileNames[file.Filename] = true
		if file.ModID != modID {
			t.Errorf("File has wrong mod ID: expected %d, got %d", modID, file.ModID)
		}
	}

	if !fileNames["file1.zip"] || !fileNames["file2.zip"] {
		t.Error("Expected both filenames")
	}
}

func TestModService_UpdateModTags(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("create new tags", func(t *testing.T) {
		tagNames := []string{"Adventure", "Magic", "Fantasy"}

		err := service.updateModTags(modID, tagNames, gameID)
		if err != nil {
			t.Fatalf("updateModTags failed: %v", err)
		}

		// Check if tags were created and associated
		tags, err := service.getModTags(modID)
		if err != nil {
			t.Fatalf("Failed to get mod tags: %v", err)
		}

		if len(tags) != 3 {
			t.Errorf("Expected 3 tags, got %d", len(tags))
		}

		tagNamesMap := make(map[string]bool)
		for _, tag := range tags {
			tagNamesMap[tag.Name] = true
		}

		for _, expectedName := range tagNames {
			if !tagNamesMap[expectedName] {
				t.Errorf("Tag %s not found", expectedName)
			}
		}
	})

	t.Run("use existing tags", func(t *testing.T) {
		// Create another mod
		mod2ID := createTestMod(t, db, "Test Mod 2", "test-mod-2", gameID, userID)

		// Use existing tags
		tagNames := []string{"Adventure", "Magic"}

		err := service.updateModTags(mod2ID, tagNames, gameID)
		if err != nil {
			t.Fatalf("updateModTags with existing tags failed: %v", err)
		}

		// Check that tags were reused, not duplicated
		var tagCount int
		err = db.QueryRow("SELECT COUNT(*) FROM tags WHERE game_id = ?", gameID).Scan(&tagCount)
		if err != nil {
			t.Fatal("Failed to count tags")
		}

		if tagCount != 3 { // Should still be 3 from previous test
			t.Errorf("Expected 3 total tags, got %d", tagCount)
		}
	})
}

func TestModService_UpdateModDependencies(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)
	dep1ID := createTestMod(t, db, "Dependency 1", "dep-1", gameID, userID)
	dep2ID := createTestMod(t, db, "Dependency 2", "dep-2", gameID, userID)

	t.Run("add dependencies", func(t *testing.T) {
		dependencyIDs := []int{dep1ID, dep2ID}

		err := service.updateModDependencies(modID, dependencyIDs)
		if err != nil {
			t.Fatalf("updateModDependencies failed: %v", err)
		}

		// Check if dependencies were added
		deps, err := service.getModDependencies(modID)
		if err != nil {
			t.Fatalf("Failed to get mod dependencies: %v", err)
		}

		if len(deps) != 2 {
			t.Errorf("Expected 2 dependencies, got %d", len(deps))
		}
	})

	t.Run("ignore self dependency", func(t *testing.T) {
		// Try to add self as dependency
		dependencyIDs := []int{modID, dep1ID}

		err := service.updateModDependencies(modID, dependencyIDs)
		if err != nil {
			t.Fatalf("updateModDependencies with self dependency failed: %v", err)
		}

		// Check that self-dependency was ignored
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM mod_dependencies WHERE mod_id = ? AND dependency_id = ?", modID, modID).Scan(&count)
		if err != nil {
			t.Fatal("Failed to check self dependency")
		}

		if count != 0 {
			t.Error("Self dependency should be ignored")
		}
	})

	t.Run("ignore non-existent dependency", func(t *testing.T) {
		dependencyIDs := []int{99999, dep1ID}

		err := service.updateModDependencies(modID, dependencyIDs)
		if err != nil {
			t.Fatalf("updateModDependencies with non-existent dependency failed: %v", err)
		}

		// Should not add non-existent dependency
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM mod_dependencies WHERE mod_id = ? AND dependency_id = ?", modID, 99999).Scan(&count)
		if err != nil {
			t.Fatal("Failed to check non-existent dependency")
		}

		if count != 0 {
			t.Error("Non-existent dependency should be ignored")
		}
	})
}

func TestModService_SimulateMalwareScan(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Check initial state
	mod, err := service.GetByID(modID)
	if err != nil {
		t.Fatalf("Failed to get mod: %v", err)
	}

	if mod.IsScanned {
		t.Error("Mod should not be scanned initially")
	}

	// Run simulation (this is called asynchronously in real code)
	service.simulateMalwareScan(modID)

	// Check that mod is now scanned and clean
	mod, err = service.GetByID(modID)
	if err != nil {
		t.Fatalf("Failed to get mod after scan: %v", err)
	}

	if !mod.IsScanned {
		t.Error("Mod should be scanned after simulation")
	}

	if mod.ScanResult != "clean" {
		t.Errorf("Expected scan result 'clean', got %s", mod.ScanResult)
	}
}

func TestModService_SimulateFileScan(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	t.Run("clean file scan", func(t *testing.T) {
		fileID := 1 // Not divisible by 13, should be clean

		service.simulateFileScan(modID, fileID)

		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("Failed to get mod after file scan: %v", err)
		}

		if mod.ScanResult != "clean" {
			t.Errorf("Expected scan result 'clean', got %s", mod.ScanResult)
		}

		if mod.IsRejected {
			t.Error("Mod should not be rejected for clean scan")
		}
	})

	t.Run("threat detection", func(t *testing.T) {
		fileID := 13 // Divisible by 13, should trigger threat detection

		service.simulateFileScan(modID, fileID)

		mod, err := service.GetByID(modID)
		if err != nil {
			t.Fatalf("Failed to get mod after threat scan: %v", err)
		}

		if mod.ScanResult != "threat" {
			t.Errorf("Expected scan result 'threat', got %s", mod.ScanResult)
		}

		if !mod.IsRejected {
			t.Error("Mod should be rejected for threat detection")
		}

		if !strings.Contains(mod.RejectionReason, "Malware detected") {
			t.Errorf("Expected malware rejection reason, got: %s", mod.RejectionReason)
		}
	})
}

func TestModService_IsLiked(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	service := NewModService(db, t.TempDir(), t.TempDir())

	userID := createTestUser(t, db, "testuser", "test@example.com")
	gameID := createTestGame(t, db, "Test Game", "test-game")
	modID := createTestMod(t, db, "Test Mod", "test-mod", gameID, userID)

	// Initially not liked
	if service.IsLiked(modID, userID) {
		t.Error("Mod should not be liked initially")
	}

	// Add like
	err := service.Like(modID, userID)
	if err != nil {
		t.Fatalf("Failed to like mod: %v", err)
	}

	// Now should be liked
	if !service.IsLiked(modID, userID) {
		t.Error("Mod should be liked after liking")
	}

	// Remove like
	err = service.Unlike(modID, userID)
	if err != nil {
		t.Fatalf("Failed to unlike mod: %v", err)
	}

	// Should not be liked again
	if service.IsLiked(modID, userID) {
		t.Error("Mod should not be liked after unliking")
	}
}

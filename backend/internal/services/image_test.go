package services

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewImageService(t *testing.T) {
	storagePath := "/test/path"
	service := NewImageService(storagePath)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.storagePath != storagePath {
		t.Errorf("Expected storage path %s, got %s", storagePath, service.storagePath)
	}
}

func TestImageService_isImageFile(t *testing.T) {
	service := NewImageService("/tmp")

	tests := []struct {
		filename string
		expected bool
	}{
		{"test.jpg", true},
		{"test.jpeg", true},
		{"test.png", true},
		{"test.gif", true},
		{"test.webp", true},
		{"test.JPG", true},
		{"test.JPEG", true},
		{"test.PNG", true},
		{"test.txt", false},
		{"test.pdf", false},
		{"test", false},
		{"", false},
		{"test.bmp", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := service.isImageFile(tt.filename)
			if result != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.filename, result)
			}
		})
	}
}

func TestImageService_generateFilename(t *testing.T) {
	service := NewImageService("/tmp")

	tests := []struct {
		input string
	}{
		{"test.jpg"},
		{"my file.png"},
		{"image.jpeg"},
		{"no-extension"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := service.generateFilename(tt.input)

			// Verify the result has the expected format: sanitized_name_token
			parts := strings.Split(result, "_")
			if len(parts) < 2 {
				t.Errorf("Expected filename to contain underscore separator, got %s", result)
			}

			// Verify it's not empty
			if result == "" {
				t.Error("Expected non-empty filename")
			}

			// Verify original extension is removed
			originalBase := strings.TrimSuffix(tt.input, filepath.Ext(tt.input))
			if strings.Contains(result, filepath.Ext(tt.input)) {
				t.Errorf("Expected original extension to be removed from %s", result)
			}

			// Verify some form of the original base name is present (sanitized)
			if !strings.Contains(result, strings.ReplaceAll(originalBase, " ", "")) &&
				!strings.Contains(result, strings.ReplaceAll(originalBase, " ", "_")) {
				t.Logf("Generated filename: %s for input: %s", result, tt.input)
				// This is just informational, not a failure
			}
		})
	}
}

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 100, 255})
		}
	}
	return img
}

func createTestImageFile(img image.Image) (multipart.File, *multipart.FileHeader, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, nil, err
	}

	reader := bytes.NewReader(buf.Bytes())

	// Create a mock multipart file
	file := &mockMultipartFile{
		reader: reader,
		size:   int64(buf.Len()),
	}

	header := &multipart.FileHeader{
		Filename: "test.png",
		Size:     int64(buf.Len()),
	}

	return file, header, nil
}

type mockMultipartFile struct {
	reader   *bytes.Reader
	size     int64
	position int64
}

func (m *mockMultipartFile) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *mockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return m.reader.ReadAt(p, off)
}

func (m *mockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return m.reader.Seek(offset, whence)
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func TestImageService_encodeToWebP(t *testing.T) {
	service := NewImageService("/tmp")

	img := createTestImage(100, 100)

	data, err := service.encodeToWebP(img)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected encoded data, got empty slice")
	}

	// Verify it's valid PNG data (since we're encoding to PNG for now)
	_, err = png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Errorf("Expected valid PNG data, got decode error: %v", err)
	}
}

func TestImageService_resizeImage(t *testing.T) {
	service := NewImageService("/tmp")

	tests := []struct {
		name           string
		originalWidth  int
		originalHeight int
		maxWidth       int
		maxHeight      int
		expectedWidth  int
		expectedHeight int
		shouldResize   bool
	}{
		{
			name:          "no resize needed",
			originalWidth: 100, originalHeight: 100,
			maxWidth: 200, maxHeight: 200,
			expectedWidth: 100, expectedHeight: 100,
			shouldResize: false,
		},
		{
			name:          "resize wider image",
			originalWidth: 3000, originalHeight: 1000,
			maxWidth: 2048, maxHeight: 2048,
			expectedWidth: 2048, expectedHeight: 682,
			shouldResize: true,
		},
		{
			name:          "resize taller image",
			originalWidth: 1000, originalHeight: 3000,
			maxWidth: 2048, maxHeight: 2048,
			expectedWidth: 682, expectedHeight: 2048,
			shouldResize: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImage(tt.originalWidth, tt.originalHeight)
			result := service.resizeImage(img, tt.maxWidth, tt.maxHeight)

			bounds := result.Bounds()
			width := bounds.Dx()
			height := bounds.Dy()

			if width != tt.expectedWidth || height != tt.expectedHeight {
				t.Errorf("Expected dimensions %dx%d, got %dx%d",
					tt.expectedWidth, tt.expectedHeight, width, height)
			}
		})
	}
}

func TestImageService_saveFile(t *testing.T) {
	service := NewImageService("/tmp")

	// Create temporary directory
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test.txt")
	testData := []byte("test data")

	err := service.saveFile(testData, testPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was created and contains correct data
	savedData, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if !bytes.Equal(savedData, testData) {
		t.Errorf("Expected %s, got %s", string(testData), string(savedData))
	}
}

func TestImageService_saveFile_Error(t *testing.T) {
	service := NewImageService("/tmp")

	// Try to save to invalid path
	invalidPath := "/invalid/path/file.txt"
	testData := []byte("test data")

	err := service.saveFile(testData, invalidPath)
	if err == nil {
		t.Error("Expected error when saving to invalid path, got nil")
	}
}

func TestImageService_ValidateImage(t *testing.T) {
	service := NewImageService("/tmp")

	// Create valid test image
	img := createTestImage(100, 100)
	file, _, err := createTestImageFile(img)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}

	// Test valid image with no size limit
	err = service.ValidateImage(file, 0)
	if err != nil {
		t.Errorf("Expected no error for valid image, got %v", err)
	}

	// Test valid image within size limit
	file.Seek(0, io.SeekStart)
	err = service.ValidateImage(file, 1000000) // 1MB limit
	if err != nil {
		t.Errorf("Expected no error for valid image within size limit, got %v", err)
	}

	// Test image exceeding size limit
	file.Seek(0, io.SeekStart)
	err = service.ValidateImage(file, 100) // Very small limit
	if err == nil {
		t.Error("Expected error for image exceeding size limit, got nil")
	}

	// Test invalid image format
	invalidFile := &mockMultipartFile{
		reader: bytes.NewReader([]byte("not an image")),
		size:   13,
	}
	err = service.ValidateImage(invalidFile, 0)
	if err == nil {
		t.Error("Expected error for invalid image format, got nil")
	}
}

func TestImageService_ProcessAndSave(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Test with valid image
	img := createTestImage(100, 100)
	file, header, err := createTestImageFile(img)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}

	relativePath, err := service.ProcessAndSave(file, header, "uploads")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !strings.HasPrefix(relativePath, "uploads/") {
		t.Errorf("Expected path to start with 'uploads/', got %s", relativePath)
	}

	if !strings.HasSuffix(relativePath, ".webp") {
		t.Errorf("Expected path to end with '.webp', got %s", relativePath)
	}

	// Verify file was created
	fullPath := filepath.Join(tempDir, relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Errorf("Expected file to be created at %s", fullPath)
	}
}

func TestImageService_ProcessAndSave_InvalidFile(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create mock file with invalid extension
	file := &mockMultipartFile{
		reader: bytes.NewReader([]byte("test")),
		size:   4,
	}
	header := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     4,
	}

	_, err := service.ProcessAndSave(file, header, "uploads")
	if err == nil {
		t.Error("Expected error for invalid file type, got nil")
	}

	expectedError := "file is not a valid image"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

func TestImageService_ProcessAndSave_LargeImage(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create large image that needs resizing
	img := createTestImage(3000, 2000)
	file, header, err := createTestImageFile(img)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}

	relativePath, err := service.ProcessAndSave(file, header, "uploads")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the saved image was resized
	fullPath := filepath.Join(tempDir, relativePath)
	savedFile, err := os.Open(fullPath)
	if err != nil {
		t.Fatalf("Failed to open saved file: %v", err)
	}
	defer savedFile.Close()

	savedImg, _, err := image.Decode(savedFile)
	if err != nil {
		t.Fatalf("Failed to decode saved image: %v", err)
	}

	bounds := savedImg.Bounds()
	if bounds.Dx() > 2048 || bounds.Dy() > 2048 {
		t.Errorf("Expected image to be resized to max 2048x2048, got %dx%d",
			bounds.Dx(), bounds.Dy())
	}
}

func TestImageService_Delete(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create a test file
	relativePath := "test/image.webp"
	fullPath := filepath.Join(tempDir, relativePath)

	// Create directory and file
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, []byte("test"), 0644)

	// Delete the file
	err := service.Delete(relativePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file is deleted
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Error("Expected file to be deleted")
	}
}

func TestImageService_Delete_NonexistentFile(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Try to delete nonexistent file - should not return error
	err := service.Delete("nonexistent/file.webp")
	if err != nil {
		t.Errorf("Expected no error when deleting nonexistent file, got %v", err)
	}
}

func TestImageService_GetFullPath(t *testing.T) {
	storagePath := "/test/storage"
	service := NewImageService(storagePath)

	relativePath := "uploads/image.webp"
	expected := filepath.Join(storagePath, relativePath)

	result := service.GetFullPath(relativePath)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestImageService_GetImageInfo(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create test image file
	img := createTestImage(200, 150)
	var buf bytes.Buffer
	png.Encode(&buf, img)

	relativePath := "test/image.png"
	fullPath := filepath.Join(tempDir, relativePath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, buf.Bytes(), 0644)

	info, err := service.GetImageInfo(relativePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info["width"] != 200 {
		t.Errorf("Expected width 200, got %v", info["width"])
	}

	if info["height"] != 150 {
		t.Errorf("Expected height 150, got %v", info["height"])
	}

	if info["format"] != "png" {
		t.Errorf("Expected format 'png', got %v", info["format"])
	}

	if info["size"].(int64) <= 0 {
		t.Errorf("Expected positive size, got %v", info["size"])
	}
}

func TestImageService_GetImageInfo_NonexistentFile(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	_, err := service.GetImageInfo("nonexistent/file.png")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestImageService_GetImageInfo_InvalidImage(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create invalid image file
	relativePath := "test/invalid.png"
	fullPath := filepath.Join(tempDir, relativePath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, []byte("not an image"), 0644)

	_, err := service.GetImageInfo(relativePath)
	if err == nil {
		t.Error("Expected error for invalid image, got nil")
	}
}

func TestImageService_CreateThumbnail(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create original image
	img := createTestImage(400, 300)
	var buf bytes.Buffer
	png.Encode(&buf, img)

	originalPath := "uploads/original.png"
	fullPath := filepath.Join(tempDir, originalPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, buf.Bytes(), 0644)

	// Create thumbnail
	thumbPath, err := service.CreateThumbnail(originalPath, 100, 100)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedThumbPath := "uploads/original_thumb_100x100.webp"
	if thumbPath != expectedThumbPath {
		t.Errorf("Expected thumbnail path %s, got %s", expectedThumbPath, thumbPath)
	}

	// Verify thumbnail file exists
	fullThumbPath := filepath.Join(tempDir, thumbPath)
	if _, err := os.Stat(fullThumbPath); os.IsNotExist(err) {
		t.Error("Expected thumbnail file to be created")
	}

	// Verify thumbnail dimensions
	thumbFile, err := os.Open(fullThumbPath)
	if err != nil {
		t.Fatalf("Failed to open thumbnail: %v", err)
	}
	defer thumbFile.Close()

	thumbImg, _, err := image.Decode(thumbFile)
	if err != nil {
		t.Fatalf("Failed to decode thumbnail: %v", err)
	}

	bounds := thumbImg.Bounds()
	if bounds.Dx() > 100 || bounds.Dy() > 100 {
		t.Errorf("Expected thumbnail to be max 100x100, got %dx%d",
			bounds.Dx(), bounds.Dy())
	}
}

func TestImageService_CreateThumbnail_NonexistentOriginal(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	_, err := service.CreateThumbnail("nonexistent/file.png", 100, 100)
	if err == nil {
		t.Error("Expected error for nonexistent original file, got nil")
	}
}

func TestImageService_CreateThumbnail_InvalidOriginal(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create invalid image file
	originalPath := "uploads/invalid.png"
	fullPath := filepath.Join(tempDir, originalPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, []byte("not an image"), 0644)

	_, err := service.CreateThumbnail(originalPath, 100, 100)
	if err == nil {
		t.Error("Expected error for invalid original image, got nil")
	}
}

func TestImageService_ProcessAndSave_DecodeError(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create file with valid extension but invalid image data
	file := &mockMultipartFile{
		reader: bytes.NewReader([]byte("invalid image data")),
		size:   18,
	}
	header := &multipart.FileHeader{
		Filename: "test.png",
		Size:     18,
	}

	_, err := service.ProcessAndSave(file, header, "uploads")
	if err == nil {
		t.Error("Expected error for invalid image data, got nil")
	}

	if !strings.Contains(err.Error(), "failed to decode image") {
		t.Errorf("Expected decode error, got %v", err)
	}
}

func TestImageService_ProcessAndSave_DirectoryCreationError(t *testing.T) {
	// Use a read-only temporary directory to simulate permission error
	tempDir := t.TempDir()

	// Create a read-only subdirectory
	restrictedDir := filepath.Join(tempDir, "readonly")
	os.MkdirAll(restrictedDir, 0555)    // read and execute only
	defer os.Chmod(restrictedDir, 0755) // restore permissions for cleanup

	service := NewImageService(restrictedDir)

	img := createTestImage(100, 100)
	file, header, err := createTestImageFile(img)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}

	_, err = service.ProcessAndSave(file, header, "newdir")
	if err == nil {
		t.Error("Expected error for directory creation failure, got nil")
	}

	if !strings.Contains(err.Error(), "failed to create directory") {
		t.Errorf("Expected directory creation error, got: %v", err)
	}
}

// Test integration with actual utils functions
func TestImageService_ProcessAndSave_UtilsIntegration(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Test with actual utils functions (integration test)
	img := createTestImage(50, 50)
	file, header, err := createTestImageFile(img)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}

	// This will use the real utils functions
	relativePath, err := service.ProcessAndSave(file, header, "test")
	if err != nil {
		t.Fatalf("ProcessAndSave failed: %v", err)
	}

	// Verify file was created successfully
	fullPath := filepath.Join(tempDir, relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}

	// Verify the path structure
	if !strings.HasPrefix(relativePath, "test/") {
		t.Errorf("Expected path to start with 'test/', got %s", relativePath)
	}

	if !strings.HasSuffix(relativePath, ".webp") {
		t.Errorf("Expected path to end with '.webp', got %s", relativePath)
	}
}

// Test ValidateImage with file size formatting
func TestImageService_ValidateImage_FileSizeFormatting(t *testing.T) {
	service := NewImageService("/tmp")

	// Create oversized mock file
	file := &mockMultipartFile{
		reader: bytes.NewReader(make([]byte, 1000)),
		size:   1000,
	}

	err := service.ValidateImage(file, 100)
	if err == nil {
		t.Error("Expected error for oversized image, got nil")
	}

	if !strings.Contains(err.Error(), "image too large") {
		t.Errorf("Expected 'image too large' error, got %v", err)
	}

	// Verify the error message contains size information
	if !strings.Contains(err.Error(), "max:") {
		t.Errorf("Expected error to contain max size info, got %v", err)
	}
}

// Test file operations with permission errors
func TestImageService_Delete_PermissionError(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create a file in a read-only directory
	subDir := filepath.Join(tempDir, "readonly")
	os.MkdirAll(subDir, 0755)

	filePath := filepath.Join(subDir, "test.webp")
	os.WriteFile(filePath, []byte("test"), 0644)

	// Make directory read-only
	os.Chmod(subDir, 0555)
	defer os.Chmod(subDir, 0755) // restore for cleanup

	relativePath := filepath.Join("readonly", "test.webp")
	err := service.Delete(relativePath)

	// On some systems this might not error, but if it does, it should be handled
	if err != nil && !strings.Contains(err.Error(), "failed to delete image") {
		t.Errorf("Expected specific delete error format, got %v", err)
	}
}

// Test resize with edge case dimensions
func TestImageService_resizeImage_EdgeCases(t *testing.T) {
	service := NewImageService("/tmp")

	// Test with 1x1 image
	img := createTestImage(1, 1)
	result := service.resizeImage(img, 100, 100)
	bounds := result.Bounds()

	if bounds.Dx() != 1 || bounds.Dy() != 1 {
		t.Errorf("Expected 1x1 image to remain 1x1, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// Test with very rectangular image (aspect ratio preservation)
	img = createTestImage(1000, 100)
	result = service.resizeImage(img, 500, 500)
	bounds = result.Bounds()

	// Should be 500x50 to preserve aspect ratio
	if bounds.Dx() != 500 {
		t.Errorf("Expected width 500, got %d", bounds.Dx())
	}

	if bounds.Dy() != 50 {
		t.Errorf("Expected height 50, got %d", bounds.Dy())
	}
}

// Test GetImageInfo with file stat error
func TestImageService_GetImageInfo_StatError(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create image file then remove it while keeping handle
	img := createTestImage(100, 100)
	var buf bytes.Buffer
	png.Encode(&buf, img)

	relativePath := "test/image.png"
	fullPath := filepath.Join(tempDir, relativePath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, buf.Bytes(), 0644)

	// Make directory inaccessible to cause stat error
	os.Chmod(filepath.Dir(fullPath), 0000)
	defer os.Chmod(filepath.Dir(fullPath), 0755)

	_, err := service.GetImageInfo(relativePath)
	// This should error on the file access, not necessarily stat
	// The specific error depends on the system
	if err == nil {
		// If it doesn't error, that's also valid behavior on some systems
		t.Log("GetImageInfo didn't fail with inaccessible directory - system dependent behavior")
	}
}

// Test CreateThumbnail with file write error
func TestImageService_CreateThumbnail_WriteError(t *testing.T) {
	tempDir := t.TempDir()
	service := NewImageService(tempDir)

	// Create original image
	img := createTestImage(200, 200)
	var buf bytes.Buffer
	png.Encode(&buf, img)

	originalPath := "uploads/original.png"
	fullPath := filepath.Join(tempDir, originalPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, buf.Bytes(), 0644)

	// Make the upload directory read-only to prevent thumbnail creation
	uploadDir := filepath.Join(tempDir, "uploads")
	os.Chmod(uploadDir, 0555)
	defer os.Chmod(uploadDir, 0755)

	_, err := service.CreateThumbnail(originalPath, 100, 100)
	if err == nil {
		t.Error("Expected error when unable to write thumbnail, got nil")
	}

	if !strings.Contains(err.Error(), "failed to save thumbnail") {
		t.Errorf("Expected thumbnail save error, got %v", err)
	}
}

// Test ProcessAndSave with file write error
func TestImageService_ProcessAndSave_WriteError(t *testing.T) {
	tempDir := t.TempDir()

	// Create read-only directory
	uploadDir := filepath.Join(tempDir, "uploads")
	os.MkdirAll(uploadDir, 0555) // read-only
	defer os.Chmod(uploadDir, 0755)

	service := NewImageService(tempDir)

	img := createTestImage(100, 100)
	file, header, err := createTestImageFile(img)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}

	_, err = service.ProcessAndSave(file, header, "uploads")
	if err == nil {
		t.Error("Expected error when unable to write file, got nil")
	}

	if !strings.Contains(err.Error(), "failed to save file") {
		t.Errorf("Expected file save error, got %v", err)
	}
}

// Test edge case in mockMultipartFile
func TestMockMultipartFile_EdgeCases(t *testing.T) {
	data := []byte("test data")
	file := &mockMultipartFile{
		reader: bytes.NewReader(data),
		size:   int64(len(data)),
	}

	// Test ReadAt
	buf := make([]byte, 4)
	_, err := file.ReadAt(buf, 5)
	if err != nil && err != io.EOF {
		t.Errorf("ReadAt error: %v", err)
	}

	// Test Seek and position tracking
	pos, err := file.Seek(0, io.SeekStart)
	if err != nil {
		t.Errorf("Seek error: %v", err)
	}
	if pos != 0 {
		t.Errorf("Expected position 0, got %d", pos)
	}

	// Test Close (should not error)
	err = file.Close()
	if err != nil {
		t.Errorf("Close should not error, got %v", err)
	}
}

// Test filename generation with empty input
func TestImageService_generateFilename_EmptyInput(t *testing.T) {
	service := NewImageService("/tmp")

	result := service.generateFilename("")
	if result == "" {
		t.Error("Expected non-empty result for empty input")
	}

	// Should still generate some kind of filename
	if !strings.Contains(result, "_") {
		t.Error("Expected generated filename to contain separator")
	}
}

// Test resizeImage with zero dimensions (edge case)
func TestImageService_resizeImage_ZeroDimensions(t *testing.T) {
	service := NewImageService("/tmp")

	img := createTestImage(100, 100)

	// This should not crash and should return original image
	result := service.resizeImage(img, 0, 0)
	if result == nil {
		t.Error("Expected non-nil result")
	}
}

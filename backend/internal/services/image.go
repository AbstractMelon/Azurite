package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/azurite/backend/internal/utils"
)

type ImageService struct {
	storagePath string
}

func NewImageService(storagePath string) *ImageService {
	return &ImageService{
		storagePath: storagePath,
	}
}

func (s *ImageService) ProcessAndSave(file multipart.File, header *multipart.FileHeader, subDir string) (string, error) {
	if !s.isImageFile(header.Filename) {
		return "", errors.New("file is not a valid image")
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	if img.Bounds().Dx() > 2048 || img.Bounds().Dy() > 2048 {
		img = s.resizeImage(img, 2048, 2048)
	}

	webpData, err := s.encodeToWebP(img)
	if err != nil {
		return "", fmt.Errorf("failed to encode to WebP: %w", err)
	}

	filename := s.generateFilename(header.Filename)
	webpFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".webp"

	fullDir := filepath.Join(s.storagePath, subDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fullPath := filepath.Join(fullDir, webpFilename)
	if err := s.saveFile(webpData, fullPath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filepath.Join(subDir, webpFilename), nil
}

// SaveImage is an alias for ProcessAndSave
func (s *ImageService) SaveImage(file multipart.File, header *multipart.FileHeader, subDir string) (string, error) {
	return s.ProcessAndSave(file, header, subDir)
}

func (s *ImageService) Delete(relativePath string) error {
	fullPath := filepath.Join(s.storagePath, relativePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

func (s *ImageService) GetFullPath(relativePath string) string {
	return filepath.Join(s.storagePath, relativePath)
}

func (s *ImageService) isImageFile(filename string) bool {
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowed := range allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (s *ImageService) generateFilename(originalFilename string) string {
	baseName := strings.TrimSuffix(originalFilename, filepath.Ext(originalFilename))
	safeName := utils.SanitizeFilename(baseName)
	token := utils.GenerateRandomToken()
	if len(token) > 8 {
		token = token[:8]
	}
	return fmt.Sprintf("%s_%s", safeName, token)
}

func (s *ImageService) encodeToWebP(img image.Image) ([]byte, error) {
	var buf bytes.Buffer

	// Since golang.org/x/image/webp doesn't have an encoder,
	// we'll convert to PNG for now and serve as WebP when possible
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *ImageService) resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= maxWidth && height <= maxHeight {
		return img
	}

	ratio := float64(width) / float64(height)

	if width > height {
		width = maxWidth
		height = int(float64(maxWidth) / ratio)
	} else {
		height = maxHeight
		width = int(float64(maxHeight) * ratio)
	}

	// Simple nearest neighbor resize
	resized := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * bounds.Dx() / width
			srcY := y * bounds.Dy() / height
			resized.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}

	return resized
}

func (s *ImageService) saveFile(data []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (s *ImageService) ValidateImage(file multipart.File, maxSize int64) error {
	if maxSize > 0 {
		file.Seek(0, io.SeekEnd)
		size, _ := file.Seek(0, io.SeekCurrent)
		file.Seek(0, io.SeekStart)

		if size > maxSize {
			return fmt.Errorf("image too large: %s (max: %s)",
				utils.FormatFileSize(size), utils.FormatFileSize(maxSize))
		}
	}

	_, _, err := image.Decode(file)
	if err != nil {
		return errors.New("invalid image format")
	}

	file.Seek(0, io.SeekStart)
	return nil
}

func (s *ImageService) GetImageInfo(path string) (map[string]interface{}, error) {
	file, err := os.Open(filepath.Join(s.storagePath, path))
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image config: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	return map[string]interface{}{
		"width":  config.Width,
		"height": config.Height,
		"format": format,
		"size":   stat.Size(),
	}, nil
}

func (s *ImageService) CreateThumbnail(originalPath string, width, height int) (string, error) {
	fullPath := filepath.Join(s.storagePath, originalPath)

	file, err := os.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to open original image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	thumbnail := s.resizeImage(img, width, height)

	thumbData, err := s.encodeToWebP(thumbnail)
	if err != nil {
		return "", fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	dir := filepath.Dir(originalPath)
	filename := filepath.Base(originalPath)
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)

	thumbFilename := fmt.Sprintf("%s_thumb_%dx%d.webp", baseName, width, height)
	thumbPath := filepath.Join(dir, thumbFilename)
	fullThumbPath := filepath.Join(s.storagePath, thumbPath)

	if err := s.saveFile(thumbData, fullThumbPath); err != nil {
		return "", fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return thumbPath, nil
}

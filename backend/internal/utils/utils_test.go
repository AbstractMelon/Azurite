package utils

import (
	"bytes"
	"mime/multipart"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	if hash == password {
		t.Error("Hash should not be the same as password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("Expected password to match hash")
	}

	if CheckPasswordHash(wrongPassword, hash) {
		t.Error("Expected wrong password to not match hash")
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test-Game_Mod", "test-game_mod"},
		{"Special!@#$%Characters", "specialcharacters"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"---dashes---", "dashes"},
		{"", "item-"}, // Will have timestamp appended
		{"UPPERCASE", "uppercase"},
		{"Mixed_Case-With123Numbers", "mixed_case-with123numbers"},
	}

	for _, test := range tests {
		result := GenerateSlug(test.input)
		if test.input == "" {
			// For empty string, check it starts with "item-" and has timestamp
			if !strings.HasPrefix(result, "item-") {
				t.Errorf("Empty input should generate slug starting with 'item-', got %s", result)
			}
		} else {
			if result != test.expected {
				t.Errorf("GenerateSlug(%q) = %q, expected %q", test.input, result, test.expected)
			}
		}
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token1 := GenerateRandomToken()
	token2 := GenerateRandomToken()

	if token1 == "" {
		t.Error("Expected non-empty token")
	}

	if len(token1) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("Expected token length 64, got %d", len(token1))
	}

	if token1 == token2 {
		t.Error("Expected different tokens to be generated")
	}
}

func TestIsAllowedFileType(t *testing.T) {
	allowedTypes := []string{".zip", ".rar", ".jar"}

	tests := []struct {
		filename string
		expected bool
	}{
		{"test.zip", true},
		{"test.ZIP", true}, // Case insensitive
		{"test.rar", true},
		{"test.jar", true},
		{"test.exe", false},
		{"test.txt", false},
		{"noextension", false},
	}

	for _, test := range tests {
		result := IsAllowedFileType(test.filename, allowedTypes)
		if result != test.expected {
			t.Errorf("IsAllowedFileType(%q) = %v, expected %v", test.filename, result, test.expected)
		}
	}
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.zip", "application/zip"},
		{"test.ZIP", "application/zip"}, // Case insensitive
		{"test.png", "image/png"},
		{"test.jpg", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.json", "application/json"},
		{"test.unknown", "application/octet-stream"},
		{"noextension", "application/octet-stream"},
	}

	for _, test := range tests {
		result := GetMimeType(test.filename)
		if result != test.expected {
			t.Errorf("GetMimeType(%q) = %q, expected %q", test.filename, result, test.expected)
		}
	}
}

func TestContainsInt(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if !ContainsInt(slice, 3) {
		t.Error("Expected slice to contain 3")
	}

	if ContainsInt(slice, 6) {
		t.Error("Expected slice to not contain 6")
	}

	if ContainsInt([]int{}, 1) {
		t.Error("Expected empty slice to not contain 1")
	}
}

func TestContainsString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	if !ContainsString(slice, "banana") {
		t.Error("Expected slice to contain 'banana'")
	}

	if ContainsString(slice, "grape") {
		t.Error("Expected slice to not contain 'grape'")
	}

	if ContainsString([]string{}, "apple") {
		t.Error("Expected empty slice to not contain 'apple'")
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"Hello World", 5, "Hello..."},
		{"Short", 10, "Short"},
		{"", 5, ""},
		{"Exactly5", 8, "Exactly5"},
		{"TruncateThis", 8, "Truncate..."},
	}

	for _, test := range tests {
		result := TruncateString(test.input, test.length)
		if result != test.expected {
			t.Errorf("TruncateString(%q, %d) = %q, expected %q", test.input, test.length, result, test.expected)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.org", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"test@.com", false},
		{"", false},
	}

	for _, test := range tests {
		result := ValidateEmail(test.email)
		if result != test.expected {
			t.Errorf("ValidateEmail(%q) = %v, expected %v", test.email, result, test.expected)
		}
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"validuser", true},
		{"user123", true},
		{"user_name", true},
		{"user-name", true},
		{"ab", false},                    // Too short
		{"", false},                      // Empty
		{strings.Repeat("a", 51), false}, // Too long
		{"user name", false},             // Space not allowed
		{"user@name", false},             // @ not allowed
	}

	for _, test := range tests {
		result := ValidateUsername(test.username)
		if result != test.expected {
			t.Errorf("ValidateUsername(%q) = %v, expected %v", test.username, result, test.expected)
		}
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal_file-name.txt", "normal_file-name.txt"},
		{"file with spaces.txt", "file_with_spaces.txt"},
		{"file@#$%^&*().txt", "file_.txt"},
		{"multiple___underscores", "multiple_underscores"},
		{"___leading_trailing___", "leading_trailing"},
		{"", ""},
	}

	for _, test := range tests {
		result := SanitizeFilename(test.input)
		if result != test.expected {
			t.Errorf("SanitizeFilename(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{0, "0 B"},
	}

	for _, test := range tests {
		result := FormatFileSize(test.bytes)
		if result != test.expected {
			t.Errorf("FormatFileSize(%d) = %q, expected %q", test.bytes, result, test.expected)
		}
	}
}

func TestStringSliceContains(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	if !StringSliceContains(slice, "banana") {
		t.Error("Expected slice to contain 'banana'")
	}

	if StringSliceContains(slice, "grape") {
		t.Error("Expected slice to not contain 'grape'")
	}

	// Test with empty slice
	if StringSliceContains([]string{}, "apple") {
		t.Error("Expected empty slice to not contain 'apple'")
	}
}

func TestRemoveStringFromSlice(t *testing.T) {
	tests := []struct {
		slice    []string
		remove   string
		expected []string
	}{
		{[]string{"a", "b", "c"}, "b", []string{"a", "c"}},
		{[]string{"a", "b", "c"}, "d", []string{"a", "b", "c"}},
		{[]string{"a"}, "a", []string{}},
		{[]string{}, "a", []string{}},
		{[]string{"a", "a", "b"}, "a", []string{"b"}}, // Remove all occurrences
	}

	for _, test := range tests {
		result := RemoveStringFromSlice(test.slice, test.remove)
		if len(result) != len(test.expected) {
			t.Errorf("RemoveStringFromSlice length mismatch: got %d, expected %d", len(result), len(test.expected))
			continue
		}

		for i, v := range result {
			if v != test.expected[i] {
				t.Errorf("RemoveStringFromSlice(%v, %q) = %v, expected %v", test.slice, test.remove, result, test.expected)
				break
			}
		}
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		remoteAddr    string
		xForwardedFor string
		xRealIP       string
		expected      string
	}{
		{"192.168.1.1:8080", "", "10.0.0.1", "10.0.0.1"},                    // X-Real-IP takes priority
		{"192.168.1.1:8080", "203.0.113.1,198.51.100.1", "", "203.0.113.1"}, // X-Forwarded-For first IP
		{"192.168.1.1:8080", "", "", "192.168.1.1"},                         // Remote addr without port
		{"", "", "", "unknown"},                                             // No IP available
		{"192.168.1.1", "", "", "192.168.1.1"},                              // Remote addr without port
	}

	for _, test := range tests {
		result := GetClientIP(test.remoteAddr, test.xForwardedFor, test.xRealIP)
		if result != test.expected {
			t.Errorf("GetClientIP(%q, %q, %q) = %q, expected %q",
				test.remoteAddr, test.xForwardedFor, test.xRealIP, result, test.expected)
		}
	}
}

func TestGenerateFileHash(t *testing.T) {
	// Create a test file in memory
	content := "test file content"
	_ = bytes.NewBuffer([]byte(content))

	// Create a multipart.File mock
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("file", "test.txt")
	fw.Write([]byte(content))
	w.Close()

	reader := bytes.NewReader(b.Bytes())
	r := multipart.NewReader(reader, w.Boundary())

	part, err := r.NextPart()
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hash, err := GenerateFileHash(part)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	if len(hash) != 64 { // SHA256 produces 64 character hex string
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}
}

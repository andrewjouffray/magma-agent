package parsing

import (
	"os"
	"testing"
)

func TestReadMagmaFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data to the temporary file
	content := `# This is a comment
line1
line2

# Another comment
line3`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Call the function with the path to the temporary file
	lines, err := ReadMagmaFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadMagmaFile returned an error: %v", err)
	}

	// Verify the output
	expectedLines := []string{"line1", "line2", "line3"}
	if len(lines) != len(expectedLines) {
		t.Fatalf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Expected line %d to be %q, got %q", i, expectedLines[i], line)
		}
	}
}

func TestWriteTrack(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "trackfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Test data to write
	lines := []string{"line1", "line2", "line3"}

	// Call the function with the path to the temporary file
	err = WriteTrack(lines, tmpFile.Name())
	if err != nil {
		t.Fatalf("WriteTrack returned an error: %v", err)
	}

	// Read the file content to verify
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	// Verify the output
	expectedContent := "line1\nline2\nline3\n"
	if string(content) != expectedContent {
		t.Errorf("Expected file content to be %q, got %q", expectedContent, string(content))
	}
}

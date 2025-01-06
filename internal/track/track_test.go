package track

import (
	"os"
	"testing"
)

func TestAddPath(t *testing.T) {
	// Create a temporary file to act as the track file
	tmpFile, err := os.CreateTemp("", "trackfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Test case: Add a new path that does not exist
	newPath := "/new/path"
	err = AddPath(newPath, tmpFile.Name())
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Test case: Add a new path that does exist
	tmpFile2, err := os.CreateTemp("", "test_to_track")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())
	newPath = tmpFile2.Name()
	err = AddPath(newPath, tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the path was added
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	if string(content) != newPath+"\n" {
		t.Errorf("Expected %s, got %s", newPath, string(content))
	}

	// Test case: Add the same path again
	err = AddPath(newPath, tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the path was not duplicated
	content, err = os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	if string(content) != newPath+"\n" {
		t.Errorf("Expected %s, got %s", newPath, string(content))
	}

	// Test case: Add a non-existent path
	nonExistentPath := "/non/existent/path"
	err = AddPath(nonExistentPath, tmpFile.Name())
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRemovePath(t *testing.T) {
	// Create a temporary file to act as the track file
	tmpFile, err := os.CreateTemp("", "trackfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Add a path to the track file
	tmpFile2, err := os.CreateTemp("", "test_to_track")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())
	newPath := tmpFile2.Name()
	err = AddPath(newPath, tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to add path: %v", err)
	}

	// Test case: Remove the path
	err = RemovePath(newPath, tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the path was removed
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	if string(content) != "" {
		t.Errorf("Expected empty file, got %s", string(content))
	}

	// Test case: Remove a non-existent path
	nonExistentPath := "/non/existent/path"
	err = RemovePath(nonExistentPath, tmpFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

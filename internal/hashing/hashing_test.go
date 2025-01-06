package hashing

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHashFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some data to the file
	data := []byte("hello world")
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Calculate the expected hash
	hasher := sha256.New()
	hasher.Write(data)
	expectedHash := fmt.Sprintf("%x", hasher.Sum(nil))

	// Call the hashFile function
	hash, err := hashFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("hashFile returned an error: %v", err)
	}

	// Check if the hash matches the expected hash
	if hash != expectedHash {
		t.Errorf("hashFile returned %s, expected %s", hash, expectedHash)
	}
}

func TestHashFile_FileNotFound(t *testing.T) {
	// Call the hashFile function with a non-existent file
	_, err := hashFile("non_existent_file.txt")
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestHashNodeList(t *testing.T) {
	nodes := []Node{
		{Hash: "abc123"},
		{Hash: "def456"},
		{Hash: "ghi789"},
	}

	expectedHash := hashString("abc123def456ghi789")

	hash := hashNodeList(nodes)

	if hash != expectedHash {
		t.Errorf("hashNodeList returned %s, expected %s", hash, expectedHash)
	}
}

func TestHashNodeList_EmptyNodes(t *testing.T) {
	nodes := []Node{}

	expectedHash := hashString("")

	hash := hashNodeList(nodes)

	if hash != expectedHash {
		t.Errorf("hashNodeList returned %s, expected %s", hash, expectedHash)
	}
}

func TestHashPath_File(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some data to the file
	data := []byte("hello world")
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Call the HashPath function
	node, err := HashPath(tmpfile.Name())
	if err != nil {
		t.Fatalf("HashPath returned an error: %v", err)
	}

	// Calculate the expected hash
	hasher := sha256.New()
	hasher.Write(data)
	expectedHash := fmt.Sprintf("%x", hasher.Sum(nil))

	// Check if the hash matches the expected hash
	if node.Hash != expectedHash {
		t.Errorf("HashPath returned %s, expected %s", node.Hash, expectedHash)
	}
}

func TestHashPath_Directory(t *testing.T) {
	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir) // clean up

	// Create a temporary file in the directory
	tmpfile, err := os.CreateTemp(tmpdir, "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some data to the file
	data := []byte("hello world")
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Call the HashPath function
	node, err := HashPath(tmpdir)
	if err != nil {
		t.Fatalf("HashPath returned an error: %v", err)
	}

	// Check if the directory hash is not empty
	if node.Hash == "" {
		t.Errorf("HashPath returned an empty hash for the directory")
	}
}

func TestHashPath_Symlink(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Create a temporary symlink
	symlink := tmpfile.Name() + "_symlink"
	if err := os.Symlink(tmpfile.Name(), symlink); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(symlink) // clean up

	// Call the HashPath function
	node, err := HashPath(symlink)
	if err != nil {
		t.Fatalf("HashPath returned an error: %v", err)
	}

	// Calculate the expected hash
	expectedHash := hashString(tmpfile.Name())

	// Check if the hash matches the expected hash
	if node.Hash != expectedHash {
		t.Errorf("HashPath returned %s, expected %s", node.Hash, expectedHash)
	}
}

func TestHashPath_FileNotFound(t *testing.T) {
	// Call the HashPath function with a non-existent file
	_, err := HashPath("non_existent_file.txt")
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestHashPath_IgnoreList_StartsWith(t *testing.T) {
	// Add a pattern to the ignore list
	ignoreList = append(ignoreList, "**/ignore_*")

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "ignore_me.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Call the HashPath function
	node, err := HashPath(tmpfile.Name())
	if err != nil {
		t.Fatalf("HashPath returned an error: %v", err)
	}

	// Check if the file was skipped
	if node.Hash != "skipped" {
		t.Errorf("HashPath did not skip the file, returned %s", node.Hash)
	}
}

func TestHashPath_IgnoreList_EndsWith(t *testing.T) {

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example_*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Add a pattern to the ignore list
	ignoreList = append(ignoreList, "**/*.log")

	// Call the HashPath function
	node, err := HashPath(tmpfile.Name())
	if err != nil {
		t.Fatalf("HashPath returned an error: %v", err)
	}

	// Check if the file was skipped
	if node.Hash != "skipped" {
		t.Errorf("HashPath did not skip the file, returned %s", node.Hash)
	}
}

func TestHashPath_IgnoreList_Directory(t *testing.T) {

	// Create a temporary directory
	tmpdir := t.TempDir()

	// Add a pattern to the ignore list
	ignoreList = append(ignoreList, tmpdir+"/**")
	defer os.RemoveAll(tmpdir) // clean up

	// Create a temporary file in the directory
	tmpfile, err := os.CreateTemp(tmpdir, "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Call the HashPath function
	node, err := HashPath(tmpdir)
	if err != nil {
		t.Fatalf("HashPath returned an error: %v", err)
	}

	// Check if the directory was skipped
	if node.Hash != "skipped" {
		t.Errorf("HashPath did not skip the directory, returned %s", node.Hash)
	}
}

func TestHashPath_IgnoreList_Subdirectory(t *testing.T) {
	// Add a pattern to the ignore list
	ignoreList = append(ignoreList, "**/ignore_subdir/*")

	// Create a temporary directory
	tmpdir, err := os.MkdirTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir) // clean up

	// Create a subdirectory
	subdir := filepath.Join(tmpdir, "ignore_subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a temporary file in the subdirectory
	tmpfile, err := os.CreateTemp(subdir, "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Call the HashPath function
	node, err := HashPath(subdir)
	if err != nil {
		t.Fatalf("HashPath returned an error: %v", err)
	}

	// Check if the subdirectory was skipped
	if node.Hash != "skipped" {
		t.Errorf("HashPath did not skip the subdirectory, returned %s", node.Hash)
	}
}

func TestSnapShot(t *testing.T) {
	// Create a temporary directory for the snapshot
	snapshotDir, err := os.MkdirTemp("", "snapshot")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(snapshotDir) // clean up

	// Create a temporary file to track
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some data to the file
	data := []byte("hello world")
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Call the SnapShot function
	err = SnapShot(snapshotDir, []string{tmpfile.Name()})
	if err != nil {
		t.Fatalf("SnapShot returned an error: %v", err)
	}

	// Check if the snapshot file was created
	files, err := os.ReadDir(snapshotDir)
	if err != nil {
		t.Fatalf("Failed to read snapshot directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No snapshot file created")
	}

	// Check if the snapshot file has the correct content
	snapshotFile := filepath.Join(snapshotDir, files[0].Name())
	content, err := os.ReadFile(snapshotFile)
	if err != nil {
		t.Fatalf("Failed to read snapshot file: %v", err)
	}

	var root Node
	if err := json.Unmarshal(content, &root); err != nil {
		t.Fatalf("Failed to unmarshal snapshot JSON: %v", err)
	}

	if root.Path != "root" {
		t.Errorf("Expected root path to be 'root', got %s", root.Path)
	}

	if len(root.Children) != 1 {
		fmt.Println(root.Children)
		t.Errorf("Expected 1 child node, got %d", len(root.Children))
	}
}

func TestSnapShot_WithTags(t *testing.T) {
	// Create a temporary directory for the snapshot
	snapshotDir, err := os.MkdirTemp("", "snapshot")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(snapshotDir) // clean up

	// Create a temporary file to track
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some data to the file
	data := []byte("hello world")
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Call the SnapShot function with tags
	err = SnapShot(snapshotDir, []string{tmpfile.Name()}, "tag1", "tag2")
	if err != nil {
		t.Fatalf("SnapShot returned an error: %v", err)
	}

	// Check if the snapshot file was created
	files, err := os.ReadDir(snapshotDir)
	if err != nil {
		t.Fatalf("Failed to read snapshot directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No snapshot file created")
	}

	// Check if the snapshot file name contains the tags
	snapshotFile := files[0].Name()
	if !strings.Contains(snapshotFile, "tag1") || !strings.Contains(snapshotFile, "tag2") {
		t.Errorf("Snapshot file name does not contain tags, got %s", snapshotFile)
	}
}

func TestSnapShot_EmptyTrackPaths(t *testing.T) {
	// Create a temporary directory for the snapshot
	snapshotDir, err := os.MkdirTemp("", "snapshot")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(snapshotDir) // clean up

	// Call the SnapShot function with empty track paths
	err = SnapShot(snapshotDir, []string{})
	if err != nil {
		t.Fatalf("SnapShot returned an error: %v", err)
	}

	// Check if the snapshot file was created
	files, err := os.ReadDir(snapshotDir)
	if err != nil {
		t.Fatalf("Failed to read snapshot directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No snapshot file created")
	}

	// Check if the snapshot file has the correct content
	snapshotFile := filepath.Join(snapshotDir, files[0].Name())
	content, err := os.ReadFile(snapshotFile)
	if err != nil {
		t.Fatalf("Failed to read snapshot file: %v", err)
	}

	var root Node
	if err := json.Unmarshal(content, &root); err != nil {
		t.Fatalf("Failed to unmarshal snapshot JSON: %v", err)
	}

	if root.Path != "root" {
		t.Errorf("Expected root path to be 'root', got %s", root.Path)
	}

	if len(root.Children) != 0 {
		t.Errorf("Expected 0 child nodes, got %d", len(root.Children))
	}
}

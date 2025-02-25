package hashing

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"magma/internal/config"
	"magma/internal/parsing"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

// Node represents a node in the file tree, the root node is returned by the HashPath function
// and can be parsed into JSON.
type Node struct {
	Path     string `json:"path"`     // The path of the file or directory
	Hash     string `json:"hash"`     // The hash value of this node
	Children []Node `json:"children"` // Child nodes
}

var ignoreList []string

func init() {
	var err error
	ignoreList, err = parsing.ReadMagmaFile(config.IgnoreFile)
	if err != nil {
		fmt.Println("Error reading ignore file:", err)
	}
}

// hashFile computes the SHA-256 hash of the file at the given filepath.
// It returns the hash as a hexadecimal string or an error if any occurs during the process.
//
// Parameters:
//   - filepath: The path to the file to be hashed.
//
// Returns:
//   - string: The hexadecimal representation of the file's SHA-256 hash.
//   - error: An error if the file cannot be opened or read.
func hashFile(filepath string) (string, error) {
	// opens the file at the given path
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	// ensures the file is closed after the function returns regardless of the outcome
	defer file.Close()

	hasher := sha256.New()

	// write bytes from the file to the hasher and checks for any errors
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// returns the hash as a string
	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	return hash, nil
}

// hashNodeList takes a slice of Node objects and a path string, concatenates
// the Hash field of each Node into a single string, and returns the hash of
// the concatenated string.
//
// Parameters:
//
//	nodes - a slice of Node objects, each containing a Hash field
//
// Returns:
//
//	A string representing the hash of the concatenated Hash fields of the input nodes.
func hashNodeList(nodes []Node) string {

	// concatenates all strings in the input list
	concatenated := ""
	for _, node := range nodes {
		concatenated += node.Hash
	}

	// hashes the concatenated string
	hash := hashString(concatenated)
	return hash

}

func hashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// recursively hashes all files in the given directory and subdirectories
// HashPath computes a hash for the given file or directory path.
// It returns a Node struct containing the hash and any child nodes.
//
// If the path is a symlink, it resolves the symlink and hashes the resolved path.
// If the path is a directory, it recursively hashes all files and directories within it.
// If the path is a file, it hashes the file content.
//
// The function also checks if the path is in the ignore list and skips hashing if it is.
//
// Parameters:
//   - path: The file or directory path to hash.
//
// Returns:
//   - Node: A Node struct containing the hash and any child nodes.
//   - error: An error if any occurred during hashing.
func HashPath(path string) (node Node, error error) {

	var localNode Node

	// check if the path is a file
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return localNode, err
	}

	// check if the path is in the ignore list
	for _, ignore := range ignoreList {
		if match, _ := doublestar.Match(ignore, path); match {
			localNode.Hash = "skipped"
			return localNode, nil
		}
	}

	// Check if the path is a symlink
	if fileInfo.Mode()&os.ModeSymlink != 0 {

		// Handle the symlink
		resolvedPath, err := os.Readlink(path)
		if err != nil {
			return localNode, fmt.Errorf("failed to resolve symlink: %w", err)
		}
		// If the resolved path is relative, make it absolute based on the symlink's parent directory
		if !filepath.IsAbs(resolvedPath) {
			resolvedPath = filepath.Join(filepath.Dir(path), resolvedPath)
		}

		// hash the resolved path string
		hash := hashString(resolvedPath)
		localNode.Hash = hash
		localNode.Children = nil
		localNode.Path = path
		return localNode, nil

	}

	if fileInfo.IsDir() {
		var nodes []Node
		// get all files and directories in the given path
		files, err := os.ReadDir(path)
		if err != nil {
			return localNode, err
		}
		for _, file := range files {
			// recursively hash all files in the directory
			child, err := HashPath(path + "/" + file.Name())
			if err != nil {
				return localNode, err
			}
			nodes = append(nodes, child)
		}
		// hash all the hashes of the files in the directory
		nodeHash := hashNodeList(nodes)
		localNode.Hash = nodeHash
		localNode.Children = nodes
		localNode.Path = path
		return localNode, nil

	} else if fileInfo.Mode()&os.ModeType != 0 {
		localNode.Hash = "skipped"
		return localNode, nil

	}

	// hash the file
	hash, err := hashFile(path)
	if err != nil {
		return localNode, err
	}
	localNode.Hash = hash
	localNode.Children = nil
	localNode.Path = path
	return localNode, nil

}

// SnapShot creates a snapshot of the given tracked paths and saves it as a JSON file.
// The snapshot includes the hash of each tracked path and their hierarchical structure.
//
// Parameters:
//   - SnapshotPath: The directory where the snapshot JSON file will be saved.
//   - trackPaths: A list of paths to be tracked and hashed.
//   - tags: Optional tags to be appended to the snapshot filename.
//
// Returns:
//   - error: An error if any occurs during the snapshot creation or file writing process.
func SnapShot(SnapshotPath string, trackPaths []string, tags ...string) error {

	// for each tracked path, create a root node
	nodes := []Node{}

	for _, path := range trackPaths {
		node, err := HashPath(path)
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
	}

	root := Node{
		Path:     "root",
		Hash:     hashNodeList(nodes),
		Children: nodes,
	}

	jsonData, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	// truncate the hash to 8 characters
	fileName := root.Hash
	if len(fileName) > 8 {
		fileName = fileName[:8]
	}

	// if there are tags, append them to the snapshot filename
	if len(tags) > 0 {
		for _, tag := range tags {
			fileName += "_" + tag
		}
	}

	savePath := filepath.Join(SnapshotPath, fileName)
	savePath = savePath + ".json"

	// Write the JSON to a file
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Println("Snapshot saved to", savePath)

	return nil
}

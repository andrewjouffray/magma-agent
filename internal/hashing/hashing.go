package hashing

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type Node struct {
	Hash     string `json:"hash"`     // The hash value of this node
	Children []Node `json:"children"` // Child nodes
}

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
	fmt.Println(hash)
	return hash, nil
}

func hashNodeList(nodes []Node) string {

	// concatenates all strings in the input list
	concatenated := ""
	for _, node := range nodes {
		concatenated += node.Hash
	}

	// hashes the concatenated string
	hasher := sha256.New()
	hasher.Write([]byte(concatenated))
	hash := fmt.Sprintf("%x", hasher.Sum(nil))
	fmt.Println(hash)
	return hash

}

// recursively hashes all files in the given directory and subdirectories
func HashPath(path string) (node Node, error error) {

	var localNode Node

	// check if the path is a file
	fileInfo, err := os.Stat(path)
	if err != nil {
		return localNode, err
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
		return localNode, nil

	} else {
		// hash the file
		hash, err := hashFile(path)
		if err != nil {
			return localNode, err
		}
		localNode.Hash = hash
		localNode.Children = nil
		return localNode, nil
	}

}

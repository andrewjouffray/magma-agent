package main

import (
	"encoding/json"
	"fmt"
	"magma/internal/hashing"
	"magma/internal/initialize"
	"os"
)

func main() {

	asciiArt := `
                                       
    _____ _____     ____   _____ _____   
   /     \\__  \   / ___\ /     \\__  \  
  |  Y Y  \/ __ \_/ /_/  >  Y Y  \/ __ \_
  |__|_|  (____  /\___  /|__|_|  (____  /
	\/     \//_____/       \/     \/ 
  `
	fmt.Println(asciiArt)

	// Ensure at least one positional argument (command) is provided
	if len(os.Args) < 1 {
		fmt.Println("please provide a command (e.g., hash)")
		return
	}

	// Get the command
	command := os.Args[1]

	switch {
	case command == "hash":
		// Hash the files in the given directory
		if len(os.Args) < 2 {
			fmt.Println("Please provide a path to hash, e.g., magma hash /path/to/directory")
			return
		}
		path := os.Args[2]
		fmt.Println("Hashing files in", path, "...")
		root, err := hashing.HashPath(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonData, err := json.MarshalIndent(root, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		// fmt.Println(string(jsonData))

		// Write the JSON to a file
		file, err := os.Create("output.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Println("JSON written to output.json")
	case command == "init":
		// Initialize the magma directory
		err := initialize.Initialize()
		if err != nil {
			fmt.Println("Error initializing magma:", err)
			return
		}
	default:
		fmt.Println("Unknown command ", os.Args[1])
	}
}

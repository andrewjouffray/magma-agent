package main

import (
	"fmt"
	"magma/internal/config"
	"magma/internal/hashing"
	"magma/internal/initialize"
	"magma/internal/parsing"
	"magma/internal/track"
	"os"
)

// main is the entry point of the magma-agent application. It displays an ASCII art banner,
// checks for at least one positional argument (command), and executes the corresponding
// command. Supported commands are:
// - "snap [tag1] [tag2] ...": Creates a new cryptographic snapshot for all tracked files and directories.
// - "track [path]": Adds a new path to the track file.
// - "untrack [path]": Removes a path from the track file.
// - "init": Initializes the magma directory.
// If an unknown command is provided, it prints an error message.
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
		fmt.Println("please provide a command (e.g., snap)")
		return
	}

	// Get the command
	command := os.Args[1]

	switch {

	// creates a new cryptographic snapshot for all tracked files and directories
	case command == "snap":

		// Get the paths to track
		trackPaths, err := parsing.ReadMagmaFile(config.TrackFile)
		if err != nil {
			fmt.Println("Error reading magma file:", err)
			return
		}

		if len(trackPaths) == 0 {
			fmt.Println("No paths to track")
			return
		}

		// get all the optional tags for the snapshot
		tags := os.Args[2:]

		// Create a snapshot
		err = hashing.SnapShot(config.SnapshotsDir, trackPaths, tags...)
		if err != nil {
			fmt.Println("Error creating snapshot:", err)
			return
		}

	case command == "track":

		// Ensure at least one positional argument (path) is provided
		if len(os.Args) < 3 {
			fmt.Println("please provide a path to track")
			return
		}

		// Get the path to track
		path := os.Args[2]

		// Track the path
		err := track.AddPath(path, config.TrackFile)
		if err != nil {
			fmt.Println("Error tracking path:", err)
			return
		}

		fmt.Println("Path tracked")

	case command == "untrack":

		// Ensure at least one positional argument (path) is provided
		if len(os.Args) < 3 {
			fmt.Println("please provide a path to untrack")
			return
		}

		// Get the path to untrack
		path := os.Args[2]

		// Untrack the path
		err := track.RemovePath(path, config.TrackFile)
		if err != nil {
			fmt.Println("Error untracking path:", err)
			return
		}

		fmt.Println("Path untracked")

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

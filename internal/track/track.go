package track

import (
	"fmt"
	"magma/internal/config"
	"magma/internal/parsing"
	"os"
)

// adds new path to the track file
// AddPath adds a new path to the track file if it does not already exist.
//
// Parameters:
//   - newPath: The new path to be added.
//   - trackFilePath: The path to the track file.
//
// Returns:
//   - error: An error if the new path does not exist, if there is an error reading the track file,
//     or if there is an error writing to the track file. Returns nil if the path is successfully added
//     or if the path already exists in the track file.
func AddPath(newPath string, trackFilePath string) error {

	// check if the new path exists
	_, err := os.Stat(newPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("path %s does not exist", newPath)
	}

	// read the current paths from the track file
	currentPaths, err := parsing.ReadMagmaFile(config.TrackFile)
	if err != nil {
		return err
	}

	// check if the new path is already in the track file
	for _, path := range currentPaths {
		if path == newPath {
			println("Path already exists in the track file")
			return nil
		}
	}

	// append the new path to the current paths
	currentPaths = append(currentPaths, newPath)

	// write the new paths to the track file
	err = parsing.WriteTrack(currentPaths, trackFilePath)
	if err != nil {
		return err
	}

	return nil
}

// removes path from the track file
// RemovePath removes a specified path from the tracking file.
//
// It reads the current paths from the tracking file, checks if the path to be removed
// is present, removes it if found, and then writes the updated paths back to the tracking file.
//
// Parameters:
//   - pathToRemove: The path that needs to be removed from the tracking file.
//   - trackFilePath: The file path of the tracking file.
//
// Returns:
//   - error: An error if there is an issue reading from or writing to the tracking file, otherwise nil.
func RemovePath(pathToRemove string, trackFilePath string) error {

	// read the current paths from the track file
	currentPaths, err := parsing.ReadMagmaFile(trackFilePath)
	if err != nil {
		return err
	}

	fmt.Println(currentPaths)
	fmt.Println(pathToRemove)

	// check if the path to remove is in the track file and remove it
	for i, path := range currentPaths {
		if path == pathToRemove {
			currentPaths = append(currentPaths[:i], currentPaths[i+1:]...)
			break
		}
	}

	fmt.Println(currentPaths)

	// write the new paths to the track file
	err = parsing.WriteTrack(currentPaths, trackFilePath)
	if err != nil {
		return err
	}

	return nil
}

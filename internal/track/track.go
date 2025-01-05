package track

import (
	"fmt"
	"magma/internal/config"
	"magma/internal/parsing"
	"os"
)

// adds new path to the track file
func AddPath(newPath string) error {

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
	err = parsing.WriteTrack(currentPaths, config.TrackFile)
	if err != nil {
		return err
	}

	return nil
}

// removes path from the track file
func RemovePath(pathToRemove string) error {

	// check if the new path exists
	_, err := os.Stat(pathToRemove)
	if os.IsNotExist(err) {
		return fmt.Errorf("path %s does not exist", pathToRemove)
	}

	// read the current paths from the track file
	currentPaths, err := parsing.ReadMagmaFile(config.TrackFile)
	if err != nil {
		return err
	}

	// check if the new path is already in the track file
	for _, path := range currentPaths {
		if path == pathToRemove {
			println("Path already exists in the track file")
			return nil
		}
	}

	// append the new path to the current paths
	currentPaths = append(currentPaths, pathToRemove)

	// write the new paths to the track file
	err = parsing.WriteTrack(currentPaths, config.TrackFile)
	if err != nil {
		return err
	}

	return nil
}

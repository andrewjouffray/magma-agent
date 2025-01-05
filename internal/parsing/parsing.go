package parsing

import (
	"bufio"
	"fmt"
	"os"
)

// ReadMagmaFile reads a file from the given path and returns a slice of strings,
// each representing a line from the file. It skips empty lines and lines that
// start with a '#' character (comments).
//
// Parameters:
//   - path: The path to the file to be read.
//
// Returns:
//   - []string: A slice of strings containing the
func ReadMagmaFile(path string) ([]string, error) {
	// open the ignore file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		// skip empty lines and comments
		if scanner.Text() == "" || scanner.Text()[0] == '#' {
			continue
		}
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

// WriteTrack writes a slice of strings to a specified file, each string on a new line.
// The file is opened in append mode, so new lines are added to the end of the file.
//
// Parameters:
//
//	lines - a slice of strings to be written to the file
//	trackFilePath - the path to the file where the lines will be written
//
// Returns:
//
//	error - an error if there is an issue opening the file, writing to it, or flushing the buffer
func WriteTrack(lines []string, trackFilePath string) error {
	// open the track file in append mode
	file, err := os.OpenFile(trackFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open track file: %w", err)
	}
	defer file.Close()

	// create a writer
	writer := bufio.NewWriter(file)

	for _, path := range lines {
		line := path + "\n" // Append newline directly
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write to track file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}

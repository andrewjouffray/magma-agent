package parsing

import (
	"bufio"
	"fmt"
	"os"
)

// provices functions to read and parse the ignore and track files
func ReadIgnore() ([]string, error) {
	// open the ignore file
	file, err := os.Open("/etc/magma/ignore")
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

func ReadTrack() ([]string, error) {
	// open the track file
	file, err := os.Open("/etc/magma/track")
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

func WriteTrack(paths []string) error {
	// open the track file in append mode
	file, err := os.OpenFile("/etc/magma/track", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open track file: %w", err)
	}
	defer file.Close()

	// create a writer
	writer := bufio.NewWriter(file)

	for _, path := range paths {
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

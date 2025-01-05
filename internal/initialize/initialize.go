package initialize

import (
	"bufio"
	"fmt"
	"os"
)

// initializes the /etc/magma directory, track file and snapshots directory
func Initialize() error {
	// check the existence of the /etc/magma directory
	_, err := os.Stat("/etc/magma")
	if os.IsNotExist(err) {
		// create the /etc/magma directory
		err = os.Mkdir("/etc/magma", 0755)
		if err != nil {
			return err
		}
	}

	// check the existence of the track file
	_, err = os.Stat("/etc/magma/track")
	if os.IsNotExist(err) {
		// create the track file
		_, err = os.Create("/etc/magma/track")
		if err != nil {
			return err
		}
	}

	// check the existence of the snapshots directory
	_, err = os.Stat("/etc/magma/snapshots")
	if os.IsNotExist(err) {
		// create the snapshots directory
		err = os.Mkdir("/etc/magma/snapshots", 0755)
		if err != nil {
			return err
		}
	}

	// check the existence of the /etc/magma/ignore file
	_, err = os.Stat("/etc/magma/ignore")
	if os.IsNotExist(err) {
		// create the /etc/magma/ignore file
		fileName := "/etc/magma/ignore"

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return err
		}
		defer file.Close()

		// Create a writer
		writer := bufio.NewWriter(file)

		lines := []string{
			"# self directory",
			"/etc/magma",
			"# any hidden files or directories",
			".*",
		}

		// Write each line to the file
		for _, line := range lines {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return err
			}
		}

		// Flush the writer to ensure all data is written to the file
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer:", err)
			return err
		}

		// Print a success message
		print("Successfully initialized magma directory\n")
		fmt.Println("Ignore file created at", fileName)
		fmt.Println("By default, the following paths are ignored:")
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	return nil
}

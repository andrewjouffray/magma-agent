package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"magma/internal/hashing"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		fmt.Println("First argument:", os.Args[1])
	} else {
		fmt.Println("No arguments provided.")
	}

	// hashes the files in the given directory
	path := flag.String("path", "/", "path to compute hashes")

	flag.Parse()

	switch {
	case os.Args[1] == "hash":
		// Hash the files in the given directory
		root, err := hashing.HashPath(*path)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonData, err := json.MarshalIndent(root, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonData))

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
	}
}

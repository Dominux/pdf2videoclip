package internal

import (
	"fmt"
	"log"
	"os"
)

func ReadFile() {
	dir, err := os.Getwd()
	if err != nil {
		// Log the error and exit if getting the directory fails
		log.Fatal(err)
	}
	// Print the working directory
	fmt.Println(dir)

	content, err := os.ReadFile("./example.pdf") // Use os.ReadFile as of Go 1.16+
	if err != nil {
		log.Fatal(err) // Handle potential errors
	}
	fmt.Println(string(content))
}

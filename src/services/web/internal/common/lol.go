package common

import (
	"fmt"
	"log"
	"os"
)

func ReadFile() {
	content, err := os.ReadFile("./example.pdf") // Use os.ReadFile as of Go 1.16+
	if err != nil {
		log.Fatal(err) // Handle potential errors
	}
	fmt.Println(string(content))
}

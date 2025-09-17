package pkg

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
)

func HashCSS() string {
	content, err := os.ReadFile("static/styles/output.css")
	if err != nil {
		log.Fatal(err)
	}

	sum := sha256.Sum256(content)
	hash := fmt.Sprintf("%x", sum)[:8] // take first 8 characters
	return hash
}

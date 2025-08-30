package pkg

import (
	"strings"
)

// SlugToTitle converts "novel-name" -> "Novel Name"
func SlugToTitle(slug string) string {
	words := strings.Split(slug, "-")
	for i, word := range words {
		if len(word) > 0 {
			// Uppercase first letter, lowercase rest
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// TitleToSlug converts "Novel Name" -> "novel-name"
func TitleToSlug(title string) string {
	// Trim and normalize spaces
	title = strings.TrimSpace(title)
	// Replace multiple spaces with one
	title = strings.Join(strings.Fields(title), " ")
	// Lowercase and replace spaces with hyphens
	return strings.ReplaceAll(strings.ToLower(title), " ", "-")
}

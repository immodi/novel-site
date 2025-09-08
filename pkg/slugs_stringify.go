package pkg

import (
	"strings"
	"unicode"
)

// SlugToTitle converts a slug like "novel-name" -> "Novel Name"
// Works correctly with Unicode
func SlugToTitle(slug string) string {
	words := strings.Split(slug, "-")
	for i, word := range words {
		if len(word) == 0 {
			continue
		}
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// TitleToSlug converts "Novel Name" -> "novel-name"
// Preserves non-Latin characters and removes unwanted symbols
func TitleToSlug(title string) string {
	title = strings.TrimSpace(title)
	words := strings.Fields(title)
	var slugParts []string
	for _, word := range words {
		var clean []rune
		for _, r := range word {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				clean = append(clean, r)
			}
		}
		if len(clean) > 0 {
			slugParts = append(slugParts, strings.ToLower(string(clean)))
		}
	}
	return strings.Join(slugParts, "-")
}

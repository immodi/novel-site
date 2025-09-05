package pkg

import (
	"math/rand"
	"slices"
	"strings"
	"time"
)

func RandomTags(n int) []string {
	allTags := []string{
		"Academy",
		"Adapted to Anime",
		"Adapted to Manga",
		"Alternate World",
		"Anti-Magic",
		"Apathetic Protagonist",
		"Beautiful Female Lead",
		"Brother Complex",
		"Calm Protagonist",
		"Clever Protagonist",
		"Devoted Love Interests",
		"Discrimination",
		"Elemental Magic",
		"Engineer",
		"Eye Powers",
		"Familial Love",
		"Genius Protagonist",
		"Hard-Working Protagonist",
		"Incest",
		"Magic",
		"Magic Formations",
		"Magical Technology",
		"Male Protagonist",
		"Military",
		"Modern Day",
		"Nobles",
		"Overpowered Protagonist",
		"Politics",
		"Protagonist Strong from the Start",
		"Schemes And Conspiracies",
		"Scientists",
		"Sister Complex",
		"Special Abilities",
		"Spirits",
		"Stoic Characters",
		"Strong Love Interests",
		"Terrorists",
		"Tsundere",
		"Unconditional Love",
		"Unique Weapon User",
		"Wars",
	}

	if n <= 0 {
		return nil
	}
	if n > len(allTags) {
		n = len(allTags)
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	// shuffle copy
	shuffled := slices.Clone(allTags)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// return first n, lowercased
	result := make([]string, n)
	for i := range n {
		result[i] = strings.ToLower(shuffled[i])
	}

	return result
}

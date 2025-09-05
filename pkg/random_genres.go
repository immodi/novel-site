package pkg

import (
	sql "immodi/novel-site/internal/db/schema"
	"math/rand"
	"time"
)

func RandomGenres(n int) []string {
	if n > len(sql.AllGenres) {
		n = len(sql.AllGenres)
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	// make a copy to shuffle
	genres := make([]sql.Genre, len(sql.AllGenres))
	copy(genres, sql.AllGenres)

	rand.Shuffle(len(genres), func(i, j int) {
		genres[i], genres[j] = genres[j], genres[i]
	})

	result := make([]string, n)
	for i := range n {
		result[i] = string(genres[i])
	}

	return result
}

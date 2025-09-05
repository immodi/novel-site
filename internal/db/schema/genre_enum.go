package sql

import "slices"

type Genre string

const (
	GenreAction       Genre = "action"
	GenreAdventure    Genre = "adventure"
	GenreFantasy      Genre = "fantasy"
	GenreRomance      Genre = "romance"
	GenreMystery      Genre = "mystery"
	GenreSupernatural Genre = "supernatural"
	GenreDrama        Genre = "drama"
)

var AllGenres = []Genre{
	GenreAction,
	GenreAdventure,
	GenreFantasy,
	GenreRomance,
	GenreMystery,
	GenreSupernatural,
	GenreDrama,
}

func (g Genre) IsValidGenre() bool {
	return slices.Contains(AllGenres, g)
}

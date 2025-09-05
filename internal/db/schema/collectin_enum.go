package sql

import "slices"

type Collection string

const (
	CollectionHot       Collection = "hot"
	CollectionLatest    Collection = "latest"
	CollectionCompleted Collection = "completed"
	CollectionOnGoing   Collection = "ongoing"
)

var AllCollections = []Collection{
	CollectionHot,
	CollectionLatest,
	CollectionCompleted,
	CollectionOnGoing,
}

func (c Collection) IsValidCollection() bool {
	return slices.Contains(AllCollections, c)
}

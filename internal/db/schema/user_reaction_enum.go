package sql

type UserReaction string

const (
	UserReactionLike    UserReaction = "like"
	UserReactionDislike UserReaction = "dislike"
)

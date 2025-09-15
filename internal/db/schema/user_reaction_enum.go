package sql

import (
	"errors"
	"strings"
)

type UserReaction string

const (
	UserReactionLike    UserReaction = "like"
	UserReactionDislike UserReaction = "dislike"
)

func ParseUserReaction(s string) (UserReaction, error) {
	switch strings.ToLower(s) {
	case string(UserReactionLike):
		return UserReactionLike, nil
	case string(UserReactionDislike):
		return UserReactionDislike, nil
	default:
		return "", errors.New("invalid reaction: must be \"like\" or \"dislike\"")
	}
}

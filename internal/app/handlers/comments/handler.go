package comments

import (
	comments "immodi/novel-site/internal/app/services/comments"
	"immodi/novel-site/internal/app/services/profile"
)

type CommentHandler struct {
	commentService comments.CommentService
	profileService profile.ProfileService
}

func NewCommentHandler(service comments.CommentService, profileService profile.ProfileService) *CommentHandler {
	return &CommentHandler{commentService: service, profileService: profileService}
}

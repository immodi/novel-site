package chaptercomments

import (
	chaptercomments "immodi/novel-site/internal/app/services/chapter_comments"
	"immodi/novel-site/internal/app/services/profile"
)

type ChapterCommentHandler struct {
	commentService chaptercomments.ChapterCommentService
	profileService profile.ProfileService
}

func NewChapterCommentHandler(service chaptercomments.ChapterCommentService, profileService profile.ProfileService) *ChapterCommentHandler {
	return &ChapterCommentHandler{commentService: service, profileService: profileService}
}

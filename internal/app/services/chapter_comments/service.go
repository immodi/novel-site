package chaptercomments

import (
	c_sql "database/sql"
	"immodi/novel-site/internal/db/repositories"
	sql "immodi/novel-site/internal/db/schema"
)

type ChapterCommentService interface {
	// Comment CRUD
	GetCommentById(commentID int) (repositories.ChapterComment, error)
	CreateComment(chapterID, userID int, parentID *int, content string) (repositories.ChapterComment, error)
	GetCommentsByChapter(chapterID int) ([]repositories.GetCommentsByChapterRow, error)
	UpdateComment(commentID int, content string) (repositories.ChapterComment, error)
	GetCommentReplies(parentID c_sql.NullInt64) ([]repositories.GetRepliesByChapterCommentRow, error)
	DeleteComment(commentID, userID int) error

	// Reactions
	AddOrUpdateOrRemoveReaction(userID, commentID int, reaction sql.UserReaction) error
	RemoveReaction(userID, commentID int) error
	GetUserReaction(userID, commentID int) (string, error)
	GetUserReactionForComments(userID int, commentIDs []int64) ([]repositories.GetUserReactionsForChapterCommentsRow, error)
	CountReactions(commentID int) (repositories.CountChapterReactionsRow, error)

	GetCommentsByChapterWithUserReactions(userID, chapterID int) ([]repositories.GetCommentsByChapterWithUserReactionsRow, error)
	GetCommentRepliesWithUserReactions(userID int, chapterID c_sql.NullInt64) ([]repositories.GetRepliesByChapterCommentWithUserReactionsRow, error)
}

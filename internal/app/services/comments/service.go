package commments

import (
	c_sql "database/sql"
	"immodi/novel-site/internal/db/repositories"
	sql "immodi/novel-site/internal/db/schema"
)

type CommentService interface {
	// Comment CRUD
	GetCommentById(commentID int) (repositories.Comment, error)
	CreateComment(novelID, userID int, parentID *int, content string) (repositories.Comment, error)
	GetCommentsByNovel(novelID int) ([]repositories.GetCommentsByNovelRow, error)
	UpdateComment(commentID int, content string) (repositories.Comment, error)
	GetCommentReplies(parentID c_sql.NullInt64) ([]repositories.GetRepliesByCommentRow, error)
	DeleteComment(commentID, userID int) error

	// Reactions
	AddOrUpdateOrRemoveReaction(userID, commentID int, reaction sql.UserReaction) error
	RemoveReaction(userID, commentID int) error
	GetUserReaction(userID, commentID int) (string, error)
	GetUserReactionForComments(userID int, commentIDs []int64) ([]repositories.GetUserReactionsForCommentsRow, error)
	CountReactions(commentID int) (repositories.CountReactionsRow, error)
}

package commments

import (
	"context"
	"database/sql"
	"errors"

	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/db/repositories"
	c_sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
)

type commentService struct {
	db *db.DBService
}

func NewCommentService(db *db.DBService) CommentService {
	return &commentService{db: db}
}

func (c *commentService) AddOrUpdateOrRemoveReaction(
	userID int,
	commentID int,
	reaction c_sql.UserReaction,
) error {
	return db.ExecuteTx(c.db, func(ctx context.Context, q *repositories.Queries) error {
		_, err := q.DeleteReactionIfSame(ctx, repositories.DeleteReactionIfSameParams{
			UserID:    int64(userID),
			CommentID: int64(commentID),
			Reaction:  string(reaction),
		})

		switch {
		case err == nil:
			// A row was deleted -> reaction removed, nothing else to do.
			return nil

		case errors.Is(err, sql.ErrNoRows):
			// Nothing deleted -> user is changing or adding reaction, so upsert.
			return q.UpsertReaction(ctx, repositories.UpsertReactionParams{
				UserID:      int64(userID),
				CommentID:   int64(commentID),
				Reaction:    string(reaction),
				LastUpdated: pkg.GetCurrentTimeRFC3339(),
			})

		default:
			// Any other database error
			return err
		}
	})
}

func (c *commentService) CountReactions(commentID int) (repositories.CountReactionsRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.CountReactionsRow, error) {
		return q.CountReactions(ctx, int64(commentID))
	})
}

func (c *commentService) GetCommentById(commentID int) (repositories.Comment, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.Comment, error) {
		return q.GetCommentById(ctx, int64(commentID))
	})
}

func (c *commentService) CreateComment(novelID int, userID int, parentID *int, content string) (repositories.Comment, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.Comment, error) {
		var parent sql.NullInt64
		if parentID != nil {
			parent = db.NewNullInt64(int64(*parentID))
		} else {
			parent = sql.NullInt64{Valid: false}
		}

		return q.CreateComment(ctx, repositories.CreateCommentParams{
			NovelID:     int64(novelID),
			UserID:      int64(userID),
			Content:     content,
			ParentID:    parent,
			LastUpdated: pkg.GetCurrentTimeRFC3339(),
		})
	})
}

func (c *commentService) UpdateComment(commentID int, content string) (repositories.Comment, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.Comment, error) {
		return q.UpdateComment(ctx, repositories.UpdateCommentParams{
			ID:          int64(commentID),
			Content:     content,
			LastUpdated: pkg.GetCurrentTimeRFC3339(),
		})
	})
}

func (c *commentService) DeleteComment(commentID int, userID int) error {
	return db.Execute(c.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.DeleteComment(ctx, repositories.DeleteCommentParams{
			ID:     int64(commentID),
			UserID: int64(userID),
		})
	})
}

func (c *commentService) GetCommentsByNovel(novelID int) ([]repositories.GetCommentsByNovelRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetCommentsByNovelRow, error) {
		return q.GetCommentsByNovel(ctx, int64(novelID))
	})
}

func (c *commentService) GetCommentReplies(parentID sql.NullInt64) ([]repositories.GetRepliesByCommentRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetRepliesByCommentRow, error) {
		replies, err := q.GetRepliesByComment(ctx, parentID)
		if err != nil {
			return nil, err
		}
		return replies, nil
	})
}

func (c *commentService) GetUserReaction(userID int, commentID int) (string, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (string, error) {
		return q.GetUserReaction(ctx, repositories.GetUserReactionParams{
			UserID:    int64(userID),
			CommentID: int64(commentID),
		})
	})
}

func (c *commentService) GetUserReactionForComments(userID int, commentIDs []int64) ([]repositories.GetUserReactionsForCommentsRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetUserReactionsForCommentsRow, error) {
		return q.GetUserReactionsForComments(ctx, repositories.GetUserReactionsForCommentsParams{
			UserID:     int64(userID),
			CommentIds: commentIDs,
		})
	})
}

func (c *commentService) RemoveReaction(userID int, commentID int) error {
	return db.Execute(c.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.RemoveReaction(ctx, repositories.RemoveReactionParams{
			UserID:    int64(userID),
			CommentID: int64(commentID),
		})
	})
}

package chaptercomments

import (
	"context"
	"database/sql"
	"errors"

	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/db/repositories"
	c_sql "immodi/novel-site/internal/db/schema"
	"immodi/novel-site/pkg"
)

type chapterCommentService struct {
	db *db.DBService
}

func NewChapterCommentService(db *db.DBService) ChapterCommentService {
	return &chapterCommentService{db: db}
}

func (c *chapterCommentService) AddOrUpdateOrRemoveReaction(
	userID int,
	commentID int,
	reaction c_sql.UserReaction,
) error {
	return db.ExecuteTx(c.db, func(ctx context.Context, q *repositories.Queries) error {
		_, err := q.DeleteChapterReactionIfSame(ctx, repositories.DeleteChapterReactionIfSameParams{
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
			return q.UpsertChapterReaction(ctx, repositories.UpsertChapterReactionParams{
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

func (c *chapterCommentService) CountReactions(commentID int) (repositories.CountChapterReactionsRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.CountChapterReactionsRow, error) {
		return q.CountChapterReactions(ctx, int64(commentID))
	})
}

func (c *chapterCommentService) GetCommentById(commentID int) (repositories.ChapterComment, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.ChapterComment, error) {
		return q.GetChapterCommentById(ctx, int64(commentID))
	})
}

func (c *chapterCommentService) CreateComment(chapterID int, userID int, parentID *int, content string) (repositories.ChapterComment, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.ChapterComment, error) {
		var parent sql.NullInt64
		if parentID != nil {
			parent = db.NewNullInt64(int64(*parentID))
		} else {
			parent = sql.NullInt64{Valid: false}
		}

		return q.CreateChapterComment(ctx, repositories.CreateChapterCommentParams{
			ChapterID:   int64(chapterID),
			UserID:      int64(userID),
			Content:     content,
			ParentID:    parent,
			LastUpdated: pkg.GetCurrentTimeRFC3339(),
		})
	})
}

func (c *chapterCommentService) UpdateComment(commentID int, content string) (repositories.ChapterComment, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (repositories.ChapterComment, error) {
		return q.UpdateChapterComment(ctx, repositories.UpdateChapterCommentParams{
			ID:          int64(commentID),
			Content:     content,
			LastUpdated: pkg.GetCurrentTimeRFC3339(),
		})
	})
}

func (c *chapterCommentService) DeleteComment(commentID int, userID int) error {
	return db.Execute(c.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.DeleteChapterComment(ctx, repositories.DeleteChapterCommentParams{
			ID:     int64(commentID),
			UserID: int64(userID),
		})
	})
}

func (c *chapterCommentService) GetCommentsByChapter(chapterID int) ([]repositories.GetCommentsByChapterRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetCommentsByChapterRow, error) {
		return q.GetCommentsByChapter(ctx, int64(chapterID))
	})
}

func (c *chapterCommentService) GetCommentReplies(parentID sql.NullInt64) ([]repositories.GetRepliesByChapterCommentRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetRepliesByChapterCommentRow, error) {
		replies, err := q.GetRepliesByChapterComment(ctx, parentID)
		if err != nil {
			return nil, err
		}
		return replies, nil
	})
}

func (c *chapterCommentService) GetCommentRepliesWithUserReactions(userID int, chapterID sql.NullInt64) ([]repositories.GetRepliesByChapterCommentWithUserReactionsRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetRepliesByChapterCommentWithUserReactionsRow, error) {
		return q.GetRepliesByChapterCommentWithUserReactions(ctx, repositories.GetRepliesByChapterCommentWithUserReactionsParams{
			UserID:   int64(userID),
			ParentID: chapterID,
		})
	})
}

func (c *chapterCommentService) GetUserReaction(userID int, commentID int) (string, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) (string, error) {
		return q.GetUserChapterReaction(ctx, repositories.GetUserChapterReactionParams{
			UserID:    int64(userID),
			CommentID: int64(commentID),
		})
	})
}

func (c *chapterCommentService) GetUserReactionForComments(userID int, commentIDs []int64) ([]repositories.GetUserReactionsForChapterCommentsRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetUserReactionsForChapterCommentsRow, error) {
		return q.GetUserReactionsForChapterComments(ctx, repositories.GetUserReactionsForChapterCommentsParams{
			UserID:     int64(userID),
			CommentIds: commentIDs,
		})
	})
}

func (c *chapterCommentService) RemoveReaction(userID int, commentID int) error {
	return db.Execute(c.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.RemoveChapterReaction(ctx, repositories.RemoveChapterReactionParams{
			UserID:    int64(userID),
			CommentID: int64(commentID),
		})
	})
}

func (c *chapterCommentService) GetCommentsByChapterWithUserReactions(userID, chapterID int) ([]repositories.GetCommentsByChapterWithUserReactionsRow, error) {
	return db.ExecuteWithResult(c.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.GetCommentsByChapterWithUserReactionsRow, error) {
		return q.GetCommentsByChapterWithUserReactions(ctx, repositories.GetCommentsByChapterWithUserReactionsParams{
			UserID:    int64(userID),
			ChapterID: int64(chapterID),
		})
	})
}

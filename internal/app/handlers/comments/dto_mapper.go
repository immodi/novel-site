package comments

import (
	sql "database/sql"
	comments "immodi/novel-site/internal/app/services/comments"
	"immodi/novel-site/internal/db/repositories"
	commentsdtostructs "immodi/novel-site/internal/http/structs/comments"
	"log"
)

func DbCommentsToCommentDtoMapper(
	comments []repositories.GetCommentsByNovelRow,
	commentsReactions []repositories.GetUserReactionsForCommentsRow,
	userID int,
	service comments.CommentService,
) []commentsdtostructs.CommentDTO {

	// Build a quick lookup: comment_id -> reaction
	reactionMap := make(map[int64]string, len(commentsReactions))
	for _, r := range commentsReactions {
		reactionMap[r.CommentID] = r.Reaction
	}

	var dtos []commentsdtostructs.CommentDTO
	if len(comments) == 0 {
		return dtos
	}

	// Track which comments are replies
	isReply := make(map[int64]bool, len(comments))
	// Cache replies for each comment
	repliesCache := make(map[int64][]repositories.GetRepliesByCommentRow, len(comments))

	for _, c := range comments {
		replRows, err := service.GetCommentReplies(sql.NullInt64{Int64: c.ID, Valid: true})
		if err != nil {
			log.Println(err)
			replRows = nil
		}
		repliesCache[c.ID] = replRows
		for _, r := range replRows {
			isReply[r.ID] = true
		}
	}

	// Recursive helper to build replies
	var buildReplyDTO func(r repositories.GetRepliesByCommentRow) commentsdtostructs.CommentDTO
	buildReplyDTO = func(r repositories.GetRepliesByCommentRow) commentsdtostructs.CommentDTO {
		nestedRows := repliesCache[r.ID]
		var nestedDtos []commentsdtostructs.CommentDTO
		for _, nr := range nestedRows {
			nestedDtos = append(nestedDtos, buildReplyDTO(nr))
		}

		var userReaction *string
		if react, ok := reactionMap[r.ID]; ok {
			userReaction = &react
		}

		return commentsdtostructs.CommentDTO{
			ID:           int(r.ID),
			Content:      r.Content,
			LastUpdated:  r.LastUpdated,
			UserID:       int(r.UserID),
			PictureURL:   r.PictureUrl,
			UserName:     r.Username,
			Likes:        int(r.Likes),
			Dislikes:     int(r.Dislikes),
			Replies:      nestedDtos,
			UserReaction: userReaction,
		}
	}

	// Top-level comments only
	for _, c := range comments {
		if isReply[c.ID] {
			continue
		}

		var replyDtos []commentsdtostructs.CommentDTO
		for _, r := range repliesCache[c.ID] {
			replyDtos = append(replyDtos, buildReplyDTO(r))
		}

		var userReaction *string
		if react, ok := reactionMap[c.ID]; ok {
			userReaction = &react
		}

		dtos = append(dtos, commentsdtostructs.CommentDTO{
			ID:           int(c.ID),
			Content:      c.Content,
			LastUpdated:  c.LastUpdated,
			UserID:       int(c.UserID),
			PictureURL:   c.PictureUrl,
			UserName:     c.Username,
			Likes:        int(c.Likes),
			Dislikes:     int(c.Dislikes),
			Replies:      replyDtos,
			UserReaction: userReaction,
		})
	}

	return dtos
}

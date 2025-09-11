package comments

import (
	sql "database/sql"
	comments "immodi/novel-site/internal/app/services/comments"
	"immodi/novel-site/internal/db/repositories"
	commentsdtostructs "immodi/novel-site/internal/http/structs/comments"
	"log"
)

func DbCommentsToCommentDtoMapper(comments []repositories.GetCommentsByNovelRow, service comments.CommentService) []commentsdtostructs.CommentDTO {
	var commentDtos []commentsdtostructs.CommentDTO
	if len(comments) == 0 {
		return commentDtos
	}

	// Track IDs that are replies (so we don't add them as top-level later)
	isReply := make(map[int64]bool, len(comments))

	// Cache replies we fetched to avoid refetching where possible:
	repliesCache := make(map[int64][]repositories.GetRepliesByCommentRow, len(comments))

	// Pre-fetch direct replies for every item in the input slice and mark reply IDs
	for _, c := range comments {
		replRows, err := service.GetCommentReplies(sql.NullInt64{Int64: c.ID, Valid: true})
		if err != nil {
			log.Println(err.Error(), "pre-caching replies - continuing")
			replRows = []repositories.GetRepliesByCommentRow{}
		}
		repliesCache[c.ID] = replRows
		for _, r := range replRows {
			isReply[r.ID] = true
		}
	}

	// Recursive builder for reply rows (re-uses cache, falls back to service)
	var buildReplyDTO func(r repositories.GetRepliesByCommentRow) commentsdtostructs.CommentDTO
	buildReplyDTO = func(r repositories.GetRepliesByCommentRow) commentsdtostructs.CommentDTO {
		// mark as reply (defensive)
		isReply[r.ID] = true

		// get nested replies from cache or fetch
		nestedRows, ok := repliesCache[r.ID]
		if !ok {
			var err error
			nestedRows, err = service.GetCommentReplies(sql.NullInt64{Int64: r.ID, Valid: true})
			if err != nil {
				log.Println(err.Error(), "fetching nested replies")
				nestedRows = []repositories.GetRepliesByCommentRow{}
			}
			repliesCache[r.ID] = nestedRows
			for _, nr := range nestedRows {
				isReply[nr.ID] = true
			}
		}

		nestedDtos := make([]commentsdtostructs.CommentDTO, 0, len(nestedRows))
		for _, nr := range nestedRows {
			nestedDtos = append(nestedDtos, buildReplyDTO(nr))
		}

		return commentsdtostructs.CommentDTO{
			ID:         int(r.ID),
			Content:    r.Content,
			CreatedAt:  r.CreatedAt,
			UserID:     int(r.UserID),
			PictureURL: r.PictureUrl,
			UserName:   r.Username,
			Likes:      int(r.Likes),
			Dislikes:   int(r.Dislikes),
			Replies:    nestedDtos,
		}
	}

	// Build top-level DTOs, skipping any rows that are actually replies
	for _, c := range comments {
		if isReply[c.ID] {
			// This row is a reply to some other comment, skip as a top-level item
			continue
		}

		// get direct replies (from cache)
		replRows := repliesCache[c.ID]
		replyDtos := make([]commentsdtostructs.CommentDTO, 0, len(replRows))
		for _, r := range replRows {
			replyDtos = append(replyDtos, buildReplyDTO(r))
		}

		commentDtos = append(commentDtos, commentsdtostructs.CommentDTO{
			ID:         int(c.ID),
			Content:    c.Content,
			CreatedAt:  c.CreatedAt,
			UserID:     int(c.UserID),
			PictureURL: c.PictureUrl,
			UserName:   c.Username,
			Likes:      int(c.Likes),
			Dislikes:   int(c.Dislikes),
			Replies:    replyDtos,
		})
	}

	return commentDtos
}

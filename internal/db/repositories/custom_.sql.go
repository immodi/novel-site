package repositories

//
// import (
// 	"context"
// 	"fmt"
// 	"strings"
// )
//
// func BuildNovelsQuery(
// 	tags []string,
// 	tagsExclude []string,
// 	tagCondition string,
// 	genres []string,
// 	genreCondition string,
// ) string {
// 	baseQuery := `
// SELECT
//     n.id,
//     n.title,
//     n.slug,
//     n.release_year,
//     n.is_completed,
//     n.update_time,
//     n.view_count,
//     n.cover_image,
//     n.author,
//     COALESCE((SELECT COUNT(*) FROM chapters WHERE novel_id = n.id), 0) AS chapter_count,
//     COALESCE((SELECT COUNT(DISTINCT user_id) FROM user_bookmarks WHERE novel_id = n.id), 0) AS bookmark_count,
//     COALESCE((SELECT COUNT(*) FROM comments WHERE novel_id = n.id), 0) AS comment_count,
//     COALESCE((SELECT GROUP_CONCAT(DISTINCT genre) FROM novel_genres WHERE novel_id = n.id), '') AS genres,
//     COALESCE((SELECT GROUP_CONCAT(DISTINCT tag) FROM novel_tags WHERE novel_id = n.id), '') AS tags
// FROM novels n
// WHERE (? = -1 OR n.is_completed = ?)`
//
// 	// Add filtering conditions
// 	var whereConditions []string
//
// 	// --- TAGS filter ---
// 	if len(tags) > 0 {
// 		placeholders := make([]string, len(tags))
// 		for i := range tags {
// 			placeholders[i] = "?"
// 		}
// 		inClause := strings.Join(placeholders, ",")
//
// 		switch tagCondition {
// 		case "and":
// 			whereConditions = append(whereConditions, fmt.Sprintf(`
//     (SELECT COUNT(DISTINCT tag) FROM novel_tags WHERE novel_id = n.id AND tag IN (%s)) = %d`, inClause, len(tags)))
// 		default: // "or" behavior
// 			whereConditions = append(whereConditions, fmt.Sprintf(`
//     EXISTS (SELECT 1 FROM novel_tags WHERE novel_id = n.id AND tag IN (%s))`, inClause))
// 		}
// 	}
//
// 	if len(tagsExclude) > 0 {
// 		excludePlaceholders := make([]string, len(tagsExclude))
// 		for i := range tagsExclude {
// 			excludePlaceholders[i] = "?"
// 		}
// 		excludeInClause := strings.Join(excludePlaceholders, ",")
// 		whereConditions = append(whereConditions, fmt.Sprintf(`
//     NOT EXISTS (SELECT 1 FROM novel_tags WHERE novel_id = n.id AND tag IN (%s))`, excludeInClause))
// 	}
//
// 	// --- GENRES filter ---
// 	if len(genres) > 0 {
// 		genrePlaceholders := make([]string, len(genres))
// 		for i := range genres {
// 			genrePlaceholders[i] = "?"
// 		}
// 		genreInClause := strings.Join(genrePlaceholders, ",")
//
// 		switch genreCondition {
// 		case "and":
// 			whereConditions = append(whereConditions, fmt.Sprintf(`
//     (SELECT COUNT(DISTINCT genre) FROM novel_genres WHERE novel_id = n.id AND genre IN (%s)) = %d`, genreInClause, len(genres)))
// 		case "exclude":
// 			whereConditions = append(whereConditions, fmt.Sprintf(`
//     NOT EXISTS (SELECT 1 FROM novel_genres WHERE novel_id = n.id AND genre IN (%s))`, genreInClause))
// 		default: // "or" behavior
// 			whereConditions = append(whereConditions, fmt.Sprintf(`
//     EXISTS (SELECT 1 FROM novel_genres WHERE novel_id = n.id AND genre IN (%s))`, genreInClause))
// 		}
// 	}
//
// 	// Add all WHERE conditions
// 	for _, condition := range whereConditions {
// 		baseQuery += " AND " + condition
// 	}
//
// 	// Add chapter count filter
// 	baseQuery += `
//     AND (SELECT COUNT(*) FROM chapters WHERE novel_id = n.id) BETWEEN ? AND ?
// ORDER BY n.title;`
//
// 	return baseQuery
// }
//
// type FilterNovelsParams struct {
// 	MinChapters    int64
// 	MaxChapters    int64
// 	Tags           []string
// 	TagsExclude    []string
// 	TagCondition   string
// 	Genres         []string
// 	GenreCondition string
// 	IsCompleted    int // -1 = ignore, 0 = incomplete, 1 = completed
// }
//
// type FilterNovelsRows struct {
// 	ID            int64
// 	Title         string
// 	Slug          string
// 	ReleaseYear   int64
// 	IsCompleted   bool
// 	UpdateTime    string
// 	ViewCount     int64
// 	CoverImage    string
// 	Author        string
// 	ChapterCount  int64
// 	BookmarkCount int64
// 	CommentCount  int64
// 	Genres        string // comma-separated
// 	Tags          string // comma-separated
// }
//
// func (q *Queries) FilterNovels(ctx context.Context, arg FilterNovelsParams) ([]FilterNovelsRows, error) {
// 	query := BuildNovelsQuery(
// 		arg.Tags,
// 		arg.TagsExclude,
// 		arg.TagCondition,
// 		arg.Genres,
// 		arg.GenreCondition,
// 	)
//
// 	// Build the args in the CORRECT order based on how the query is constructed
// 	args := []any{
// 		arg.IsCompleted, // First ? in: (? = -1 OR n.is_completed = ?)
// 		arg.IsCompleted, // Second ? in: (? = -1 OR n.is_completed = ?)
// 	}
//
// 	// Add tag values (these come after the base WHERE clause)
// 	for _, t := range arg.Tags {
// 		args = append(args, t)
// 	}
//
// 	// Add excluded tag values
// 	for _, t := range arg.TagsExclude {
// 		args = append(args, t)
// 	}
//
// 	// Add genre values
// 	for _, g := range arg.Genres {
// 		args = append(args, g)
// 	}
//
// 	// Finally add chapter count filters (these come at the very end)
// 	args = append(args, arg.MinChapters, arg.MaxChapters)
//
// 	rows, err := q.db.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var items []FilterNovelsRows
// 	for rows.Next() {
// 		var i FilterNovelsRows
// 		if err := rows.Scan(
// 			&i.ID,
// 			&i.Title,
// 			&i.Slug,
// 			&i.ReleaseYear,
// 			&i.IsCompleted,
// 			&i.UpdateTime,
// 			&i.ViewCount,
// 			&i.CoverImage,
// 			&i.Author,
// 			&i.ChapterCount,
// 			&i.BookmarkCount,
// 			&i.CommentCount,
// 			&i.Genres,
// 			&i.Tags,
// 		); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, i)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return items, nil
// }

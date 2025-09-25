package repositories

import (
	"context"
	"fmt"
	"strings"
)

const (
	SortByBookmark = "bookmark"
	SortByComment  = "review"
	SortByABC      = "abc"
	SortByCBA      = "cba"
	SortByDate     = "date"
	SortByChapter  = "chapter-count-most"
)

func BuildNovelsQuery(
	tags []string,
	tagsExclude []string,
	tagCondition string,
	genres []string,
	genreCondition string,
) string {
	baseQuery := `
SELECT
    n.id,
    n.title,
    n.slug,
    n.release_year,
    n.is_completed,
    n.update_time,
    n.view_count,
    n.cover_image,
    n.author,
    COALESCE((SELECT COUNT(*) FROM chapters WHERE novel_id = n.id), 0) AS chapter_count,
    COALESCE((SELECT COUNT(DISTINCT user_id) FROM user_bookmarks WHERE novel_id = n.id), 0) AS bookmark_count,
    COALESCE((SELECT COUNT(*) FROM comments WHERE novel_id = n.id), 0) AS comment_count,
    COALESCE((SELECT GROUP_CONCAT(DISTINCT genre) FROM novel_genres WHERE novel_id = n.id), '') AS genres,
    COALESCE((SELECT GROUP_CONCAT(DISTINCT tag) FROM novel_tags WHERE novel_id = n.id), '') AS tags
FROM novels n
WHERE (? = -1 OR n.is_completed = ?)`

	// Add filtering conditions
	var whereConditions []string

	// --- TAGS filter ---
	if len(tags) > 0 {
		placeholders := make([]string, len(tags))
		for i := range tags {
			placeholders[i] = "?"
		}
		inClause := strings.Join(placeholders, ",")

		switch tagCondition {
		case "and":
			whereConditions = append(whereConditions, fmt.Sprintf(`
    (SELECT COUNT(DISTINCT tag) FROM novel_tags WHERE novel_id = n.id AND tag IN (%s)) = %d`, inClause, len(tags)))
		default: // "or" behavior
			whereConditions = append(whereConditions, fmt.Sprintf(`
    EXISTS (SELECT 1 FROM novel_tags WHERE novel_id = n.id AND tag IN (%s))`, inClause))
		}
	}

	if len(tagsExclude) > 0 {
		excludePlaceholders := make([]string, len(tagsExclude))
		for i := range tagsExclude {
			excludePlaceholders[i] = "?"
		}
		excludeInClause := strings.Join(excludePlaceholders, ",")
		whereConditions = append(whereConditions, fmt.Sprintf(`
    NOT EXISTS (SELECT 1 FROM novel_tags WHERE novel_id = n.id AND tag IN (%s))`, excludeInClause))
	}

	// --- GENRES filter ---
	if len(genres) > 0 {
		genrePlaceholders := make([]string, len(genres))
		for i := range genres {
			genrePlaceholders[i] = "?"
		}
		genreInClause := strings.Join(genrePlaceholders, ",")

		switch genreCondition {
		case "and":
			whereConditions = append(whereConditions, fmt.Sprintf(`
    (SELECT COUNT(DISTINCT genre) FROM novel_genres WHERE novel_id = n.id AND genre IN (%s)) = %d`, genreInClause, len(genres)))
		case "exclude":
			whereConditions = append(whereConditions, fmt.Sprintf(`
    NOT EXISTS (SELECT 1 FROM novel_genres WHERE novel_id = n.id AND genre IN (%s))`, genreInClause))
		default: // "or" behavior
			whereConditions = append(whereConditions, fmt.Sprintf(`
    EXISTS (SELECT 1 FROM novel_genres WHERE novel_id = n.id AND genre IN (%s))`, genreInClause))
		}
	}

	// Add all WHERE conditions
	for _, condition := range whereConditions {
		baseQuery += " AND " + condition
	}

	// Add chapter count filter
	baseQuery += `
    AND (SELECT COUNT(*) FROM chapters WHERE novel_id = n.id) BETWEEN ? AND ?`

	return baseQuery
}

// New function to build count query (same filters, just COUNT(*))
func BuildNovelsCountQuery(
	tags []string,
	tagsExclude []string,
	tagCondition string,
	genres []string,
	genreCondition string,
) string {
	baseQuery := `
SELECT COUNT(*)
FROM novels n
WHERE (? = -1 OR n.is_completed = ?)`

	// Add filtering conditions (same logic as main query)
	var whereConditions []string

	// --- TAGS filter ---
	if len(tags) > 0 {
		placeholders := make([]string, len(tags))
		for i := range tags {
			placeholders[i] = "?"
		}
		inClause := strings.Join(placeholders, ",")

		switch tagCondition {
		case "and":
			whereConditions = append(whereConditions, fmt.Sprintf(`
    (SELECT COUNT(DISTINCT tag) FROM novel_tags WHERE novel_id = n.id AND tag IN (%s)) = %d`, inClause, len(tags)))
		default: // "or" behavior
			whereConditions = append(whereConditions, fmt.Sprintf(`
    EXISTS (SELECT 1 FROM novel_tags WHERE novel_id = n.id AND tag IN (%s))`, inClause))
		}
	}

	if len(tagsExclude) > 0 {
		excludePlaceholders := make([]string, len(tagsExclude))
		for i := range tagsExclude {
			excludePlaceholders[i] = "?"
		}
		excludeInClause := strings.Join(excludePlaceholders, ",")
		whereConditions = append(whereConditions, fmt.Sprintf(`
    NOT EXISTS (SELECT 1 FROM novel_tags WHERE novel_id = n.id AND tag IN (%s))`, excludeInClause))
	}

	// --- GENRES filter ---
	if len(genres) > 0 {
		genrePlaceholders := make([]string, len(genres))
		for i := range genres {
			genrePlaceholders[i] = "?"
		}
		genreInClause := strings.Join(genrePlaceholders, ",")

		switch genreCondition {
		case "and":
			whereConditions = append(whereConditions, fmt.Sprintf(`
    (SELECT COUNT(DISTINCT genre) FROM novel_genres WHERE novel_id = n.id AND genre IN (%s)) = %d`, genreInClause, len(genres)))
		case "exclude":
			whereConditions = append(whereConditions, fmt.Sprintf(`
    NOT EXISTS (SELECT 1 FROM novel_genres WHERE novel_id = n.id AND genre IN (%s))`, genreInClause))
		default: // "or" behavior
			whereConditions = append(whereConditions, fmt.Sprintf(`
    EXISTS (SELECT 1 FROM novel_genres WHERE novel_id = n.id AND genre IN (%s))`, genreInClause))
		}
	}

	// Add all WHERE conditions
	for _, condition := range whereConditions {
		baseQuery += " AND " + condition
	}

	// Add chapter count filter
	baseQuery += `
    AND (SELECT COUNT(*) FROM chapters WHERE novel_id = n.id) BETWEEN ? AND ?`

	return baseQuery
}

type FilterNovelsParams struct {
	MinChapters    int64
	MaxChapters    int64
	Tags           []string
	TagsExclude    []string
	TagCondition   string
	Genres         []string
	SortBy         string
	GenreCondition string
	IsCompleted    int   // -1 = ignore, 0 = incomplete, 1 = completed
	Limit          int64 // Number of results to return
	Offset         int64 // Number of results to skip
}

type FilterNovelsRows struct {
	ID            int64
	Title         string
	Slug          string
	ReleaseYear   int64
	IsCompleted   bool
	UpdateTime    string
	ViewCount     int64
	CoverImage    string
	Author        string
	ChapterCount  int64
	BookmarkCount int64
	CommentCount  int64
	Genres        string // comma-separated
	Tags          string // comma-separated
}

type FilterNovelsResult struct {
	Novels     []FilterNovelsRows
	TotalCount int64
}

func (q *Queries) FilterNovels(ctx context.Context, arg FilterNovelsParams) (*FilterNovelsResult, error) {
	// Set default pagination values
	if arg.Limit <= 0 {
		arg.Limit = 20 // Default page size
	}
	if arg.Offset < 0 {
		arg.Offset = 0
	}

	// Build args (same for both queries)
	args := []any{
		arg.IsCompleted, // First ? in: (? = -1 OR n.is_completed = ?)
		arg.IsCompleted, // Second ? in: (? = -1 OR n.is_completed = ?)
	}

	// Add tag values (these come after the base WHERE clause)
	for _, t := range arg.Tags {
		args = append(args, t)
	}

	// Add excluded tag values
	for _, t := range arg.TagsExclude {
		args = append(args, t)
	}

	// Add genre values
	for _, g := range arg.Genres {
		args = append(args, g)
	}

	// Add chapter count filters (these come at the very end)
	args = append(args, arg.MinChapters, arg.MaxChapters)

	// First, get the total count
	countQuery := BuildNovelsCountQuery(
		arg.Tags,
		arg.TagsExclude,
		arg.TagCondition,
		arg.Genres,
		arg.GenreCondition,
	)

	var totalCount int64
	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Then get the paginated data
	dataQuery := BuildNovelsQuery(
		arg.Tags,
		arg.TagsExclude,
		arg.TagCondition,
		arg.Genres,
		arg.GenreCondition,
	) + " " + getOrderByClause(arg.SortBy) + " LIMIT ? OFFSET ?"

	// Add pagination params to args
	dataArgs := append(args, arg.Limit, arg.Offset)

	rows, err := q.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var items []FilterNovelsRows
	for rows.Next() {
		var i FilterNovelsRows
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Slug,
			&i.ReleaseYear,
			&i.IsCompleted,
			&i.UpdateTime,
			&i.ViewCount,
			&i.CoverImage,
			&i.Author,
			&i.ChapterCount,
			&i.BookmarkCount,
			&i.CommentCount,
			&i.Genres,
			&i.Tags,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return &FilterNovelsResult{
		Novels:     items,
		TotalCount: totalCount,
	}, nil
}

func getOrderByClause(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case SortByBookmark:
		return "ORDER BY bookmark_count DESC"
	case SortByComment:
		return "ORDER BY comment_count DESC"
	case SortByABC:
		return "ORDER BY title ASC"
	case SortByCBA:
		return "ORDER BY title DESC"
	case SortByDate:
		return "ORDER BY update_time DESC"
	case SortByChapter:
		return "ORDER BY chapter_count DESC"
	default:
		return "ORDER BY title ASC" // default fallback
	}
}

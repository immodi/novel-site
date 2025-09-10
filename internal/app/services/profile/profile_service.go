package profile

import (
	"context"
	"immodi/novel-site/internal/app/services/db"
	"immodi/novel-site/internal/app/services/novels"
	"immodi/novel-site/internal/db/repositories"
	"immodi/novel-site/pkg"
)

type profileService struct {
	db           *db.DBService
	novelService novels.NovelService
}

func NewProfileService(db *db.DBService, novelService novels.NovelService) ProfileService {
	return &profileService{db: db, novelService: novelService}
}

func (p *profileService) CountBookMarkedNovels(userId int64) (int64, error) {
	return db.ExecuteWithResult(p.db, func(ctx context.Context, q *repositories.Queries) (int64, error) {
		return q.CountUserBookmarks(ctx, userId)
	})
}

func (p *profileService) GetUserById(id int64) (repositories.User, error) {
	return db.ExecuteWithResult(p.db, func(ctx context.Context, q *repositories.Queries) (repositories.User, error) {
		return q.GetUserByID(ctx, id)
	})
}

func (p *profileService) ListBookMarkedNovelsPaginated(userId int64, offset int, limit int) ([]repositories.Novel, error) {
	return db.ExecuteWithResult(p.db, func(ctx context.Context, q *repositories.Queries) ([]repositories.Novel, error) {
		return q.ListUserBookmarksPaginated(ctx, repositories.ListUserBookmarksPaginatedParams{
			UserID: userId,
			Offset: int64(offset),
			Limit:  int64(limit),
		})
	})
}

func (p *profileService) UpdateUserPartial(params repositories.UpdateUserPartialParams) (repositories.User, error) {
	return db.ExecuteWithResult(p.db, func(ctx context.Context, q *repositories.Queries) (repositories.User, error) {
		return q.UpdateUserPartial(ctx, params)
	})
}

func (p *profileService) AddUserBookMark(userId, novelID int64) error {
	return db.Execute(p.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.AddUserBookmark(ctx, repositories.AddUserBookmarkParams{
			UserID:    userId,
			NovelID:   novelID,
			CreatedAt: pkg.GetCurrentTimeRFC3339(),
		})
	})
}

func (p *profileService) RemoveUserBookmark(userId, novelID int64) error {
	return db.Execute(p.db, func(ctx context.Context, q *repositories.Queries) error {
		return q.RemoveUserBookmark(ctx, repositories.RemoveUserBookmarkParams{
			UserID:  userId,
			NovelID: novelID,
		})
	})
}

func (p *profileService) GetNovelBySlug(slug string) (repositories.Novel, error) {
	return p.novelService.GetNovelBySlug(slug)
}

package profile

import "immodi/novel-site/internal/db/repositories"

type ProfileService interface {
	GetNovelBySlug(slug string) (repositories.Novel, error)

	AddUserBookMark(userId, novelID int64) error
	RemoveUserBookmark(userId, novelID int64) error

	ListBookMarkedNovelsPaginated(userId int64, offset, limit int) ([]repositories.Novel, error)
	CountBookMarkedNovels(userId int64) (int64, error)

	UpdateUserPartial(params repositories.UpdateUserPartialParams) (repositories.User, error)
	GetUserById(id int64) (repositories.User, error)
}

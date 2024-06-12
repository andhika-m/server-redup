package bookmark

import (
	"redup/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetBookmarks(userID string) ([]model.Bookmark, error)
	AddBookmark(userID, videoID string) error
	RemoveBookmark(userID, bookmarkID string) error
}

type bookmarkRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &bookmarkRepo{
		db: db,
	}
}

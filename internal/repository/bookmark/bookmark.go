package bookmark

import (
	"redup/internal/model"
	"time"

	"github.com/google/uuid"
)

func (br *bookmarkRepo) GetBookmarks(userID string) ([]model.Bookmark, error) {
	var bookmarks []model.Bookmark
	err := br.db.Where("user_id = ?", userID).Find(&bookmarks).Error
	return bookmarks, err
}

func (br *bookmarkRepo) AddBookmark(userID, videoID string) error {
	bookmark := model.Bookmark{
		ID:        uuid.New().String(),
		UserID:    userID,
		VideoID:   videoID,
		CreatedAt: time.Now(),
	}
	return br.db.Create(&bookmark).Error
}

func (br *bookmarkRepo) RemoveBookmark(userID, bookmarkID string) error {
	return br.db.Where("user_id = ? AND id = ?", userID, bookmarkID).Delete(&model.Bookmark{}).Error
}

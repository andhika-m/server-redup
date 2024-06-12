package history

import (
	"redup/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *historyRepo) GetHistory(userID string) ([]model.History, error) {
	var history []model.History
	err := r.db.Where("user_id = ?", userID).Find(&history).Error
	return history, err
}

func (r *historyRepo) AddHistory(userID, videoID string) error {
	var history model.History

	err := r.db.Where("user_id = ? AND video_id = ?", userID, videoID).First(&history).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			history = model.History{
				ID:        uuid.New().String(),
				UserID:    userID,
				VideoID:   videoID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			return r.db.Create(&history).Error
		}
		return err
	}

	return r.db.Model(&history).Update("updated_at", time.Now()).Error
}

func (r *historyRepo) RemoveHistory(userID, historyID string) error {
	return r.db.Where("user_id = ? AND id = ?", userID, historyID).Delete(&model.History{}).Error
}

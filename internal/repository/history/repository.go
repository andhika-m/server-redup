package history

import (
	"redup/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetHistory(userID string) ([]model.History, error)
	AddHistory(userID, videoID string) error
	RemoveHistory(userID, historyID string) error
}

type historyRepo struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) Repository {
	return &historyRepo{
		db: db,
	}
}

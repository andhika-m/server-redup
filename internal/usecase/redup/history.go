package redup

import "redup/internal/model"

func (ru *redupUsecase) GetHistory(userID string) ([]model.History, error) {
	return ru.historyRepo.GetHistory(userID)
}

func (ru *redupUsecase) AddHistory(userID, videoID string) error {
	return ru.historyRepo.AddHistory(userID, videoID)
}

func (ru *redupUsecase) RemoveHistory(userID, historyID string) error {
	return ru.historyRepo.RemoveHistory(userID, historyID)
}

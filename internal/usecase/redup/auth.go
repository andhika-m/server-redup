package redup

import (
	"redup/internal/model"
)

func (ru *redupUsecase) VerifySession(token string) (string, error) {
	return ru.userRepo.VerifySession(token)
}


func (ru *redupUsecase) CheckSession(data model.UserSession) (userID string, err error) {
	userID, err = ru.userRepo.CheckSession(data)
	if err != nil {
		return "", err
	}

	return userID, nil
}
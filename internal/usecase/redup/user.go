package redup

import (
	"errors"
	"redup/internal/model"

	"redup/internal/model/constant"

	"github.com/google/uuid"
)

func (ru *redupUsecase) RegisterUser(request model.RegisterRequest) (model.User, error) {
	userRegistered, err := ru.userRepo.CheckRegistered(request.Email)
	if err != nil {
		return model.User{}, err
	}
	if userRegistered {
		return model.User{}, errors.New("email already registered")
	}

	userHash, err := ru.userRepo.GenerateUserHash(request.Password)
	if err != nil {
		return model.User{}, err
	}

	role := constant.RoleCategoryStudent
	if request.Role == constant.RoleCategoryTeacher {
		role = constant.RoleCategoryTeacher
	} else if request.Role != "" && request.Role != constant.RoleCategoryStudent {
		return model.User{}, errors.New("invalid role specified")
	}

	userData, err := ru.userRepo.RegisterUser(model.User{
		ID:     uuid.New().String(),
		Name:   request.Name,
		School: request.School,
		Email:  request.Email,
		Hash:   userHash,
		Role:   role,
	})

	if err != nil {
		return model.User{}, err
	}

	return userData, nil
}

func (ru *redupUsecase) Login(request model.LoginRequest) (model.UserSession, error) {
	userData, err := ru.userRepo.GetUserData(request.Email)
	if err != nil {
		return model.UserSession{}, err
	}

	verified, err := ru.userRepo.VerifyLogin(request.Email, request.Password, userData)
	if err != nil {
		return model.UserSession{}, err
	}

	if !verified {
		return model.UserSession{}, errors.New("can't verify user login")
	}

	userSession, err := ru.userRepo.CreateUserSession(userData.ID)
	if err != nil {
		return model.UserSession{}, err
	}

	return userSession, nil
}

func (ru *redupUsecase) GetUserByID(userID string) (model.User, error) {
	return ru.userRepo.GetUserByID(userID)
}

func (ru *redupUsecase) UpdateUser(userID string, updatedData model.User) (model.User, error) {

	currentUser, err := ru.userRepo.GetUserByID(userID)
	if err != nil {
		return model.User{}, err
	}

	if updatedData.Name != "" {
		currentUser.Name = updatedData.Name
	}
	if updatedData.Email != "" {
		currentUser.Email = updatedData.Email
	}
	if updatedData.School != "" {
		currentUser.School = updatedData.School
	}
	if updatedData.Hash != "" {
		newHash, err := ru.userRepo.GenerateUserHash(updatedData.Hash)
		if err != nil {
			return model.User{}, err
		}
		currentUser.Hash = newHash
	}

	updatedUser, err := ru.userRepo.UpdateUser(currentUser)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (ru *redupUsecase) DeleteUser(userID string) error {
	return ru.userRepo.DeleteUser(userID)
}

func (ru *redupUsecase) Logout(userID string, token string) error {
	return ru.userRepo.DeleteUserSession(userID, token)
}

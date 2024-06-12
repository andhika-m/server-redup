package user

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"redup/internal/model"
	"time"

	"gorm.io/gorm"
)

type userRepo struct {
	db        *gorm.DB
	gcm       cipher.AEAD
	time      uint32
	memory    uint32
	threads   uint8
	keyLen    uint32
	signKey   *rsa.PrivateKey
	accessExp time.Duration
}

func GetRepository(
	db *gorm.DB,
	secret string,
	time uint32,
	memory uint32,
	threads uint8,
	keyLen uint32,
	signKey *rsa.PrivateKey,
	accessExp time.Duration,

) (Repository, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &userRepo{
		db:        db,
		gcm:       gcm,
		time:      time,
		memory:    memory,
		threads:   threads,
		keyLen:    keyLen,
		signKey:   signKey,
		accessExp: accessExp,
	}, nil
}

func (ur *userRepo) RegisterUser(userData model.User) (model.User, error) {
	if err := ur.db.Create(&userData).Error; err != nil {
		return model.User{}, err
	}

	return userData, nil
}

func (ur *userRepo) CheckRegistered(email string) (bool, error) {
	var userData model.User

	if err := ur.db.Where(model.User{Email: email}).First(&userData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return userData.ID != "", nil
}

func (ur *userRepo) GetUserData(email string) (model.User, error) {
	var userData model.User

	if err := ur.db.Where(model.User{Email: email}).First(&userData).Error; err != nil {
		return userData, err
	}

	return userData, nil
}

func (ur *userRepo) VerifyLogin(email, password string, userData model.User) (bool, error) {
	if email != userData.Email {
		return false, nil
	}

	verified, err := ur.comparePassword(password, userData.Hash)
	if err != nil {
		return false, err
	}

	return verified, nil
}

func (ur *userRepo) GetUserByID(userID string) (model.User, error) {
	var user model.User
	if err := ur.db.First(&user, "id = ?", userID).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ur *userRepo) UpdateUser(user model.User) (model.User, error) {
	if err := ur.db.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (ur *userRepo) DeleteUser(userID string) error {
	if err := ur.db.Delete(&model.User{}, "id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepo) DeleteUserSession(userID string, token string) error {
	return ur.db.Where("user_id = ? AND jwt_token = ?", userID, token).Delete(&model.UserSession{}).Error
}

func (ur *userRepo) VerifySession(token string) (string, error) {
	var session model.UserSession
	if err := ur.db.First(&session, "jwt_token = ?", token).Error; err != nil {
		return "", errors.New("invalid session")
	}
	return session.UserID, nil
}

func ParseToken(token string) (map[string]interface{}, error) {
	signingKey := []byte("your-signing-key")

	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

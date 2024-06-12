package redup

import (
	"redup/internal/model"
	"redup/internal/repository/bookmark"
	"redup/internal/repository/user"
	"redup/internal/repository/video"
	"redup/internal/repository/history"
)

type redupUsecase struct {
	userRepo     user.Repository
	videoRepo    video.Repository
	bookmarkRepo bookmark.Repository
	historyRepo history.Repository
}

func RedupUsecase(videoRepo video.Repository, userRepo user.Repository, bookmarkRepo bookmark.Repository, historyRepo history.Repository) Usecase {
	return &redupUsecase{
		videoRepo:    videoRepo,
		userRepo:     userRepo,
		bookmarkRepo: bookmarkRepo,
		historyRepo: historyRepo,
	}
}

type Usecase interface {
	// user
	RegisterUser(request model.RegisterRequest) (model.User, error)
	Login(request model.LoginRequest) (model.UserSession, error)
	GetUserByID(userID string) (model.User, error)
	UpdateUser(userID string, updatedData model.User) (model.User, error)
	DeleteUser(userID string) error
	Logout(userID string, token string) error

	// bookmark
	GetBookmarks(userID string) ([]model.Bookmark, error)
	RemoveBookmark(userID, bookmarkID string) error
	AddBookmark(userID, videoID string) error

	// auth
	CheckSession(data model.UserSession) (userID string, err error)
	VerifySession(token string) (string, error)

	// video
	VideoDescription(request model.VideoDescription) (model.VideoDescription, error)
	VideoCreate(request model.VideoFile) (model.VideoFile, error)
	GetVideos(request model.VideoFile) (model.VideoDescription, error)
	GetVideoList(videoKategori, videoKelas string) ([]model.VideoDescription, error)
	GetVideoByID(videoID string) (model.VideoDescription, error)
	GetVideoFilePathByID(videoID string) (string, error)
	GetFileNameByID(id string) (string, error)
	UpdateVideo(video model.VideoDescription) (model.VideoDescription, error)
	DeleteVideo(id string) error
	SearchVideos(query string) ([]model.VideoDescription, error)

	// history
	GetHistory(userID string) ([]model.History, error)
	AddHistory(userID, videoID string) error
	RemoveHistory(userID, historyID string) error
}

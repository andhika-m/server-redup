package video

import (
	"redup/internal/model"
)

type Repository interface {
	CreateVideo(video model.VideoFile) (model.VideoFile, error)
	CreateVideoDescription(video model.VideoDescription) (model.VideoDescription, error)
	GetVideosWithDescriptions(videoFileID string) (model.VideoDescription, error)

	// use
	GetVideoList(videoKategori, videoKelas string) ([]model.VideoDescription, error)
	GetVideoFileByID(videoFileID string) (model.VideoFile, error)
	GetVideoDataByID(videoID string) (model.VideoDescription, error)

	//try
	GetVideoFilePathByID(videoID string) (string, error)
	GetFileNameByID(id string) (string, error)

	GetVideoByID(videoID string) (model.VideoDescription, error)
	UpdateVideo(video model.VideoDescription) (model.VideoDescription, error)
	DeleteVideo(id string) error
	SearchVideos(query string) ([]model.VideoDescription, error)
	GetAllVideos() ([]model.VideoDescription, error)
}

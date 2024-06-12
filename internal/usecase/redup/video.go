package redup

import (
	"errors"
	"redup/internal/model"
)

func (ru *redupUsecase) GetVideoByID(videoID string) (model.VideoDescription, error) {
	return ru.videoRepo.GetVideoDataByID(videoID)
}

func (ru *redupUsecase) GetVideoList(videoKategori, videoKelas string) ([]model.VideoDescription, error) {
	videoData, err := ru.videoRepo.GetVideoList(videoKategori, videoKelas)
	if err != nil {
		return nil, err
	}

	for i, video := range videoData {
		videoFile, err := ru.videoRepo.GetVideoFileByID(video.VideoFileID)
		if err != nil {
			return nil, err
		}

		videoData[i].VideoFile = []model.VideoFile{videoFile}
	}

	return videoData, nil
}

func (ru *redupUsecase) VideoDescription(request model.VideoDescription) (model.VideoDescription, error) {

	createdVideoDescription, err := ru.videoRepo.CreateVideoDescription(request)
	if err != nil {
		return model.VideoDescription{}, err
	}

	return createdVideoDescription, nil
}

func (ru *redupUsecase) GetVideoFilePathByID(videoID string) (string, error) {
	return ru.videoRepo.GetVideoFilePathByID(videoID)
}

func (ru *redupUsecase) GetFileNameByID(id string) (string, error) {
	return ru.videoRepo.GetFileNameByID(id)
}

func (ru *redupUsecase) VideoCreate(request model.VideoFile) (model.VideoFile, error) {

	createdVideoFile, err := ru.videoRepo.CreateVideo(request)
	if err != nil {
		return model.VideoFile{}, err
	}

	return createdVideoFile, nil
}

func (ru *redupUsecase) UpdateVideo(video model.VideoDescription) (model.VideoDescription, error) {
	return ru.videoRepo.UpdateVideo(video)
}

func (ru *redupUsecase) DeleteVideo(id string) error {
	return ru.videoRepo.DeleteVideo(id)
}

func (ru *redupUsecase) GetVideos(request model.VideoFile) (model.VideoDescription, error) {
	videoFile, err := ru.videoRepo.GetVideosWithDescriptions(request.ID)
	if err != nil {
		return videoFile, err
	}

	if videoFile.VideoFileID != request.ID {
		return model.VideoDescription{}, errors.New("unauthorized")
	}

	return videoFile, nil
}

func (ru *redupUsecase) SearchVideos(query string) ([]model.VideoDescription, error) {
	return ru.videoRepo.SearchVideos(query)
}
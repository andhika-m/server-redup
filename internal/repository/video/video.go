package video

import (
	"errors"
	"os"
	"path/filepath"
	"redup/internal/model"
	"redup/internal/model/constant"

	"gorm.io/gorm"
)

type videoRepo struct {
	storagePath string
	db          *gorm.DB
}

func GetRepository(db *gorm.DB, storagePath string) Repository {
	return &videoRepo{
		storagePath: storagePath,
		db:          db,
	}
}

func (vr *videoRepo) GetVideoList(videoKategori, videoKelas string) ([]model.VideoDescription, error) {
	videoData := make([]model.VideoDescription, 0)

	if videoKategori == "" && videoKelas == "" {
		return nil, errors.New("invalid video category and class")
	}

	if err := vr.db.Where(model.VideoDescription{Kategori: constant.VideoCategory(videoKategori), Kelas: constant.ClassCategory(videoKelas)}).Find(&videoData).Error; err != nil {
		return nil, err
	}

	return videoData, nil
}

func (vr *videoRepo) GetVideoFileByID(videoFileID string) (model.VideoFile, error) {
	var videoFile model.VideoFile
	if err := vr.db.Where("id = ?", videoFileID).First(&videoFile).Error; err != nil {
		return model.VideoFile{}, err
	}

	return videoFile, nil
}

func (vr *videoRepo) GetFileNameByID(id string) (string, error) {
	var videoFile model.VideoFile

	if err := vr.db.Model(&model.VideoFile{}).Where("id = ?", id).First(&videoFile).Error; err != nil {
		return "", err
	}

	return videoFile.FileName, nil
}

func (vr *videoRepo) GetVideoFilePathByID(videoFileID string) (string, error) {
	var videoFile model.VideoFile
	if err := vr.db.Where("id = ?", videoFileID).First(&videoFile).Error; err != nil {
		return "", err
	}

	filePath := filepath.Join(vr.storagePath, videoFile.FileName)

	return filePath, nil
}

func (vr *videoRepo) GetVideoByID(videoID string) (model.VideoDescription, error) {
	var videoData model.VideoDescription

	if err := vr.db.Preload("VideoFile").First(&videoData, "id = ?", videoID).Error; err != nil {
		return videoData, err
	}

	return videoData, nil
}

func (vr *videoRepo) GetVideoDataByID(videoID string) (model.VideoDescription, error) {
	var video model.VideoDescription
	if err := vr.db.Where("id = ?", videoID).First(&video).Error; err != nil {
		return model.VideoDescription{}, err
	}

	videoFile, err := vr.GetVideoFileByID(video.VideoFileID)
	if err != nil {
		return model.VideoDescription{}, err
	}

	video.VideoFile = []model.VideoFile{videoFile}

	return video, nil
}

func (vr *videoRepo) CreateVideo(video model.VideoFile) (model.VideoFile, error) {
	if err := vr.db.Create(&video).Error; err != nil {
		return video, err
	}

	return video, nil
}

func (vr *videoRepo) CreateVideoDescription(video model.VideoDescription) (model.VideoDescription, error) {
	if err := vr.db.Create(&video).Error; err != nil {
		return video, err
	}

	return video, nil
}

func (vr *videoRepo) UpdateVideo(video model.VideoDescription) (model.VideoDescription, error) {
	if err := vr.db.Save(video).Error; err != nil {
		return video, err
	}

	return video, nil
}

func (vr *videoRepo) DeleteVideo(id string) error {
	var videoDescription model.VideoDescription
	if err := vr.db.First(&videoDescription, "id = ?", id).Error; err != nil {
		return err
	}

	videoFileID := videoDescription.VideoFileID

	var videoFile model.VideoFile
	if err := vr.db.First(&videoFile, "id = ?", videoFileID).Error; err != nil {
		return err
	}

	fileName := videoFile.FileName
	filePath := filepath.Join(vr.storagePath, fileName)

	if err := os.Remove(filePath); err != nil {
		return err
	}

	if err := vr.db.Delete(&videoFile).Error; err != nil {
		return err
	}

	if err := vr.db.Delete(&videoDescription).Error; err != nil {
		return err
	}

	return nil
}

func (vr *videoRepo) GetVideosWithDescriptions(videoFileID string) (model.VideoDescription, error) {
	var videos model.VideoDescription

	if err := vr.db.Where(model.VideoDescription{Kategori: constant.VideoCategory(videoFileID)}).Preload("VideoDescription").Find(&videos).Error; err != nil {
		return videos, err
	}

	return videos, nil
}

func (vr *videoRepo) SearchVideos(query string) ([]model.VideoDescription, error) {
	var videos []model.VideoDescription
	searchQuery := "%" + query + "%"

	err := vr.db.Where("judul ILIKE ?", searchQuery).Find(&videos).Error
	if err != nil {
		return nil, err
	}

	for i := range videos {
		var videoFiles []model.VideoFile
		err := vr.db.Where("id = ?", videos[i].VideoFileID).Find(&videoFiles).Error
		if err != nil {
			return nil, err
		}
		videos[i].VideoFile = videoFiles
	}

	return videos, nil
}

func (vr *videoRepo) GetAllVideos() ([]model.VideoDescription, error) {
	var videos []model.VideoDescription
	err := vr.db.Find(&videos).Error
	return videos, err
}

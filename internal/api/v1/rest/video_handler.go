package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"redup/internal/model"
	"redup/internal/model/constant"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *handler) CreateVideo(c echo.Context) error {
	file, err := c.FormFile("video")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to get video file"})
	}

	newUUID := uuid.New().String()

	uploadDir := "./public/upload/videos"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
	}

	fileName := strings.ReplaceAll(file.Filename, " ", "-")

	uniqueFileName := newUUID + "_" + fileName

	videoPath := filepath.Join(uploadDir, uniqueFileName)

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open video file"})
	}
	defer src.Close()

	dst, err := os.Create(videoPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create video file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save video file"})
	}

	videoFileData := model.VideoFile{
		ID:       newUUID,
		FileName: uniqueFileName,
	}

	createdVideo, err := h.redupUsecase.VideoCreate(videoFileData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create video"})
	}

	videoDescriptionData := model.VideoDescription{
		ID:          uuid.New().String(),
		Judul:       c.FormValue("judul"),
		Kategori:    constant.VideoCategory(c.FormValue("kategori")),
		Kelas:       constant.ClassCategory(c.FormValue("kelas")),
		Deskripsi:   c.FormValue("deskripsi"),
		VideoFileID: createdVideo.ID,
		VideoFile:   []model.VideoFile{createdVideo},
	}

	createdVideoDescription, err := h.redupUsecase.VideoDescription(videoDescriptionData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create video description"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": createdVideoDescription,
	})
}

func (h *handler) GetVideos(c echo.Context) error {
	videoKategori := c.FormValue("video_kategori")
	videoKelas := c.FormValue("video_kelas")

	videoData, err := h.redupUsecase.GetVideoList(videoKategori, videoKelas)
	if err != nil {
		fmt.Printf("got error: %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": videoData,
	})
}

func (h *handler) GetVideoByID(c echo.Context) error {
	videoID := c.Param("id")

	video, err := h.redupUsecase.GetVideoByID(videoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": video,
	})
}

func (h *handler) DownloadVideo(c echo.Context) error {
	videoID := c.Param("id")
	fileName, err := h.redupUsecase.GetFileNameByID(videoID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Video not found"})
	}

	filePath := filepath.Join("./public/upload/videos/", fileName)

	removeUuid := strings.LastIndex(fileName, "_")

	downloadFile := fileName[removeUuid+1:]

	c.Response().Header().Set("Content-Disposition", "attachment; filename="+downloadFile)
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	err = c.File(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to download video"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Video successfully downloaded"})
}

func (h *handler) EditVideo(c echo.Context) error {
	videoID := c.Param("id")

	existingVideo, err := h.redupUsecase.GetVideoByID(videoID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Video not found"})
	}

	newJudul := c.FormValue("judul")
	newKategori := constant.VideoCategory(c.FormValue("kategori"))
	newKelas := constant.ClassCategory(c.FormValue("kelas"))
	newDeskripsi := c.FormValue("deskripsi")

	if newJudul != "" {
		existingVideo.Judul = newJudul
	}
	if newKategori != "" {
		existingVideo.Kategori = newKategori
	}
	if newKelas != "" {
		existingVideo.Kelas = newKelas
	}
	if newDeskripsi != "" {
		existingVideo.Deskripsi = newDeskripsi
	}

	updatedVideo, err := h.redupUsecase.UpdateVideo(existingVideo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update video"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": updatedVideo,
	})
}

func (h *handler) DeleteVideo(c echo.Context) error {
	videoID := c.Param("id")

	if err := h.redupUsecase.DeleteVideo(videoID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete video"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Video deleted successfully"})
}

func (h *handler) SearchVideos(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "query parameter is required"})
	}

	videos, err := h.redupUsecase.SearchVideos(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"data": videos,
	}

	return c.JSON(http.StatusOK, response)
}

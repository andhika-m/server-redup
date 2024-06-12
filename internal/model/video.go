package model

import (
	"redup/internal/model/constant"
)

type VideoDescription struct {
	ID          string                 `gorm:"primaryKey" json:"id"`
	VideoFileID string                 `gorm:"index" json:"video_file_id"`
	VideoFile   []VideoFile            `json:"video_file" gorm:"-"`
	Judul       string                 `json:"judul"`
	Kategori    constant.VideoCategory `json:"kategori"`
	Kelas       constant.ClassCategory `json:"kelas"`
	Deskripsi   string                 `json:"deskripsi"`
}

type VideoFile struct {
	ID       string `json:"id" gorm:"primaryKey"`
	FileName string `json:"file_name"`
}

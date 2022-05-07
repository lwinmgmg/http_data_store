package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name      string `gorm:"not null;size:256;index;uniqueIndex:unique_idx_name_folder" json:"name"`
	FolderID  uint   `gorm:"not null;uniqueIndex:unique_idx_name_folder" json:"folder_id"`
	Path      string `gorm:"not null;unique;size:55" json:"path"`
	MimeType  string `gorm:"not null;default:application/octet-stream;size:256" json:"mime_type"`
	FirstSize int64  `gorm:"not null" json:"first_size"`
	LastSize  int64  `gorm:"not null" json:"last_size"`
	Folder    Folder
}

func (file *File) TableName() string {
	return *prefix + "_files"
}

func (file *File) Create() (*File, error) {
	err := db.Create(file).Error
	return file, err
}

func GetAllFile[T any](folder_id uint) ([]T, error) {
	var result []T
	err := db.Model(&File{}).Where("folder_id=?", folder_id).Find(&result).Error
	return result, err
}

func GetFileById[T any](folder_id, id uint) (*T, error) {
	var result T
	err := db.Model(&File{}).Where("folder_id=? AND id=?", folder_id, id).Take(&result).Error
	return &result, err
}

func DeleteFileById[T any](folder_id, id uint) (*T, error) {
	var result T
	res := db.Model(&File{}).Where("folder_id=? AND id=?", folder_id, id).Take(&result)
	if res.Error != nil {
		return &result, res.Error
	}
	return &result, res.Unscoped().Delete(&File{}).Error
}

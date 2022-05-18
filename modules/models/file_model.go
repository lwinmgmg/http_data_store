package models

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func (file *File) GetID() uint {
	return file.ID
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

func GetFileByName[T any](folder_id uint, filename string) (*T, error) {
	var result T
	err := db.Model(&File{}).Where("folder_id=? AND name=?", folder_id, filename).Take(&result).Error
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

func GetFileByNameForUpdate[T any](folder_id uint, filename string) (*T, *gorm.DB, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var result T
	var id uint
	err := db.Model(&File{}).Select("id").Where("folder_id=? AND name=?", folder_id, filename).Take(&id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, tx, nil
		}
		return nil, tx, err
	}
	file_ctx, cancel := context.WithTimeout(common_ctx, 5*time.Second)
	defer cancel()
	err = tx.WithContext(file_ctx).Model(&File{}).Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where("id=?", id).Take(&result).Error
	return &result, tx, err
}

func UpdateThroughTransactionById[T any](id uint, data map[string]any, tx *gorm.DB) (*T, error) {
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	defer tx.Commit()
	var result T
	res := tx.Model(&File{}).Where("id=?", id).Updates(data)
	if res.Error != nil {
		return nil, res.Error
	}
	err := res.Take(&result).Error
	return &result, err
}

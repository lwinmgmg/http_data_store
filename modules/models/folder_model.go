package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model
	Name   string `gorm:"size:256;not null;uniqueIndex:unique_idx_name,uid" json:"name"`
	UserID uint   `gorm:"not null;uniqueIndex:unique_idx_name" json:"user_id"`
	User   User
	Key    string `json:"key"`
}

func (folder *Folder) TableName() string {
	return *prefix + "_folders"
}

func (folder *Folder) GetID() uint {
	return folder.ID
}

func (folder *Folder) Create() (*Folder, error) {
	err := db.Create(folder).Error
	return folder, err
}

func GetAllFolder[T any](uid uint, dest *[]T) ([]T, error) {
	err := db.Model(&Folder{}).Where("user_id=?", uid).Find(dest).Error
	return nil, err
}

func GetFolderById[T any](uid, id uint, dest *T) error {
	return db.Model(&Folder{}).Where("id=? AND user_id=?", id, uid).Take(dest).Error
}

func GetFolderByName[T any](name string, dest *T) error {
	return db.Model(&Folder{}).Where("name=?", name).Take(dest).Error
}

func UpdateFolderById[T any](uid, id uint, data map[string]any) (*T, error) {
	var result T
	res := db.Model(&Folder{}).Where("id=? AND user_id=?", id, uid).Updates(data)
	if res.Error != nil {
		return &result, res.Error
	}
	err := res.Take(&result).Error
	return &result, err
}

func DeleteFolderById[T any](uid, id uint) (*T, error) {
	var result T
	res := db.Model(&Folder{}).Where("id=? AND user_id=?", id, uid).Take(&result)
	fmt.Println(result)
	if res.Error != nil {
		return &result, res.Error
	}
	return &result, res.Unscoped().Delete(&Folder{}).Error
}

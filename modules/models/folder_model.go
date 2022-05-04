package models

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Name   string `gorm:"size:256;not null;uniqueIndex:unique_idx_name,uid"`
	UserID int    `gorm:"not null;uniqueIndex:unique_idx_name" json:"user_id"`
	User   User
	Key    string `json:"key"`
}

func (folder *Folder) TableName() string {
	return *prefix + "_folders"
}

func (folder *Folder) Create() (*Folder, error) {
	err := db.Create(folder).Error
	return folder, err
}

package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	UserName string `gorm:"index;unique;size:50;column:username" json:"username"`
	Password string `gorm:"not null;size:150"`
	Key      string `gorm:"size:256" json:"key"`
}

func (user *User) Create() (*User, error) {
	err := db.Create(user).Error
	return user, err
}

func GetAllUser() ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func GetUserById(id uint) (*User, error) {
	var user User
	return &user, db.Where("id=?", id).Find(&user).Error
}

func DeleteUserById(id uint) (*User, error) {
	var user User
	res := db.Where("id=?", id).Find(&user)
	if res.Error != nil {
		return &user, res.Error
	}
	return &user, res.Unscoped().Delete(&user).Error
}

func UpdateUserById(id uint, data map[string]interface{}) (*User, error) {
	var user User
	err := db.Where("id=?", id).Find(&user).Updates(data).Error
	return &user, err
}

func GetUserForUpdateById(id uint) (*User, *gorm.DB) {
	var user User
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	tx.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
	).Where("id=?", id).Find(&user)
	return &user, tx
}

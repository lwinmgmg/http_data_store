package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	UserName string `gorm:"" json:"username"`
	Password string
}

func (user *User) Create() *User {
	db.Create(user)
	return user
}

func GetAllUser() []User {
	var users []User
	db.Find(&users)
	return users
}

func GetUserById(id uint) *User {
	var user User
	db.Where("id=?", id).Find(&user)
	return &user
}

func DeleteById(id uint) *User {
	var user User
	db.Where("id=?", id).Delete(&user)
	return &user
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

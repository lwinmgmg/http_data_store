package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"" json:"username"`
	Password string `json:"password"`
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

func GetUserById(id uint) (*User, *gorm.DB) {
	var user User
	db := db.Where("id=?", id).Find(&user)
	return &user, db
}

func DeleteById(id uint) *User {
	var user User
	db.Where("id=?", id).Delete(&user)
	return &user
}

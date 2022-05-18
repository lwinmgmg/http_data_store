package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	UserName string `gorm:"index;unique;size:50;column:username" json:"username"`
	Password string `gorm:"not null;size:150" json:"password"`
	Key      string `gorm:"size:256" json:"key"`
}

func (user *User) TableName() string {
	return *prefix + "_users"
}

func (user *User) GetID() uint {
	return user.ID
}

func (user *User) Create(output any) error {
	err := db.Model(user).Create(user).Error
	if err != nil {
		return err
	}
	return db.Model(user).Take(output, user.ID).Error
}

func (user *User) GetAll(users any) error {
	return db.Model(user).Find(users).Error
}

func (user *User) GetByID(id int, output any) error {
	return db.Model(user).First(output, id).Error
}

// func (user *User) GetByIDs()

func GetAllUser() ([]User, error) {
	var users []User
	err := db.Select("id", "username", "key").Find(&users).Error
	return users, err
}

func GetUserById(id uint) (*User, error) {
	var user User
	return &user, db.Select("id", "username", "key").Where("id=?", id).First(&user).Error
}

func DeleteUserById(id uint) (*User, error) {
	var user User
	res := db.Select("id").Where("id=?", id).Find(&user)
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

func GetUserIdByUserName(username string) uint {
	var user User
	db.Select("id").Where("username=?", username).First(&user)
	return user.ID
}

func GetUserByUserName(username string) *User {
	var user User
	db.Select("id", "password").Where("username=?", username).First(&user)
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

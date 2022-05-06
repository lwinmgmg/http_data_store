package models

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lwinmgmg/http_data_store/helper"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

var userMap map[string]string = map[string]string{
	"username": "lwinmgmg",
	"password": "defaultpassword",
}

func TestCreate(t *testing.T) {
	username := userMap["username"]
	password := userMap["password"]
	user := views.UserCreate{
		UserName: &username,
		Password: &password,
	}
	if err := user.Validate(); err != nil {
		t.Error("Error on user validating :", err)
	}
	user1 := &User{
		UserName: *user.UserName,
		Password: *user.Password,
	}
	user1, err := user1.Create()
	if err != nil {
		t.Error(err)
	}
	if user1.ID == 0 {
		t.Errorf("User ID should be greater than 0, resulting '%v'", user1.ID)
	}
}

func TestGetAllUser(t *testing.T) {
	users, err := GetAllUser()
	fmt.Println(len(users))
	if err != nil {
		t.Error(err)
	}
	if len(users) == 0 {
		t.Error("User just added but no user")
	}
	if users[0].UserName != userMap["username"] {
		t.Errorf("Expecting username as '%v', getting '%v'", userMap["username"], users[0].UserName)
	}
}

func TestGetUserIdByUserName(t *testing.T) {
	users, err := GetAllUser()
	if err != nil {
		t.Error(err)
	}
	id := users[0].ID
	uid := GetUserIdByUserName(userMap["username"])
	if id != uid {
		t.Errorf("Expecting user id %v, getting %v", id, uid)
	}
}

func TestGetUserById(t *testing.T) {
	id := GetUserIdByUserName(userMap["username"])
	user, err := GetUserById(id)
	if err != nil {
		t.Error(err)
	}
	if user.ID != id {
		t.Errorf("User ID must be equal %v, getting %v", id, user.ID)
	}
}

func TestUpdateUserById(t *testing.T) {
	username := "lwinmgmg1"
	password := "password"
	key := "abcd"
	userUpdate1 := views.UserUpdate{
		UserName: &username,
		Password: &password,
		Key:      &key,
	}
	updateMap, err := userUpdate1.Validate()
	if err != nil {
		t.Error("Getting error on UserUpdate view validate :", err)
	}
	id := GetUserIdByUserName(userMap["username"])
	user, err := UpdateUserById(id, updateMap)
	if err != nil {
		t.Error("Error on update :", err)
	}
	if user.ID != id {
		t.Errorf("Expecting ID %v, getting %v", id, user.ID)
	}
	if user.UserName != username {
		t.Errorf("Expecting username %v, gettint %v", username, user.UserName)
	}
	if user.Password != helper.HexString(password) {
		t.Errorf("Expecting username %v, gettint %v", user.Password, helper.HexString(password))
	}
	userMap["username"] = username
	userMap["password"] = password
}

func TestGetUserForUpdateById(t *testing.T) {
	var user User
	uid := GetUserIdByUserName(userMap["username"])
	_, tx := GetUserForUpdateById(uid)
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	context, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	if err := db.WithContext(context).Where("id=?", uid).Delete(&user).Error; err == nil {
		t.Error("Locking does not work")
	}
	tx.Commit()
}

func TestDeleteUserById(t *testing.T) {
	users, err := GetAllUser()
	if err != nil {
		t.Error(err)
	}
	id := users[0].ID
	user, err := DeleteUserById(id)
	if err != nil {
		t.Error(err)
	}
	if user.ID != id {
		t.Errorf("Expecting user id as %v, getting %v", id, user.ID)
	}
	if _, err = GetUserById(id); err == nil {
		t.Errorf("User is already deleted but not getting error on get by id")
	}
}

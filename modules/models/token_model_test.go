package models

import (
	"testing"

	"github.com/lwinmgmg/http_data_store/modules/views"
)

func TestGenerateToken(t *testing.T) {
	username := userMap["username"] + "token"
	password := userMap["password"]
	userCreate := views.UserCreate{
		UserName: &username,
		Password: &password,
	}
	if err := userCreate.Validate(); err != nil {
		t.Error("Error on creating user for token :", err)
	}
	user := User{
		UserName: *userCreate.UserName,
		Password: *userCreate.Password,
	}
	user.Create()
	token, err := GenerateToken(user.UserName)
	if err != nil {
		t.Error("error on generating token :", token)
	}
	id, err := ValidateToken(token)
	if err != nil {
		t.Error("Error on validating token :", err)
	}
	if _, err := DeleteUserById(id); err != nil {
		t.Error("Error on deleting created user :", err)
	}
}

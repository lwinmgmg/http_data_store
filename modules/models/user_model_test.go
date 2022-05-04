package models

import (
	"testing"
)

func TestGetAllUser(t *testing.T) {
	user1 := &User{UserName: "lwinmgmg"}
	user1, err := user1.Create()
	if err != nil {
		t.Error(err)
	}
	users, err := GetAllUser()
	if err != nil {
		t.Error(err)
	}
	if len(users) == 0 {
		t.Error("User just added but no user")
	}
}

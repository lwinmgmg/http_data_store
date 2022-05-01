package models

import "testing"

func Test_GetAllUser(t *testing.T) {
	users := GetAllUser()
	if len(users) > 0 {
		t.Errorf("Expected no user, got %v users", len(users))
	}
	user := User{
		UserName: "LMM",
		Password: "ABC",
	}
	user.Create()
	users = GetAllUser()
	if len(users) != 1 {
		t.Errorf("Expected one user, got %v users", len(users))
	}
	DeleteById(user.ID)
}

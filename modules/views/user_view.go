package views

import (
	"regexp"

	"github.com/lwinmgmg/http_data_store/helper"
)

type UserRead struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Key      string `json:"key"`
}

type UserCreate struct {
	UserName *string `json:"username" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

func (user *UserCreate) Validate() error {
	if user.UserName == nil {
		return helper.NewCustomError("Username can't be empty", helper.ValidationError)
	}
	isMatch, _ := regexp.MatchString("^[a-z]*[0-9a-z]*", *user.UserName)
	if !isMatch {
		return helper.NewCustomError("Username can be use lowercase character and number", helper.ValidationError)
	}
	if user.Password == nil {
		return helper.NewCustomError("Password can't be empty", helper.ValidationError)
	}
	if len(*(user.Password)) < 8 {
		return helper.NewCustomError("Password must have at least 8 character", helper.ValidationError)
	}
	*user.Password = helper.HexString(*user.Password)
	return nil
}

type UserUpdate struct {
	UserName *string `json:"username"`
	Password *string `json:"password"`
	Key      *string `json:"key"`
}

func (user *UserUpdate) Validate() (map[string]interface{}, error) {
	userMap := make(map[string]interface{}, 3)
	if user.UserName != nil {
		isMatch, _ := regexp.MatchString("^[a-z]*[0-9a-z]*", *user.UserName)
		if !isMatch {
			return nil, helper.NewCustomError("Username can be use lowercase character and number", helper.ValidationError)
		}
		userMap["username"] = *user.UserName
	}
	if user.Password != nil {
		if len(*user.Password) < 8 {
			return nil, helper.NewCustomError("Password must have at least 8 character", helper.ValidationError)
		}
		userMap["password"] = helper.HexString(*user.Password)
	}
	if user.Key != nil {
		userMap["key"] = *user.Key
	}
	return userMap, nil
}

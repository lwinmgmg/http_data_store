package views

import (
	"regexp"

	"github.com/lwinmgmg/http_data_store/helper"
)

var (
	UserReg *regexp.Regexp
)

type UserRead struct {
	ID       uint   `json:"id"`
	UserName string `gorm:"column:username" json:"username"`
	Key      string `json:"key"`
}

type UserCreate struct {
	UserName *string `json:"username" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

func (user *UserCreate) Validate() error {
	if user.UserName == nil {
		return helper.NewCustomError(helper.ValidationError, "Username can't be empty")
	}
	UserReg = regexp.MustCompile(`^[a-z][a-z0-9]+$`)
	isMatch := UserReg.MatchString(*user.UserName)
	if !isMatch {
		return helper.NewCustomError(helper.ValidationError, "Username can be use lowercase character and number")
	}
	if user.Password == nil {
		return helper.NewCustomError(helper.ValidationError, "Password can't be empty")
	}
	if len(*(user.Password)) < 8 {
		return helper.NewCustomError(helper.ValidationError, "Password must have at least 8 character")
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
		UserReg = regexp.MustCompile(`^[a-z]*[0-9a-z]*`)
		isMatch := UserReg.MatchString(*user.UserName)
		if !isMatch {
			return nil, helper.NewCustomError(helper.ValidationError, "Username can be use lowercase character and number")
		}
		userMap["username"] = *user.UserName
	}
	if user.Password != nil {
		if len(*user.Password) < 8 {
			return nil, helper.NewCustomError(helper.ValidationError, "Password must have at least 8 character")
		}
		userMap["password"] = helper.HexString(*user.Password)
	}
	if user.Key != nil {
		userMap["key"] = *user.Key
	}
	return userMap, nil
}

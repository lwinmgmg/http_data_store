package views

import (
	"regexp"

	"github.com/lwinmgmg/http_data_store/helper"
)

type UserCreate struct {
	UserName *string `json:"username"`
	Password *string `json:"password"`
}

func (user *UserCreate) Validate() error {
	if user.UserName == nil {
		return helper.NewCustomError("Username can't be empty", helper.ValidationError)
	}
	isMatch, _ := regexp.MatchString("^[a-z]*[0-9a-z]*", "peach")
	if !isMatch {
		return helper.NewCustomError("Username can be use lowercase character and number", helper.ValidationError)
	}
	if user.Password == nil {
		return helper.NewCustomError("Password can't be empty", helper.ValidationError)
	}
	if len(*(user.Password)) < 8 {
		return helper.NewCustomError("Password must have at least 8 character", helper.ValidationError)
	}
	return nil
}

package models

import (
	"time"

	"github.com/lwinmgmg/http_data_store/helper"
)

type Claim struct {
	Issuer   string `json:"iss,omitempty"`
	Audience string `json:"aud,omitempty"`
	IssuedAt int64  `json:"iat,omitempty"`
}

func (claim *Claim) Valid() (uint, error) {
	id := GetUserIdByUserName(claim.Audience)
	if id == 0 {
		return 0, helper.NewCustomError("Unknown issuer", helper.AuthenticationError)
	}
	if claim.IssuedAt < time.Now().Unix()+int64(env.HDS_TOKEN_DEFAULT_TIMEOUT) {
		return 0, helper.NewCustomError("Token has already expired", helper.AuthenticationError)
	}
	if claim.Audience != env.HDS_TOKEN_AUDIENCE {
		return 0, helper.NewCustomError("Unknown audience", helper.AuthenticationError)
	}
	return id, nil
}

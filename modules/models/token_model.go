package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lwinmgmg/http_data_store/helper"
)

type Claim struct {
	Issuer   string `json:"iss,omitempty"`
	Audience string `json:"aud,omitempty"`
	IssuedAt int64  `json:"iat,omitempty"`
}

func NewClaim(issuer string) Claim {
	return Claim{
		Issuer:   issuer,
		Audience: env.HDS_TOKEN_AUDIENCE,
		IssuedAt: time.Now().Unix(),
	}
}

func (claim *Claim) GetIssuer() (uint, error) {
	id := GetUserIdByUserName(claim.Issuer)
	if id == 0 {
		return 0, helper.NewCustomError("Unknown issuer", helper.AuthenticationError)
	}
	return id, nil
}

func (claim *Claim) Valid() error {
	if claim.IssuedAt > time.Now().Unix()+int64(env.HDS_TOKEN_DEFAULT_TIMEOUT) {
		return helper.NewCustomError("Token has already expired", helper.AuthenticationError)
	}
	if claim.Audience != env.HDS_TOKEN_AUDIENCE {
		return helper.NewCustomError("Unknown audience", helper.AuthenticationError)
	}
	return nil
}

func KeyFunc() jwt.Keyfunc {
	return func(tkn *jwt.Token) (interface{}, error) {
		return []byte(env.HDS_TOKEN_KEY), nil
	}
}

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = &Claim{
		Issuer:   username,
		Audience: env.HDS_TOKEN_AUDIENCE,
		IssuedAt: time.Now().Unix(),
	}
	return token.SignedString([]byte(env.HDS_TOKEN_KEY))
}

func ValidateToken(tokenStr string) (uint, error) {
	claim := Claim{}
	if _, err := jwt.ParseWithClaims(tokenStr, &claim, KeyFunc()); err != nil {
		fmt.Println(claim)
		return 0, err
	}
	return claim.GetIssuer()
}

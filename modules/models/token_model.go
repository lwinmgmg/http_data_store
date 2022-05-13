package models

import (
	"context"
	"fmt"
	"strconv"
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
		return 0, helper.NewCustomError(helper.AuthenticationError, "Unknown issuer")
	}
	return id, nil
}

func (claim *Claim) Valid() error {
	if claim.IssuedAt > time.Now().Unix()+int64(env.HDS_TOKEN_DEFAULT_TIMEOUT) {
		return helper.NewCustomError(helper.AuthenticationError, "Token has already expired")
	}
	if claim.Audience != env.HDS_TOKEN_AUDIENCE {
		return helper.NewCustomError(helper.AuthenticationError, "Unknown audience")
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
	ctx, cancel := context.WithTimeout(common_ctx, time.Millisecond*100)
	defer cancel()
	v, err := redis_client.Get(ctx, tokenStr).Result()
	if err == nil {
		out, err := strconv.Atoi(v)
		if err == nil {
			return uint(out), nil
		}
	}
	claim := Claim{}
	if _, err := jwt.ParseWithClaims(tokenStr, &claim, KeyFunc()); err != nil {
		return 0, err
	}
	if err := claim.Valid(); err != nil {
		return 0, err
	}
	uid, err := claim.GetIssuer()
	if err != nil {
		return 0, err
	}
	var duration int64 = claim.IssuedAt + int64(env.HDS_TOKEN_DEFAULT_TIMEOUT) - time.Now().Unix()
	if duration < 0 {
		return 0, helper.NewCustomError(helper.AuthenticationError, "Token has already expired")
	}
	if err := redis_client.Set(ctx, tokenStr, uid, time.Second*time.Duration(duration)).Err(); err != nil {
		fmt.Println(err, "Error on key set")
	}
	return uid, nil
}

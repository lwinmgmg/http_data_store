package middlewares

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/helper"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

const (
	BEARER_SCHEMA string = "Bearer"
	BASIC_SCHEMA  string = "Basic"
)

func AuthParser(value string) (string, string, error) {
	tokens := strings.Split(value, " ")
	if len(tokens) != 2 {
		return "", "", fmt.Errorf("wrong value in authorization value : %v", value)
	}
	return tokens[0], tokens[1], nil
}

func BasicDecoder(value string) (views.UserCreate, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return views.UserCreate{}, err
	}
	decodedStrList := strings.SplitN(string(decoded), ":", 2)
	if len(decodedStrList) != 2 {
		return views.UserCreate{}, fmt.Errorf("can't split string from base64 : %v", value)
	}
	return views.UserCreate{
		UserName: &decodedStrList[0],
		Password: &decodedStrList[1],
	}, nil
}

func InternalServerErrorHandler(ctx *gin.Context, err interface{}) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
		"detail": err,
	})
}

func BasicAuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		basicKey := ctx.GetHeader("Authorization")
		if basicKey != "" {
			tokenType, tokenStr, err := AuthParser(basicKey)
			if err != nil || tokenType != BASIC_SCHEMA {
				ctx.Header("WWW-Authenticate", "Authorization Required")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
					"detail": "Authorization Required!",
				})
				return
			}
			userInput, err := BasicDecoder(tokenStr)
			if err != nil {
				ctx.Header("WWW-Authenticate", "Authorization Required")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
					"detail": err.Error(),
				})
				return
			}
			user := models.GetUserByUserName(*userInput.UserName)
			if user.ID == 0 {
				ctx.Header("WWW-Authenticate", "Authorization Required")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
					"detail": "Unauthorize user!",
				})
				return
			}
			if user.Password != helper.HexString(*userInput.Password) {
				ctx.Header("WWW-Authenticate", "Authorization Required")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
					"detail": "Authentication failed!",
				})
				return
			}
			ctx.Set("uid", user.ID)
			ctx.Set("username", *userInput.UserName)
			return
		}
		ctx.Header("WWW-Authenticate", "Authorization Required")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"detail": "Authorization Required!",
		})
	}
}

func JWTAuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerKey := ctx.GetHeader("Authorization")
		if bearerKey != "" {
			tokenType, tokenStr, err := AuthParser(bearerKey)
			if err != nil || tokenType != BEARER_SCHEMA {
				ctx.Header("WWW-Authenticate", "Authorization Required")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
					"detail": "Authorization Required!",
				})
				return

			}
			uid, err := models.ValidateToken(tokenStr)
			if err != nil {
				ctx.Header("WWW-Authenticate", "Authorization Required")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
					"detail": err.Error(),
				})
				return
			}
			ctx.Set("uid", uid)
			return
		}
		ctx.Header("WWW-Authenticate", "Authorization Required")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"detail": "Authorization Required!",
		})
	}
}

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

func (cmgr *ControllerManager) GetAllUser(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	if uid != 1 {
		ctx.Header("WWW-Authenticate", "Authorization Required")
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "User not allow"})
		return
	}
	users, err := models.GetAllUser()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, users)
}

func (cmgr *ControllerManager) GetUserById(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	if uid != 1 {
		ctx.Header("WWW-Authenticate", "Authorization Required")
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "User not allow"})
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
	}
	user, err := models.GetUserById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": "User not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

func (cmgr *ControllerManager) Create(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	if uid != 1 {
		ctx.Header("WWW-Authenticate", "Authorization Required")
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "User not allow"})
		return
	}
	var user views.UserCreate
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"detail": err.Error()})
		return
	}
	if err := user.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	dbUser := &models.User{
		UserName: *user.UserName,
		Password: *user.Password,
	}
	dbUser, err := dbUser.Create()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, dbUser)
}

func (cmgr *ControllerManager) DeleteUserById(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	if uid != 1 {
		ctx.Header("WWW-Authenticate", "Authorization Required")
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "User not allow"})
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	user, err := models.DeleteUserById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": "User not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

func (cmgr *ControllerManager) UpdateUserById(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	if uid != 1 {
		ctx.Header("WWW-Authenticate", "Authorization Required")
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "User not allow"})
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	var updateUser views.UserUpdate
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"detail": err.Error()})
		return
	}
	data, err := updateUser.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	user, err := models.UpdateUserById(uint(id), data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": "User not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

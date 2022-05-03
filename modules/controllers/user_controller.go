package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

func (cmgr *CManager) GetAllUser(ctx *gin.Context) {
	users, err := models.GetAllUser()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, users)
}

func (cmgr *CManager) GetUserById(ctx *gin.Context) {
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

func (cmgr *CManager) Create(ctx *gin.Context) {
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

func (cmgr *CManager) DeleteById(ctx *gin.Context) {
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

func (cmgr *CManager) UpdateById(ctx *gin.Context) {
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

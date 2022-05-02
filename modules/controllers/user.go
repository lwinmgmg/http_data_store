package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

func (cmgr *CManager) GetAllUser(ctx *gin.Context) {
	var users []models.User = models.GetAllUser()
	ctx.IndentedJSON(http.StatusOK, users)
}

func (cmgr *CManager) GetUserById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
	}
	user := models.GetUserById(uint(id))
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": "User not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

func (cmgr *CManager) Create(ctx *gin.Context) {
	user := &views.UserCreate{}
	ctx.ShouldBindJSON(user)
	if err := user.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	dbUser := &models.User{
		UserName: *user.UserName,
		Password: *user.Password,
	}
	dbUser = dbUser.Create()
	ctx.IndentedJSON(http.StatusOK, dbUser)
}

func (cmgr *CManager) DeleteById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
	}
	user := models.DeleteById(uint(id))
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": "User not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (cMgr ControllerManager) GetAllFile(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	idStr := ctx.Param("folder_id")
	folder_id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	fmt.Println(uid, folder_id)
	// models.GetAllFile[]()
}

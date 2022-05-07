package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

func (cMgr *ControllerManager) GetAllFolder(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	var folders []views.FolderRead
	if _, err := models.GetAllFolder(uid, &folders); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, folders)
}

func (cMgr *ControllerManager) CreateFolder(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	var folderCreate views.FolderCreate
	if err := ctx.ShouldBind(&folderCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	folder := &models.Folder{
		Name:   folderCreate.Name,
		Key:    folderCreate.Key,
		UserID: uid,
	}
	folder, err := folder.Create()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, folder)
}

func (cMgr *ControllerManager) GetFolderById(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	var folder views.FolderRead
	if err := models.GetFolderById(uid, uint(id), &folder); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, &folder)
}

func (cMgr *ControllerManager) UpdateFolderById(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	var folder views.FolderUpdate
	if err := ctx.ShouldBind(&folder); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	data := folder.ToMap()
	updatedFolder, err := models.UpdateFolderById[views.FolderRead](uid, uint(id), data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, updatedFolder)
}

func (cMgr *ControllerManager) DeleteFolderById(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(uint)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	folder, err := models.DeleteFolderById[views.FolderRead](uid, uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, folder)
}

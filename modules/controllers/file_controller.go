package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/cron"
	"github.com/lwinmgmg/http_data_store/helper"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

func GetFileRequirements(ctx *gin.Context) (uint, uint, error) {
	uid := ctx.MustGet("uid").(uint)
	idStr := ctx.Param("folder_id")
	folder_id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, 0, helper.NewCustomError(helper.ValidationError, "Error on converting integer : %v", err)
	}
	var folderRead views.FolderRead
	err = models.GetFolderById(uid, uint(folder_id), &folderRead)
	if err != nil {
		return 0, 0, helper.NewCustomError(helper.DatabaseError, "Error on getting folder by id : %v", err)
	}
	if folderRead.ID == 0 {
		return 0, 0, helper.NewCustomError(helper.DataDoesNotExistError, "Folder ID %v not found", folder_id)
	}
	return uid, uint(folder_id), nil
}

func (cMgr ControllerManager) GetAllFile(ctx *gin.Context) {
	_, folder_id, err := GetFileRequirements(ctx)
	if err != nil {
		ctx.JSON(helper.ErrorResponse(err))
		return
	}
	files, err := models.GetAllFile[views.FileRead](uint(folder_id))
	if err != nil {
		newErr := helper.NewCustomError(helper.DatabaseError, "Error on fetching files : %v", err)
		ctx.JSON(helper.ErrorResponse(newErr))
		return
	}
	ctx.IndentedJSON(http.StatusOK, files)
}

func (cMgr ControllerManager) CreateFile(ctx *gin.Context) {
	_, folder_id, err := GetFileRequirements(ctx)
	if err != nil {
		ctx.JSON(helper.ErrorResponse(err))
		return
	}
	myForm := struct {
		FileName string `form:"filename" binding:"required"`
		MimeType string `form:"mime_type"`
	}{MimeType: "application/octet-stream"}
	if err := ctx.Bind(&myForm); err != nil {
		newErr := helper.NewCustomError(helper.ValidationError, "Error on binding form data : %v", err)
		ctx.JSON(helper.ErrorResponse(newErr))
		return
	}
	oldFile, tx, err := models.GetFileByNameForUpdate[models.File](folder_id, myForm.FileName)
	if err != nil {
		ctx.JSON(helper.ErrorResponse(err))
		return
	}
	fileInput, err := ctx.FormFile("file")
	if err != nil {
		newErr := helper.NewCustomError(helper.ValidationError, "Error on getting file : %v", err)
		ctx.JSON(helper.ErrorResponse(newErr))
		return
	}
	reader, err := fileInput.Open()
	if err != nil {
		ctx.JSON(helper.ErrorResponse(err))
		return
	}
	defer reader.Close()
	writer := helper.NewWriterManager(reader, fileInput.Size)
	size1, size2, err := writer.WriteOriginal()
	if err != nil {
		ctx.JSON(helper.ErrorResponse(err))
		return
	}
	gc := cron.GetGarbageChannelWriter()
	if oldFile != nil {
		data := map[string]any{
			"name":       myForm.FileName,
			"mime_type":  myForm.MimeType,
			"path":       writer.FileName,
			"first_size": size1,
			"last_size":  size2,
		}
		UpdatedFolder, err := models.UpdateThroughTransactionById[views.FileRead](oldFile.ID, data, tx)
		if err != nil {
			ctx.JSON(helper.ErrorResponse(err))
			gc <- writer.FileName
			return
		}
		gc <- oldFile.Path
		ctx.JSON(http.StatusOK, UpdatedFolder)
	} else {
		file := &models.File{
			FolderID:  folder_id,
			Name:      myForm.FileName,
			Path:      writer.FileName,
			MimeType:  myForm.MimeType,
			FirstSize: size1,
			LastSize:  size2,
		}
		file, err = file.Create()
		if err != nil {
			newErr := helper.NewCustomError(helper.DatabaseError, "Error on file record create : %v", err)
			ctx.JSON(helper.ErrorResponse(newErr))
			gc <- writer.FileName
			return
		}
		ctx.JSON(http.StatusCreated, file)
	}
}

func (cMgr *ControllerManager) GetFileById(ctx *gin.Context) {
	_, folder_id, err := GetFileRequirements(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	idStr := ctx.Param("file_id")
	file_id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	file, err := models.GetFileById[models.File](folder_id, uint(file_id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	reader, closers, err := helper.ReadFile(file.Path, file.FirstSize, file.LastSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"detail": err.Error()})
		return
	}
	defer helper.MultiCloser(closers...)
	contentLength := file.FirstSize + file.LastSize
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%v"`, file.Name),
	}
	ctx.DataFromReader(http.StatusOK, contentLength, file.MimeType, reader, extraHeaders)
}

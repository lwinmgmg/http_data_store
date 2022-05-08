package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/helper"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

func (ctr *ControllerManager) ServeTempUrl(ctx *gin.Context) {
	var query views.TempUrlQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": err.Error()})
		return
	}
	if time.Now().Unix() > query.ExpireTime {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "expired"})
		return
	}
	keyList := make([]string, 0, 2)
	var folder views.FolderRead
	if err := models.GetFolderByName(query.FolderName, &folder); err != nil {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": err.Error()})
		return
	}
	var signal bool = false
	if folder.Key == "" {
		user, err := models.GetUserById(folder.UserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": err.Error()})
			return
		}
		if user.Key == "" {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "Unauthorize"})
			return
		}
		keyList = append(keyList, user.Key)
	} else {
		keyList = append(keyList, folder.Key)
		signal = true
	}
	for i := 0; i < 2; i++ {
		values1 := helper.HexString(fmt.Sprintf("%v:%v:%v:%v", query.File, query.FolderName, query.ExpireTime, keyList[i]))
		if values1 == query.HashKey {
			file, err := models.GetFileByName[models.File](folder.ID, query.File)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": err.Error()})
				return
			}
			reader, closers, err := helper.ReadFile(file.Path, file.FirstSize, file.LastSize)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, map[string]string{"detail": err.Error()})
				return
			}
			defer helper.MultiCloser(closers...)
			contentLength := file.FirstSize + file.LastSize
			var filename string
			if query.FileName != "" {
				filename = query.FileName
			} else {
				filename = file.Name
			}
			extraHeaders := map[string]string{
				"Content-Disposition": fmt.Sprintf(`attachment; filename="%v"`, filename),
				"Cache-Control":       "private,max-age=31536000,immutable",
			}
			ctx.DataFromReader(http.StatusOK, contentLength, file.MimeType, reader, extraHeaders)
			return
		}
		if signal {
			user, err := models.GetUserById(folder.UserID)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": err.Error()})
				return
			}
			if user.Key == "" {
				ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "Unauthorize"})
				return
			}
			keyList = append(keyList, user.Key)
		}
	}
	ctx.JSON(http.StatusUnauthorized, map[string]string{"detail": "invalid hashkey"})
}

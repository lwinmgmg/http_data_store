package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/middlewares"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

func RegisterFileRoutes(router *gin.RouterGroup, cMgr *controllers.ControllerManager) {
	router.Use(middlewares.JWTAuthenticationMiddleware())
	router.GET("/folders/:folder_id/files", cMgr.GetAllFile)
	router.POST("/folders/:folder_id/files", cMgr.CreateFile)
	router.GET("/folders/:folder_id/files/:file_id", cMgr.GetFileById)
}

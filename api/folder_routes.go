package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/middlewares"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

func RegisterFolderRoutes(router *gin.RouterGroup, cMgr *controllers.ControllerManager) {
	router.Use(middlewares.JWTAuthenticationMiddleware())
	router.GET("/folders", cMgr.GetAllFolder)
	router.POST("/folders", cMgr.CreateFolder)
	router.GET("/folders/:folder_id", cMgr.GetFolderById)
	router.DELETE("/folders/:folder_id", cMgr.DeleteFolderById)
	router.PUT("/folders/:folder_id", cMgr.UpdateFolderById)
}

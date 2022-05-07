package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/middlewares"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

func RegisterRoutes(app *gin.Engine) {
	cMgr := &controllers.ControllerManager{}

	//Admin Routes
	router := app.Group("/api")
	router.Use(middlewares.BasicAuthenticationMiddleware())
	router.GET("/token", cMgr.GetToken)
	router.GET("/users", cMgr.GetAllUser)
	router.POST("/users", cMgr.Create)
	router.GET("/users/:id", cMgr.GetUserById)
	router.PUT("/users/:id", cMgr.UpdateUserById)
	router.DELETE("/users/:id", cMgr.DeleteUserById)

	clientRouter := app.Group("/api/client")
	clientRouter.Use(middlewares.JWTAuthenticationMiddleware())

	//Client File Routers
	clientRouter.GET("/folders/:folder_id/files", cMgr.GetAllFile)
	clientRouter.POST("/folders/:folder_id/files", cMgr.CreateFile)
	clientRouter.GET("/folders/:folder_id/files/:file_id", cMgr.GetFileById)

	//Client Folder Routes
	clientRouter.GET("/folders", cMgr.GetAllFolder)
	clientRouter.POST("/folders", cMgr.CreateFolder)
	clientRouter.GET("/folders/:folder_id", cMgr.GetFolderById)
	clientRouter.DELETE("/folders/:folder_id", cMgr.DeleteFolderById)
	clientRouter.PUT("/folders/:folder_id", cMgr.UpdateFolderById)
}

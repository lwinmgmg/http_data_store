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

	//Client Routes
	clientRouter := app.Group("/api/client")
	clientRouter.Use(middlewares.JWTAuthenticationMiddleware())
	clientRouter.GET("/folders", cMgr.GetAllFolder)
	clientRouter.POST("/folders", cMgr.CreateFolder)
	clientRouter.GET("/folders/:id", cMgr.GetFolderById)
	clientRouter.DELETE("/folders/:id", cMgr.DeleteFolderById)
	clientRouter.PUT("/folders/:id", cMgr.UpdateFolderById)
}

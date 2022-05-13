package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/middlewares"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

var app *gin.Engine

func init() {
	if app == nil {
		app = gin.New()
	}
	app.Use(gin.Logger(), gin.CustomRecovery(middlewares.InternalServerErrorHandler))
}

func GetApp() *gin.Engine {
	return app
}

func RegisterRoutes() {
	cMgr := &controllers.ControllerManager{}
	v1Router := app.Group("/v1")
	//Admin Routes
	router := v1Router.Group("/api")
	RegisterUserRoutes(router, cMgr)

	//Token Routes
	RegisterTokenRoutes(router, cMgr)

	clientRouter := v1Router.Group("/api/client")
	//Client Folder Routes
	RegisterFolderRoutes(clientRouter, cMgr)

	//Client File Routes
	RegisterFileRoutes(clientRouter, cMgr)

	//Temp URL route
	tempRouter := v1Router.Group("/temp")
	RegisterTempUrlRoutes(tempRouter, cMgr)

}

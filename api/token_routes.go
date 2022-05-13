package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/middlewares"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

func RegisterTokenRoutes(router *gin.RouterGroup, cMgr *controllers.ControllerManager) {
	router.Use(middlewares.BasicAuthenticationMiddleware())
	router.GET("/token", cMgr.GetToken)
}

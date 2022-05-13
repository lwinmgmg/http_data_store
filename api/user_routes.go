package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/middlewares"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

func RegisterUserRoutes(router *gin.RouterGroup, cMgr *controllers.ControllerManager) {
	router.Use(middlewares.BasicAuthenticationMiddleware())
	router.GET("/users", cMgr.GetAllUser)
	router.POST("/users", cMgr.Create)
	router.GET("/users/:id", cMgr.GetUserById)
	router.PUT("/users/:id", cMgr.UpdateUserById)
	router.DELETE("/users/:id", cMgr.DeleteUserById)
}

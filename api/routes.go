package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

func RegisterRoutes(app *gin.Engine) {
	router := app.Group("/api")
	cMgr := &controllers.CManager{}
	router.GET("/users", cMgr.GetAllUser)
	router.POST("/users", cMgr.Create)
	router.GET("/users/:id", cMgr.GetUserById)
	router.PUT("/users/:id", cMgr.UpdateById)
	router.DELETE("/users/:id", cMgr.DeleteById)
}

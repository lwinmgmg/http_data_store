package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/modules/controllers"
)

func RegisterTempUrlRoutes(router *gin.RouterGroup, cMgr *controllers.ControllerManager) {
	router.GET("", cMgr.ServeTempUrl)
}

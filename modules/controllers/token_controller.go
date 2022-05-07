package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/modules/models"
	"github.com/lwinmgmg/http_data_store/modules/views"
)

func (ctr *ControllerManager) GetToken(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)
	tokenStr, _ := models.GenerateToken(username)
	token := views.NewBearerTokenRead(tokenStr)
	ctx.IndentedJSON(http.StatusOK, &token)
}

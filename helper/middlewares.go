package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalServerErrorHandler(ctx *gin.Context, err interface{}) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
		"detail": err,
	})
}

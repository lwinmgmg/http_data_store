package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(app *gin.Engine) {
	router := app.Group("/api")
}

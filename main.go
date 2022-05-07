package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/api"
	"github.com/lwinmgmg/http_data_store/environ"
	"github.com/lwinmgmg/http_data_store/middlewares"
	"github.com/lwinmgmg/http_data_store/modules/models"
)

var (
	HOST string           = "localhost"
	PORT int              = 8000
	env  *environ.Environ = environ.GetAllEnv()
)

func main() {
	models.SettingUp(env.HDS_TABLE_PREFIX)
	app := gin.New()
	app.Use(gin.Logger(), gin.CustomRecovery(middlewares.InternalServerErrorHandler))
	api.RegisterRoutes(app)
	err := app.Run(fmt.Sprintf("%v:%v", HOST, PORT))
	if err != nil {
		fmt.Println("Error on running server :", err)
		return
	}
	fmt.Println("Server started on PORT :", PORT)

}

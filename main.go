package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lwinmgmg/http_data_store/api"
)

var (
	HOST string = "localhost"
	PORT int    = 8000
)

func main() {
	app := gin.New()
	api.RegisterRoutes(app)
	err := app.Run()
	if err != nil {
		fmt.Println("Error on running server :", err)
		return
	}
	fmt.Println("Server started on PORT :", PORT)

}

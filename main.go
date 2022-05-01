package main

import (
	"fmt"

	"github.com/lwinmgmg/http_data_store/modules/models"
)

func main() {
	fmt.Println("Started")
	user, _ := models.GetUserById(3)
	if user.ID != 0 {
		models.DeleteById(user.ID)
	}
	fmt.Printf("user: %v\n", user)
}

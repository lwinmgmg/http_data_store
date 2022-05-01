package models

import (
	"github.com/lwinmgmg/http_data_store/dbm"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	if db == nil {
		if err := dbm.ConnectSqlite(); err != nil {
			panic(err)
		}
		db = dbm.GetSqliteDatabase()
	}
	db.AutoMigrate(&User{})
}

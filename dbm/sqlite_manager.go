package dbm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectSqlite() error {
	conn, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = conn
	return nil
}

func GetSqliteDatabase() *gorm.DB {
	return db
}

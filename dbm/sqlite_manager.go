package dbm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectSqlite() error {
	dsn := "lwinmgmg:password@tcp(127.0.0.1:3306)/ds?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = conn
	return nil
}

func GetSqliteDatabase() *gorm.DB {
	return db
}

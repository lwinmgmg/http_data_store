package dbm

import (
	"fmt"

	"github.com/lwinmgmg/http_data_store/environ"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var env *environ.Environ

func ConnectSqlite() error {
	env = environ.GetAllEnv()
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		env.HDS_DB_USER, env.HDS_DB_PASSWORD, env.HDS_DB_HOST, env.HDS_DB_PORT, env.HDS_DB_NAME)
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

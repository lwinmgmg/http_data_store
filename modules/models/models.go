package models

import (
	"fmt"

	"github.com/lwinmgmg/http_data_store/dbm"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	prefix *string
)

func init() {
	fmt.Println("ABCD")
	if db == nil {
		if err := dbm.ConnectSqlite(); err != nil {
			panic(err)
		}
		db = dbm.GetSqliteDatabase()
	}
}

func SettingUp(pref string) {
	prefix = &pref
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Folder{})
}

func TearDown(pref string) {
	prefix = &pref
	db.Migrator().DropTable(&Folder{})
	db.Migrator().DropTable(&User{})
}

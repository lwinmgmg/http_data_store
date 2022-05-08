package models

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lwinmgmg/http_data_store/dbm"
	"github.com/lwinmgmg/http_data_store/environ"
	"github.com/lwinmgmg/http_data_store/helper"
	"gorm.io/gorm"
)

var (
	db           *gorm.DB
	redis_client *redis.Client
	common_ctx    context.Context = context.Background()
	prefix       *string
	env          *environ.Environ = environ.GetAllEnv()
)

func init() {
	if db == nil {
		if err := dbm.ConnectMySql(); err != nil {
			panic(err)
		}
		db = dbm.GetMySqlDatabase()
	}
	if redis_client == nil {
		dbm.ConnectRedis()
		redis_client = dbm.GetRedisConnection()
	}
	
}

func SettingUp(pref string) {
	prefix = &pref
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Folder{})
	db.AutoMigrate(&File{})
	users, err := GetAllUser()
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		user := User{
			UserName: "admin",
			Password: helper.HexString(env.HDS_ADMIN_PASSWORD),
		}
		if _, err := user.Create(); err != nil {
			panic(err)
		}
	}
}

func TearDown(pref string) {
	prefix = &pref
	db.Migrator().DropTable(&Folder{})
	db.Migrator().DropTable(&User{})
	db.Migrator().DropTable(&File{})
}

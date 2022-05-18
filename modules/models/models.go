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
	common_ctx   context.Context = context.Background()
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
		if err := user.Create(&user); err != nil {
			panic(err)
		}
	}
}

func CreateTable(table Table) {
	db.AutoMigrate(table)
}

func DropTable(table Table) {
	db.Migrator().DropTable(table)
}

func TearDown(pref string) {
	prefix = &pref
	db.Migrator().DropTable(&Folder{})
	db.Migrator().DropTable(&User{})
	db.Migrator().DropTable(&File{})
}

type ModelManager interface {
	// Getting Data
	GetAll(any) error
	GetByID(uint, any) error
	GetByIDs([]uint, any) error

	// Creating Data
	Create(Table, any) error

	// Updating Data
	UpdateByID(uint, Table, ...any) error
	UpdateByIDs([]uint, any, ...any) error
}

type Table interface {
	TableName() string
	GetID() uint
}

type HdsModel struct {
	Table Table
}

func (model *HdsModel) GetAll(dest any) error {
	return db.Model(model.Table).Find(dest).Error
}

func (model *HdsModel) GetByID(id uint, dest any) error {
	return db.Model(model.Table).First(dest, id).Error
}

func (model *HdsModel) GetByIDs(ids []uint, dest any) error {
	return db.Model(model.Table).Find(dest, ids).Error
}

func (model *HdsModel) Create(data Table, dest ...any) error {
	switch length := len(dest); length {
	case 0:
		return db.Model(model.Table).Create(data).Error
	default:
		if err := db.Model(model.Table).Create(data).Error; err != nil{
			return err
		}
		return model.GetByID(data.GetID(), dest[0])
	}
}

func (model *HdsModel) UpdateByID(id uint, data any, dest ...any) error {
	switch length := len(dest); length {
	case 0:
		return db.Model(model.Table).Create(data).Error
	default:
		if err := db.Model(model.Table).Create(data).Error; err != nil{
			return err
		}
		return model.GetByID(id, dest[0])
	}
}

func (model *HdsModel) UpdateByIDs(ids []uint, data any, dest ...any) error{
	switch length := len(dest); length {
	case 0:
		return db.Model(model.Table).Create(data).Error
	default:
		if err := db.Model(model.Table).Create(data).Error; err != nil{
			return err
		}
		return model.GetByIDs(ids, dest[0])
	}
}


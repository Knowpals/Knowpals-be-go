package ioc

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/Knowpals/Knowpals-be-go/repository/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(conf *config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(conf.Mysql.Dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	err = model.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

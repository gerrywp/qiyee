package models

import (
	"gerry.wang/qiyee/configs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open(configs.GetAbsPath("pai.db3")), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}

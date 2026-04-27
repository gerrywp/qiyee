package models

import (
	"gerry.wang/qiyee/configs"
	"github.com/glebarez/sqlite"
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

	// 自动迁移所有模型
	DB.AutoMigrate(
		&Banner{},
		&User{},
		&About{},
		&Site{},
		&Prod{},
		&News{},
	)
}

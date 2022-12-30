package common

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"poc/graph/customTypes"
)

func InitDb() (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&customTypes.Todo{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

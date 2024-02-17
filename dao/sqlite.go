package dao

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSqlite() (err error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

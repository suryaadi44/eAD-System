package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(dbHost string, dbPort string, dbUsername string, dbPassword string, dbName string, retries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	conString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err = gorm.Open(mysql.Open(conString), &gorm.Config{})
	for err != nil {
		if retries == 0 {
			return nil, err
		}
		retries--

		time.Sleep(5 * time.Second)
		db, err = gorm.Open(mysql.Open(conString), &gorm.Config{})
	}

	return db, err
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate()
}

package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	databaseFile = "sandissa.db"
	tempRange    = 60
)

var (
	DB *gorm.DB
)

// Temperature is the DB model used to store temperatures
type Temperature struct {
	gorm.Model

	Value float64
}

// init opens the database.
func initDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(databaseFile), &gorm.Config{})
	if err != nil {
		lerr("Failed to init storage", err, params{
			"filename": databaseFile,
		})
		return err
	}
	DB.AutoMigrate(&Temperature{})
	if err != nil {
		return err
	}
	lf("Initialized storage", params{"file": databaseFile})
	return nil
}

// addTempDB logs the temperature value
func addTempDB(value float64) error {
	return DB.Create(&Temperature{Value: value}).Error
}

// getTempDB returns the last temperature value
func getTempDB() (*Temperature, error) {
	val := &Temperature{}
	return val, DB.Model(&Temperature{}).Order("ID desc").First(val).Error
}

// getTempsDB returns the last temperature value
func getTempsDB() ([]Temperature, error) {
	val := make([]Temperature, 0, tempRange)
	err := DB.Model(&Temperature{}).Order("ID desc").Limit(tempRange).Find(&val).Error
	return val, err
}

func closeDB() error {
	db, err := DB.DB()
	if err != nil {
		lerr("Failed to get database while closing", err, params{})
		return err
	}
	err = db.Close()
	if err != nil {
		lerr("Failed to close the database", err, params{})
		return err
	}
	lf("Closed database", params{"file": databaseFile})
	return nil
}

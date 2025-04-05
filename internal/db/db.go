package db

import (
	calculationservice "kalc/internal/calculationService"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// подкючаемся к БД
func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=pass dbname=postgres port=5432 sslmode=disable" // data source name
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // пустой конфиг

	if err != nil {
		log.Fatal("Не смог подключиться к БД: ", err)
	}

	// автомиграция
	if err := db.AutoMigrate(&calculationservice.Calculation{}); err != nil {
		log.Fatal("Поломался на миграции: ", err)
	}

	return db, nil
}

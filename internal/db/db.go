package db

import (
	calculationservice "kalc/internal/calculationService"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// подкючаемся к БД
func InitDB() (*gorm.DB, error) {
	// dsn := "host=localhost user=postgres password=pass dbname=kalc port=5432 sslmode=disable"

	// ожидаем что переменные окружения лежат в контейнере
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// data source name
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

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

package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"steamtrade.shop/parser/pkg"
)

func Conn() *gorm.DB {
	err := godotenv.Load("parser/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Some database environment variables are not set")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		defer panic("failed to connect database")
	}
	db.AutoMigrate(&model.Product{})
	return db
	// Миграция схемы
	//db.AutoMigrate(&model.Product{}, &model.ProductPrice{})
}
func Contains(name string, db *gorm.DB) (bool, error, model.Product) {
	var item model.Product
	result := db.Where("name = ?", name).First(&item)
	if result.Error != nil {
		return false, result.Error, model.Product{}
	} else {
		//item.Price = price
		//db.Create(&item)
		return true, nil, item
	}
}

//func UpdateItems(names string, prices float64, db *gorm.DB)

func InsertItem(item *model.Product, db *gorm.DB) error {
	result := db.Create(&item)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func InsertItems(items []*model.Product, db *gorm.DB) error {
	result := db.Create(&items)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

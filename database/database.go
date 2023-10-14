package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=Password@2001 dbname=gin-commerce port=5432 sslmode=disable"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database...")
	}

	DB = connection

	DB.AutoMigrate(&Product{}, &User{}, &Order{}, &OrderItem{}, &Cart{}, &CartItem{})
}

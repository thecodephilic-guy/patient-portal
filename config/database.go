package config

import (
	"hospital-management/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB //creating and exporting the instance of the database connection

func ConnectDatabase() {
	er := godotenv.Load()
	if er != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, // Improves performance for migrations
	})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("✅ Successfully Connected To The Database!")
}

func MigrateDB() {
	err := DB.AutoMigrate(&models.User{}, &models.Patient{}) // Creates tables if not exists
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migration completed successfully!")
}

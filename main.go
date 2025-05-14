package main

import (
	"hospital-management/config"
	"hospital-management/models"
	"hospital-management/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	config.ConnectDatabase()

	// Run migrations
	// config.MigrateDB()

	// Create default admin user if not exists
	var adminCount int64
	config.DB.Model(&models.User{}).Count(&adminCount)
	if adminCount == 0 {
		// Create seed data
		createSeedData()
	}

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func createSeedData() {
	log.Println("Creating seed data...")

	// Hash a default password
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password: ", err)
	}

	// Create a receptionist user
	receptionist := models.User{
		Email:    "receptionist@hospital.com",
		Password: string(hashedPassword),
		Name:     "Reception Staff",
		Role:     "receptionist",
	}
	config.DB.Create(&receptionist)

	// Create a doctor user
	doctor := models.User{
		Email:    "doctor@hospital.com",
		Password: string(hashedPassword),
		Name:     "Dr. John Smith",
		Role:     "doctor",
	}
	config.DB.Create(&doctor)

	log.Println("Seed data created successfully")
}

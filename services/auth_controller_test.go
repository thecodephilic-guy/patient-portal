package services

import (
	"bytes"
	"encoding/json"
	"hospital-management/config"
	"hospital-management/controllers"
	"hospital-management/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTestDB() *gorm.DB {
	// Use SQLite for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Run migrations
	db.AutoMigrate(&models.User{}, &models.Patient{})

	// Set global DB
	config.DB = db
	return db
}

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestAuthMain(m *testing.M) {
	// Load env for JWT secret
	godotenv.Load("../.env.test")

	// Setup test environment
	setupTestDB()

	// Run tests
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

func TestLoginUserAuthController(t *testing.T) {
	// Setup
	db := setupAuthTestDB()
	r := setupAuthRouter()

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	testUser := models.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
		Role:     "receptionist",
	}
	db.Create(&testUser)

	// Define login endpoint
	r.POST("/login", controllers.LoginUser)

	// Test case 1: Valid login
	t.Run("Valid Login", func(t *testing.T) {
		// Create login request
		loginJSON, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"password": "testpassword",
		})

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// Assert status code is 200 OK
		assert.Equal(t, http.StatusOK, resp.Code)

		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check if token exists
		token, exists := response["token"]
		assert.True(t, exists)
		assert.NotEmpty(t, token)

		// Check if user data exists
		user, exists := response["user"]
		assert.True(t, exists)
		assert.NotNil(t, user)
	})

	// Test case 2: Invalid password
	t.Run("Invalid Password", func(t *testing.T) {
		loginJSON, _ := json.Marshal(map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		})

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// Assert status code is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	// Test case 3: Non-existent user
	t.Run("Non-existent User", func(t *testing.T) {
		loginJSON, _ := json.Marshal(map[string]string{
			"email":    "nonexistent@example.com",
			"password": "testpassword",
		})

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		// Assert status code is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

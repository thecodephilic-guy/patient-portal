package controllers

import (
	"hospital-management/config"
	"hospital-management/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreatePatient handles patient registration by receptionists
func CreatePatient(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "receptionist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var input struct {
		FirstName      string `json:"first_name" binding:"required"`
		LastName       string `json:"last_name" binding:"required"`
		Email          string `json:"email" binding:"required,email"`
		Phone          string `json:"phone"`
		DateOfBirth    string `json:"date_of_birth"`
		Gender         string `json:"gender"`
		Address        string `json:"address"`
		MedicalHistory string `json:"medical_history"`
		Allergies      string `json:"allergies"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if patient already exists
	var existingPatient models.Patient
	if err := config.DB.Where("email = ?", input.Email).First(&existingPatient).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient with this email already exists"})
		return
	}

	// Parse date of birth
	var dob time.Time
	var err error
	if input.DateOfBirth != "" {
		dob, err = time.Parse("2006-01-02", input.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create patient
	patient := models.Patient{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Phone:          input.Phone,
		DateOfBirth:    dob,
		Gender:         input.Gender,
		Address:        input.Address,
		MedicalHistory: input.MedicalHistory,
		Allergies:      input.Allergies,
		CreatedBy:      userID.(uint),
		UpdatedBy:      userID.(uint),
	}

	if err := config.DB.Create(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register patient"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient registered successfully",
		"patient": patient,
	})
}

// GetAllPatients returns all patients
func GetAllPatients(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "doctor" && role != "receptionist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var patients []models.Patient
	if err := config.DB.Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}

// GetPatient returns a specific patient by ID
func GetPatient(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "doctor" && role != "receptionist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	id := c.Param("id")

	var patient models.Patient
	if err := config.DB.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patient": patient})
}

// UpdatePatient updates patient information
func UpdatePatient(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "doctor" && role != "receptionist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	id := c.Param("id")

	var patient models.Patient
	if err := config.DB.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var input struct {
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		DateOfBirth    string `json:"date_of_birth"`
		Gender         string `json:"gender"`
		Address        string `json:"address"`
		MedicalHistory string `json:"medical_history"`
		Allergies      string `json:"allergies"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.FirstName != "" {
		patient.FirstName = input.FirstName
	}
	if input.LastName != "" {
		patient.LastName = input.LastName
	}
	if input.Email != "" {
		patient.Email = input.Email
	}
	if input.Phone != "" {
		patient.Phone = input.Phone
	}
	if input.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		patient.DateOfBirth = dob
	}
	if input.Gender != "" {
		patient.Gender = input.Gender
	}
	if input.Address != "" {
		patient.Address = input.Address
	}
	if input.MedicalHistory != "" {
		patient.MedicalHistory = input.MedicalHistory
	}
	if input.Allergies != "" {
		patient.Allergies = input.Allergies
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	patient.UpdatedBy = userID.(uint)

	// Save updated patient
	if err := config.DB.Save(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient updated successfully",
		"patient": patient,
	})
}

// DeletePatient removes a patient from the database
func DeletePatient(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "receptionist" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	// Check if patient exists
	var patient models.Patient
	if err := config.DB.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Delete patient
	if err := config.DB.Delete(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}

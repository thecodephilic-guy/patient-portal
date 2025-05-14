package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `gorm:"not null" json:"name"`
	Role     string `gorm:"not null" json:"role"` // "doctor" or "receptionist"
}

type Patient struct {
	gorm.Model
	FirstName      string    `gorm:"not null" json:"first_name"`
	LastName       string    `gorm:"not null" json:"last_name"`
	Email          string    `gorm:"unique;not null" json:"email"`
	Phone          string    `json:"phone"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Gender         string    `json:"gender"`
	Address        string    `json:"address"`
	MedicalHistory string    `json:"medical_history"`
	Allergies      string    `json:"allergies"`
	CreatedBy      uint      `json:"created_by"`
	UpdatedBy      uint      `json:"updated_by"`
}

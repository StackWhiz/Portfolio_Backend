package models

import (
	"time"

	"gorm.io/gorm"
)

// Profile represents the main profile information
type Profile struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Location  string    `json:"location"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Phone     string    `json:"phone"`
	Telegram  string    `json:"telegram"`
	GitHub    string    `json:"github"`
	LinkedIn  string    `json:"linkedin"`
	Summary   string    `json:"summary" gorm:"type:text"`
	Avatar    string    `json:"avatar"`
	ResumeURL string    `json:"resume_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Experience represents work experience entries
type Experience struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Company      string     `json:"company" gorm:"not null"`
	Position     string     `json:"position" gorm:"not null"`
	Location     string     `json:"location"`
	StartDate    time.Time  `json:"start_date" gorm:"not null"`
	EndDate      *time.Time `json:"end_date"`
	Current      bool       `json:"current" gorm:"default:false"`
	Description  string     `json:"description" gorm:"type:text"`
	Achievements []string   `json:"achievements" gorm:"type:json"`
	Technologies []string   `json:"technologies" gorm:"type:json"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Skill represents technical skills
type Skill struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Category    string    `json:"category" gorm:"not null"` // Languages, Frameworks, Tools, etc.
	Level       int       `json:"level" gorm:"default:5"`   // 1-10 scale
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Project represents portfolio projects
type Project struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Description     string    `json:"description" gorm:"type:text"`
	LongDescription string    `json:"long_description" gorm:"type:text"`
	Technologies    []string  `json:"technologies" gorm:"type:json"`
	GitHubURL       string    `json:"github_url"`
	LiveURL         string    `json:"live_url"`
	ImageURL        string    `json:"image_url"`
	Featured        bool      `json:"featured" gorm:"default:false"`
	Category        string    `json:"category"`                          // Blockchain, Backend, Full-stack, etc.
	Status          string    `json:"status" gorm:"default:'completed'"` // completed, in-progress, planned
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Contact represents contact form submissions
type Contact struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Status    string    `json:"status" gorm:"default:'new'"` // new, read, replied
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents admin users
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // Hidden from JSON
	Role      string    `json:"role" gorm:"default:'admin'"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Hash password before creating user
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// BeforeUpdate hook for User
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Hash password before updating if it's being changed
	if tx.Statement.Changed("Password") {
		hashedPassword, err := HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

package domain

import (
	"time"
)

// Users slice of User
type Users []User

// User User is model of users
type User struct {
	ID        uint       `gorm:"primary_key"json:"id"`
	Name      string     `gorm:"type:varchar(255);not null;"json:"name"`
	UID       string     `gorm:"type:varchar(255);not null;unique"json:"uid"`
	Pasword   string     `gorm:"size:255;not null"json:"password"`
	Role      int        `gorm:"not null"json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

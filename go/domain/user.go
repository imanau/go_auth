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
	UID       string     `gorm:"type:varchar(255);not null;unique_index"json:"uid"`
	Password  string     `gorm:"size:255;not null"json:"password"`
	Role      int        `gorm:"not null"json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// UserForGeneral api返却用のユーザー情報を格納する
type UserForGeneral struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// PasswordInfo パスワード更新用の情報を格納する
type PasswordInfo struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	NewPassword          string `json:"new_password"`
}

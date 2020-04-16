package model

import (
	"go_auth/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // porstgres driver
)

var db *gorm.DB

// ConnectDB return *gorm.DB
func ConnectDB() (*gorm.DB, error) {
	var err error
	db, err = gorm.Open("postgres", "host=db port=5432 user=postgres dbname=auth_service password=postgres sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.User{})
	return db, err
}

// AllUser return all of users
func AllUser(db *gorm.DB) *gorm.DB {
	users := new(domain.Users)
	result := db.Find(&users)
	return result
}

// FindUser return User
func FindUser(db *gorm.DB, user *domain.User) domain.User {
	db.Where("uid = ?", user.UID).First(&user)
	return *user
}

// CreateUser model Create
func CreateUser(db *gorm.DB, user *domain.User) {
	db.Create(user)
}

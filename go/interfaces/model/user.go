package model

import (
	"go_auth/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=db port=5432 user=postgres dbname=auth_service password=postgres sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.User{})
}

// AllUser return all of users
func AllUser() *gorm.DB {
	users := new(domain.Users)
	result := db.Find(&users)
	return result
}

package model

import (
	"fmt"
	"go_auth/domain"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // porstgres driver
	"github.com/joho/godotenv"
)

var db *gorm.DB

// ConnectDB return *gorm.DB
func ConnectDB() (*gorm.DB, error) {
	// envファイル呼び込み
	err := godotenv.Load(fmt.Sprintf("envfiles/%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		panic("Error loading .env file")
	}
	dbhost := os.Getenv("dbhost")
	dbport := os.Getenv("dbport")
	dbuser := os.Getenv("dbuser")
	dbname := os.Getenv("dbname")
	dbpassword := os.Getenv("dbpassword")
	sslmode := os.Getenv("sslmode")
	dbconfig := "host=" + dbhost + " port=" + dbport + " user=" + dbuser + " dbname=" + dbname + " password=" + dbpassword + " sslmode=" + sslmode
	db, err = gorm.Open("postgres", dbconfig)
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

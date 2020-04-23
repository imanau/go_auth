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
	err := godotenv.Load(fmt.Sprintf("/usr/src/go_auth/envfiles/%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		panic(err)
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
	if user.ID == 0 {
		db.Where("uid = ?", user.UID).First(&user)
	} else {
		db.Where("id = ?", user.ID).First(&user)
	}
	return *user
}

// CreateUser model Create
func CreateUser(db *gorm.DB, user *domain.User) {
	db.Create(user)
}

// UpdateUser model Update withot password
func UpdateUser(db *gorm.DB, user *domain.User) {
	db.Model(&user).Updates(map[string]interface{}{"name": user.Name, "uid": user.UID, "role": user.Role})
}

// ChangePassword model Update with password
func ChangePassword(db *gorm.DB, user *domain.User) {
	db.Model(&user).Update("password", user.Password)
}

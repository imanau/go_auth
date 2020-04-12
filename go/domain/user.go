// User User is model of users
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255);not null;"`
	UID       string `gorm:"type:varchar(255);not null;unique"`
	Pasword   string `gorm:"size:255;not null"`
	Role      int    `gorm:"not null"`
}
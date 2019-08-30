package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Password    string
	LineId		string `gorm:"unique"`
	Investments []Investment `gorm:"many2many:user_investments;"`
}

type Product struct {
	gorm.Model
	Name     string
	ImageURL string
	Price    int64
}

type Investment struct {
	gorm.Model
	Name    string
	Product *Product
	Goal    int64
	Year    int
	Current int64
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Product{}, &Investment{})
}

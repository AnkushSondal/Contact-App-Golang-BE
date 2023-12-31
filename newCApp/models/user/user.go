package user

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"Username" gorm:"unique;type:varchar(100)"`
	Password string `json:"Password" gorm:"type:varchar(100)"`
	FullName string `json:"FullName" gorm:"type:varchar(100)"`
	IsAdmin  string   `json:"isAdmin" gorm:"type:boolean;default:0" `
}

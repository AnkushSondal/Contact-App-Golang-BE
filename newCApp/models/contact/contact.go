package contact

import (
	"contactApp/models/user"

	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	User        user.User `gorm:"foreignkey:UserRefer"` // use UserRefer as foreign key
	UserRefer   uint      `json:"UserRefer" `
	ContactName string    `json:"ContactName" gorm:"type:varchar(100)"`
}

package contactinfo

import (
	"contactApp/models/contact"

	"github.com/jinzhu/gorm"
)

type ContactInfo struct {
	gorm.Model
	Contact      contact.Contact `gorm:"foreignkey:ContactRefer"`
	ContactRefer uint            `json:"ContactRefer" `
	CIType       string          `json:"CIType" gorm:"type:varchar(100)"`
	CIValue      string          `json:"CIValue" gorm:"type:varchar(100)"`
}

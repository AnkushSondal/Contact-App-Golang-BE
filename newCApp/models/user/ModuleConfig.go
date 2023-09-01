package user

import (
	"contactApp/components/log"
	"sync"

	"github.com/jinzhu/gorm"
)

type ModuleConfig struct {
	DB *gorm.DB
}

// NewTestModuleConfig Create New Test Module Config
func NewUserModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {
	// Table List
	var models []interface{} = []interface{}{
		&User{},
	}
	// Table Migrant
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.GetLogger().Print("Auto Migration ==>", err)
		}
	}
	// if err := config.DB.Model(&User{}).
	// 	AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.GetLogger().Print("Auto Migration ==> %s", err)
	// }
	log.GetLogger().Print("Test Module Configured.")
}

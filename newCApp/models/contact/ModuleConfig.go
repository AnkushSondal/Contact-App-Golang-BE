package contact

import (
	"contactApp/components/log"
	"sync"

	"github.com/jinzhu/gorm"
)

type ModuleConfig struct {
	DB *gorm.DB
}

// NewTestModuleConfig Create New Test Module Config
func NewContactModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {
	// Table List
	var models []interface{} = []interface{}{
		&Contact{},
	}
	// Table Migrant
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.GetLogger().Print("Auto Migration ==>", err)
		}
	}
	if err := config.DB.Model(&Contact{}).
		AddForeignKey("user_refer", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.GetLogger().Print("Auto Migration ==>", err)
		log.GetLogger().Print("Auto Migration ................==>", err)
	}
	log.GetLogger().Print("Test Module Configured.")
}

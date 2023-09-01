package module

import (
	"contactApp/app"
	"contactApp/models/contact"
	"contactApp/models/contactinfo"
	"contactApp/models/user"
)

// Configure will configure all modules
func Configure(appObj *app.App) {
	userModule := user.NewUserModuleConfig(appObj.DB)
	contactModule := contact.NewContactModuleConfig(appObj.DB)
	contactInfoModule := contactinfo.NewContactInfoModuleConfig(appObj.DB)

	appObj.MigrateTables([]app.ModuleConfig{userModule,contactModule, contactInfoModule })
	
}

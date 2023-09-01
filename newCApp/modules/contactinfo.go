package module

import (
	"contactApp/app"
	"contactApp/components/contact_info/controller"
	"contactApp/components/contact_info/service"
	"contactApp/repository"
)

func registerContactInfoRoutes(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	contactInfoService := service.NewContactInfoService(appObj.DB, repository)

	contactInfoController := controller.NewContactInfoController(contactInfoService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactInfoController,
	})
}
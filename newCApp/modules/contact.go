package module

import (
	"contactApp/app"
	"contactApp/components/contact/controller"
	"contactApp/components/contact/service"
	"contactApp/repository"
)

func registerContactRoutes(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	contactService := service.NewContactService(appObj.DB, repository)

	contactController := controller.NewContactController(contactService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactController,
	})
}
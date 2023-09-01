package module

import (
	"contactApp/app"
	"contactApp/components/user/controller"
	"contactApp/components/user/service"
	"contactApp/repository"
)

func registerUserRoutes(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	userService := service.NewUserService(appObj.DB, repository)

	userController := controller.NewUserController(userService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		userController,
	})
}

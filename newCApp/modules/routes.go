package module

import (
	"contactApp/app"
	"contactApp/repository"
)

func RegisterModuleRoutes(app *app.App, repository repository.Repository) {
	log := app.Log
	log.Print("============RegisterModuleRoutes.go==============")
	app.WG.Add(3)
	go registerUserRoutes(app, repository)
	go registerContactRoutes(app, repository)
	go registerContactInfoRoutes(app, repository)
	
	app.WG.Wait()
}

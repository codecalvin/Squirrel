package main

import (
	"github.com/astaxie/beego"
	
	_ "squirrelchuckle/routers"
	"squirrelchuckle/database"
	"squirrelchuckle/services"
	"squirrelchuckle/settings"
)

var serviceManager *services.ServiceManager

func main() {
	defer dispose()
	setup()
	beego.Run()
}

func setup() {
	serviceManager = services.GetManager()
	serviceManager.Initialize()
}

func dispose() {
	settings.Serialize()
	serviceManager.UnInitialize()
	database.Close()
}


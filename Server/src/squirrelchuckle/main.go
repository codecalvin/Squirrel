package main

import (
	"github.com/astaxie/beego"
	
	_ "squirrelchuckle/routers"
	"squirrelchuckle/database"
	"squirrelchuckle/services"
)

var serviceManager *services.ServiceManager

func main() {
	defer dispose()
	beego.Run()
}

func setup() {
	serviceManager = services.New()
	serviceManager.Initialize()
}

func dispose() {
	serviceManager.UnInitialize()
	database.Close()
}


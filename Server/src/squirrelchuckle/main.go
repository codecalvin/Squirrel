package main

import (
	"github.com/astaxie/beego"

	_ "squirrelchuckle/routers"
	"squirrelchuckle/services"
)

var serviceManager *services.ServiceManager
var appSetting *services.AppSetting
var database *services.Database
func main() {
	defer dispose()
	setup()
	beego.Run()
}

func setup() {
	var instances []services.ServiceInterface = make([]services.ServiceInterface, 10)

	appSetting = services.AppSettingInstance()
	database = services.DatabaseInstance()
	instances = append(instances, database)
	err := appSetting.Initialize()
	if err != nil {
		panic(err)
	}
	
	serviceManager = services.GetManager()
	serviceManager.Initialize()
}

func dispose() {
	// make sure tear down order
	serviceManager.UnInitialize()
	database.UnInitialize()
	appSetting.Serialize()
}


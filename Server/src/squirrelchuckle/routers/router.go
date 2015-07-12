package routers

import (
	"github.com/astaxie/beego"

	"squirrelchuckle/controllers"
	"squirrelchuckle/core"
	"squirrelchuckle/services"
)

func router(rootPath, info string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
	controllers.EndPoints[info] = rootPath
	return beego.Router(rootPath, c, mappingMethods...)
}

func init() {
	controllers.EndPoints = make(map[string]string)
	registerServices()

	router("/", "home_page", &controllers.MainController{})
	router("/api", "api_page", &controllers.ApiController{})
	router("/api/admin", "admin_page", &controllers.AdminController{})
	
	router("/ws", "websocket_example:Join", &controllers.WebSocketController{})
	router("/ws/join", "websocket_example:Join", &controllers.WebSocketController{}, "get:Join")

	router("/users", "users", &controllers.UsersController{})
	router("/users/:id", "user_url", &controllers.UsersController{})

	router("/user", "current_user", &controllers.UserController{})

	router("/device", "device_token", &controllers.DeviceTokenController{})
}

func registerServices() {
	core.SquirrelApp.RegisterService(&services.ApplePushService{})
	core.SquirrelApp.RegisterService(&services.DeviceTokenService{})
	core.SquirrelApp.RegisterService(&services.UserService{})
	
}
package routers

import (
	"github.com/astaxie/beego"

	"squirrelchuckle/controllers"
)

func router(rootPath, info string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
	controllers.EndPoints[info] = rootPath
	return beego.Router(rootPath, c, mappingMethods...)
}

func init() {
	controllers.EndPoints = make(map[string]string)

	router("/", "home_page", &controllers.MainController{})
	router("/api", "api_page", &controllers.ApiController{})
	
	router("/apns/test", "apns_tester", &controllers.APNPushTestController{})
	
	router("/ws", "websocket_example:Join", &controllers.WebSocketController{})
	router("/ws/join", "websocket_example:Join", &controllers.WebSocketController{}, "get:Join")

	router("/users", "users", &controllers.UsersController{})
	router("/users/:id", "user_url", &controllers.UsersController{})

	router("/user", "current_user", &controllers.UserController{})

	router("/classes", "classes", &controllers.ClassController{})
	router("/classes/:id", "class_url", &controllers.ClassController{})
}

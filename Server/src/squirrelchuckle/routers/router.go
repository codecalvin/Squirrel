package routers

import (
	//"fmt"
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

	router("/device", "device_token", &controllers.DeviceTokenController{})
	router("/classes", "classes", &controllers.ClassController{})
	router("/classes/:id", "class_url", &controllers.ClassController{})

	router("/API1/Post", "add/update a class by admin", &controllers.PostController{})
	router("/API1/Classes", "get classes", &controllers.ClassController{})
	router("/API1/OneClass/:classKey", "get one class", &controllers.OneClassController{})
	router("/API1/Classes/Register", "a user register a class", &controllers.RegisterController{})
	router("/API1/Classes/UnRegister", "a user unregister a class", &controllers.UnRegisterController{})

	router("/API1/Classes/User/:userKey", "get classes of a user", &controllers.UsersController{})
	router("/API1/OneClassUsers/:classKey", "get users of a class", &controllers.OneClassUsersController{})
	router("/API1/OneClassUserCount/:classKey", "get user count information of a class", &controllers.OneClassUserCountController{})
	router("/API1/Classes/QueryRegisterStatus", "query register status", &controllers.RegisterStatusController{})
}

func registerServices() {
	core.SquirrelApp.RegisterService(&services.ApplePushService{})
	core.SquirrelApp.RegisterService(&services.DeviceTokenService{})
	core.SquirrelApp.RegisterService(&services.UserService{})
}

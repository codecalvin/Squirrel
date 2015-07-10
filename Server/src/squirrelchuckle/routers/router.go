package routers

import (
	//"fmt"
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

package routers

import (
	"github.com/astaxie/beego"

	"squirrelchuckle/controllers"
)

func router(rootPath, info string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
	controllers.EndPoint[info] = rootPath
	return beego.Router(rootPath, c, mappingMethods...)
}

func init() {
	controllers.EndPoint = make(map[string]string)

	router("/", "home_page", &controllers.MainController{})
	router("/api", "api_page", &controllers.ApiController{})
	router("/users", "users", &controllers.UsersController{})
	router("/users/:id", "user_url", &controllers.UsersController{})

	router("/user", "current_user", &controllers.UserController{})

	router("/classes", "classes", &controllers.ClassController{})
	router("/classes/:id", "class_url", &controllers.ClassController{})
}

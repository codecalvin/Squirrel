package routers

import (
	"squirrelchuckle/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

	beego.Router("/users", &controllers.UsersController{})
	//beego.Router("/users/create", &controllers.UsersController{})
	//beego.Router("/users/delete", &controllers.UsersController{})

	beego.Router("/user/:id", &controllers.UserController{})
}

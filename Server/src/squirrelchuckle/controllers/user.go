package controllers


import (

	"github.com/astaxie/beego"
	
	"squirrelchuckle/services"
)

type UsersController struct {
	beego.Controller
}

type SignUpController struct {
	beego.Controller
}

var userService *services.UserService

func (this *UsersController) Get() {
}

func (this *UsersController) Post() {
//	input := this.Input()
//	name := input.Get("name")
//	id := input.Get("id")
//
//	p := ProfileInfo{UserName:name, UserId:id}
//	if cInfo, err := userService.Upsert(bson.M{"userid": id}, p); err != nil {
//		this.Ctx.Output.Body([]byte(err.Error()))
//	} else {
//		this.Ctx.Output.Body([]byte(fmt.Sprintf("updated %v, raw name %v, raw id %v", cInfo, name, id)))
//	}
}


// -----------------------------------------------
// User Management
type UserController struct {
	beego.Controller
}

func (this *SignUpController) Post() {
	input 		:= this.Input()
	name 		:= input.Get("name")
	password 	:= input.Get("password")
	deviceToken := input.Get("device_token")

	if _, ok := userService.Users[name]; ok {
		this.CustomAbort(400, "User already registed")
		return
	}

	user := &services.User {
		AdsName: name,
		Password:password,
	}

	if err := userService.AddWithDevice(user, deviceToken); err != nil {
		this.Ctx.Output.Body([]byte(err.Error()))
	}
}
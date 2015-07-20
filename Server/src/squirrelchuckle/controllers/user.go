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

type SignInController struct {
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
	adsName := input.Get("ads_name")
	adsPass := input.Get("ads_pass")
	deviceToken := input.Get("device_token")

	if _, ok := userService.Users[adsName]; ok {
		this.CustomAbort(400, "User already registed")
		return
	}

	user := &services.User {
		AdsName:    adsName,
		AdsPass:    adsPass,
	}

	if newUser, err := userService.AddWithDevice(user, deviceToken); newUser == nil && err != nil {
		this.CustomAbort(400, err.Error())
	} else {
		this.Data["password"] = newUser.Password
	}
	this.ServeJson()
}

func (this *SignInController) Post() {
	input 		:= this.Input()
	adsName	 	:= input.Get("ads_name")
	adsPass 	:= input.Get("ads_pass")
	password 	:= input.Get("password")
	deviceToken := input.Get("device_token")

	if user, ok := userService.Users[adsName]; !ok {
		this.CustomAbort(400, "User doesn't exist")
		return
	} else {
		if password == user.Password {
			// authorized login
		} else if len(adsPass) > 0 {
			if userService.Auth(&user.AdsName, &user.AdsPass) {
				// authorized login
			} else {
				this.CustomAbort(400, "User name or password error")
				return
			}
		} else {
			this.CustomAbort(400, "User name or password error")
			return
		}

		this.Data["password"] = user.Password

		// transfer & touch device token
		deviceTokenService.Add(user.AdsName, deviceToken)
		this.ServeJson()
	}
}

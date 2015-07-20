package controllers

import (
	"github.com/astaxie/beego"
	"squirrelchuckle/services"
	"fmt"
	"squirrelchuckle/core"
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

type RegisterItem struct {
	ElementType_UniqueKey     string
	ElementType_ClassName     string
	ElementType_UserUniqueKey string
	ElementType_UserName      string
}

func (this *UsersController) Get() {
	userKey := this.Ctx.Input.Param(":userKey")
	fmt.Println(userKey, "========")
	if user, ok := userService.Users[userKey]; ok {
		this.Data["json"] = user.Classes
	}
	this.ServeJson()
}

// -----------------------------------------------
// User Management
type UserController struct {
	beego.Controller
}

func (this *SignUpController) Post() {
	if userService == nil {
		userService, _ = core.SquirrelApp.GetServiceByName("UserService").(*services.UserService)
	}
	if userService == nil {
		this.Data["json"] = fmt.Sprintf("service is down a %v, u %v, d %v", apnService==nil, userService == nil, deviceTokenService==nil)
		this.ServeJson()
		return
	}
	input 		:= this.Input()
	adsName := input.Get("ads_name")
	adsPass := input.Get("ads_pass")
	deviceToken := input.Get("device_token")

	fmt.Println("111111")
	if _, ok := userService.Users[adsName]; ok {
		fmt.Println("User already registed")
		this.CustomAbort(400, "User already registed")
		return
	}

	user := &services.User {
		AdsName:    adsName,
		AdsPass:    adsPass,
	}

	fmt.Println("222222")
	if newUser, err := userService.AddWithDevice(user, deviceToken); newUser == nil {
		fmt.Println(newUser, err)
		this.CustomAbort(400, err.Error())
	} else {
		this.Data["json"] = []string { newUser.Password }
		fmt.Println("OK", newUser, err)
	}
	fmt.Println("333333")
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

		this.Data["password"] = user

		// transfer & touch device token
		deviceTokenService.Add(user.AdsName, deviceToken)
		this.ServeJson()
	}
}

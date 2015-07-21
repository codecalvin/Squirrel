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

type ListUserController struct {
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

func (this *ListUserController) Get() {


	if ! initUserController() {
		this.ServeJson()
		return
	}

	i := 0
	users := make([]*services.User, len(userService.Users))
	for _, v := range userService.Users {
		users[i] = v
		i += 1
	}
	this.Data["json"] = users
	this.ServeJson()
}

func (this *ListUserController) Post() {
	if ! initUserController() {
		this.ServeJson()
		return
	}

	input := this.Input()
	if name := input.Get("name"); name != "" {
		if user, ok := userService.Users[name]; ok {
			message := input.Get("msg")
			if message == "" {
				message = "mock up"
			}
			userService.NotifyUser(user, message, input.Get("badge"))
		}
	}
	this.ServeJson()
}

func (this *UsersController) Get() {
	userKey := this.Ctx.Input.Param(":userKey")
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

const (
	SUCC_CODE	int = iota
	NEW_AUTH
	PASS_ERROR
)

func initUserController() bool {
	if userService == nil {
		userService, _ = core.SquirrelApp.GetServiceByName("UserService").(*services.UserService)
	}
	if deviceTokenService == nil {
		deviceTokenService, _ = core.SquirrelApp.GetServiceByName("DeviceTokenService").(*services.DeviceTokenService)
	}
	if userService == nil || deviceTokenService == nil ||
		!userService.Alive() || !deviceTokenService.Alive() {
		return false
	}

	return true
}

func (this *SignUpController) Post() {
	if !initUserController() {
		this.ServeJson()
		return
	}

	input	:= this.Input()
	adsName := input.Get("ads_name")
	adsPass := input.Get("ads_pass")
	deviceToken := input.Get("device_token")

	var ok bool
	var err error
	var user *services.User

	status := SUCC_CODE
	if user, ok = userService.Users[adsName]; ok {
		if user.AdsPass != adsPass || adsPass == "" {
			status = PASS_ERROR
		}
	} else {
		if user, err = userService.Add(adsName, adsPass); err != nil {
			status = PASS_ERROR
		}
	}

	switch status {
	case SUCC_CODE:
		if deviceToken != "" {
			userService.AddDevice(user, deviceToken)
		}
		this.Data["json"] = []string {user.Password}

		userService.NotifyUser(user, fmt.Sprintf("hello %v, login from a device", user.AdsName), "-1")

	case PASS_ERROR:
		this.Ctx.Output.SetStatus(Unauth)
	}

	this.ServeJson()
}

func (this *SignInController) Post() {
	input 		:= this.Input()
	adsName	 	:= input.Get("ads_name")
	password 	:= input.Get("password")
	deviceToken := input.Get("device_token")

	if user, ok := userService.Users[adsName]; !ok {
		this.Ctx.Output.SetStatus(Unauth)
	} else {
		if password == user.Password  {
			if deviceToken != "" {
				// add & update device token
				userService.AddDevice(user, deviceToken)

				// TODO test
				userService.NotifyUser(user, fmt.Sprintf("hello %v, login from a device", user.AdsName), "-1")
			}
		} else {
			this.Ctx.Output.SetStatus(Unauth)
		}
		this.ServeJson()
	}
}

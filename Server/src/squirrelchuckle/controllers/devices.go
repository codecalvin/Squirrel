package controllers
import (
	"github.com/astaxie/beego"

	"squirrelchuckle/core"
	"squirrelchuckle/services"
	"encoding/hex"
//	"encoding/binary"
)


type DeviceTokenController struct {
	beego.Controller
}

var apnService *services.ApplePushService
var deviceTokenService *services.DeviceTokenService

func initController() bool {
	if apnService == nil {
		apnService, _ = core.SquirrelApp.GetServiceByName("ApplePushService").(*services.ApplePushService)
	}
	if userService == nil {
		userService, _ = core.SquirrelApp.GetServiceByName("UserService").(*services.UserService)
	}
	if deviceTokenService == nil {
		deviceTokenService, _ = core.SquirrelApp.GetServiceByName("DeviceTokenService").(*services.DeviceTokenService)
	}
	if apnService == nil || userService == nil || deviceTokenService == nil ||
		!apnService.Alive() || !userService.Alive() || !deviceTokenService.Alive() {
		return false
	}

	return true
}

func (this *DeviceTokenController) Get() {
	if !initController() {
		this.Data["json"] = "service is down"
	} else {
		deviceTokenService.TestAPN()
		this.Data["json"] = "Test APN"
	}

	this.ServeJson()
}


func (this *DeviceTokenController) Post() {
	if !initController() {
		this.Data["json"] = "service is down"
		this.ServeJson()
		return
	}

	input := this.Input()
	email := input.Get("email")
	deviceToken := input.Get("device_token")
	if _, ok := userService.Users[email]; !ok {
		if _, err := hex.DecodeString(deviceToken); err == nil {
			//user.Devices = append(user.Devices, services.Device{ UserID : deviceTokenId, })
//			deviceTokenService.Add([]uint32{ binary.BigEndian.Uint32(deviceTokenId) })
		} else {
			this.Ctx.Output.Body([]byte(err.Error()))
		}
		this.Ctx.Output.Body([]byte("not found user"))
	}
}

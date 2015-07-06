package controllers
import (
	"github.com/astaxie/beego"

	"squirrelchuckle/core"
	"squirrelchuckle/services"
)


type DeviceTokenController struct {
	beego.Controller
}

var apnService *services.ApplePushService

func (this *DeviceTokenController) Get() {
	if apnService == nil || !apnService.Alive() {
		apnService, _ = core.SquirrelApp.GetServiceByName("APNService").(*services.ApplePushService)
	}

	result := apnService.TestAPN()

	this.Data["json"] = result
	this.ServeJson()
}

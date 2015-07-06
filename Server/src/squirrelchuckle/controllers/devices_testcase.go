package controllers
import (
	"github.com/astaxie/beego"

	"squirrelchuckle/core"
	"squirrelchuckle/services"
)


type APNPushTestController struct {
	beego.Controller
}

var tapnService *services.ApplePushService

func (this *APNPushTestController) Get() {
	if apnService == nil || !apnService.Alive() {
		tapnService, _ = core.SquirrelApp.GetServiceByName("APN_Service").(*services.ApplePushService)
	}

	result := tapnService.TestAPN()
	
	this.Data["json"] = result
	this.ServeJson()
}

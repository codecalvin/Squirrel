package controllers
import (
	"github.com/astaxie/beego"

	"squirrelchuckle/services"
)


type APNPushTestController struct {
	beego.Controller
}

var tapnService *services.ApplePushService

func (this *APNPushTestController) Get() {
	if apnService == nil || !apnService.Alive() {
		tapnService, _ = services.GetManager().GetServiceByName("APN_Service").(*services.ApplePushService)
	}

	result := tapnService.TestAPN()
	
	this.Data["json"] = result
	this.ServeJson()
}

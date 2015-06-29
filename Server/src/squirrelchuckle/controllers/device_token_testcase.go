package controllers
import (
	"github.com/astaxie/beego"

	"squirrelchuckle/services"
)


type APNPushTestController struct {
	beego.Controller
}

var apnService *services.APNService

func (this *APNPushTestController) Get() {
	if apnService == nil || !apnService.Alive() {
		apnService, _ = services.GetManager().GetServiceByName("APN_Service").(*services.APNService)
	}

	result := apnService.TestAPN()
	
	this.Data["json"] = result
	this.ServeJson()
}

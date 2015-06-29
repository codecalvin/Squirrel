package controllers
import (
	"github.com/astaxie/beego"

	"squirrelchuckle/services"
)


type APNPushController struct {
	beego.Controller
}

var apnService *services.APNService

func (this *APNPushController) Get() {

	if apnService == nil || !apnService.Alive() {
		apnService, _ = services.GetManager().GetServiceByName("APN_Service").(*services.APNService)
	}

	result := apnService.TestAPN()
	
	this.Data["json"] = result
	this.ServeJson()
}

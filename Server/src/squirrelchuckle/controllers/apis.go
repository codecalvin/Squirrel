package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"squirrelchuckle/core"
)

type ApiController struct {
	beego.Controller
}

var EndPoints map[string]string

func (this *ApiController) Get() {
	fmt.Print("ApiController get")
	fmt.Print(EndPoints)
	this.Data["json"] = EndPoints
	this.ServeJson()
}

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Get() {
	core.CloseChan <- true
	this.Ctx.WriteString("Service down")
	this.ServeJson()
}
package controllers

import (
	"github.com/astaxie/beego"
	"squirrelchuckle/core"
)

type ApiController struct {
	beego.Controller
}

var EndPoints map[string]string

func (this *ApiController) Get() {
	this.Data["json"] = EndPoints
	this.ServeJson()
}

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Get() {
	core.CloseChan <- true
	this.ServeJson()
}
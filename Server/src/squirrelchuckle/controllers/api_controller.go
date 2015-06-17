package controllers

import (
	"github.com/astaxie/beego"
)

type ApiController struct {
	beego.Controller
}

var EndPoint map[string]string

func (this *ApiController) Get() {
	this.Data["json"] = EndPoint
	this.ServeJson()
}
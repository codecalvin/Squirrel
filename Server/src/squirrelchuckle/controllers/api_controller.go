package controllers

import (
	"github.com/astaxie/beego"
)

type ApiController struct {
	beego.Controller
}

var EndPoints map[string]string

func (this *ApiController) Get() {
	this.Data["json"] = EndPoints
	this.ServeJson()
}
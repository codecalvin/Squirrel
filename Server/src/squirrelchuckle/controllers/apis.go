package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
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
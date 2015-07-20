package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"squirrelchuckle/core"
	"squirrelchuckle/services"
)

type OneClassUsersController struct {
	beego.Controller
}

func (this *OneClassUsersController) Get() {
	fmt.Print("OneClassUsersController.Get():")

	classKey := this.Ctx.Input.Param(":classKey")
	fmt.Println(classKey)

	classCollection := core.SquirrelApp.DB("squirrel").C("class")
	classResult := services.ClassItem{}
	err := classCollection.Find(bson.M{"elementtype_uniquekey": classKey}).One(&classResult)
	if err != nil {
		fmt.Println("error1")
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	fmt.Println("Results All: ", classResult.RegisterUsers)
	this.Data["json"] = classResult.RegisterUsers
	this.ServeJson()
}




package controllers

import (
	"fmt"
	"strconv"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"squirrelchuckle/database"
)

type OneClassUserCountController struct {
	beego.Controller
}

func (this *OneClassUserCountController) Get() {
	fmt.Print("OneClassUserCountController.Get():")

	classKey := this.Ctx.Input.Param(":classKey")
	fmt.Println(classKey)

	classCollection := database.MSession.DB("squirrel").C("class")
	classResult := ClassItem{}
	err := classCollection.Find(bson.M{"elementtype_uniquekey": classKey}).One(&classResult)
	if err != nil {
		fmt.Println("error1")
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	registeredCount := len(classResult.RegisterUsers)
	registeredCountString := strconv.Itoa(registeredCount)
    userCountInfo := map[string]string{"ElementType_ClassStudent":classResult.ElementType_ClassStudent,
    							"ElementType_RegisteredStudentCount":registeredCountString}
        
	fmt.Println("Results All: ", userCountInfo)
	this.Data["json"] = userCountInfo
	this.ServeJson()
}




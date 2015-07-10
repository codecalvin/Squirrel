package controllers


import (
	"fmt"
	"github.com/astaxie/beego"
	"squirrelchuckle/database"
	"gopkg.in/mgo.v2/bson"
)

type OneClassController struct {
	beego.Controller
}

func (this *OneClassController) Get(){
	fmt.Print("OneClassController.Get():")
	
	classKey := this.Ctx.Input.Param(":classKey")
	fmt.Println(classKey)
	
	c := database.MSession.DB("squirrel").C("class")
	
	result := ClassItem{}
	err := c.Find(bson.M{"elementtype_uniquekey": classKey}).One(&result)
	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	fmt.Println("Results All: ", result)
	this.Data["json"] = result
	this.ServeJson()
}

func (this *OneClassController) Delete(){
	// not used yet
	fmt.Print("OneClassController.Delete():")
	
	classKey := this.Ctx.Input.Param(":classKey")
	fmt.Println(classKey)
	
	c := database.MSession.DB("squirrel").C("class")

	_, err := c.RemoveAll(bson.M{"elementtype_uniquekey": classKey})
	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	this.Data["json"] = map[string]string{"result": "delete success"}
	this.ServeJson()
}
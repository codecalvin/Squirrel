package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"

	"squirrelchuckle/core"
)

type ClassController struct {
	beego.Controller
}

var database *core.Database
func (this *ClassController) Get() {
	if database == nil {
		if service := core.SquirrelApp.GetServiceByName("Database"); service == nil || !service.Alive() {
			this.Ctx.Redirect(300, "/")
		} else {
			database = service.(*core.Database)
		}
	}

	var result [] ProfileInfo
	c := database.MSession.DB("squirrel").C("user")
	q := c.Find(nil)
	iterator := q.Iter()
	_ = iterator.All(&result)

	this.Data["json"] = result
	this.ServeJson()
}

func (this *ClassController) Post() {
	input := this.Input()
	name := input.Get("name")
	id := input.Get("id")

	c := database.MSession.DB("squirrel").C("test_user")

	p := ProfileInfo{UserName:name, UserId:id}
	if cInfo, err := c.Upsert(bson.M{"userid": id}, p); err != nil {
		this.Ctx.Output.Body([]byte(err.Error()))
	} else {
		this.Ctx.Output.Body([]byte(fmt.Sprintf("updated %v, raw name %v, raw id %v", cInfo, name, id)))
	}
}

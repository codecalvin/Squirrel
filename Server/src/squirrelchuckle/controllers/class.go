package controllers


import (
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"

	"squirrelchuckle/services"
)

type ClassController struct {
	beego.Controller
}

func (this *ClassController) Get() {
	var result [] ProfileInfo
	c := services.DatabaseInstance().MSession.DB("squirrel").C("user")
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

	c := services.DatabaseInstance().MSession.DB("squirrel").C("test_user")

	p := ProfileInfo{UserName:name, UserId:id}
	cinfo, err := c.Upsert(bson.M{"userid": id}, p)

	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}

	this.Ctx.Output.Body([]byte(fmt.Sprintf("updated %v, raw name %v, raw id %v", cinfo, name, id)))
}

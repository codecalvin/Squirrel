package controllers


import (
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"

	"squirrelchuckle/services"
	"squirrelchuckle/core"
	"golang.org/x/tools/godoc/redirect"
)

type ClassController struct {
	beego.Controller
}

var database *core.Database
func (this *ClassController) Get() {
	if database == nil {
		if database = core.SquirrelApp.GetServiceByName("Database"); database == nil || database.Alive() {
			database = nil
			// Todo
			redirect.Handler("Service down")
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
	cinfo, err := c.Upsert(bson.M{"userid": id}, p)

	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}

	this.Ctx.Output.Body([]byte(fmt.Sprintf("updated %v, raw name %v, raw id %v", cinfo, name, id)))
}

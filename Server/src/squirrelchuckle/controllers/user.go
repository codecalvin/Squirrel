package controllers


import (
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	
	"squirrelchuckle/core"
)

type UsersController struct {
	beego.Controller
}

type ProfileInfo struct {
	UserName string
	UserId   string
}

func (this *UsersController) Get() {
	var result [] ProfileInfo
	c := core.SquirrelApp.MSession.DB("squirrel").C("user")
	q := c.Find(nil)
	iterator := q.Iter()
	_ = iterator.All(&result)

	this.Data["json"] = result
	this.ServeJson()
}

func (this *UsersController) Post() {
	input := this.Input()
	name := input.Get("name")
	id := input.Get("id")

	c := core.SquirrelApp.MSession.DB("squirrel").C("user")

	p := ProfileInfo{UserName:name, UserId:id}
	if cInfo, err := c.Upsert(bson.M{"userid": id}, p); err != nil {
		this.Ctx.Output.Body([]byte(err.Error()))
	} else {
		this.Ctx.Output.Body([]byte(fmt.Sprintf("updated %v, raw name %v, raw id %v", cInfo, name, id)))
	}
}


// -----------------------------------------------
// User Management
type UserController struct {
	beego.Controller
}


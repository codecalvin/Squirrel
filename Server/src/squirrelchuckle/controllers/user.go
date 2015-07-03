package controllers


import (
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	
	"squirrelchuckle/services"
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
	c := services.DatabaseInstance().MSession.DB("squirrel").C("user")
	q := c.Find(nil)
	iterator := q.Iter()
	print(iterator)
	_ = iterator.All(&result)

	this.Data["json"] = result
	this.ServeJson()
}

func (this *UsersController) Post() {
	input := this.Input()
	name := input.Get("name")
	id := input.Get("id")

	c := services.DatabaseInstance().MSession.DB("squirrel").C("user")

	p := ProfileInfo{UserName:name, UserId:id}
	cinfo, err := c.Upsert(bson.M{"userid": id}, p)

	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}

	this.Ctx.Output.Body([]byte(fmt.Sprintf("updated %v, raw name %v, raw id %v", cinfo, name, id)))
}


// -----------------------------------------------
// User Management
type UserController struct {
	beego.Controller
}


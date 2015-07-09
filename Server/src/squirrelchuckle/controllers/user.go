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

var userService *services.UserService

func (this *UsersController) Get() {
	var result [] ProfileInfo
	q := userService.Find(nil)
	iterator := q.Iter()
	_ = iterator.All(&result)

	this.Data["json"] = result
	this.ServeJson()
}

func (this *UsersController) Post() {
	input := this.Input()
	name := input.Get("name")
	id := input.Get("id")

	p := ProfileInfo{UserName:name, UserId:id}
	if cInfo, err := userService.Upsert(bson.M{"userid": id}, p); err != nil {
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


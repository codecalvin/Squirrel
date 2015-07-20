package controllers


import (
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"squirrelchuckle/core"
	"squirrelchuckle/services"
)

type UnRegisterController struct {
	beego.Controller
}


func (this *UnRegisterController) Post() {
	fmt.Println("RegisterController Post:")

	input := this.Input()
	eventUniqueKey := input.Get("ElementType_UniqueKey")
	eventName := input.Get("ElementType_ClassName")
	userUniqueKey := input.Get("ElementType_UserUniqueKey")
	userName := input.Get("ElementType_UserName")
	
	fmt.Println(eventUniqueKey)
	fmt.Println(eventName)
	fmt.Println(userUniqueKey)
	fmt.Println(userName)
	
	// register user to class
	c := core.SquirrelApp.DB("squirrel").C("class")
	result := services.ClassItem{}
	err := c.Find(bson.M{"elementtype_uniquekey": eventUniqueKey}).One(&result)
	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
		this.Data["json"] = map[string]string{"result": "not found"}
		this.ServeJson()
		return
	} else{
		newClassItemResult := result
		_, err = c.RemoveAll(bson.M{"elementtype_uniquekey": eventUniqueKey})
		if err != nil{
			this.Ctx.Output.Body([]byte(err.Error()))
		}
		delete(newClassItemResult.RegisterUsers, userUniqueKey)
		err = c.Insert(&newClassItemResult)
		fmt.Println("user in class Results All: ", newClassItemResult.RegisterUsers)
	}
	
	// register class to user
	if user, ok := userService.Users[userUniqueKey]; ok {
		delete(user.Classes, eventUniqueKey)
	}
	this.Data["json"] = map[string]string{"result": "OK"}
	this.ServeJson()
}

package controllers


import (
	"fmt"
	"github.com/astaxie/beego"
	"squirrelchuckle/database"
	"gopkg.in/mgo.v2/bson"
)

type RegisterStatusController struct {
	beego.Controller
}

func (this *RegisterStatusController) Get(){
	fmt.Println("RegisterStatusController Get():")
	
	input := this.Input()
	eventUniqueKey := input.Get("ElementType_UniqueKey")
	eventName := input.Get("ElementType_ClassName")
	userUniqueKey := input.Get("ElementType_UserUniqueKey")
	userName := input.Get("ElementType_UserName")
	
	fmt.Println("input: ", eventUniqueKey)
	fmt.Println(eventName)
	fmt.Println(userUniqueKey)
	fmt.Println(userName)
	
	// register user to class
	c := database.MSession.DB("squirrel").C("class")
	result := ClassItem{}
	err := c.Find(bson.M{"elementtype_uniquekey": eventUniqueKey}).One(&result)
	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	} 
	
	status := "NO"
	_, existingKey := result.RegisterUsers[userUniqueKey]
	if existingKey {
		status = "YES"
	}
	
	fmt.Println(userUniqueKey, status)
	this.Data["json"] = map[string]string{userUniqueKey: status}
	this.ServeJson()
}


package controllers


import (
	"fmt"
	"strconv"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"squirrelchuckle/core"
	"squirrelchuckle/services"
)

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) Get(){
	// not used yet
	fmt.Print("RegisterController Get:")

	this.Data["json"] = map[string]string{"result": "to do"}
	this.ServeJson()
}

func (this *RegisterController) Post() {
	fmt.Print("RegisterController Post:")

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
	
	registeredCount := len(result.RegisterUsers)
	maxCount, _:= strconv.Atoi(result.ElementType_ClassStudent)
	if registeredCount >= maxCount{
		this.Data["json"] = map[string]string{"result": "registered count is full, no more seat"}
		this.ServeJson()
		return 
		}
	
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
		newClassItemResult.RegisterUsers[userUniqueKey]=userName
		err = c.Insert(&newClassItemResult)
		fmt.Println("user in class Results All: ", newClassItemResult.RegisterUsers)
	}
	
	// register class to user
	var user *services.User
	var ok bool
	if user, ok = userService.Users[userUniqueKey]; !ok {
		this.Data["json"] = map[string]string{"result": "fail"}
		this.ServeJson()
		return
	}

	if user.Classes == nil {
		user.Classes = make(map[string] *services.ClassBriefItem)
	}

	user.Classes[eventUniqueKey] = &services.ClassBriefItem{ElementType_UniqueKey:result.ElementType_UniqueKey,
		ElementType_ClassName:result.ElementType_ClassName,
		ElementType_ClassTime:result.ElementType_ClassTime,
	}

    this.Data["json"] = map[string]string{"result": "OK"}
	this.ServeJson()
}

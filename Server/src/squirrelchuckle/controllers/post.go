package controllers


import (
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
	"squirrelchuckle/core"
	"squirrelchuckle/services"
)

type PostController struct {
	beego.Controller
}

func (this *PostController) Post() {
	fmt.Print("post go:")

	input := this.Input()
	eventUniqueKey := input.Get("ElementType_UniqueKey")
	eventName := input.Get("ElementType_ClassName")
	eventTime := input.Get("ElementType_ClassTime")
	teacher := input.Get("ElementType_ClassTeacher")
	studentCount := input.Get("ElementType_ClassStudent")
	eventDescription := input.Get("ElementType_ClassDescription")
	
	fmt.Println(eventUniqueKey)
	fmt.Println(eventName)
	fmt.Println(eventTime)
	fmt.Println(teacher)
	fmt.Println(studentCount)
	fmt.Println(eventDescription)
	
	c := core.SquirrelApp.DB("squirrel").C("class")
	result := services.ClassItem{}
	err := c.Find(bson.M{"elementtype_uniquekey": eventUniqueKey}).One(&result)
	if err != nil{
		err = c.Insert(&services.ClassItem{eventUniqueKey, eventName, eventTime,teacher,studentCount,eventDescription, make(map[string]string)})
	} else {
			newResult := result
			_, err = c.RemoveAll(bson.M{"elementtype_uniquekey": eventUniqueKey})
			if err != nil{
				this.Ctx.Output.Body([]byte(err.Error()))
			}
			newResult.ElementType_UniqueKey = eventUniqueKey
			newResult.ElementType_ClassName = eventName
			newResult.ElementType_ClassTime = eventTime
			newResult.ElementType_ClassTeacher = teacher
			newResult.ElementType_ClassStudent = studentCount
			newResult.ElementType_ClassDescription = eventDescription
			err = c.Insert(&newResult)
		}
	
	var resultAll [] services.ClassItem
	err = c.Find(nil).All(&resultAll)
	fmt.Println("Results All: ", resultAll)
	
    this.Data["json"] = map[string]string{"result": "OK"}
	this.ServeJson()
}

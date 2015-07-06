package controllers


import (
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"

	"squirrelchuckle/database"
)

type ClassController struct {
	beego.Controller
}

func (this *ClassController) Get() {
	fmt.Print("\nClassController::Get()\n") 
        
	var result [] ClassItem
	c := database.MSession.DB("squirrel").C("class")
	err := c.Find(nil).All(&result)
	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	
	briefItems := map[string]string{}
	for index := 0; index < len(result); index++{
		if len(result[index].ElementType_UniqueKey) < 1{
			continue
			}
		briefItems[result[index].ElementType_UniqueKey] = result[index].ElementType_ClassName
	}
	
	fmt.Print(briefItems)
	fmt.Print("\n")
	this.Data["json"] = briefItems
	this.ServeJson()
}

func (this *ClassController) Post() {
	fmt.Print("post go:")
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

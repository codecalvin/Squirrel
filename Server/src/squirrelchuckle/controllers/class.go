package controllers


import (
	"fmt"
	"github.com/astaxie/beego"
	"squirrelchuckle/core"
	"squirrelchuckle/services"
)

type ClassController struct {
	beego.Controller
}

func (this *ClassController) Get() {
	var result [] services.ClassItem
	c := core.SquirrelApp.DB("squirrel").C("class")
	err := c.Find(nil).All(&result)
	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	
	briefItems := map[string]services.ClassBriefItem{}
	for index := 0; index < len(result); index++{
		if len(result[index].ElementType_UniqueKey) < 1{
			continue
			}
		briefItem := services.ClassBriefItem{ElementType_UniqueKey:result[index].ElementType_UniqueKey,
    							ElementType_ClassName:result[index].ElementType_ClassName,
    							ElementType_ClassTime:result[index].ElementType_ClassTime}
		briefItems[result[index].ElementType_UniqueKey] = briefItem
	}
	
	fmt.Print(briefItems)
	fmt.Print("\n")
	this.Data["json"] = briefItems
	this.ServeJson()
}

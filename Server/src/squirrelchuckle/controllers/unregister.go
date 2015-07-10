package controllers


import (
	"fmt"
	"github.com/astaxie/beego"
	"squirrelchuckle/database"
	"gopkg.in/mgo.v2/bson"
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
	c := database.MSession.DB("squirrel").C("class")
	result := ClassItem{}
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
	fmt.Println(database.DataBaseNameString)
	fmt.Println(database.DataBaseUserCollectionNameString)
	userCollection := database.MSession.DB("squirrel").C("user")
	fmt.Println("flag 0")
	userResult := UserItem{}
	err = userCollection.Find(bson.M{"elementtype_useruniquekey": userUniqueKey}).One(&userResult)
	if err != nil {
		fmt.Println(" error1")
		err = userCollection.Insert(&UserItem{userUniqueKey, userName, map[string]ClassBriefItem{}})
		err = userCollection.Find(bson.M{"elementtype_useruniquekey": userUniqueKey}).One(&userResult)
		if err != nil {
			fmt.Println(" error2")
		}
	} 
	fmt.Println("flag 3")
	newUserResult := userResult
	delete(newUserResult.Classes, eventUniqueKey)
	fmt.Println("Class in a user Results All: ", newUserResult.Classes)
	
	// to do: improve it. Workaround by delete and insert for update
	//_, err = userCollection.Upsert(bson.M{"elementtype_useruniquekey": userUniqueKey}, newUserResult)
	_, err = userCollection.RemoveAll(bson.M{"elementtype_useruniquekey": userUniqueKey})
	if err != nil{
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	err = userCollection.Insert(&newUserResult)

	fmt.Println("user:", userResult)
    this.Data["json"] = map[string]string{"result": "OK"}
	this.ServeJson()
}

package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"

	"squirrelchuckle/database"
)

type UsersController struct {
	beego.Controller
}

type ProfileInfo struct {
	UserName string
	UserId   string
}

type ClassBriefItem struct {
	ElementType_UniqueKey string
	ElementType_ClassName string
}

type ClassItem struct {
	ElementType_UniqueKey        string
	ElementType_ClassName        string
	ElementType_ClassTime        string
	ElementType_ClassTeacher     string
	ElementType_ClassStudent     string
	ElementType_ClassDescription string
	RegisterUsers                map[string]string
}

type UserItem struct {
	ElementType_UserUniqueKey string
	ElementType_UserName      string
	Classes                   map[string]string
}

type RegisterItem struct {
	ElementType_UniqueKey     string
	ElementType_ClassName     string
	ElementType_UserUniqueKey string
	ElementType_UserName      string
}

func (this *UsersController) Get() {

	// not used yet
	fmt.Print("UsersController.Get():")

	userKey := this.Ctx.Input.Param(":userKey")
	fmt.Println(userKey)

	userCollection := database.MSession.DB("squirrel").C("user")
	var userResult = UserItem{}
	err := userCollection.Find(bson.M{"elementtype_useruniquekey": userKey}).One(&userResult)
	if err != nil {
		fmt.Println("error1")
		this.Ctx.Output.Body([]byte(err.Error()))
	}
	fmt.Println("Results All: ", userResult.Classes)
	this.Data["json"] = userResult.Classes
	this.ServeJson()
}

func (this *UsersController) Post() {
	input := this.Input()
	name := input.Get("name")
	id := input.Get("id")

	c := database.MSession.DB("squirrel").C("user")

	p := ProfileInfo{UserName: name, UserId: id}
	cinfo, err := c.Upsert(bson.M{"userid": id}, p)

	if err != nil {
		this.Ctx.Output.Body([]byte(err.Error()))
	}

	this.Ctx.Output.Body([]byte(fmt.Sprintf("updated %v, raw name %v, raw id %v", cinfo, name, id)))
}

// -----------------------------------------------
// User Management
type UserController struct {
	beego.Controller
}

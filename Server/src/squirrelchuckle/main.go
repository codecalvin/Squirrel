package main

import (
	_ "squirrelchuckle/routers"
	database "squirrelchuckle/database"
	"github.com/astaxie/beego"
)

func main() {
	defer database.Close()
	beego.Run()
}


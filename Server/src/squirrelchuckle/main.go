package main

import (
	"github.com/astaxie/beego"

	"squirrelchuckle/core"
	_ "squirrelchuckle/routers"
)

func main() {
	core.Run()
	beego.Run()
}

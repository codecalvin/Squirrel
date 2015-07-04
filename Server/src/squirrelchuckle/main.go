package main

import (
	"github.com/astaxie/beego"

	"squirrelchuckle/core"
)

func main() {
	defer core.CloseApp()
	beego.Run()
}

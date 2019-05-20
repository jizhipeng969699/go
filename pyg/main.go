package main

import (
	"github.com/astaxie/beego"
	_ "pyg/models"
	_ "pyg/routers"
)

func main() {
	beego.Run()
}

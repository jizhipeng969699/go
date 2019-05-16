package main

import (
	"github.com/astaxie/beego"
	_ "pinyg/models"
	_ "pinyg/routers"
)

func main() {
	beego.Run()
}

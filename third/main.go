package main

import (
	//"third/controllers"
	_"third/controllers"
	"github.com/astaxie/beego"
	_ "third/routers"
)

func main() {
	//controllers.Rfunc()

	beego.AddFuncMap("Prepage", Prepage)
	beego.AddFuncMap("Nextpage", Nextpage)
	beego.Run()
}

func Prepage(pagenum int) int {
	if pagenum <= 1 {
		return 1
	}
	return pagenum - 1
}
func Nextpage(pagenum int, pagecount float64) int {
	if pagenum >= int(pagecount) {
		return int(pagecount)
	}
	return pagenum + 1
}

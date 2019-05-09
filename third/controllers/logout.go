package controllers

import "github.com/astaxie/beego"

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Logout() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}

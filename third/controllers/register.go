package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"third/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

func (this *UserController) Register() {
	userName := this.GetString("userName")
	password := this.GetString("password")

	if userName == "" || password == "" {
		beego.Error("数据不能为空")
		this.TplName = "register.html"
		return
	}

	o := orm.NewOrm()
	var user models.Thirduser
	user.Name = userName
	user.Pwd = password

	n, err := o.Insert(&user)
	if err != nil {
		beego.Error("插入失败")
		this.TplName = "register.html"
		return
	}
	beego.Info(n)
	this.Redirect("/login", 302)
}

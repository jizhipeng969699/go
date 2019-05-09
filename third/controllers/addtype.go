package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"third/models"
)

type AddtypeController struct {
	beego.Controller
}

func (this *AddtypeController) ShowAddtype() {
	username := this.GetSession("userName")
	if username == nil {
		this.Redirect("/login", 302)
		return
	}
	//这是类型断言
	this.Data["username"] = username.(string) //返回值是一个接口的时候 最好使用 类型断言后在使用

	o := orm.NewOrm()
	//txttypes := []models.Txttype{}
	var txttypes []models.Txttype
	o.QueryTable("Txttype").All(&txttypes)

	this.Data["txttypes"] = txttypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

func (this *AddtypeController) Addtype() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		beego.Info("数据不完整")
		this.Redirect("/zzz/addtype", 302)
		return
	}
	o := orm.NewOrm()
	addtype := models.Txttype{Typename: typeName}
	_, err := o.Insert(&addtype)
	if err != nil {
		beego.Info("插入失败")
		this.Redirect("/zzz/addtype", 302)
		return
	}
	this.Redirect("/zzz/addtype", 302)

}

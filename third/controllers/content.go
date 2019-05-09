package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"third/models"
)

type ContentController struct {
	beego.Controller
}

func (this *ContentController) ShowContent() {
	username := this.GetSession("userName")
	if username == nil {
		this.Redirect("/login", 302)
		return
	}

	this.Data["username"] = username.(string)

	id, err := this.GetInt("Id")
	if err != nil {
		beego.Info("获取id失败")
		this.Redirect("/zzz/index", 302)
		return
	}

	o := orm.NewOrm()
	txt := models.Txt{Id: id}
	err = o.Read(&txt)
	if err != nil {
		beego.Info("Content Read(&txt)", err)
		this.Redirect("/zzz/index", 302)
		return
	}
	txt.Readcount += 1

	//多表查询yi
	//o.LoadRelated(&txt,"Thirduser")
	//多表查询er
	var users []models.Thirduser
	//高级查询二               字段名————表明————字段  ，要比较的字段
	o.QueryTable("Thirduser").Filter("Txt__Txt__Id", id).Distinct().All(&users)
	this.Data["users"] = users

	o.Update(&txt)

	this.Data["txt"] = txt

	userName := this.GetSession("userName")
	var user models.Thirduser
	user.Name = userName.(string)
	o.Read(&user, "Name")

	m2m := o.QueryM2M(&user, "Txt")

	m2m.Add(txt)

	this.Layout = "layout.html"
	this.TplName = "content.html"
}

//func (this *ContentController) Content() {
//	id, err := this.GetInt("Id")
//	if err != nil {
//		beego.Info("获取id失败")
//		this.Redirect("/index", 302)
//		return
//	}
//
//	o := orm.NewOrm()
//	txt := models.Txt{Id: id}
//	err = o.Read(&txt)
//	if err != nil {
//		beego.Info("Content Read(&txt)", err)
//		this.Redirect("/index", 302)
//		return
//	}
//
//	this.Data["txt"] = txt
//}

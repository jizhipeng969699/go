package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"path"
	"third/models"
	"time"
)

type AddController struct {
	beego.Controller
}

func (this *AddController) ShowAdd() {
	username := this.GetSession("userName")
	if username == nil {
		this.Redirect("/login", 302)
		return
	}

	this.Data["username"]=username.(string)

	o := orm.NewOrm()
	//txttypes := []models.Txttype{}
	var txttypes []models.Txttype
	o.QueryTable("Txttype").All(&txttypes)

	this.Data["txttypes"] = txttypes

	this.Layout = "layout.html"
	this.TplName = "add.html"
}
func (this *AddController) Add() {

	articleName := this.GetString("articleName")
	content := this.GetString("content")

	if articleName == "" || content == "" {
		beego.Error("数据不完整")
		this.Redirect("/zzz/add", 302)
		return
	}

	file, head, err := this.GetFile("uploadname")
	if err != nil {
		beego.Error("上传图片错误")
		this.Redirect("/zzz/add", 302)
		return
	}
	defer file.Close()

	if head.Size > 5000000 {
		beego.Error("上传图片过大")
		this.Redirect("/zzz/add", 302)
		return
	}
	//获得文件后缀名
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" {
		beego.Error("上传图片格式错误")
		this.Redirect("/zzz/add", 302)
		return
	}

	fname := time.Now().Format("201601021504052222") + ext

	err = this.SaveToFile("uploadname", "./static/img/"+fname)
	if err != nil {
		beego.Error("图片存储失败")
		this.Redirect("/zzz/add", 302)
		return
	}

	o := orm.NewOrm()
	var txt models.Txt
	txt.Title = articleName
	txt.Content = content
	txt.Img = "/static/img/" + fname

	typename := this.GetString("select")
	var txttype models.Txttype
	txttype.Typename = typename
	o.Read(&txttype, "Typename")

	txt.Txttype = &txttype

	n, err := o.Insert(&txt)
	if err != nil {
		beego.Error("插入失败", err)
		this.Redirect("/zzz/add", 302)
		return
	}

	this.Redirect("/zzz/index", 302)
	beego.Info(n, "index", 302)
}

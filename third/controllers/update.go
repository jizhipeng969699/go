package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"path"
	"strconv"
	"third/models"
	"time"
)

type UpdateController struct {
	beego.Controller
}

func (this *UpdateController) ShowUpdate() {
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

	this.Data["txt"] = txt
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

func (this *UpdateController) Update() {
	id, err := this.GetInt("Id")
	if err != nil {
		beego.Info("获取id失败")
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		return
	}
	o := orm.NewOrm()
	txt := models.Txt{Id: id}
	err = o.Read(&txt)
	if err != nil {
		beego.Info("Content Read(&txt)", err)
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		return
	}
	//<input value="{{.txt.Id}}" name="id" hidden="hidden">
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	if articleName == "" || content == "" {
		beego.Error("数据不完整")
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		this.Data["errmsg"] = "数据不完整"
		return
	}

	file, head, err := this.GetFile("uploadname")
	if err != nil {
		beego.Error("上传图片错误")
		this.Data["errmsg"] = "上传图片错误"
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		return
	}
	defer file.Close()

	if head.Size > 5000000 {
		beego.Error("上传图片过大")
		this.Data["errmsg"] = "上传图片过大"
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		return
	}
	//获得文件后缀名
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" {
		beego.Error("上传图片格式错误")
		this.Data["errmsg"] = "上传图片格式错误"
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		return
	}

	fname := time.Now().Format("201601021504052222") + ext

	err = this.SaveToFile("uploadname", "./static/img/"+fname)
	if err != nil {
		beego.Error("图片存储失败")
		this.Data["errmsg"] = "图片存储失败"
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		return
	}

	txt.Title = articleName
	txt.Content = content
	txt.Img = "/static/img/" + fname

	_, err = o.Update(&txt)
	if err != nil {
		beego.Error("更新失败")
		this.Data["errmsg"] = "更新失败"
		this.Redirect("/zzz/update?Id="+strconv.Itoa(id), 302)
		return
	}

	this.Redirect("/zzz/index", 302)

}

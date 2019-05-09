package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"third/models"
)

type DeleteController struct {
	beego.Controller
}

func (this *DeleteController) Delete() {
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

	filepath := txt.Img
	err = os.Remove("." + filepath) //删除存储的图片
	if err != nil {
		beego.Info("删除失败", err)
		this.Redirect("/zzz/index", 302)
		return
	}

	n, err := o.Delete(&txt)
	if err != nil {
		beego.Info("删除失败", err)
		this.Redirect("/zzz/index", 302)
		return
	}
	beego.Info(n)
	this.Redirect("/zzz/index", 302)
}

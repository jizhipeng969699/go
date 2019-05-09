package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"third/models"
)

type DeletetypeController struct {
	beego.Controller
}

func (this *DeletetypeController) Deletetype() {
	id, err := this.GetInt("Id")
	if err != nil {
		beego.Info("获取id失败")
		this.Redirect("/zzz/addtype", 302)
		return
	}

	o := orm.NewOrm()
	var txttype models.Txttype
	txttype.Id = id

	o.Delete(&txttype) //一对多删除 主表 信息时  默认是 级联删除
	this.Redirect("/zzz/addtype", 302)

}

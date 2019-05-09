package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"third/models"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	//这里使用 base64 进行解密
	dec, _ := base64.StdEncoding.DecodeString(userName)
	if userName != "" {
		this.Data["userName"] = string(dec)
		this.Data["checked"] = "checked"
	} else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}

	this.TplName = "login.html"

}

func (this *LoginController) Login() {
	userName := this.GetString("userName")
	password := this.GetString("password")

	if userName == "" || password == "" {
		beego.Error("用户名或密码不能为空")
		this.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	var user models.Thirduser
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("没有找到用户名")
		this.TplName = "login.html"
		return
	}
	if password != user.Pwd {
		beego.Error("密码错误")
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")
	//beego的cookie 默认不支持中文 要注册中文用户名  需要先进行 加密
	enc := base64.StdEncoding.EncodeToString([]byte(userName))

	if remember == "on" {
		this.Ctx.SetCookie("userName", enc, 60)
	} else {
		this.Ctx.SetCookie("userName", enc, -1)
	}

	this.SetSession("userName", userName) //保存登陆状态  session 在beego中需要手动开启 在conf中

	this.Redirect("/zzz/index", 302)
	//this.Ctx.WriteString("ok")

}

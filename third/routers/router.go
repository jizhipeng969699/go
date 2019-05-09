package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"third/controllers"
	_ "third/models"
)

func init() {

	//插入过滤器       参数1 正则表达式  参数2 过滤器的位置   参数3 是一个回调函数 要干什么事
	beego.InsertFilter("/zzz/*", beego.BeforeExec, Filterfunc)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:Register")
	beego.Router("/login", &controllers.LoginController{}, "get:ShowLogin;post:Login")
	beego.Router("/zzz/index", &controllers.IndexController{}, "get:ShowIndex")
	beego.Router("/zzz/add", &controllers.AddController{}, "get:ShowAdd;post:Add")

	beego.Router("/zzz/addtype", &controllers.AddtypeController{}, "get:ShowAddtype;post:Addtype")

	beego.Router("/zzz/content", &controllers.ContentController{}, "get:ShowContent")

	beego.Router("/zzz/delete", &controllers.DeleteController{}, "get:Delete")

	beego.Router("/zzz/update", &controllers.UpdateController{}, "get:ShowUpdate;post:Update")

	beego.Router("/zzz/logout", &controllers.LogoutController{}, "get:Logout")

	beego.Router("/zzz/deletetype", &controllers.DeletetypeController{}, "get:Deletetype")

}

func Filterfunc(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == "" {
		ctx.Redirect(302, "/login") //等同与  this.Redirect("/login",302)  他的参数正好相干
		return
	}
}

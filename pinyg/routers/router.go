package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"pinyg/controllers"
)

func init() {
	beego.InsertFilter("/user/*", beego.BeforeExec, fnFunc) //路由过滤

	//注册
	beego.Router("/register", &controllers.UserController{}, "get:Showregister;post:Register")
	//处理邮箱注册
	beego.Router("/register-email", &controllers.UserController{}, "get:Showremail;post:Email")
	//发送短信																	//这里的post请求时通过ajas发送的 .post请求
	beego.Router("/sendMsg", &controllers.UserController{}, "post:Sendmsg")
	//激活用户
	beego.Router("/active", &controllers.UserController{}, "get:Active")
	//用户登录
	beego.Router("/login", &controllers.UserController{}, "get:Showlogin;post:Login")
	//用户退出
	beego.Router("/user/logout", &controllers.UserController{}, "get:Logout")

	//主界面
	beego.Router("/index", &controllers.GoolsController{}, "get:Showindex")
	//展示用户中心
	beego.Router("/user/user_center_info", &controllers.GoolsController{}, "get:ShowCenterinfo")
	//展示用户order中心
	beego.Router("/user/user_center_order", &controllers.GoolsController{}, "get:ShowCenterorder")
	//展示用户site中心
	beego.Router("/user/user_center_site", &controllers.GoolsController{}, "get:ShowCentersite;post:CenterSite")
	//展示生鲜模块
	beego.Router("/user/index_sx", &controllers.GoolsController{}, "get:ShowIndex_sx")

	//商品详细信息页
	beego.Router("/user/detail", &controllers.GoolsController{}, "get:Showgoodsdetail")
	//商品类型列表页
	beego.Router("/user/list", &controllers.GoolsController{}, "get:Showgoodslist")
}
func fnFunc(ctx *context.Context) {
	//过滤路由
	username := ctx.Input.Session("username")
	if username == nil {
		ctx.Redirect(302, "/index")
		return
	}
}

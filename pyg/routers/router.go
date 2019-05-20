package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"pyg/controllers"
)

func init() {
	//路由过滤器
	beego.InsertFilter("/user/*", beego.BeforeExec, guolcFunc)

	beego.Router("/", &controllers.MainController{})
	//用户注册陆游
	beego.Router("/register", &controllers.Usercontrollers{}, "get:Showregister;post:Register")
	//处理验证码
	beego.Router("/sendmsg", &controllers.Usercontrollers{}, "post:Sendmsg")
	//注册邮箱
	beego.Router("/register-email", &controllers.Usercontrollers{}, "get:Showemail;post:Email")
	//激活用户
	beego.Router("/active", &controllers.Usercontrollers{}, "get:Active")
	//用户登录
	beego.Router("/login", &controllers.Usercontrollers{}, "get:Showlogin;post:Login")
	//退出登陆
	beego.Router("/user/logout", &controllers.Usercontrollers{}, "get:Logout")

	//展示用户信息中心
	beego.Router("/user/userinfo", &controllers.Usercontrollers{}, "get:ShowUserinfo")
	//展示用户收货地质
	beego.Router("/user/usersite", &controllers.Usercontrollers{}, "get:ShowUsersite;post:Usersite")
	//展示用户订单
	beego.Router("/user/order", &controllers.Usercontrollers{}, "get:ShowUserorder")

	//主页信息
	beego.Router("/index", &controllers.Gooscontrollers{}, "get:Showindex")
	//主页生鲜模块
	beego.Router("/index_sx", &controllers.Gooscontrollers{}, "get:Showindex_sx")
	//展示商品详细
	beego.Router("/user/goodsdetail", &controllers.Gooscontrollers{}, "get:Showgoodsdetail")
	//统一类型的商品
	beego.Router("/user/goodslist", &controllers.Gooscontrollers{}, "get:Showgoodslist")
	//模糊查找
	beego.Router("/serch", &controllers.Gooscontrollers{}, "get,post:Showserch")

	//添加购物车
	beego.Router("/addCart", &controllers.Cartcontrollers{}, "post:Cart")
	//展示购物车
	beego.Router("/user/cart", &controllers.Cartcontrollers{}, "get:Showcart")
	//更新redis count数据加 减
	beego.Router("/upCart", &controllers.Cartcontrollers{}, "post:UpCart")
	//删除redis中的数据
	beego.Router("/delete", &controllers.Cartcontrollers{}, "post:Delete")
	//提交订单 展示订单页面
	beego.Router("/user/place_order", &controllers.Ordercontrollers{}, "post:Showorder")
	//处理订单 添加订单
	beego.Router("/user/pushOrder", &controllers.Ordercontrollers{}, "post:Pushorder")
	//支付
	beego.Router("/pay", &controllers.Ordercontrollers{}, "get:Pay")

}
func guolcFunc(ctx *context.Context) {
	//过滤校验
	name := ctx.Input.Session("name")
	if name == nil {
		ctx.Redirect(302, "/login") //状态马 不能 乱写
		return
	}

}

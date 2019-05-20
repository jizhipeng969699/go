package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/gomodule/redigo/redis"
	"math"
	"math/rand"
	"pyg/models"
	"regexp"
	"time"
)

type Usercontrollers struct {
	beego.Controller
}

func (this *Usercontrollers) Showregister() {
	this.TplName = "register.html"
}

type Message struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
} //接收短信格式  json

//传输json格式
func Respfunc(this *beego.Controller, resp map[string]interface{}) {
	//把容器传递给前端
	this.Data["json"] = resp //给前端说明  是一个json格式的数据
	//指定传输的方式
	this.ServeJSON()
}

//发送短信
func (this *Usercontrollers) Sendmsg() {
	phone := this.GetString("phone")
	resp := make(map[string]interface{})
	defer Respfunc(&this.Controller, resp) //面向对象
	if phone == "" {
		beego.Error("获取手机号错误")

		//通过json格式 传值   定义map类型的容器   需要先make 初始化
		//resp := make(map[string]interface{})

		resp["errno"] = 1
		resp["errmsg"] = "获取手机号错误"
		//	指定传输的方式  是服务器json方式

		//Respfunc(this, resp)
		return
	}

	//检查电话的格式是否正确
	reg, _ := regexp.Compile(`^1[3-9][0-9]{9}$`)

	result := reg.FindString(phone)

	if result == "" {
		beego.Error("电话格式错误")

		//通过json格式 传值   定义map类型的容器   需要先make 初始化
		resp["errno"] = 2
		resp["errmsg"] = "电话格式错误"
		//	指定传输的方式  是服务器json方式
		//Respfunc(this, resp)
		return
	}
	//发送短信   SDK调用
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "LTAIu4sh9mfgqjjr", "sTPSi0Ybj0oFyqDTjQyQNqdq9I9akE")
	if err != nil {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 3
		resp["errmsg"] = "初始化短信错误"
		return
	}
	//生成6位数随机数
	rand.Seed(time.Now().UnixNano()) //设置随机数种子
	//aa := rand.Intn(899999) + 100000
	aa := fmt.Sprintf("%06d", rand.Intn(899999)+100000)

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = "品优购"
	request.QueryParams["TemplateCode"] = "SMS_164275022"
	request.QueryParams["TemplateParam"] = "{\"code\":" + aa + "}"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		beego.Error("短信发送失败", err)
		//2.给容器赋值
		resp["errno"] = 4
		resp["errmsg"] = "短信发送失败"
		return
	}

	//json数据解析
	var message Message
	json.Unmarshal(response.GetHttpContentBytes(), &message)
	if message.Message != "OK" {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 6
		resp["errmsg"] = message.Message
		return
	}
	resp["errno"] = 5
	resp["errmsg"] = "发送成功"
	resp["num"] = aa
}

//处理登陆信息
func (this *Usercontrollers) Register() {
	phone := this.GetString("phone")
	password := this.GetString("password")
	repassword := this.GetString("repassword")

	if phone == "" || password == "" || repassword == "" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "register.html"
		return
	}
	if password != repassword {
		beego.Error("两次密码输入不一致")
		this.Data["errmsg"] = "两次密码输入不一致"
		this.TplName = "register.html"
		return
	}

	//处理和插入数据
	o := orm.NewOrm()
	var user models.User
	user.Name = phone
	user.Phone = phone
	user.Pwd = password
	_, err := o.Insert(&user)
	if err != nil {
		beego.Error("插入失败")
		this.Data["errmsg"] = "插入失败"
		this.TplName = "register.html"
		return
	}

	//把用户名保存到本地 cookie中
	this.Ctx.SetCookie("userName", user.Name, 60*10)
	//跳转界面
	this.Redirect("register-email", 302)

}

//展示注册邮箱界面
func (this *Usercontrollers) Showemail() {
	this.TplName = "register-email.html"
}

//注册邮箱
func (this *Usercontrollers) Email() {
	//获得数据
	email := this.GetString("email")
	password := this.GetString("password")
	repassword := this.GetString("repassword")
	//比较数据
	if email == "" || password == "" || repassword == "" {
		beego.Error("输入的数据不完整")
		this.Data["errmsg"] = "输入的数据不完整"
		this.TplName = "register-email.html"
		return
	}
	if password != repassword {
		beego.Error("两次的密码不一致")
		this.Data["errmsg"] = "两次的密码不一致"
		this.TplName = "register-email.html"
		return
	}
	//校验邮箱的格式
	rep, err := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	if err != nil {
		beego.Error("正则错误", err)
		this.Data["errmsg"] = "两次的密码不一致"
		this.TplName = "register-email.html"
		return
	}
	result := rep.FindString(email)
	if result == "" {
		beego.Error("邮箱格式错误")
		this.Data["errmsg"] = "邮箱格式错误"
		this.TplName = "register-email.html"
		return
	}

	//处理数据
	//发送邮件
	//utils     全局通用接口  工具类  邮箱配置
	//参数 1  邮箱的用户名   参数2 是  授权玛                   参数3  smtp  格式			参数4  port端口号
	config := `{"username":"2297134578@qq.com","password":"xjvmmqpyyfuydifb","host":"smtp.qq.com","port":587}`

	emailReg := utils.NewEMail(config) //建立一个新的email 对象
	//内容配置
	emailReg.Subject = "品优购用户激活"        //设置邮件的标题
	emailReg.From = "2297134578@qq.com" //邮件发出处
	emailReg.To = []string{email}       //邮件要发到哪里
	userName := this.Ctx.GetCookie("userName")
	emailReg.HTML = `<a href="http://127.0.0.1:8080/active?userName=` + userName + `"> 点击激活该用户</a>`

	//发送
	err = emailReg.Send()
	if err != nil {
		beego.Info(err, "----")
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err = o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名不存在", err)
		this.Data["errmsg"] = "用户名不存在"
		this.TplName = "register-email.html"
		return
	}

	user.Email = email
	_, err = o.Update(&user)
	if err != nil {
		beego.Error("更新邮箱失败")
		this.Data["errmsg"] = "更新邮箱失败"
		this.TplName = "register-email.html"
		return
	}

	//返回数据
	this.Ctx.WriteString("邮件已发送，请去目标邮箱激活用户！")
}

//激活邮箱
func (this *Usercontrollers) Active() {
	userName := this.GetString("userName")
	if userName == "" {
		beego.Error("获取用户名失败")
		this.Data["errmsg"] = "获取用户名失败"
		this.Redirect("/register-email", 302)
		return
	}

	//通过用户名去更新他的  Active 为 true
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("获取用户名失败")
		this.Data["errmsg"] = "获取用户名失败"
		this.Redirect("/register-email", 302)
		return
	}

	user.Active = true //激活
	_, err = o.Update(&user)
	if err != nil {
		beego.Error("激活失败")
		this.Data["errmsg"] = "激活失败"
		this.Redirect("/register-email", 302)
		return
	}

	this.Redirect("/login", 302)
	beego.Info("集获成功")
}

//展示登陆界面
func (this *Usercontrollers) Showlogin() {
	name := this.Ctx.GetCookie("logname")

	//var buffer bytes.Buffer
	//dec := gob.NewDecoder(&buffer)
	//dec.Decode(name)

	if name == "" {
		this.Data["checked"] = ""
	} else {
		this.Data["checked"] = "checked"
	}
	this.Data["name"] = name

	this.TplName = "login.html"
}

//处理登陆界面
func (this *Usercontrollers) Login() {
	username := this.GetString("username")
	password := this.GetString("password")

	if username == "" || password == "" {
		beego.Error("用户名或者密码为空")
		this.Data["errmsg"] = "用户名或者密码为空"
		this.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = username
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名不存在")
		this.Data["errmsg"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}

	if password != user.Pwd {
		beego.Error("密码错误")
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}

	if user.Active == false {
		beego.Error("没有激活")
		this.Data["errmsg"] = "没有激活"
		this.TplName = "login.html"
		return
	}

	m1 := this.GetString("m1")
	//var buffer bytes.Buffer
	//enc := gob.NewEncoder(&buffer)
	//err = enc.Encode(&username)
	//if err != nil {
	//	beego.Error("加密失败")
	//	this.Redirect("/login", 302)
	//	return
	//}

	if m1 == "on" {
		this.Ctx.SetCookie("logname", username, 60*60)
	} else {
		this.Ctx.SetCookie("logname", username, -1)
	}

	//将登陆信息保存到seesion中
	this.SetSession("name", username)

	beego.Info("登陆成功")
	this.Redirect("/index", 302)
}

//退出登陆
func (this *Usercontrollers) Logout() {
	this.DelSession("name")
	this.Redirect("/login", 302)
}

//展示用户中心
func (this *Usercontrollers) ShowUserinfo() {
	//通过session查找
	name := this.GetSession("name")
	if name == nil {
		beego.Error("获取session name 错误")
		this.Data["errmsg"] = "获取session name 错误"
		this.Redirect("/index", 302)
		return
	}
	this.Data["name"] = name //给上方的名字赋值

	o := orm.NewOrm()
	var user models.User
	user.Name = name.(string)
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名不存在")
		this.Data["errmsg"] = "用户名不存在"
		this.Redirect("/index", 302)
		return
	}

	var addr models.Address
	err = o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name.(string)).Filter("Isdefault", true).One(&addr)
	if err == nil {
		this.Data["addr"] = addr
	} else {
		this.Data["addr"] = ""
	}

	//最近浏览
	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		beego.Error("redis连接失败")
		this.Redirect("/index", 302)
		return
	}
	defer conn.Close()

	goodsid, err := redis.Ints(conn.Do("lrange", "history"+name.(string), 0, 3))
	if err != nil {
		beego.Error("商品id获取失败")
		this.Redirect("/index", 302)
		return
	}

	var goodsskus []models.GoodsSKU

	for _, v := range goodsid {
		var goods models.GoodsSKU
		goods.Id = v
		o.Read(&goods)
		goodsskus = append(goodsskus, goods)
	}

	this.Data["goodsskus"] = goodsskus
	//给user  附加信息   地址 没有插入
	this.Data["user"] = user

	this.Layout = "user_center_layout.html"
	this.Data["contrl"] = 1
	this.TplName = "user_center_info.html"
}

//展示用户收货地质
func (this *Usercontrollers) ShowUsersite() {
	//展示默认地址
	//通过session获取用户名
	name := this.GetSession("name")

	o := orm.NewOrm()
	var addr models.Address //通过session获取到的用户名获取用户的默认地址
	qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name.(string))
	err := qs.Filter("Isdefault", true).One(&addr)
	if err == nil {
		this.Data["addr"] = addr
	} else {
		this.Data["addr"] = ""
	}

	this.Layout = "user_center_layout.html"
	this.Data["contrl"] = 3
	this.Data["name"] = name
	this.TplName = "user_center_site.html"
}

//处理用户收货地址
func (this *Usercontrollers) Usersite() {
	reciver := this.GetString("receiver")
	addr := this.GetString("addr")
	postCode := this.GetString("postCode")
	phone := this.GetString("phone")

	if reciver == "" || addr == "" || postCode == "" || phone == "" {
		beego.Error("收货地址不完整")
		this.Data["errmsg"] = "收货地址不完整"
		this.TplName = "user_center_site.html"
		return
	}

	o := orm.NewOrm()
	var useraddr models.Address
	useraddr.Receiver = reciver
	useraddr.Addr = addr
	useraddr.PostCode = postCode
	useraddr.Phone = phone

	name := this.GetSession("name")
	var user models.User
	user.Name = name.(string)
	o.Read(&user, "Name")
	//这里意思是  那个用户的地址
	useraddr.User = &user
	//查询看 有没有默认地址 如果有 把默认地址修改为非默认
	// 如果没有 直接插入默认地址
	//查询当前用户是否有默认地址
	var oldaddr models.Address
	qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name.(string))
	err := qs.Filter("Isdefault", true).One(&oldaddr)
	//err 为空说明找到了  不为空说明没有找到
	//if err != nil {
	//	useraddr.Isdefault = true //如果找到了就把找到的isdefault赋值为false 再更新  在把要插入的地址设置为默认地址
	//} else { //如果没有找到 就直接把 要插入的地址设置为默认地址
	//	oldaddr.Isdefault = false
	//	o.Update(&oldaddr, "Isdefault")
	//	useraddr.Isdefault = true
	//}
	if err == nil {
		oldaddr.Isdefault = false
		_, err = o.Update(&oldaddr, "Isdefault")
		if err != nil {
			beego.Error("更新失败")
			this.TplName = "user_center_site.html"
			return
		}
	}
	useraddr.Isdefault = true

	_, err = o.Insert(&useraddr)
	if err != nil {
		beego.Error("插入失败", err)
		this.TplName = "user_center_site.html"
		return
	}

	this.Redirect("/user/usersite", 302)

}

//展示用户订单
func (this *Usercontrollers) ShowUserorder() {
	//通过session查找
	name := this.GetSession("name")
	if name == nil {
		beego.Error("获取session name 错误")
		this.Data["errmsg"] = "获取session name 错误"
		this.Redirect("/index", 302)
		return
	}
	//展示订单的信息
	o := orm.NewOrm()
	var OrderInfo []models.OrderInfo
	qs := o.QueryTable("OrderInfo").RelatedSel("User").Filter("User__Name", name.(string))

	count, _ := qs.Count()
	pagesize := 1
	pagecount := int(math.Ceil(float64(count) / float64(pagesize)))
	pageindex, _ := this.GetInt("pageindex")
	if pageindex == 0 {
		pageindex = 1
	}

	pages := PageMode(pagecount, pageindex)
	this.Data["pages"] = pages
	this.Data["pageindex"] = pageindex

	var prepage, nextpage int
	if pageindex-1 <= 0 {
		prepage = 1
	} else {
		prepage = pageindex - 1
	}
	if pageindex+1 > pagecount {
		nextpage = pagecount
	} else {
		nextpage = pageindex + 1
	}

	this.Data["prepage"] = prepage
	this.Data["nextpage"] = nextpage

	qs.Limit(pagesize, pagesize*(pageindex-1)).All(&OrderInfo) //截取

	//定义大容器
	var orderinfos []map[string]interface{}

	for _, v := range OrderInfo {
		//定义行容器
		t := make(map[string]interface{})
		var OrderGoods []models.OrderGoods
		o.QueryTable("OrderGoods").RelatedSel("OrderInfo", "GoodsSKU").Filter("OrderInfo__Id", v.Id).All(&OrderGoods)

		t["OrderInfo"] = v
		t["OrderGoods"] = OrderGoods
		orderinfos = append(orderinfos, t)
	}

	this.Data["orderinfos"] = orderinfos
	this.Data["name"] = name.(string) //给上方的名字赋值
	this.Layout = "user_center_layout.html"
	this.Data["contrl"] = 2
	this.TplName = "user_center_order.html"
}

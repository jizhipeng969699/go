package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"math/rand"
	"pinyg/models"
	"regexp"
	"strings"
	"time"
)

type UserController struct {
	beego.Controller
}

//显示登陆界面
func (this *UserController) Showregister() {
	this.TplName = "register.html"
}

//发送短信
func (this *UserController) Sendmsg() {
	////接受数据
	//phone := this.GetString("phone")
	////返回json格式数据
	//
	////校验数据
	//if phone == "" {
	//	beego.Error("获取电话失败")
	//	//定义一个给ajax传递数据的容器
	//	//接收的时json格式 返回的也是json格式
	//	resp := make(map[string]interface{})
	//	//给容器赋值
	//	resp["errno"] = 1         //错误玛
	//	resp["errmsg"] = "获取电话失败" //错误信息
	//	//把容器传递给前段
	//	this.Data["json"] = resp
	//	this.ServeJSON()
	//}
	//beego.Info("------ok")
	resp := make(map[string]interface{})
	defer Senderr(this, resp)

	phone := this.GetString("phone")
	if phone == "" {
		beego.Error("获取电话失败")

		//resp := make(map[string]interface{})
		resp["errno"] = 1
		resp["errmsg"] = "获取电话失败"

		//this.Data["json"] = resp
		//this.ServeJSON()
		return
	}
	reg, _ := regexp.Compile(`^1[3-9][0-9]{9}$`)
	result := reg.FindString(phone) //通过正则找 找到了返回首字母 没有找到返回空
	//bll := reg.MatchString(phone)   符合返回true  不符合返回false
	if result == "" {
		beego.Error("电话格式错误")
		resp["errno"] = 2
		resp["errmsg"] = "电话格式错误"
		return
	}

	//初始化客户端  需要accessKey  需要开通申请
	client, err := sdk.NewClientWithAccessKey("default", "LTAIULAvXMZQS2XZ", "uLaaURM2cvD9L2m6Lie9XmeaYvdaKx")
	if err != nil {
		resp["errno"] = 3
		resp["errmsg"] = "阿里云客户端初始化失败"
		return
	}
	//获取6位数随机码
	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	//vcode := fmt.Sprintf("%06d", rnd.Int31n(1000000))
	rand.Seed(time.Now().UnixNano())
	vcode := fmt.Sprint(rand.Intn(899999) + 100000)

	//初始化请求对象
	request := requests.NewCommonRequest()
	request.Method = "POST"                                         //设置请求方法
	request.Scheme = "https"                                        // https | http   //设置请求协议
	request.Domain = "dysmsapi.aliyuncs.com"                        //域名
	request.Version = "2017-05-25"                                  //版本号
	request.ApiName = "SendSms"                                     //api名称
	request.QueryParams["PhoneNumbers"] = phone                     //需要发送的电话号码
	request.QueryParams["SignName"] = "品优购"                         //签名名称   需要申请
	request.QueryParams["TemplateCode"] = "SMS_164275022"           //模板号   需要申请
	request.QueryParams["TemplateParam"] = `{"code":` + vcode + `}` //发送短信验证码

	response, err := client.ProcessCommonRequest(request) //发送短信
	if err != nil {
		resp["errno"] = 4
		resp["errmsg"] = "发送短信失败"
		return
	}

	var msg MSG
	//
	//json.Unmarshal(response.GetHttpContentBytes(), &msg) //解析发送结果
	//if msg.Message != "OK" {
	//	resp["errno"] = 6
	//	resp["errmsg"] = "短信发送失败"
	//	return
	//}

	//json 数据格式解析//解析阿里云发短信返回的数据
	json.Unmarshal(response.GetHttpContentBytes(), &msg)
	if msg.Message != "OK" {
		beego.Error("短信发送失败")
		resp["errno"] = 6
		resp["errmsg"] = msg.Message
		return
	}

	resp["errno"] = 5
	resp["errmsg"] = "短信发送成功"
	resp["vcode"] = vcode
}

//解析阿里云发短信返回的数据
type MSG struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}

//发短信    发送函数
func Senderr(this *UserController, resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//处理注册界面
func (this *UserController) Register() {
	phone := this.GetString("phone")
	password := this.GetString("password")
	repassword := this.GetString("repassword")
	//code := this.GetString("code")   验证马  没有检验

	if phone == "" || password == "" || repassword == "" {
		beego.Error("用户名或者密码为空")
		this.Data["errmsg"] = "用户名或者密码为空"
		this.TplName = "register.html"
		return
	}
	if password != repassword {
		beego.Error("两次密码不一致")
		this.Data["errmsg"] = "两次密码不一致"
		this.TplName = "register.html"
		return
	}

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

	beego.Info("成功")
	this.Ctx.SetCookie("userName", user.Name, 60*60)
	this.Redirect("/register-email", 302)
}

//展示邮箱界面
func (this *UserController) Showremail() {
	this.TplName = "register-email.html"
}

//处理注册邮箱业务
func (this *UserController) Email() {
	email := this.GetString("email")
	password := this.GetString("password")
	repassword := this.GetString("repassword")

	if email == "" || password == "" || repassword == "" {
		beego.Error("邮箱信息不全")
		this.Data["errmsg"] = "邮箱信息不全"
		this.TplName = "register-email.html"
		return
	}
	if password != repassword {
		beego.Error("两次密码不一致")
		this.Data["errmsg"] = "两次密码不一致"
		this.TplName = "register-email.html"
		return
	}

	newemail := strings.ToLower(email)
	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(newemail)
	if result == "" {
		beego.Error("邮箱格式错误")
		this.Data["errmsg"] = "邮箱格式错误"
		this.TplName = "register-email.html"
		return
	}

	//utils   全剧通用接口 工具类    stmp
	//配置
	config := `{"username":"2297134578@qq.com","password":"xjvmmqpyyfuydifb","host":"smtp.qq.com","port":587}`
	eml := utils.NewEMail(config) //创建email对像
	//配置邮箱的内容
	eml.Subject = "pingyg邮箱激活"     //邮箱的标题
	eml.From = "2297134578@qq.com" //邮箱的发起地址  发件人地址
	eml.To = []string{email}       //发给谁  是个切片 可以多发 群法

	userName := this.Ctx.GetCookie("userName")
	eml.HTML = `<a href="http://192.168.0.104:8080/active?userName=` + userName + `">点击激活该用户</a>` //href 中的url整体必须要“” 阔起来
	err := eml.Send()                                                                             //发送
	if err != nil {
		beego.Error("邮箱发送失败")
		this.Data["errmsg"] = "邮箱发送失败"
		this.TplName = "register-email.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err = o.Read(&user, "Name")
	if err != nil {
		beego.Error("没有找到用户名")
		this.Redirect("/register-email", 302)
		return
	}

	user.Email = email
	_, err = o.Update(&user)
	if err != nil {
		beego.Error("邮箱插入失败")
		this.Redirect("/register-email", 302)
		return
	}

	this.Ctx.WriteString("邮箱发送成功")
}

//激活用户
func (this *UserController) Active() {
	userName := this.GetString("userName")
	if userName == "" {
		beego.Error("用户名不存在")
		this.Redirect("/register-email", 302)
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("没有找到用户名")
		this.Redirect("/register-email", 302)
		return
	}

	user.Active = true
	_, err = o.Update(&user)
	if err != nil {
		beego.Error("激活失败")
		this.Redirect("/register-email", 302)
		return
	}

	this.Redirect("/login", 302)
}

//显示登陆界面
func (this *UserController) Showlogin() {
	name := this.Ctx.GetCookie("logname")
	if name == "" {
		this.Data["checked"] = ""
	} else {
		this.Data["checked"] = "checked"
	}

	this.Data["name"] = name
	this.TplName = "login.html"
}

//处理登陆界面
func (this *UserController) Login() {
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

	if user.Pwd != password {
		beego.Error("密码错误")
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}

	if user.Active != true {
		beego.Error("没有激活")
		this.Data["errmsg"] = "没有激活"
		this.TplName = "login.html"
		return
	}

	beego.Info("登陆成功")

	m1 := this.GetString("m1")
	if m1 == "on" {
		this.Ctx.SetCookie("logname", user.Name, 60*60)
		this.Data["checked"] = "checked"
	} else {
		this.Ctx.SetCookie("logname", user.Name, -1)
	}
	//beego.Info("1111111111111111111111")
	this.SetSession("username", user.Name)
	//beego.Info("22222222222222222222222222222")
	this.Redirect("/index", 302)
	//beego.Info("33333333333333333")
}

//用户退出
func (this *UserController) Logout() {
	this.DelSession("username")
	this.Redirect("/index",302)
}

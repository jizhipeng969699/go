package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"pyg/models"
)

type Cartcontrollers struct {
	beego.Controller
}

//将订单信息放到redis中去
func (this *Cartcontrollers) Cart() {

	//获取数据
	goodsId, err := this.GetInt("goodsId") //商品的id
	num, err1 := this.GetInt("num")        //商品的数量

	//返回ajax步骤
	//定义一个map容器
	resp := make(map[string]interface{})

	//面向对象的 封装 继承 多态
	defer Respfunc(&this.Controller, resp)

	if err != nil || err1 != nil {
		resp["errno"] = 1
		resp["errmsg"] = "输入的数据不完整"
		return
	}

	//校验登陆状态
	name := this.GetSession("name")
	if name == nil {
		resp["errno"] = 2
		resp["errmsg"] = "当前用户没有登陆，不能添加购物车"
		return
	}

	//go操作redis

	//conn, err := redis.Dial("tcp", ":6379")
	beego.Info(goodsId, num, "----------------")

	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		resp["errno"] = 3
		resp["errmsg"] = "redis连接失败"
		return
	}

	defer conn.Close()

	//存储到redis
	reply, err := conn.Do("hset", "cart_"+name.(string), goodsId, num)
	if err != nil {
		resp["errno"] = 4
		resp["errmsg"] = "存储失败"
		return
	}
	beego.Info(reply)
	resp["errno"] = 5
	resp["errmsg"] = "ok"

}

//显示购物车
func (this *Cartcontrollers) Showcart() {
	//从redis中取数据
	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		beego.Error("连接redis失败")
		this.Redirect("/index_sx", 302)
		return
	}
	defer conn.Close()
	//从session中获取用户名
	name := this.GetSession("name").(string)
	//hgetall  获取redis中所有的key 中的 所有key和value
	reply, err := redis.Ints(conn.Do("hgetall", "cart_"+name))
	if err != nil {
		beego.Error("查询失败")
		this.Redirect("/index_sx", 302)
		return
	}
	//定义大容器
	var goods []map[string]interface{}
	//定义总价
	totlprice := 0
	//定义总剑术
	goodscount := 0
	o := orm.NewOrm()
	for i := 0; i < len(reply); i += 2 {
		//reply[i]时商品id
		//replay[i+1]是商品数量

		//定义但行容器
		t := make(map[string]interface{})
		var GoodsSKU models.GoodsSKU
		GoodsSKU.Id = reply[i]
		o.Read(&GoodsSKU)

		littleprice := reply[i+1] * GoodsSKU.Price
		goodscount += reply[i+1]
		totlprice += littleprice

		t["GoodsSKU"] = GoodsSKU
		t["count"] = reply[i+1]
		t["littleprice"] = littleprice
		//将行容器放到大容器中
		goods = append(goods, t)
	}

	this.Data["totlprice"] = totlprice

	this.Data["goodscount"] = goodscount
	this.Data["goods"] = goods

	this.TplName = "cart.html"
}

//更新count 加
func (this *Cartcontrollers) UpCart() {
	//获取ajax传来的 json值
	count, err := this.GetInt("count")
	goodsId, err1 := this.GetInt("goodsId")

	//定义容器 用来返回ajax数据
	recep := make(map[string]interface{})
	defer Respfunc(&this.Controller, recep) //555555555555

	if err != nil || err1 != nil {
		beego.Info(err, "-----------11--------", err1)
		recep["errno"] = 1
		recep["errmsg"] = "获取数据失败"
		return
	}

	//在redis中村的是hash类型  是key key val 类型
	name := this.GetSession("name").(string) //获取hash的第一个key
	if name == "" {
		beego.Info("-------22---------", name)
		recep["errno"] = 2
		recep["errmsg"] = "获取用户名失败"
		return
	}
	//处理数据
	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		beego.Info(err, "---------33----------")
		recep["errno"] = 3
		recep["errmsg"] = "redis连接失败"
		return
	}
	defer conn.Close()

	_, err = conn.Do("hset", "cart_"+name, goodsId, count)
	if err != nil {
		beego.Info(err, "-------444------")
		recep["errno"] = 4
		recep["errmsg"] = "redis更新失败"
		return
	}
	recep["errno"] = 5
	recep["errmsg"] = "OK"
}

//更新count减
func (this *Cartcontrollers) Minus() {
	count, err := this.GetInt("count")
	goodsId, err1 := this.GetInt("goodsId")

	//创建大容器
	recep := make(map[string]interface{})

	defer Respfunc(&this.Controller, recep) //面向对象

	if err != nil || err1 != nil {
		recep["errno"] = 1
		recep["errmsg"] = "获取数据失败"
		return
	}

	//处理数据
	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		beego.Error("redis连接失败")
		recep["errno"] = 2
		recep["errmsg"] = "redis连接失败"
		return
	}
	defer conn.Close()

	name := this.GetSession("name").(string)
	if name == "" {
		beego.Error("获取用户名失败")
		recep["errno"] = 3
		recep["errmsg"] = "获取用户名失败"
		return
	}

	_, err = conn.Do("hset", "cart_"+name, goodsId, count)
	if err != nil {
		beego.Error("更新失败")
		recep["errno"] = 4
		recep["errmsg"] = "更新失败"
		return
	}

	recep["errno"] = 5
	recep["errmsg"] = "OK"
}

//删除redis中的数据
func (this *Cartcontrollers) Delete() {

	//定义返回给ajax的json数据
	recp := make(map[string]interface{})
	//指定是以json个数传输  面向对象调用
	defer Respfunc(&this.Controller, recp)
	//获取数据
	goodsId, err := this.GetInt("goodsId")
	if err != nil {
		beego.Info(err, "dddddddddddddd", goodsId)
		recp["errno"] = 1
		recp["errmsg"] = "获取数据失败"
		return
	}

	//校验登陆状态
	name := this.GetSession("name")
	if name == nil {
		recp["errno"] = 2
		recp["errmsg"] = "用户没有登陆"
		return
	}

	//处理数据 连接redis
	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		recp["errno"] = 3
		recp["errmsg"] = "redis连接失败"
		return
	}
	defer conn.Close()

	_, err = conn.Do("hdel", "cart_"+name.(string), goodsId)
	if err != nil {
		recp["errno"] = 4
		recp["errmsg"] = "redis数据删除失败"
		return
	}

	recp["errno"] = 5
	recp["errmsg"] = "OK"

}

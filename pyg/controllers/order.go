package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"github.com/smartwalle/alipay"
	"pyg/models"
	"strconv"
	"strings"
	"time"
)

type Ordercontrollers struct {
	beego.Controller
}

//展示订单界面
func (this *Ordercontrollers) Showorder() {
	//只有getstrings 接受的是切片
	//提交切片或者多个值  是 用标签的value提交的   获取的花是 标签的name
	goodsIds := this.GetStrings("checkGoods")

	o := orm.NewOrm()
	//显示所有地址
	name := this.GetSession("name")
	var addrs []models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name.(string)).All(&addrs)

	this.Data["addrs"] = addrs

	conn, _ := redis.Dial("tcp", "192.168.11.135:6379")

	//获取商品,获取总价和总件数
	var goods []map[string]interface{}
	var totalPrice, totalCount int

	for _, v := range goodsIds {
		temp := make(map[string]interface{})
		id, _ := strconv.Atoi(v)
		var goodsSku models.GoodsSKU
		goodsSku.Id = id
		o.Read(&goodsSku)

		//获取商品数量
		count, _ := redis.Int(conn.Do("hget", "cart_"+name.(string), id))

		//计算小计
		littlePrice := count * goodsSku.Price

		//把商品信息放到行容器
		temp["goodsSku"] = goodsSku
		temp["count"] = count
		temp["littlePrice"] = littlePrice

		totalPrice += littlePrice
		totalCount += count

		goods = append(goods, temp)

	}

	//返回数据
	this.Data["totalPrice"] = totalPrice
	this.Data["totalCount"] = totalCount
	this.Data["truePrice"] = totalPrice + 10
	this.Data["goods"] = goods
	this.Data["goodsIds"] = goodsIds
	this.TplName = "place_order.html"

}

//处理订单  添加订单
func (this *Ordercontrollers) Pushorder() {
	//获取数据
	addrId, err := this.GetInt("addrId")

	payId, err1 := this.GetInt("payId")
	//虽然是个切片 但是js从页面获取的数据都是string类型
	goodsIds := this.GetString("goodsIds")

	totalCount, err2 := this.GetInt("totalCount")

	totalPrice, err3 := this.GetInt("totalPrice")

	recp := make(map[string]interface{})

	defer Respfunc(&this.Controller, recp)

	if err != nil || err1 != nil || err2 != nil || err3 != nil || goodsIds == "" {
		recp["errno"] = 1
		recp["errmsg"] = "获取数据错误"
		return
	}
	//插入数据
	//通过session获取用户名 获得用户对象
	o := orm.NewOrm()
	var user models.User
	name := this.GetSession("name")
	user.Name = name.(string)
	o.Read(&user, "Name")
	//通过地址id获得地址对象
	var address models.Address
	address.Id = addrId
	o.Read(&address)
	//给订单信息表然后插入订单订单信息表
	var orderInfo models.OrderInfo
	orderInfo.OrderId = time.Now().Format("20060102150405" + strconv.Itoa(user.Id))
	orderInfo.User = &user
	orderInfo.Address = &address
	orderInfo.PayMethod = payId
	orderInfo.TotalCount = totalCount
	orderInfo.TotalPrice = totalPrice
	orderInfo.TransitPrice = 10

	//开启事务防止脏读 幻读
	o.Begin()

	//插入订单信息表
	_, err = o.Insert(&orderInfo)
	if err != nil {
		o.Rollback() //出错回滚
		recp["errno"] = 2
		recp["errmsg"] = "订单表插入失败"
		return
	}

	//插入订单商品表
	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		o.Rollback() //出错回滚
		recp["errno"] = 3
		recp["errmsg"] = "redis连接失败"
		return
	}
	defer conn.Close()

	//插入订单商品表
	//goodsIds  [2 3 5]
	goodslice := strings.Split(goodsIds[1:len(goodsIds)-1], " ")
	for _, v := range goodslice {
		goodsid, _ := strconv.Atoi(v)
		//if err != nil {
		//	beego.Info(err, "----------------")
		//}
		var ordergoods models.OrderGoods
		var goodssku models.GoodsSKU
		goodssku.Id = goodsid
		o.Read(&goodssku)

		oldstock := goodssku.Stock
		beego.Info("原始库存", oldstock)

		ordergoods.GoodsSKU = &goodssku
		ordergoods.OrderInfo = &orderInfo
		//count  在redis中所以要到redis中查找
		//获取商品的数量
		count, err := redis.Int(conn.Do("hget", "cart_"+name.(string), goodsid))
		if err != nil {
			o.Rollback() //出错回滚
			recp["errno"] = 4
			recp["errmsg"] = "获取商品数量失败"
			return
		}
		ordergoods.Count = count

		//计算小计
		ordergoods.Price = count * goodssku.Price

		//插入之前需要更新商品的库存和销量
		if goodssku.Stock < count {
			o.Rollback() //出错回滚
			recp["errno"] = 6
			recp["errmsg"] = "库存不足"
			return
		}
		//goodssku.Stock -= count
		//goodssku.Sales += count
		//_, err = o.Update(&goodssku, "Stock", "Sales")

		o.Read(&goodssku)

		qs := o.QueryTable("GoodsSKU").Filter("Id", goodsid).Filter("Stock", oldstock)
		_, err = qs.Update(orm.Params{"Stock": goodssku.Stock - count, "Sales": goodssku.Sales + count})
		if err != nil {
			o.Rollback() //出错回滚
			recp["errno"] = 7
			recp["errmsg"] = "库存和销量更新失败"
			return
		}
		beego.Info(goodssku.Stock, "-------------")

		_, err = o.Insert(&ordergoods)
		if err != nil {
			o.Rollback() //出错回滚
			recp["errno"] = 8
			recp["errmsg"] = "订单商品表插入失败"
			return
		}
	}
	o.Commit() //没有出错就提交
	recp["errno"] = 5
	recp["errmsg"] = "订单商品表插入失败"
}

//支付
func (this *Ordercontrollers) Pay() {
	orderId, err := this.GetInt("orderId")
	if err != nil {
		beego.Error("获取订单id失败")
		this.Redirect("/user/order", 302)
		return
	}

	o := orm.NewOrm()
	var orderInfo models.OrderInfo
	orderInfo.Id = orderId
	o.Read(&orderInfo)
	aliPublicKey := `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvDLKRBQjrDY2V5hBiELW
PTVt4/moRNKqYPxxn5C4+7M64KZsO9Au/imXvlpo33gzk+tg/VcU6gpWTPzioz26
BxSrJb8lUwL+PYmWzCyn6zIAiezXvc/a0YNYnZOPo2insAaKfQa1GLjHvrejOrjs
ev7oJyIEnJP6QNCAF0DZkpmlBabNleFiSvl/pNhuUylStyQsWu0a0lPZNHJAYg2F
VIqF6RcYgHL8U1VDhCr40hZ+s3Q8/7qVcNCRNQNtXfyHvJLYF9Pwf+T6s8fAbsKV
gA9Yy9fKzt7Ehgo9L38Xc5uEko22SWrOSdUpvq7DCPvlwJgx2HUf7BIHjM4rq4ue
KwIDAQAB`
	privateKey := `MIIEowIBAAKCAQEAvDLKRBQjrDY2V5hBiELWPTVt4/moRNKqYPxxn5C4+7M64KZs
O9Au/imXvlpo33gzk+tg/VcU6gpWTPzioz26BxSrJb8lUwL+PYmWzCyn6zIAiezX
vc/a0YNYnZOPo2insAaKfQa1GLjHvrejOrjsev7oJyIEnJP6QNCAF0DZkpmlBabN
leFiSvl/pNhuUylStyQsWu0a0lPZNHJAYg2FVIqF6RcYgHL8U1VDhCr40hZ+s3Q8
/7qVcNCRNQNtXfyHvJLYF9Pwf+T6s8fAbsKVgA9Yy9fKzt7Ehgo9L38Xc5uEko22
SWrOSdUpvq7DCPvlwJgx2HUf7BIHjM4rq4ueKwIDAQABAoIBADYVmo8qAn5xLfjK
ZrrUCmGYwzXq/3KX4CLzKWwj2SVsgpUg/qXJ1FdyeItZzB28m4X89OiZaAdCZT7G
xhMHoDh8thqo7f8HMy1YJyASDUfIHOOGTxHbdBUV3Bec2oCpiNdLae04Sj8Ki0zB
1YUAs+U88FaerhiZZifOQ6FrH5vDUcFT+BOIMz+XcXeForUUAh52BM+MbWf8Ybhm
Td+AixXAIKEz0+3I961KsZUJM2e7+T7KQUt3IUmf1/T2ydU2rMxTHgWsEOUWSXep
m4U7Q/eHVIXr81fpPq93ig9cGZNMSFtubDJCDn+1sq5mOege5Y1+nz7ccMmYtYmh
x79Ug1ECgYEA3WW7dxHUkZzt1tsHmLjwtC3ljWWsrloLVcBlda3GuoHvLtDtC/iO
4aPo5ZbP+RpPRfzg+HLGnFvm8rww0GToIPkHUOJtJLZsVgO7Iv5gHqRa71pF1rS7
MSXHpXTJtPzPPSBmGvj3rLpErf+1UVcaEo0SvL1Fb9SCxEVK5d/XbQcCgYEA2Zy+
DljEmfyZto7WNsBhBuR5pjHwTAFoGzdDpa5YNbWqZT1Y0ZlpwJ5DvaYX35mOm/tj
EV2Cd0N9+PyT4wdDZ2cXJR1lZuNjuwjLFVcQmCYeinvCDbZO34LdCVKdjtaEke+2
CIiqJJvSdK012EI5ENL1UCtBvk8o14qBbriW4L0CgYApt6qsVArG6VgTnTAAQhDx
EpjnnNn9/G+YV+mGVBlXOXaVTr7r+4kI0XboFBPaL2/ykyTdC9uKyPJqmLVl/y2S
UPpk7lV8jnc5efiALoh0HaoY4dy5CVMgfqrw7WG4nc57CSjUOkeJH4wjcUR8MDp0
dmqqb8uut89wJlJnoFcabwKBgFMQeZuQVrtQqHv+2dXcaSzwWV8PAOKCRvLcjX3Y
puMZjQlH+XdIUA2uW70wgTxgqQbxVkdyojJUGOnJv8mRJDF8MGGCbwpvEcp6+MoU
ickKA+5ofxIs3L6EfUrptiqnx8LM9Xccc5W1xQe0djEuVgoN/IW1fUrffH+J5w4U
d8MhAoGBAIzjroTmOBkCTb4DvgaY6TyrqC7xkVCGmtNvfTewtaf57p3cBiOP5nDE
V1WECaoslj+aeg75srB9w8A+bGiJ66bDSb1Lu6pUB+myAobGzYaf7mdw4ngIHRMc
1Jcw/FQUp2KLFp3WJ3Cx3NPtkjPgxEeXSB73UTqDn7S/zfFNmoWj`
	//支付
	//appId, aliPublicKey, privateKey string, isProduction bool
	//先定义一个 client对象          //支付宝的
	client := alipay.New("2016093000634324", aliPublicKey, privateKey, false)
	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://192.168.11.135:8080/user/order" //异步
	p.ReturnURL = "http://192.168.11.135:8080/user/order" //同步
	p.Subject = "品邮购"
	p.OutTradeNo = strconv.Itoa(orderInfo.Id)          //确定订单号
	p.TotalAmount = strconv.Itoa(orderInfo.TotalPrice) //确定钱数
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		beego.Error(err, "支付失败")
		this.Redirect("/user/order", 302)
		return
	}

	orderInfo.Orderstatus = 1
	n,err:=o.Update(&orderInfo,"Orderstatus")
	if err != nil {
		beego.Error(err, "支付失败")
		this.Redirect("/user/order", 302)
		return
	}
	beego.Info(n,"Update(&orderInfo,Orderstatus)")
	this.Redirect(url.String(), 302)

}

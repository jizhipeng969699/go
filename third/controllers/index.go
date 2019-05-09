package controllers

import (
	"bytes"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"third/models"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) ShowIndex() {
	username := this.GetSession("userName")
	if username == nil {
		this.Redirect("/login", 302)
		return
	}

	this.Data["username"] = username.(string)

	o := orm.NewOrm()
	qs := o.QueryTable("Txt")
	var txts []models.Txt
	var count int64
	typename := this.GetString("select")

	if typename == "全部" || typename == "" || typename == "aaa" { //count()  函数 返回的是 一个int64的整数  和一个错误
		count, _ = qs.RelatedSel("Txttype").Count() //宗记录书
	} else {
		count, _ = qs.RelatedSel("Txttype").Filter("Txttype__Typename", typename).Count()
	}

	txtinfo := 2 //没页现实几条记录

	pagecount := math.Ceil(float64(count) / float64(txtinfo)) //页数

	pagenum, err := this.GetInt("pagenum")
	if err != nil {
		pagenum = 1
	}

	//orm在进行查询时  是惰性查询  你不指定他不查询
	if typename == "" || typename == "全部" {
		qs.Limit(txtinfo, txtinfo*(pagenum-1)).RelatedSel("Txttype").All(&txts)
	} else {
		qs.Limit(txtinfo, txtinfo*(pagenum-1)).RelatedSel("Txttype").Filter("Txttype__Typename", typename).All(&txts)
	}
	var txttypes []models.Txttype

	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		beego.Info("连接数据库失败")
		return
	}
	defer conn.Close()

	rep, err := conn.Do("get", "newtypes")
	results, _ := redis.Bytes(rep, err)
	if len(results) == 0 {
		o.QueryTable("Txttype").All(&txttypes)

		var buffer bytes.Buffer        //设置容器  相当与缓冲区  把加密的数据放到这里来
		enc := gob.NewEncoder(&buffer) //建立加密器 参数是  放置加密数据的缓冲区
		enc.Encode(txttypes)

		conn.Do("set", "newtypes", buffer.Bytes())
		beego.Info(txttypes, "从mysql中获取数据")
	} else {
		dec := gob.NewDecoder(bytes.NewReader(results))
		dec.Decode(&txttypes) //给struct 类型 赋值时  要是用引用传递
		beego.Info(txttypes, "从redis中获取数据")

	}

	this.Data["txttypes"] = txttypes
	this.Data["Typename"] = typename

	this.Data["txts"] = txts
	this.Data["count"] = count
	this.Data["pagecount"] = pagecount

	if pagenum <= 1 {
		pagenum = 1
	}

	this.Data["pagenum"] = pagenum

	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["indexjs"] = "indexjs.html"

	this.TplName = "index.html"
}

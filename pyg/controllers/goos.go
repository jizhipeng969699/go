package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"pyg/models"
)

type Gooscontrollers struct {
	beego.Controller
}

//显示主页界面
func (this *Gooscontrollers) Showindex() {
	logname := this.GetSession("name")

	if logname != nil {
		this.Data["name"] = logname.(string)
	} else {
		this.Data["name"] = ""
	}

	//获取类型信息并传递给前端
	//获取一级菜单
	//o := orm.NewOrm()
	//var oneClass []models.TpshopCategory
	//o.QueryTable("TpshopCategory").Filter("Pid", 0).All(&oneClass)
	//
	//var types []map[string]interface{} //总容器	是一个map类型的切片
	////var types []map[interface{}]interface{}
	//for _, v := range oneClass {
	//	//行容器
	//	t := make(map[string]interface{}) //map中key方的时以及菜单  val放的是二级菜单
	//
	//	//var t []map[interface{}]interface{}  //行容器
	//
	//	var twoClass []models.TpshopCategory
	//	//var twoClass []map[string]interface{}
	//
	//	o.QueryTable("TpshopCategory").Filter("Pid", v.Id).All(&twoClass)
	//
	//	//var erji []map[interface{}]interface{}
	//	//for _, v := range twoClass {
	//	//
	//	//	t := make(map[string]interface{})
	//	//	var threeClass []models.TpshopCategory
	//	//	o.QueryTable("TpshopCategory").Filter("Pid", v.Id).All(&threeClass)
	//	//	t["t22"] = v
	//	//	t["t23"] = threeClass
	//	//	erji = append(erji, t)
	//	//}
	////	twoClass = append(twoClass, erji.(map[string]interface{}))
	//
	//	t["t1"] = v        //这时候的v  就是对应的以及标签
	//	t["t2"] = twoClass //对应的二级标签
	//	types = append(types, t)
	//}
	//
	//for _, v1 := range types {
	//	//循环获取二级菜单
	//	var erji []map[string]interface{} //定义二级容器
	//	for _, v2 := range v1["t2"].([]models.TpshopCategory) {
	//		t := make(map[string]interface{})
	//		var thirdClass []models.TpshopCategory
	//		//获取三级菜单
	//		o.QueryTable("TpshopCategory").Filter("Pid",v2.Id).All(&thirdClass)
	//		t["t22"] = v2  //二级菜单
	//		t["t23"] = thirdClass   //三级菜单
	//		erji = append(erji,t)
	//	}
	//	//把二级容器放到总容器中
	//	v1["t3"] = erji
	//}

	o := orm.NewOrm()
	var oneClass []models.TpshopCategory
	o.QueryTable("TpshopCategory").Filter("Pid", 0).All(&oneClass)

	var types []map[string]interface{}
	for _, v := range oneClass {
		t := make(map[string]interface{})

		var twoClass []models.TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid", v.Id).All(&twoClass)

		t["t1"] = v
		t["t2"] = twoClass
		types = append(types, t)
	}

	for _, v1 := range types {
		var erji []map[string]interface{}
		for _, v2 := range v1["t2"].([]models.TpshopCategory) {

			t := make(map[string]interface{})
			var three []models.TpshopCategory
			o.QueryTable("TpshopCategory").Filter("Pid", v2.Id).All(&three)

			t["t22"] = v2
			t["t23"] = three
			erji = append(erji, t)
		}
		v1["t3"] = erji
	}

	this.Data["types"] = types

	this.TplName = "index.html"
}

//商品生鲜模块
func (this *Gooscontrollers) Showindex_sx() {
	//获取生鲜主页内容
	//获取商品类型
	//获取商品所有类型
	o := orm.NewOrm()
	var goostypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goostypes)
	this.Data["goostypes"] = goostypes

	//获取轮播图
	var indexGoodsBanner []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&indexGoodsBanner)
	this.Data["indexGoodsBanner"] = indexGoodsBanner
	//获取促销商品
	var indexPromotionBanner []models.IndexPromotionBanner
	o.QueryTable("indexPromotionBanner").OrderBy("Index").All(&indexPromotionBanner)
	this.Data["indexPromotionBanner"] = indexPromotionBanner
	//获取首页商品展示

	//按照商品的分类展示  有文字商品  和 图片商品
	//总容器
	var goods []map[string]interface{}

	for _, v := range goostypes {
		var textgoods []models.IndexTypeGoodsBanner
		var imggoods []models.IndexTypeGoodsBanner
		qs := o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").Filter("GoodsType__Id", v.Id).OrderBy("Index")
		qs.Filter("DisplayType", 0).All(&textgoods) //文字 商品
		qs.Filter("DisplayType", 1).All(&imggoods)  //图片 商品

		//行容器
		t := make(map[string]interface{})
		t["goostypes"] = v
		t["textgoods"] = textgoods
		t["imggoods"] = imggoods
		goods = append(goods, t)

	}

	this.Data["goods"] = goods
	//this.Layout = "layout.html"
	this.TplName = "index_sx.html"
}

//展示商品详细
func (this *Gooscontrollers) Showgoodsdetail() {
	goodsskuid, err := this.GetInt("goodsskuid")
	if err != nil {
		beego.Error("没有获取到id")
		this.TplName = "index_sx.html"
		return
	}

	o := orm.NewOrm()
	var goodssku models.GoodsSKU
	//goodsdetail.Id = goodsid
	//err = o.Read(&goodsdetail)
	//if err != nil {
	//	beego.Error("获取商品详细失败")
	//	this.TplName = "index_sx.html"
	//	return
	//}
	conn, err := redis.Dial("tcp", "192.168.11.135:6379")
	if err != nil {
		beego.Error("redis连接失败")
		this.TplName = "index_sx.html"
		return
	}
	defer conn.Close()
	name := this.GetSession("name")
	if name == nil {
		beego.Error("用户没有登陆")
		this.TplName = "index_sx.html"
		return
	}
	_, err = conn.Do("lrem", "history"+name.(string), 0, goodsskuid) //去除重复
	beego.Info(err,"去除重复去除重复去除重复")
	_, err = conn.Do("lpush", "history"+name.(string), goodsskuid)
	if err != nil {
		beego.Error("redis存储商品id失败")
		this.TplName = "index_sx.html"
		return
	}
	//获取商品详情
	o.QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").Filter("Id", goodsskuid).One(&goodssku)

	//获取统一类型商品  新品推荐
	var newgoods []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name", goodssku.GoodsType.Name)
	qs.OrderBy("-Time").Limit(2, 0).All(&newgoods)

	this.Data["newgoods"] = newgoods

	this.Data["goodssku"] = goodssku

	this.TplName = "detail.html"
}

//展示商品类型  列表页
func (this *Gooscontrollers) Showgoodslist() {
	goostypesid, err := this.GetInt("goostypesid")
	if err != nil {
		beego.Error("没有获取到id")
		this.TplName = "index_sx.html"
		return
	}
	//显示商品列表
	o := orm.NewOrm()
	var goodsskus []models.GoodsSKU

	//按照默认 价格 销量排序
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goostypesid)
	sort := this.GetString("sort")

	count, _ := qs.Count()                                          //记录总数
	pagesize := 1                                                   //每页显示的记录数
	pagecount := int(math.Ceil(float64(count) / float64(pagesize))) //页数

	pageindex, err := this.GetInt("pageindex") //获取当前页数
	if err != nil {
		pageindex = 1
	}
	pages := PageMode(pagecount, pageindex)
	this.Data["pages"] = pages

	//上一页下一页
	var prepage, nextpage int
	if pageindex-1 <= 0 {
		prepage = 1
	} else {
		prepage = pageindex - 1
	}
	if pageindex+1 >= pagecount {
		nextpage = pagecount
	} else {
		nextpage = pageindex + 1
	}
	this.Data["prepage"] = prepage
	this.Data["nextpage"] = nextpage
	this.Data["pageindex"] = pageindex
	qs = qs.Limit(pagesize, pagesize*(pageindex-1)) //每页对应显示的记录

	if sort == "" {
		qs.All(&goodsskus)
	} else if sort == "price" {
		qs.OrderBy("-Price").All(&goodsskus)
	} else {
		qs.OrderBy("-Sales").All(&goodsskus)
	}

	this.Data["goodsskus"] = goodsskus
	this.Data["sort"] = sort

	//显示同类型 新品
	var newgoods []models.GoodsSKU
	qs.OrderBy("-Time").Limit(2, 0).All(&newgoods)
	this.Data["newgoods"] = newgoods

	this.Data["goostypesid"] = goostypesid
	//this.Layout = "layout.html"
	this.TplName = "list.html"
}

//页码分页函数
func PageMode(pagecont int, pageindex int) []int {
	//1.不足五页  有几页显示几页
	//2.大于五页 前三页   1 2 3 4 5     12345
	//3.后三页    10页       67 8  9  10
	//4.中间页码    6  10页     6-2   6-1  6  6+1  6+2
	var pages []int
	if pagecont <= 5 { //1.不足五页  有几页显示几页
		for i := 1; i <= pagecont; i++ {
			pages = append(pages, i)
		}
	} else if pageindex <= 3 { //2.大于五页 前三页   1 2 3 4 5     12345
		for i := 1; i <= 5; i++ {
			pages = append(pages, i)
		}
	} else if pageindex >= pagecont-2 { //3.后三页    10页       67 8  9  10
		for i := pagecont - 4; i <= pagecont; i++ {
			pages = append(pages, i)
		}
	} else { //4.中间页码    6  10页     6-2   6-1  6  6+1  6+2
		for i := pageindex - 2; i <= pageindex+2; i++ {
			pages = append(pages, i)
		}
	}
	return pages
}

//模糊查找
func (this *Gooscontrollers) Showserch() {
	serchcontent := this.GetString("serchcontent") //获取要查找的 东西

	o := orm.NewOrm()
	var goodsskus []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").Filter("Name__icontains", serchcontent)

	count, _ := qs.Count()                                          //得到的总记录数
	pagesize := 1                                                   //每页显示的记录数
	pagecount := int(math.Ceil(float64(count) / float64(pagesize))) //页数

	pageindex, err := this.GetInt("pageindex") //获取当前页数
	if err != nil {
		pageindex = 1
	}
	pages := PageMode(pagecount, pageindex)
	this.Data["pages"] = pages

	var prepage, nextpage int
	if pageindex-1 <= 0 {
		prepage = 1
	} else {
		prepage = pageindex - 1
	}
	if pageindex+1 >= pagecount {
		nextpage = pagecount
	} else {
		nextpage = pageindex + 1
	}

	//var prepage, nextpage int
	//if pageindex-1 <= 0 {
	//	prepage = 1
	//} else {
	//	prepage = pageindex - 1
	//}
	//if pageindex+1 >= pagecount {
	//	nextpage = pagecount
	//} else {
	//	nextpage = pageindex + 1
	//}

	this.Data["prepage"] = prepage
	this.Data["nextpage"] = nextpage

	qs = qs.Limit(pagesize, pagesize*(pageindex-1))

	qs.All(&goodsskus)
	this.Data["serchcontent"] = serchcontent
	this.Data["goodsskus"] = goodsskus

	this.TplName = "serch.html"
}

//func (this *Gooscontrollers)Showlayout(){
//	this.TplName="layout.html"
//}

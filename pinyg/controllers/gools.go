package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"pinyg/models"
)

type GoolsController struct {
	beego.Controller
}

//显示主界面
func (this *GoolsController) Showindex() {
	name := this.GetSession("username")
	if name == nil {
		this.Data["name"] = ""
	} else {
		this.Data["name"] = name.(string)
	}
	//一级菜单
	o := orm.NewOrm()
	var oneClass []models.TpshopCategory
	o.QueryTable("TpshopCategory").Filter("Pid", 0).All(&oneClass)

	//定义总容器
	var types []map[string]interface{}

	for _, v := range oneClass {
		//定义行容器
		t := make(map[string]interface{})
		//获取二级菜单
		var twoClass []models.TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid", v.Id).All(&twoClass)

		t["t1"] = v
		t["t2"] = twoClass
		types = append(types, t)
	}

	for _, v1 := range types {
		//定义二级容器
		var erji []map[string]interface{}
		for _, v2 := range v1["t2"].([]models.TpshopCategory) {
			//定义三级容器
			tt := make(map[string]interface{})
			//获取三级菜单
			var threeClass []models.TpshopCategory
			o.QueryTable("TpshopCategory").Filter("Pid", v2.Id).All(&threeClass)

			tt["t21"] = v2
			tt["t22"] = threeClass
			erji = append(erji, tt)
		}
		v1["t3"] = erji
	}

	this.Data["types"] = types
	this.TplName = "index.html"
}

//展示生鲜模块
func (this *GoolsController) ShowIndex_sx() {
	//实现类型下拉框
	o := orm.NewOrm()
	var GoodsType []models.GoodsType
	o.QueryTable("GoodsType").All(&GoodsType)
	this.Data["GoodsType"] = GoodsType

	//实现幻灯片轮播
	var IndexGoodsBanner []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&IndexGoodsBanner)
	this.Data["IndexGoodsBanner"] = IndexGoodsBanner

	//实现促销产品展示
	var IndexPromotionBanner []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&IndexPromotionBanner)
	this.Data["IndexPromotionBanner"] = IndexPromotionBanner

	//主页内容列
	//总容器
	var goods []map[string]interface{}

	for _, v := range GoodsType {
		//行容器
		t := make(map[string]interface{})

		var testTypeGoodsBanner []models.IndexTypeGoodsBanner
		var imgTypeGoodsBanner []models.IndexTypeGoodsBanner
		qs := o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").Filter("GoodsType__Id", v.Id).OrderBy("Index")
		//获取文字商品
		qs.Filter("DisplayType", 0).All(&testTypeGoodsBanner)
		//获取图片商品
		qs.Filter("DisplayType", 1).All(&imgTypeGoodsBanner)

		t["GoodsType"] = v
		t["testTypeGoodsBanner"] = testTypeGoodsBanner
		t["imgTypeGoodsBanner"] = imgTypeGoodsBanner
		goods = append(goods, t)
	}

	this.Data["goods"] = goods
	this.TplName = "index_sx.html"
}

//商品详细信息页
func (this *GoolsController) Showgoodsdetail() {
	goodsSKU, err := this.GetInt("GoodsSKU")
	if err != nil {
		beego.Error("获取商品sku id 失败")
		this.Redirect("/user/index_sx", 302)
		return
	}

	//给前端赋值 商品详情
	o := orm.NewOrm()
	var GoodsSKU models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("Goods")
	qs.Filter("Id", goodsSKU).One(&GoodsSKU)

	//商品新品推荐
	var newgoods []models.GoodsSKU
	qs.OrderBy("Time").Limit(2, 0).All(&newgoods)

	this.Data["newgoods"] = newgoods
	this.Data["GoodsSKU"] = GoodsSKU
	this.TplName = "detail.html"
}

//商品列表页
func (this *GoolsController) Showgoodslist() {
	GoodsTypeid, err := this.GetInt("GoodsTypeid")
	if err != nil {
		beego.Error("获取类型id失败")
		this.Redirect("/user/index_sx", 302)
		return
	}

	o := orm.NewOrm()
	var GoodsSKU []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", GoodsTypeid)

	count, _ := qs.Count()                                          //总记录数
	pagesize := 1                                                   //每页显示的记录数
	pagecount := int(math.Ceil(float64(count) / float64(pagesize))) //总页数
	pageindex, err := this.GetInt("pageindex")
	if err != nil {
		pageindex = 1
	}
	pages := fnPage(pagecount, pageindex)
	this.Data["pages"] = pages

	//实现按照指定的顺序排序
	sort := this.GetString("sort")
	if sort == "" {
		qs.All(&GoodsSKU)
	} else if sort == "price" {
		qs.OrderBy("-Price").All(&GoodsSKU)
	} else {
		qs.OrderBy("Sales").All(&GoodsSKU)
	}

	this.Data["GoodsSKU"] = GoodsSKU
	this.Data["GoodsTypeid"] = GoodsTypeid
	this.Data["sort"] = sort
	this.TplName = "list.html"
}

//商品列表的分页函数
func fnPage(pagecount int, pageindex int) []int {

	var pages []int
	if pagecount <= 5 {
		for i := 1; i <= pagecount; i++ {
			pages = append(pages, i)
		}
	} else if pageindex <= 3 {
		for i := 1; i < 5; i++ {
			pages = append(pages, i)
		}
	} else if pageindex >= pagecount-2 {
		for i := pagecount - 4; i <= pagecount; i++ {
			pages = append(pages, i)
		}
	} else {
		for i := pageindex - 2; i <= pageindex+2; i++ {
			pages = append(pages, i)
		}
	}

	return pages
}

//控制用户中心左上角头
func Info(this *GoolsController) {
	username := this.GetSession("username")
	if username == nil {
		this.Data["info"] = ""
	} else {
		this.Data["info"] = username.(string)
	}
}

//用户中心信息
func (this *GoolsController) ShowCenterinfo() {

	//给用户 中心 附上地址
	username := this.GetSession("username").(string)
	o := orm.NewOrm()
	var user models.User
	user.Name = username
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户没找到")
		this.Data["errmsg"] = "用户没找到"
		//this.Redirect("/user/user_center_site",302)
		this.Layout = "layout.html"
		this.Data["bb"] = 3
		Info(this)
		this.TplName = "user_center_order.html"
		return
	}

	var addr models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name", username).Filter("IsDefault", true).One(&addr)

	//手机号加密
	qian := addr.Phone[:3]
	hou := addr.Phone[7:]
	addr.Phone = qian + "****" + hou

	this.Data["name"] = username
	this.Data["addr"] = addr

	this.Layout = "layout.html"
	this.Data["bb"] = 1
	Info(this)
	this.TplName = "user_center_info.html"
}

//用户order
func (this *GoolsController) ShowCenterorder() {
	this.Layout = "layout.html"
	this.Data["bb"] = 2
	Info(this)
	this.TplName = "user_center_order.html"
}

//用户site 地址
func (this *GoolsController) ShowCentersite() {
	name := this.GetSession("username").(string)
	//展示默认地址
	o := orm.NewOrm()
	var addr models.Address
	//one 返回的错误 找到了 err 为空  传到前端的数据是有数据的    没有找到 err 有东西  传到前端的东西为空 nil
	err := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name).Filter("IsDefault", true).One(&addr)

	//手机号加密
	qian := addr.Phone[:3]
	hou := addr.Phone[7:]
	addr.Phone = qian + "****" + hou

	if err != nil {
		this.Data["addr"] = ""
	} else {
		this.Data["addr"] = addr
	}

	this.Layout = "layout.html"
	this.Data["bb"] = 3
	Info(this)
	this.TplName = "user_center_site.html"
}

//处理用户地址
func (this *GoolsController) CenterSite() {
	receiver := this.GetString("Receiver")
	addr := this.GetString("Addr")
	postCode := this.GetString("PostCode")
	phone := this.GetString("Phone")

	if receiver == "" || addr == "" || postCode == "" || phone == "" {
		beego.Error("地址信息不完整")
		this.Data["errmsg"] = "地址信息不完整"
		//this.Redirect("/user/user_center_site",302)
		this.Layout = "layout.html"
		this.Data["bb"] = 3
		Info(this)
		this.TplName = "user_center_site.html"
		return
	}

	o := orm.NewOrm()

	username := this.GetSession("username").(string)
	var user models.User
	user.Name = username
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名没找到")
		this.Data["errmsg"] = "用户名没找到"
		//this.Redirect("/user/user_center_site",302)
		this.Layout = "layout.html"
		this.Data["bb"] = 3
		Info(this)
		this.TplName = "user_center_site.html"
		return
	}
	var addrr models.Address
	addrr.Receiver = receiver
	addrr.Addr = addr
	addrr.PostCode = postCode
	addrr.Phone = phone
	addrr.User = &user
	var oldaddr models.Address
	err = o.QueryTable("Address").RelatedSel("User").Filter("User__Name", username).Filter("IsDefault", true).One(&oldaddr)
	//if err != nil {
	//	addrr.IsDefault = true
	//}else{
	//	oldaddr.IsDefault = false
	//	o.Update(&oldaddr,"IsDefault")
	//	addrr.IsDefault = true
	//}
	if err == nil {
		oldaddr.IsDefault = false
		o.Update(&oldaddr, "IsDefault")
	}
	addrr.IsDefault = true
	_, err = o.Insert(&addrr)
	if err != nil {
		beego.Error("插入失败")
		this.Data["errmsg"] = "插入失败"
		//this.Redirect("/user/user_center_site",302)
		this.Layout = "layout.html"
		this.Data["bb"] = 3
		Info(this)
		this.TplName = "user_center_site.html"
		return
	}

	this.Redirect("/user/user_center_site", 302)
}

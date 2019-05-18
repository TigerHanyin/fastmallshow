package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pyg/pyg/models"
	"math"
)

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController) ShowIndex() {
	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)

	} else {
		this.Data["name"] = ""
	}
	//一级标签
	o := orm.NewOrm()
	var oneclass []models.TpshopCategory
	o.QueryTable("TpshopCategory").Filter("Pid", 0).All(&oneclass)
	//二级标签
	//var types []map[string]interface{}
	//for _,v := range oneClass{
	//	//行容器
	//	t := make(map[string]interface{})
	//
	//	var temp []models.TpshopCategory
	//	o.QueryTable("TpshopCategory").Filter("Pid",v.Id).All(&temp)
	//	t["t1"] = v
	//	t["t2"] = temp
	//	types = append(types,t)
	//}

	var types []map[string]interface{}
	for _, v := range oneclass {
		t := make(map[string]interface{})
		var temp []models.TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid", v.Id).All(&temp)
		t["t1"] = v
		t["t2"] = temp
		types = append(types, t)

	}
	//获取第三季菜单
	for _, v1 := range types {
		//循环获取二级菜单
		var erji []map[string]interface{} //定义二级容器
		for _, v2 := range v1["t2"].([]models.TpshopCategory) {
			t := make(map[string]interface{})
			var thirdClass []models.TpshopCategory
			//获取三级菜单
			o.QueryTable("TpshopCategory").Filter("Pid", v2.Id).All(&thirdClass)
			t["t22"] = v2         //二级菜单
			t["t23"] = thirdClass //三级菜单
			erji = append(erji, t)
			//把二级容器放到总容器中
			v1["t3"] = erji
		}
	}

	this.Data["types"] = types
	this.TplName = "index.html"

}
func (this *GoodsController) ShowIndexSx() {

	o := orm.NewOrm()
	//show goodstypes
	var goodsTypes []models.GoodsType
	//find all goodstypes ->[]goodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	//send to index.html   & index.html recevie by {{.goodsType}}
	this.Data["goodsTypes"] = goodsTypes

	//轮播图
	var goodsBanners []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&goodsBanners)
	this.Data["goodsBanners"] = goodsBanners
	//促销商品
	var promotionBanners []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionBanners)
	this.Data["promotions"] = promotionBanners
	//获取首页商品展示
	var goods []map[string]interface{}
	for _, v := range goodsTypes {
		var textgoods []models.IndexTypeGoodsBanner
		var imagegoods []models.IndexTypeGoodsBanner
		qs := o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").Filter("GoodsType__Id", v.Id).OrderBy("Index")
		qs.Filter("DisplayType", 0).All(&textgoods)
		qs.Filter("DisplayType", 1).All(&imagegoods)
		//行容器
		temp := make(map[string]interface{})
		temp["goodsType"] = v
		temp["textgoods"] = textgoods
		temp["imagegoods"] = imagegoods
		goods = append(goods, temp)

	}
	this.Data["goods"] = goods
	this.TplName = "index_sx.html"
}
func (this *GoodsController) ShowDetail() {
	id, err := this.GetInt("Id")
	if err != nil {
		beego.Error("商品下架")
		this.Redirect("index_sx", 302)
		return
	}
	var goodsSKU models.GoodsSKU
	var newGoods []models.GoodsSKU
	o := orm.NewOrm()
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType", "Goods").Filter("Id", id).All(&goodsSKU)
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name", goodsSKU.GoodsType.Name)
	qs.OrderBy("-Time").Limit(2, 0).All(&newGoods)
	this.Data["goodsSKU"] = goodsSKU
	this.Data["newGoods"] = newGoods
	this.TplName = "detail.html"

}
func PageEdit(pageCount, pageIndex int) []int {
	var pages []int
	if pageCount < 5 {
		for i := 1; i <= pageCount; i++ {
			pages = append(pages, i)
		}
	} else if pageIndex <= 5 {
		for i := 1; i <= 5; i++ {
			pages = append(pages, i)
		}
	} else if pageIndex >= pageCount-2 {
		for i := pageCount - 4; i <= pageCount; i++ {
			pages = append(pages, i)
		}
	} else {
		for i := pageIndex - 2; i <= pageIndex+2; i++ {
			pages = append(pages, i)
		}
	}

	return pages

}
func (this *GoodsController) ShowList() {
	o := orm.NewOrm()
	//show goodstypes
	var goodsTypes []models.GoodsType
	//find all goodstypes ->[]goodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	//send to index.html   & index.html recevie by {{.goodsType}}
	this.Data["goodsTypes"] = goodsTypes
	//find the click type by Id
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("类型不存在")
		this.Layout = "llayout.html"
		this.TplName = "index_sx.html"
		return
	}
	this.Data["id"] = id
	//sendback Id to show in top of newgoodsshow
	goodtype := o.QueryTable("GoodsType").Filter("Goodstype__Id", id)
	//??????
	this.Data["goodtype"] = goodtype

	//recieve sort to rand goods
	sort := this.GetString("sort")

	//goods
	var goods []models.GoodsSKU
	//newgoods
	var newgoods []models.GoodsSKU
	//find goods by type through Id
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id)

	//only 2 newest
	qs.OrderBy("-Time").Limit(2, 0).All(&newgoods)
	this.Data["newgoods"] = newgoods
	this.Layout = "llayout.html"
	this.TplName = "list.html"

	//page recive
	//pagecount = count/pagesize
	count, _ := qs.Count()
	pagesize := 1
	pagecount := int(math.Ceil(float64(count) / float64(pagesize)))
	pageindex, err := this.GetInt("pageindex")
	if err != nil {
		pageindex = 1
	}
	pages := PageEdit(pagecount, pageindex)
	var prepage, nextpage int
	if pageindex-1 < 0 {
		prepage = 1
	} else {
		prepage = pageindex - 1
	}
	if pageindex+1 > pagecount {
		nextpage = pagecount
	} else {
		nextpage = pageindex + 1
	}

	this.Data["pages"] = pages
	this.Data["prepage"] = prepage
	this.Data["nextpage"] = nextpage
	qs = qs.Limit(pagesize, pagesize*(pageindex-1))
	if sort == "" {
		qs.All(&goods)

	} else if sort == "price" {
		qs.OrderBy("Price").All(&goods)
	} else {
		qs.OrderBy("-Sales").All(&goods)
	}
	this.Data["sort"] = sort
	this.Data["goods"] = goods
}

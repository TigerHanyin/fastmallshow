package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pyg/pyg/models"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type OderController struct {
	beego.Controller
}

func(this*OderController)ShowOrder(){
	goodsIds:=this.GetStrings("checkGoods")
	if len(goodsIds)==0{
		this.Redirect("/user/ShowCart",302)
		return
	}
	//处理数据

	//获取当前用户的所有收货地址
	name:=this.GetSession("name")
	o:=orm.NewOrm()
	var addrs []models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name",name).All(&addrs)
	this.Data["addrs"]=addrs
	conn,_:=redis.Dial("tcp",":6739")
	defer conn.Close()
	var goods []map[string]interface{}
	var totalPrice,totalCount int
	for _,v:=range goodsIds{
		temp:=make(map[string]interface{})
		id,_ := strconv.Atoi(v)
		var goodsSku models.GoodsSKU
		goodsSku.Id=id
		o.Read(&goodsSku)
		//获取商品数量
		count,_:=redis.Int(conn.Do("hget","cart_"+name.(string),id))
		//计算小计
		littlePrice := count * goodsSku.Price

		//把商品信息放到行容器
		temp["goodsSku"] = goodsSku
		temp["count"] = count
		temp["littlePrice"] = littlePrice
		totalPrice += littlePrice
		totalCount += 1
		goods = append(goods,temp)
	}

	this.Data["totalPrice"] = totalPrice
	this.Data["totalCount"] = totalCount
	this.Data["truePrice"] = totalPrice + 10
	this.Data["goods"] = goods
	this.TplName = "place_order.html"
}

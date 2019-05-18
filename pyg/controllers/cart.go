package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego/orm"
	"pyg/pyg/models"
)

type CartController struct {
	beego.Controller
}

func (this *CartController) HandellAddCart() {
	//获取数据
	id, err := this.GetInt("goodsId")
	num, err2 := this.GetInt("num")
	//返回ajax步骤
	//定义一个map容器
	resp := make(map[string]interface{})

	//封装，集成，多态
	defer RespFunc(&this.Controller, resp)

	//校验数据
	if err != nil || err2 != nil {
		resp["errno"] = 1
		resp["errmsg"] = "输入数据不完整"
		return
	}
	//校验登录状态
	name := this.GetSession("name")
	if name == nil {
		resp["errno"] = 2
		resp["errmsg"] = "当前用户未登录，不能添加购物车"
		return
	}
	conn, err := redis.Dial("tcp",":6379")
	if err!=nil{
		resp["errno"]=3
		resp["errmsg"]="服务器异常"
		return
	}
	defer conn.Close()

	oldNum,_ := redis.Int(conn.Do("hget","cart_"+name.(string),id))

	_,err = conn.Do("hset","cart_"+name.(string),id,oldNum + num)
	if err != nil {
		resp["errno"] = 4
		resp["errmsg"] = "添加商品到购物车失败"
		return
	}

	//返回数据
	resp["errno"] = 5
	resp["errmsg"] = "OK"
}
func(this*CartController)ShowCart(){
conn,err:=redis.Dial("tcp",":6379")
if err!=nil{
	this.Redirect("/index",302)
	return
}
defer conn.Close()
	//查询所有购物车数据
	name:=this.GetSession("name")
	result,err:=redis.Ints( conn.Do("hgetall","cart_"+name.(string)))
	if err!=nil{
		this.Redirect("/index_sx",302)
		return
	}
	//定义大容器
	var goods []map[string]interface{}
	o:=orm.NewOrm()
	litPrice:=0
	totPrice:=0
	totCount:=0
	for i:=0;i<len(result);i+=2{
		 temp:=make(map[string]interface{})
		//result[i]//goodsId    获取商品信息   商品数量
		var goodsSku models.GoodsSKU
		goodsSku.Id=result[i]
		o.Read(&goodsSku)
		//给行容器赋值
		temp["goodsSku"]=goodsSku
		temp["count"]=result[i+1]
		temp["litPrice"]=goodsSku.Price*result[i+1]
		totPrice+=litPrice
		totCount++
		//把行容器添加到大容器里面
		goods=append(goods,temp)

	}

	this.Data["goods"]=goods
	this.Data["totPrice"]=totPrice
	this.Data["totCount"]=totCount

	this.TplName="cart.html"

}
func(this*CartController)HandellUpCart(){
	count,err:=this.GetInt("count")
	Id,err2:=this.GetInt("goodsId")
	resp:=make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)
	if err!=nil||err2!=nil{
		resp["errno"]=1
		resp["errmasg"]="传输数据不完全"
		return
	}
	name:=this.GetSession("name")
	if name==nil{
		resp["errno"]=2
		resp["errmasg"]="用户未登陆"
		return
	}
	//向redis中写入购物车数量
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		resp["errno"]=3
		resp["errmasg"]="redis连接失败"
		return
	}
	defer conn.Close()
	_,err=conn.Do("hset","cart_"+name.(string),Id,count)
	if err!=nil{
		resp["errno"]=4
		resp["errmasg"]="redis插入数据失败"
		return
	}
	resp["errno"]=5
	resp["errmasg"]="ok"
}
func(this*CartController)HandellDeleteCart(){

	Id,err:=this.GetInt("goodsId")
	resp:=make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)
	if err!=nil{
		resp["errno"]=1
		resp["errmasg"]="传输数据不完全"
		return
	}
	name:=this.GetSession("name")
	if name==nil{
		resp["errno"]=2
		resp["errmasg"]="用户未登陆"
		return
	}
	//向redis中写入购物车数量
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		resp["errno"]=3
		resp["errmasg"]="redis连接失败"
		return
	}
	defer conn.Close()
	_,err=conn.Do("hdel","cart_"+name.(string),Id)
	if err!=nil{
		resp["errno"]=4
		resp["errmasg"]="数据库异常"
		return
	}
	resp["errno"]=5
	resp["errmasg"]="ok"
}
package routers

import (
	"pyg/pyg/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

)

func init() {
	//路由过滤
	beego.InsertFilter("/user/*",beego.BeforeExec,guolvFunc)
	beego.Router("/", &controllers.MainController{})
	//用户注册
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	//发送短信
	beego.Router("/sendMsg", &controllers.UserController{}, "post:HandleSendMsg")

	beego.Router("/register-email", &controllers.UserController{}, "get:ShowEmail;post:HandleEmail")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/active", &controllers.UserController{}, "get:Active")
	beego.Router("/index", &controllers.GoodsController{}, "get:ShowIndex")
	beego.Router("/user/userCenterInfo", &controllers.UserController{}, "get:ShowUserCenterInfo")
	beego.Router("/user/site",&controllers.UserController{},"get:ShowSite;post:HandleSite")
	beego.Router("/user/logout",&controllers.UserController{},"get:Logout")
	beego.Router("/index_sx",&controllers.GoodsController{},"get:ShowIndexSx")
	//商品详情
	beego.Router("/goodsDetail",&controllers.GoodsController{},"get:ShowDetail")
	beego.Router("/goodsType",&controllers.GoodsController{},"get:ShowList")
	beego.Router("/addCart",&controllers.CartController{},"post:HandellAddCart")
	beego.Router("/user/ShowCart",&controllers.CartController{},"get:ShowCart")
	beego.Router("/upCart",&controllers.CartController{},"post:HandellUpCart")
	beego.Router("/deleteCart",&controllers.CartController{},"post:HandellDeleteCart")
	beego.Router("/user/addOrder",&controllers.OderController{},"post:ShowOrder")
}


func guolvFunc(ctx*context.Context){
	name:=ctx.Input.Session("name")
	if name==nil{
		ctx.Redirect(302,"/login")
		return

	}

}
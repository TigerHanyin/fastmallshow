package controllers

import (
	"github.com/astaxie/beego"
	"regexp"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"encoding/json"
	"math/rand"
	"time"
	"fmt"
	"github.com/astaxie/beego/orm"
	"pyg/pyg/models"
	"github.com/astaxie/beego/utils"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

func RespFunc(this *beego.Controller, resp map[string]interface{}) {
	//3.把容器传递给前段
	this.Data["json"] = resp
	//4.指定传递方式
	this.ServeJSON()
}

type Message struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}

//发送短信
func (this *UserController) HandleSendMsg() {
	//接受数据
	phone := this.GetString("phone")
	resp := make(map[string]interface{})

	defer RespFunc(&this.Controller, resp)
	//返回json格式数据
	//校验数据
	if phone == "" {
		beego.Error("获取电话号码失败")
		//2.给容器赋值
		resp["errno"] = 1
		resp["errmsg"] = "获取电话号码错误"
		return
	}
	//检查电话号码格式是否正确
	reg, _ := regexp.Compile(`^1[3-9][0-9]{9}$`)
	result := reg.FindString(phone)
	if result == "" {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 2
		resp["errmsg"] = "电话号码格式错误"
		return
	}
	//发送短信   SDK调用
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "LTAIu4sh9mfgqjjr", "sTPSi0Ybj0oFyqDTjQyQNqdq9I9akE")
	if err != nil {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 3
		resp["errmsg"] = "初始化短信错误"
		return
	}
	//生成6位数随机数
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06d", rnd.Int31n(1000000))

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = "品优购"
	request.QueryParams["TemplateCode"] = "SMS_164275022"
	request.QueryParams["TemplateParam"] = `{"code":` + vcode + `}`
	response, err := client.ProcessCommonRequest(request)
	//if err != nil {
	//	beego.Error("电话号码格式错误")
	//	//2.给容器赋值
	//	resp["errno"] = 4
	//	resp["errmsg"] = "短信发送失败"
	//	return
	//}
	//json数据解析
	var message Message
	json.Unmarshal(response.GetHttpContentBytes(), &message)
	if message.Message != "OK" {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 6
		resp["errmsg"] = message.Message
		return
	}

	resp["errno"] = 5
	resp["errmsg"] = "发送成功"

}
func (this *UserController) HandleRegister() {
	phone := this.GetString("phone")
	pwd := this.GetString("password")
	rpwd := this.GetString("repassword")

	if phone == "" || pwd == "" || rpwd == "" {
		beego.Error("获取数据错误")
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "register.html"
		return
	}
	if pwd != rpwd {
		beego.Error("两次输入密码不一致")
		this.Data["errmsg"] = "两次输入密码不一致"
		this.TplName = "register.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = phone
	user.Pwd = pwd
	o.Insert(&user)
	this.Ctx.SetCookie("userName", user.Name, 60*10)
	this.Redirect("/register-email", 302)
}

func (this *UserController) ShowEmail() {
	this.TplName = "register-email.html"
}
func (this *UserController) HandleEmail() {
	email := this.GetString("email")
	pwd := this.GetString("password")
	rpwd := this.GetString("repassword")
	if email == "" || pwd == "" || rpwd == "" {
		beego.Error("输入的数据不完整")
		this.Data["errmsg"] = "输入数据不完整"
		this.TplName = "register-email.html"
		return
	}
	if pwd != rpwd {
		beego.Error("两次输入密码不一致")
		this.Data["errmsg"] = "两次输入密码不一致"
		this.TplName = "register-email.html"
		return
	}
	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(email)
	if result == "" {
		beego.Error("邮箱格式错误")
		this.Data["errmsg"] = "邮箱格式错误"
		this.TplName = "register-email.html"
		return
	}
	//处理数据 发送邮件
	config := `{"username":"czbkttsx@163.com","password":"czbkpygbj3q","host":"smtp.163.com","port":25}`
	emailReg := utils.NewEMail(config)
	emailReg.Subject = "品优购用户激活"
	emailReg.From = "czbkttsx@163.com"
	emailReg.To = []string{email}
	userName := this.Ctx.GetCookie("userName")
	emailReg.HTML = `<a href="http://127.0.0.1:8080/active?userName=` + userName + `"点击激活该用户</a>`
	err := emailReg.Send()
	beego.Error(err)

	//插入邮箱 更新邮箱字段
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")
	if err != nil {
		beego.Error("错误处理")
		return
	}
	user.Email = email
	o.Update(&user)
	this.Ctx.WriteString("邮件已发送，去激活吧")
}
func (this *UserController) Active() {
	//获取数据
	userName := this.GetString("userName")
	//校验数据
	if userName == "" {
		beego.Error("用户名错误")
		this.Redirect("/register-email", 302)
		return
	}

	//处理数据   本质上是更新active
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名不存在")
		this.Redirect("/register-email", 302)
		return
	}
	user.Active = true
	o.Update(&user, "Active")

	//返回数据
	this.Redirect("/login", 302)

}
func (this *UserController) ShowLogin() {
	name := this.Ctx.GetCookie("LoginName")
	if name == "" {
		this.Data["checked"] = ""

	} else {
		this.Data["checked"] = "checked"
	}
	this.Data["userName"] = name
	this.TplName = "login.html"
}
func (this *UserController) HandleLogin() {
	//获取数据   注册的时候要求用户名必须为字母加数字
	userName := this.GetString("userName")
	pwd := this.GetString("pwd")
	if userName == "" || pwd == "" {
		beego.Error("用户名或密码错误")
		this.Data["errmsg"] = "用户名或密码错"
		this.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	reg, _ := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	result := reg.FindString(userName)
	if result != "" {
		user.Email = userName
		err := o.Read(&user, "Email")
		if err != nil {
			this.Data["errmsg"] = "邮箱未注册"
			this.TplName = "login.html"
			return
		}
		if user.Pwd != pwd {
			this.Data["errmsg"] = "密码错误"
			this.TplName = "login.html"
			return
		}

	} else {
		user.Name = userName
		err := o.Read(&user, "Name")
		if err != nil {
			this.Data["errmsg"] = "用户名不存在"
			this.TplName = "login.html"
			return
		}

	}
	if user.Active == false {
		this.Data["errmsg"] = "当前用户未激活，请去目标邮箱激活！"
		this.TplName = "login.html"
		return
	}
	m1 := this.GetString("m1")
	if m1 == "2" {
		this.Ctx.SetCookie("LoginName", user.Name, 60*60)
	} else {
		this.Ctx.SetCookie("LoginName", user.Name, -1)
	}
	this.SetSession("name", user.Name)
	this.Redirect("/index", 302)

}
func (this *UserController) Logout() {
	this.DelSession("name")
	this.Redirect("/index", 302)
}

func (this *UserController) ShowUserCenterInfo() {
	o:=orm.NewOrm()
	var user models.User
	name:=this.GetSession("name")
	user.Name=name.(string)
	o.Read(&user,"Name")
	this.Data["user"]=user

	var addr models.Address
	qs:=o.QueryTable("Address").RelatedSel("User").Filter("User__Name",user.Name)
	qs.Filter("IsDefault",true).One(&addr)
	this.Data["addr"]=addr
	this.Data["tplName"]="个人信息"



	this.Layout = "layout.html"
	this.TplName = "user_center_info.html"
}
func (this *UserController) ShowSite() {
	o := orm.NewOrm()
	var address models.Address
	name := this.GetSession("name")
	qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name.(string))
	qs.Filter("IsDefault", true).One(&address)
	this.Data["address"] = address
	this.Layout = "layout.html"
	this.TplName = "user_center_site.html"
}
func (this *UserController) HandleSite() {
	receiver := this.GetString("receiver")
	phone := this.GetString("phone")
	code := this.GetString("code")
	addr := this.GetString("addr")

	if receiver == "" || phone == "" || code == "" || addr == "" {
		beego.Error("补全信息")
		this.Layout = "layout.html"
		this.TplName = "/user_center_site.html"
		return
	}
	//获取orm对象
	o := orm.NewOrm()
	//获取插入对象
	var Addr models.Address
	//给插入对象赋值
	Addr.Phone = phone
	Addr.PostCode = code
	Addr.Addr = addr
	Addr.Receiver = receiver
	//是哪个用户的地址
	username := this.GetSession("name")
	var user models.User
	user.Name = username.(string)
	o.Read(&user, "Name")
	Addr.User = &user
	//查询看有没有默认地址，如果有，把默认地址修改为非默认 ，如果没有，直接插入默认地址
	//查询当前用户是否有默认地址  queryseter
	var oldAddress models.Address
	qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", username.(string))
	err := qs.Filter("IsDefault", true).One(&oldAddress)
	if err == nil {
		oldAddress.IsDefault = false
		o.Update(&oldAddress, "IsDefault")
	}
	Addr.IsDefault = true
	_, err = o.Insert(&Addr)
	if err != nil {
		beego.Error("插入失败", err)
		return
	}
	this.Layout = "layout.html"
	this.Redirect("/user/site", 302)

}

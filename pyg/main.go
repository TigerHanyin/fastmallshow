package main

import (
	_ "pyg/pyg/routers"
	"github.com/astaxie/beego"
	_ "pyg/pyg/models"
	"strings"
)

func main() {
	beego.AddFuncMap("changenum", ChangeNum)
	beego.Run()
}
func ChangeNum(str string) string {
strr:=strings.Split(str,"")
	for k, _ := range strr {
		if k >= 3 && k <= 6 {
			strr[k]="*"
		}
	}
	str=strings.Join(strr,"")
	return str

}

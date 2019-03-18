package main

import (
	"github.com/astaxie/beego"
	"github.com/monsterry/openvpn-web-ui/lib"
	_ "github.com/monsterry/openvpn-web-ui/routers"
)

func main() {
	lib.AddFuncMaps()
	beego.Run()
}

package main

import (
	_ "BookCommunity/routers"
	_ "BookCommunity/sysint"
	"github.com/astaxie/beego"
)

func main() {

	beego.Run()
}


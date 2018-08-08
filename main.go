package main

import (
	_ "29q/day7/ihome/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
	"strings"
)

func main() {

	beego.SetStaticPath("/static","static")
	ignoreStaticPath()
	beego.Run()
}

func ignoreStaticPath() {
	//透明static
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	///index.html --->/static/html/index.html

	//通过这个地址进行访问127.0.0.1:8080/index.html
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
	//http://127.0.0.1:8080/static/html/index.html
}



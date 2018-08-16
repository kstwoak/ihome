package controllers

import (
	"github.com/astaxie/beego"
	"29q/day9/ihome/utils"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Retdata(resp interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}


func (c *IndexController) GetIndex() {
	//打印被调用的函数
	beego.Info("---------------- GET  /api/v1.0/houses/index GetIndex() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)


	return
}


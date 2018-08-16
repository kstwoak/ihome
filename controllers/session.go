package controllers

import (
	"github.com/astaxie/beego"
	"29q/day9/ihome/utils"
)

type SessionController struct {
	beego.Controller
}

func (c *SessionController) Retdata(resp interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}


func (c *SessionController) Getsession() {
	//打印被调用的函数
	beego.Info("---------------- GET  /api/v1.0/session GetSession() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})
	//初始化是没有用户登陆的状态

	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)

	/*从session中获取 name 字段 */
	name_map:=  make(map[string]interface{})
	name := c.GetSession("name")
	if name !=nil {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		/*如果有返回成功 并且返回 name字段*/
		name_map["name"]= name.(string)

		resp["data"] = name_map

		return
	}

	/*如果没有 初始化的时候默认返回错误*/

	return
}

func (c *SessionController) Deletesession() {
	//打印被调用的函数
	beego.Info("---------------- DELETE  /api/v1.0/session Deletesession() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})
	//初始化是没有用户登陆的状态

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)
	//c.DelSession()
	/*删除相应的session字段*/
	c.DelSession("name")
	c.DelSession("user_id")
	c.DelSession("mobile")


	return
}

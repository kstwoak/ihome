package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"29q/day7/ihome/models"

	"29q/day7/ihome/utils"
)

type AreasController struct {
	beego.Controller
}

func (c *AreasController) Retdata(resp interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}


func (c *AreasController) GetAreas() {
	//打印被调用的函数
	beego.Info("---------------- GET  api/v1.0/areas GetAreas() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)

	/*从 缓存数据库 中获取缓存信息 如果存在就发送给前端*/

	/*如果缓存中没有就 查询mysql 数据库 获取到areas的数据*/
	//创建orm句柄
	o :=orm.NewOrm()
	//创建存储查询条件的变量
	var area []models.Area
	//设置查询条件
	qs:=o.QueryTable("area")
	num ,err :=qs.All(&area)
	if err != nil{
		beego.Info("GetAreas() qs.All(&area) err",err,num)
						//调用错误码 方法函数包 中的 错误码
		resp["errno"] = utils.RECODE_DBERR
						//通过错误码获取到错误信息
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		//c.Data["json"] = &resp
		//c.ServeJSON()
		return
	}
	if num == 0{
		beego.Info("GetAreas() qs.All(&area)  无数据",num)
		//resp["errno"] = 404
		//resp["errmsg"] = "数据库没有数据"
		resp["errno"] = utils.RECODE_NODATA
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))

		//c.Data["json"] = &resp
		//c.ServeJSON()
		return

	}

	/*将获取到的数据打包成为json数据存入 缓存数据库 */


	/*将获取好的areas 发送给前端*/
	/*
	"errno": 0,
    "errmsg":"OK",
    "data": [


	*/

	//resp["errno"] = 200
	//resp["errmsg"] = "成功"


	resp["data"] = area
	//c.Data["json"] = &resp
	//c.ServeJSON()

	return
}


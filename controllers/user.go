package controllers


import (
	"github.com/astaxie/beego"
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/gomodule/redigo/redis"
	"29q/day9/ihome/utils"
	"time"
	"encoding/json"
	"reflect"
	"github.com/garyburd/redigo/redis"
	//"regexp"


	"crypto/md5"
	"encoding/hex"
	//"fmt"
	//"io/ioutil"
	//"net/http"
	"net/url"
	"strconv"
	//"strings"
	"math/rand"
	"29q/day9/ihome/models"
	"github.com/astaxie/beego/orm"
	"path"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Retdata(resp interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}


func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *UserController) GetImage() {
	//打印被调用的函数
	beego.Info("---------------- GET  /api/v1.0/imagecode/* GetImage() ------------------")
	//创建返回空间
	//resp := make(map[string]interface{})

	//resp["errno"] = utils.RECODE_SESSIONERR
	//resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	//defer c.Retdata(resp)

	/*获取前端发送过来的uuid*/
	//this.Ctx.Input.Param(":splat")
	uuid :=c.Ctx.Input.Param(":splat")
	beego.Info(uuid)
	//http://127.0.0.1:8080/api/v1.0/imagecode/bec28d5e-501b-4117-8b6f-9c1c493ca3f4
	/*生成1个随机数验证码与验证码图片*/

	//创建1个句柄
	cap := captcha.New()
	//通过句柄调用 字体文件
	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}

	//设置图片的大小
	cap.SetSize(91, 41)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	//SetFrontColor(colors ...color.Color)  这两个颜色设置的函数属于不定参函数
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	//生成图片 返回图片和 字符串(图片内容的文本形式)
	img, str := cap.Create(4, captcha.NUM)
	beego.Info(str)


	//config :=make(map[string]string)
	//config["key"] =  "ihome"
	//config["conn"] = utils.G_redis_addr


	/*将uuid与 随机数验证码对应的存储在redis缓存中*/
	//初始化缓存全局变量的对象
	bm, err := cache.NewCache("redis", `{"key":"ihome","conn":"127.0.0.1:6379","dbNum":"0"}`)
	if err !=nil{
		beego.Info("GetImage()   cache.NewCache err ",err)
	}
	//bm.Put("astaxie", 1, 10*time.Second) redis 存储的操作
	bm.Put(uuid ,str, 600 *time.Second )  //验证码进行1个小时缓存


	/*
	配置信息如下所示，redis 采用了库 redigo:
	{"key":"collectionName","conn":":6039","dbNum":"0","password":"thePassWord"}
	key: Redis collection 的名称
	conn: Redis 连接信息
	dbNum: 连接 Redis 时的 DB 编号. 默认是0.
	password: 用于连接有密码的 Redis 服务器.

	*/


	/*向前端页面返回验证码图片*/
	//将图片发送给前端的 直接发送图片
	png.Encode(c.Ctx.ResponseWriter, img)


	return
}


func (c *UserController) Getsmscode() {
	//打印被调用的函数
	beego.Info("---------------- GET  /api/v1.0/smscode/:id Getsmscode() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)

	/*获取前端发送过来的id*/
	id :=c.Ctx.Input.Param(":id")
	beego.Info(id)

	/*接收url发送过来的参数*/
	//http://127.0.0.1:8080/api/v1.0/smscode/1111?
	// text=7529  &  id=ea4ab7dd-6e9f-497e-bdf5-05bec585a40d

	var text string
	c.Ctx.Input.Bind(&text,"text")
	var uuid string
	c.Ctx.Input.Bind(&uuid ,"id")

	beego.Info(text,uuid)

	/*连接 缓存数据库 获取缓存信息进行验证 图片验证码是否正确*/
	//构建连接缓存的数据
	redis_config_map := map[string]string{
		"key":"ihome",
		//"conn":"127.0.0.1:6379",
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
	}
	beego.Info(redis_config_map)
	redis_config ,_:=json.Marshal(redis_config_map)
	beego.Info( string(redis_config) )


	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redis_config) )
	if err !=nil{
		beego.Info("GetImage()   cache.NewCache err ",err)
	}

	value :=bm.Get(uuid)
	if  value == nil{

		beego.Info("Getsmscode()bm.Get(uuid) err  ",value)
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}
	beego.Info(value,reflect.TypeOf(value))

	value_str ,_ :=redis.String(value,nil)
	/*
	第一个参数填写通过bm.get 获取到的返回值
	第二个一般情况下写nil
	*/
	beego.Info(value_str,reflect.TypeOf(value_str))
	//数据对比
	if text != value_str{
		beego.Info("图片验证码 错误 ")
		resp["errno"] = utils.RECODE_PWDERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))

	}

	/*手机号的验证
	0?(13|14|15|17|18|19)[0-9]{9}
	*/
/*
	myreg :=regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	bo :=myreg.MatchString(id) //判断传入的字符串与正则格式是否匹配
	//MatchString类似Match，但匹配对象是字符串。
	if bo == false {
		beego.Info("手机号错误")
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = "手机号错误"
		return
	}*/



	/*通过已有的短信验证码的接口 模拟发送短信*/
	//type Values map[string][]string
	v := url.Values{}  //创建1个url.values map
	//格式化当前的时间
	_now := strconv.FormatInt(time.Now().Unix(), 10)
	beego.Info(_now)
	_account := "C10921244"  //账户名 需要花钱购买的
	_password := "da3018614650fa96137bb61dd71e85d8" //查看密码请登录用户中心->验证码、通知短信->帐户及签名设置->APIKEY

	_mobile := string(id)  //手机号  通过id进行赋值
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //生成随机数

	sms_code := r.Intn(9999)
	beego.Info("短信验证码是",sms_code)
	//拼接短信发送的内容
	_content := "您的验证码是：" + strconv.Itoa(sms_code) + "。请不要把验证码泄露给其他人。"
	//添加url的内容
	v.Set("account", _account)
	v.Set("password", GetMd5String(_account+_password+_mobile+_content+_now))
	v.Set("mobile", _mobile)
	v.Set("content", _content)
	v.Set("time", _now)

	//body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	//client := &http.Client{}
	//req, _ := http.NewRequest("POST", "http://106.ihuyi.com/webservice/sms.php?method=Submit&format=json", body)
	//
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	////fmt.Printf("%+v\n", req) //看下发送的结构
	//beego.Info(req)
	//resp1, err := client.Do(req) //发送
	//defer resp1.Body.Close()     //一定要关闭resp.Body
	//data, _ := ioutil.ReadAll(resp1.Body)
	//fmt.Println(string(data), err)



	/*将短信验证码存入缓存数据库 */
	bm.Put("smscode",strconv.Itoa(sms_code) , 600 *time.Second)

	/*成功返回ok */
	return

}


func (c *UserController) Postuserret() {
	//打印被调用的函数
	beego.Info("---------------- POST  /api/v1.0/users Postuserret() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)

	/* 获得用户注册信息*/
	//err = json.Unmarshal(this.Ctx.Input.RequestBody, &ob);
	var Requestmap = make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &Requestmap)

	for key, value := range Requestmap {
		beego.Info(key, value)
	}
	 beego.Info( Requestmap["mobile"])
	 beego.Info( Requestmap["password"])
	 beego.Info( Requestmap["sms_code"])

	/*校验信息准确信*/
	if Requestmap["mobile"] == ""|| Requestmap["password"] == "" || Requestmap["sms_code"] == ""{
		beego.Info("注册数据为空")
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	}

	/*验证短信验证码 */
	//构建连接缓存的数据
	redis_config_map := map[string]string{
		"key":"ihome",
		//"conn":"127.0.0.1:6379",
		"conn":utils.G_redis_addr+":"+utils.G_redis_port,
		"dbNum":utils.G_redis_dbnum,
	}
	beego.Info(redis_config_map)
	redis_config ,_:=json.Marshal(redis_config_map)
	beego.Info( string(redis_config) )
	//连接redis数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redis_config) )
	if err !=nil{
		beego.Info("GetImage()   cache.NewCache err ",err)
	}
    //获取我们存在缓存数据库中的短信验证码
	value :=bm.Get("smscode")
	if  value == nil{

		beego.Info("Postuserret()  bm.Get(uuid) err  ",value)
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}
	beego.Info(value,reflect.TypeOf(value))

	value_str ,_ :=redis.String(value,nil)
	/*
	第一个参数填写通过bm.get 获取到的返回值
	第二个一般情况下写nil
	*/
	beego.Info(value_str,reflect.TypeOf(value_str))
	//短信验证码对比
	if  Requestmap["sms_code"].(string) != value_str{
		beego.Info("短信验证码错误 错误 ")
		resp["errno"] = utils.RECODE_PWDERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}

	/*将用户信息存储在 mysql 中 */

	//beego.Info( Requestmap["mobile"])
	//beego.Info( Requestmap["password"])
	//beego.Info( Requestmap["sms_code"])
	//创建1个mysql对象
	user:= models.User{}
	user.Name = Requestmap["mobile"].(string)
	user.Mobile = Requestmap["mobile"].(string)
	//正常情况下我们需要吧 password 转成md5 或sha256 等等密文格式   为了调试所以进行直接赋值
	//正常情况下密码的相关对比都要转成密文进行操作
	user.Password_hash =Requestmap["password"].(string)
	//操作数据库
	o :=orm.NewOrm()
	id ,err := o.Insert(&user)
	if err != nil{
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}
	beego.Info(id)

	/*添加session字段 */
	//SetSession(name string, value interface{})
	c.SetSession("name", Requestmap["mobile"].(string))
	c.SetSession("user_id",id)
	c.SetSession("mobile", Requestmap["mobile"].(string))

	/*进行返回*/

	return
}



func (c *UserController) Postlogin() {
	//打印被调用的函数
	beego.Info("---------------- POST  /api/v1.0/sessions Postlogin() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)


	/*解析前端用户发送过来的信息*/
	var Requestmap = make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &Requestmap)

	for key, value := range Requestmap {
		beego.Info(key, value)
	}
	beego.Info( Requestmap["mobile"])
	beego.Info( Requestmap["password"])

	/*校验信息准确信*/
	if Requestmap["mobile"] == ""|| Requestmap["password"] == ""{
		beego.Info("注册数据为空")
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	}

	/*查询数据库 获取 信息 */

	var user models.User

	o :=orm.NewOrm()
	//设置查询的表
	//select mobile = Requestmap["mobile"] from user  返回是1个唯一的
	qs:=o.QueryTable("user")
	err :=qs.Filter("mobile",Requestmap["mobile"]).One(&user)
	if err!=nil{
		/*如果不匹配则登陆失败*/
		beego.Info("用户名查询失败",err)
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}

	/*查询密码判断是否正确 */

	if user.Password_hash != Requestmap["password"].(string){
		/*如果不匹配则登陆失败*/
		beego.Info("密码错误")
		resp["errno"] = utils.RECODE_PWDERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}


	/*添加session字段 */

	c.SetSession("name", Requestmap["mobile"].(string))
	c.SetSession("user_id",user.Id)
	c.SetSession("mobile", Requestmap["mobile"].(string))


	return
}




func (c *UserController) Postupavatar() {
	//打印被调用的函数
	beego.Info("---------------- POST  /api/v1.0/user/avatar Postupavatar() ------------------")
	//创建返回空间
	resp := make(map[string]interface{})

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
	//延迟调用发送给前端json数据
	defer c.Retdata(resp)

	/*获取前端发过来的文件数据*/
	//GetFile(key string) (multipart.File, *multipart.FileHeader, error)
	file ,hander,err := c.GetFile("avatar")
	if err != nil{
		beego.Info("Postupavatar   c.GetFile(avatar) err" ,err)
		resp["errno"] = utils.RECODE_IOERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}
	beego.Info(file ,hander)
	beego.Info("文件大小",hander.Size)
	beego.Info("文件名",hander.Filename)
	//二进制的空间用来存储文件
	filebuffer:= make([]byte,hander.Size)
	//将文件读取到filebuffer里
	_,err = file.Read(filebuffer)
	if err !=nil{
		beego.Info("Postupavatar   file.Read(filebuffer) err" ,err)
		resp["errno"] = utils.RECODE_IOERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}


	/*获取文件的后缀名*/     //dsnlkjfajadskfksda.sadsdasd.sdasd.jpg
	beego.Info("后缀名",path.Ext(hander.Filename))




	/*存储文件到fastdfs当中并且获取 url*/
	//.jpg
	fileext :=path.Ext(hander.Filename)
	//group1 group1/M00/00/00/wKgLg1t08pmANXH1AAaInSze-cQ589.jpg

	Group,FileId ,err :=  models.UploadByBuffer(filebuffer,fileext[1:])
	if err != nil {
		beego.Info("Postupavatar  models.UploadByBuffer err" ,err)
		resp["errno"] = utils.RECODE_IOERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}

	beego.Info(Group,FileId)

	/*通过session 获取我们当前现在用户的uesr_id*/
	userid :=c.GetSession("user_id")

	beego.Info(userid ,reflect.TypeOf(userid))
	//创建数据库对象
	user :=models.User{Id:userid.(int),Avatar_url:FileId}

	//我们往数据库里存的是这一段  group1/M00/00/00/wKgLg1t08pmANXH1AAaInSze-cQ589.jpg


	/*将当前fastdfs-url 存储到我们当前用户的表中*/

	//创建数据库句柄
	o :=orm.NewOrm()
	//我们将数据更新上去 返回的id就是我们user的id也就不需要获取了
	_,err =o.Update(&user,"avatar_url")
	if err!=nil{
		beego.Info("Postupavatar  o.Update err" ,err)
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(resp["errno"].(string))
		return
	}

	/*将fastdfs-url 与我们当前的信息进行拼接*/

	allurl:=utils.AddDomain2Url(FileId)
	beego.Info(allurl)

	/*返回给前端完整的图片url*/
	url_map := make(map[string]interface{})

	url_map["avatar_url"] = allurl

	resp["data"] = url_map

	return

}


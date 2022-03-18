package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"path"

	// "math/rand"
	"net/http"
	// "time"
	"zufang/web/model"
	"zufang/web/utils"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"

	capt "zufang/web/proto/getCaptcha"
	userMicro "zufang/web/proto/user"

	"github.com/gomodule/redigo/redis"
	"github.com/tedcy/fdfs_client"
)

func GetSession(ctx *gin.Context) {
	resp := map[string]interface{}{}

	session := sessions.Default(ctx)
	userName := session.Get("userName")
	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		resp["data"] = map[string]string{"name": userName.(string)}
	}

	ctx.JSON(http.StatusOK, resp)
}

//获取图片信息
func GetImageCd(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	consulReg := consul.NewRegistry()
	microService := micro.NewService(
		micro.Registry(consulReg),
	)
	microClient := capt.NewCaptchaService("captcha", microService.Client())
	resp, err := microClient.Call(context.TODO(), &capt.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("未找到远程服务。。。")
	}

	var img captcha.Image
	json.Unmarshal(resp.Img, &img)
	// img, str := cap.Create(4, captcha.NUM)

	png.Encode(ctx.Writer, img)

	// fmt.Println(str)
}

func GetSmscd(ctx *gin.Context) {
	phone := ctx.Param("phone")
	imgCode := ctx.Query("text")
	uuid := ctx.Query("id")

	fmt.Println("-----out----:", phone, imgCode, uuid)

	/////////////---------------------------------------
	// resp := make(map[string]string)
	// //校验图片验证码 是否正确
	// result := model.CheckImgCode(uuid, imgCode)
	// if result {
	// 	//发送短信
	// 	rand.Seed(time.Now().UnixNano())
	// 	smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
	// 	err := model.SaveSmsCode(phone, smsCode)
	// 	if err != nil {
	// 		fmt.Println("存储验证码到redis失败, err :", err)

	// 		resp["errno"] = utils.RECODE_DBERR
	// 		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
	// 	}
	// 	//使用阿里短信api发送短信
	// 	fmt.Println("验证码为：", smsCode)

	// 	resp["errno"] = utils.RECODE_OK
	// 	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	// } else {
	// 	resp["errno"] = utils.RECODE_SMSERR
	// 	resp["errmsg"] = utils.RecodeText(utils.RECODE_SMSERR)
	// }

	consulReg := consul.NewRegistry()
	microService := micro.NewService(
		micro.Registry(consulReg),
	)
	microClient := userMicro.NewUserService("user", microService.Client())
	resp, err := microClient.SendSms(context.TODO(), &userMicro.Request{Phone: phone, ImgCode: imgCode, Uuid: uuid})
	if err != nil {
		fmt.Println("未找到远程服务。。。")
	}

	ctx.JSON(http.StatusOK, resp)
}

//获取注册信息
func PostRet(ctx *gin.Context) {
	// mobile := ctx.PostForm("mobile")
	// pwd := ctx.PostForm("password")
	// smsCode := ctx.PostForm("sms_code")
	// fmt.Println("mobile:", mobile, "pwd:", pwd, "smscode:", smsCode)

	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.Bind(&regData)
	fmt.Println("获取到的数据为：", regData)

	//注册用户
	microService := InitMicro()
	microClient := userMicro.NewUserService("user", microService.Client())
	resp, err := microClient.Register(context.TODO(), &userMicro.RegReq{
		Mobile: regData.Mobile, Password: regData.PassWord, SmsCode: regData.SmsCode})
	if err != nil {
		fmt.Println("未找到远程服务。。。")
	}

	ctx.JSON(http.StatusOK, resp)
}

func InitMicro() micro.Service {
	consulReg := consul.NewRegistry()
	return micro.NewService(
		micro.Registry(consulReg),
	)
}

func GetArea(ctx *gin.Context) {
	conn := model.RedisPool.Get()
	var areas []model.Area

	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))
	if len(areaData) == 0 {
		fmt.Println("从mysql中获取数据。。。")
		model.GlobalConn.Find(&areas)
		areaBuf, _ := json.Marshal(areas)
		conn.Do("set", "areaData", areaBuf)
	} else {
		fmt.Println("从redis中获取数据。。。")
		json.Unmarshal(areaData, &areas)
	}

	resp := make(map[string]interface{})
	resp["errno"] = "0"
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas

	ctx.JSON(http.StatusOK, resp)
}

func PostLogin(ctx *gin.Context) {
	var loginData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
	}
	ctx.Bind(&loginData)
	fmt.Println("获取到的数据为：", loginData)

	resp := make(map[string]interface{})
	//获取数据库数据 查询是否和数据的数据匹配
	userName, err := model.Login(loginData.Mobile, loginData.PassWord)
	if err != nil {
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
	} else {

		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		//将登陆状态保存到Session中
		session := sessions.Default(ctx)
		session.Set("userName", userName)
		session.Save()
	}

	ctx.JSON(http.StatusOK, resp)
}

func DeleteSession(ctx *gin.Context) {
	resp := map[string]interface{}{}

	session := sessions.Default(ctx)
	session.Delete("userName")
	err := session.Save()
	if err != nil {
		resp["errno"] = utils.RECODE_IOERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetUserInfo(ctx *gin.Context) {
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)

	session := sessions.Default(ctx)
	userName := session.Get("userName")

	if userName == nil { //用户没登陆，但进入该页面，恶意进入
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}

	user, err := model.GetUserInfo(userName.(string))
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = map[string]interface{}{
		"user_id":    user.ID,
		"name":       user.Name,
		"mobile":     user.Mobile,
		"real_name":  user.Real_name,
		"id_card":    user.Id_card,
		"avatar_url": "http://192.168.154.131:8888/" + user.Avatar_url,
	}
}

func PutUserInfo(ctx *gin.Context) {
	resp := map[string]interface{}{}
	defer ctx.JSON(http.StatusOK, resp)

	//获取当前用户名
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	//获取新用户名
	var nameData struct {
		Name string `json:"name"`
	}
	ctx.Bind(&nameData)

	//更新用户名
	err := model.UpdateUserName(nameData.Name, userName.(string))
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	session.Set("userName", nameData.Name)
	err = session.Save()
	if err != nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = nameData
}

func saveImagInFastDfs(buffer []byte, fileExtName string) (string, error) {
	client, err := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
	defer client.Destory()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return client.UploadByBuffer(buffer, fileExtName)
}

func PostAvatar(ctx *gin.Context) {
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)
	//上传头像
	// 单文件
	file, _ := ctx.FormFile("avatar")
	// 上传文件到指定的路径
	// err := ctx.SaveUploadedFile(file, "media/"+file.Filename)
	// fmt.Println(err)
	f, err := file.Open()
	if err != nil {
		resp["errno"] = utils.RECODE_IOERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)
		return
	}
	buf := make([]byte, file.Size)
	if _, err := f.Read(buf); err != nil {
		resp["errno"] = utils.RECODE_IOERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)
		return
	}
	fileExt := path.Ext(file.Filename)

	fileId, err := saveImagInFastDfs(buf, fileExt[1:])

	userName := sessions.Default(ctx).Get("userName")

	if err := model.UpdateAvatar(userName.(string), fileId); err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = map[string]interface{}{"avatar_url": "http://192.168.154.131:8888/" + fileId}
}

// 上传实名认证
type AuthStu struct {
	IdCard   string `json:"id_card"`
	RealName string `json:"real_name"`
}

func PutUserAuth(ctx *gin.Context) {
	resp := map[string]interface{}{}
	defer ctx.JSON(http.StatusOK, resp)

	//获取数据
	var auth AuthStu
	err := ctx.Bind(&auth)
	//校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}
	//获取当前用户名
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	//更新真实信息
	err = model.UpdateUserAuth(userName.(string), auth.IdCard, auth.RealName)
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
}

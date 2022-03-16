package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"image/png"

	// "math/rand"
	"net/http"
	// "time"
	"zufang/web/model"
	"zufang/web/utils"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"

	"github.com/gomodule/redigo/redis"
	capt "zufang/web/proto/getCaptcha"
	userMicro "zufang/web/proto/user"
)

func GetSession(ctx *gin.Context) {
	// // You also can use a struct
	// var msg struct {
	// 	Errno  string `json:"errno"`
	// 	Errmsg string `json:"errmsg"`
	// 	Number int    `json:"errno"`
	// }
	// msg.Errno = utils.RECODE_SESSIONERR
	// msg.Errmsg = utils.RecodeText(utils.RECODE_SESSIONERR)
	// msg.Number = 123
	// // Note that msg.Name becomes "user" in the JSON
	// // Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
	// c.JSON(http.StatusOK, msg)

	resp := map[string]string{}
	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
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

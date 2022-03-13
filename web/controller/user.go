package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"zufang/web/utils"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"

	capt "zufang/web/proto"
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

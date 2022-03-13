package handler

import (
	"context"
	"encoding/json"
	"image/color"

	log "github.com/micro/micro/v3/service/logger"

	"captcha/model"
	capt "captcha/proto"

	captcha "github.com/afocus/captcha"
)

type Captcha struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Captcha) Call(ctx context.Context, req *capt.Request, rsp *capt.Response) error {
	log.Info("Received Captcha.Call request")

	cap := captcha.New()
	// 可以设置多个字体 或使用cap.AddFont("xx.ttf")追加
	cap.SetFont("./conf/comic.ttf")
	// 设置验证码大小
	cap.SetSize(128, 64)
	// 设置干扰强度
	cap.SetDisturbance(captcha.NORMAL)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cap.SetFrontColor(color.RGBA{102, 139, 139, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 48, 48, 255}, color.RGBA{255, 48, 167, 255}, color.RGBA{85, 26, 139, 255})

	img, str := cap.Create(4, captcha.NUM)
	err := model.SaveImgCode(str, req.Uuid)
	if err != nil {
		return err
	}
	imgBuf, _ := json.Marshal(img)
	rsp.Img = imgBuf

	return nil
}

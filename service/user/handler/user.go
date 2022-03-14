package handler

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"user/model"
	"user/utils"

	user "user/proto"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) SendSms(ctx context.Context, req *user.Request, rsp *user.Response) error {
	//校验图片验证码 是否正确
	result := model.CheckImgCode(req.Uuid, req.ImgCode)
	if result {
		//发送短信
		rand.Seed(time.Now().UnixNano())
		smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
		err := model.SaveSmsCode(req.Phone, smsCode)
		if err != nil {
			fmt.Println("存储验证码到redis失败, err :", err)

			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		}
		//使用阿里短信api发送短信
		fmt.Println("验证码为：", smsCode)

		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	} else {
		rsp.Errno = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
	}

	return nil
}

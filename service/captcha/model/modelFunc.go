package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func SaveImgCode(code, uuid string) error {
	//链接数据库
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis dial err:", err)
		return err
	}
	defer conn.Close()

	//操作数据库
	_, err = conn.Do("setex", uuid, 60*5, code)

	//回复助手函数-----确定成具体的数据类型
	return err
}

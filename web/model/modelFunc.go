package model

import "fmt"

func CheckImgCode(uuid, imgCode string) bool {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err:", err)
		return false
	}
	defer conn.Close()

	//查询redis
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误 err:", err)
		return false
	}
	return code == imgCode
}

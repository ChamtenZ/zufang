package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

//创建全局redis连接池句柄
var RedisPool redis.Pool

//创建函数，初始化Redis连接池
func InitRedis() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

func CheckImgCode(uuid, imgCode string) bool {
	// conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	// if err != nil {
	// 	fmt.Println("redis.Dial err:", err)
	// 	return false
	// }
	// defer conn.Close()

	conn := RedisPool.Get()
	defer conn.Close()

	//查询redis
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询Imgcode错误 err:", err)
		return false
	}
	return code == imgCode
}

func SaveSmsCode(phone, code string) error {
	//从连接池获取一条连接
	conn := RedisPool.Get()
	defer conn.Close()

	//查询redis
	_, err := redis.String(conn.Do("setex", phone+"_code", 60*3, code))
	return err
}

func CheckSmsCode(phone, smsCode string) bool {
	conn := RedisPool.Get()
	defer conn.Close()

	//查询redis
	// fmt.Println("phone_code:", phone+"_code")
	code, err := redis.String(conn.Do("get", phone+"_code"))
	if err != nil {
		fmt.Println("查询Smscode错误 err:", err)
		return false
	}
	return code == smsCode
}

func RegisterUser(name, mobile, pwd string) error {
	var user User
	user.Name = name
	user.Mobile = mobile

	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))
	user.Password_hash = pwd_hash

	return GlobalConn.Create(&user).Error
}

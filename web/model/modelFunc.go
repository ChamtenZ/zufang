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

func SaveSmsCode(phone, code string) error {
	//从连接池获取一条连接
	conn := RedisPool.Get()
	defer conn.Close()

	//查询redis
	_, err := redis.String(conn.Do("setex", phone+"_code", 60*3, code))
	return err
}

func Login(mobile, pwd string) (string, error) {
	var user User

	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))
	err := GlobalConn.Select("name").Where("mobile=? and password_hash=?", mobile, pwd_hash).Find(&user).Error
	return user.Name, err
}

func GetUserInfo(userName string) (User, error) {
	var user User
	err := GlobalConn.Where("name=?", userName).First(&user).Error
	return user, err
}

func UpdateUserName(newName, oldName string) error {
	return GlobalConn.Model(new(User)).Where("name=?", oldName).Update("name", newName).Error
}

//根据用户名更新用户头像
func UpdateAvatar(userName, avatar string) error {
	return GlobalConn.Model(new(User)).Where("name=?", userName).Update("avatar_url", avatar).Error
}

func UpdateUserAuth(userName, idCard, realName string) error {
	return GlobalConn.Model(new(User)).Where("name=?", userName).Update(map[string]interface{}{
		"id_card": idCard, "real_name": realName}).Error
}

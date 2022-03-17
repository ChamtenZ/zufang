package main

import (
	"zufang/web/controller"
	"zufang/web/model"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.Default(ctx)
		userName := s.Get("userName")
		if userName == nil {
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}

func main() {

	// //初始化连接池
	model.InitRedis()
	//初始化Mysql链接池
	model.InitDb()

	//初始化路由
	router := gin.Default()

	//初始化session
	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("zufang"))
	router.Use(sessions.Sessions("mysession", store))

	//路由匹配
	// router.GET("/", func(context *gin.Context) {
	// 	context.Writer.WriteString("项目开始了。。。。")
	// })
	router.Static("/home", "view")
	// router.GET("/api/v1.0/session", controller.GetSession)
	// router.GET("api/v1.0/imagecode/:uuid", controller.GetImageCd)

	//添加路由分组
	r1 := router.Group("/api/v1.0")
	{
		r1.GET("/session", controller.GetSession)
		r1.GET("/imagecode/:uuid", controller.GetImageCd)
		r1.GET("/smscode/:phone", controller.GetSmscd)
		r1.POST("/users", controller.PostRet)
		r1.GET("/areas", controller.GetArea)
		r1.Use(LoginFilter()) //以后的路由都不需要校验session
		r1.POST("/sessions", controller.PostLogin)
		r1.DELETE("/session", controller.DeleteSession)
		r1.GET("/user", controller.GetUserInfo)
		r1.PUT("/user/name", controller.PutUserInfo)
		r1.POST("/user/avatar", controller.PostAvatar)
	}
	router.Run(":8080")
}

// func InitRedis() {
// 	panic("unimplemented")
// }

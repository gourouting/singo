package main

import (
	"gin_example/api"
	"gin_example/cache"
	"gin_example/middleware"
	"gin_example/model"
	"gin_example/util"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// 从配置文件读取配置
	if err := util.LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
		panic(err)
	}

	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(model.Database(os.Getenv("MYSQL_DSN")))
	r.Use(cache.Redis())
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", api.UserLogin)

		// 需要登录保护的
		v1.Use(middleware.AuthRequired())
		{
			// User Routing
			v1.GET("user/me", api.UserMe)
			v1.DELETE("user/logout", api.UserLogout)
		}
	}

	// 监听
	r.Run(":3000")
}

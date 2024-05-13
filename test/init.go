package test

import (
	"os"
	"singo/cache"
	"singo/conf"
	"singo/model"
	"singo/server"
	"singo/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	s *gin.Engine
)

func init() {
	// 从配置文件读取配置
	confInit()
	// API
	s = server.NewRouter()
}

// Init 初始化配置项
func confInit() {
	// 从本地读取环境变量
	godotenv.Load()

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 读取翻译文件
	if err := conf.LoadLocales("../conf/locales/zh-cn.yaml"); err != nil {
		util.Log().Panic("翻译文件加载失败", err)
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
}

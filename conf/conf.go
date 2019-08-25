package conf

import (
	"singo/cache"
	"singo/model"
	"os"

	"github.com/joho/godotenv"
)

// Init 初始化配置项
func init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 读取翻译文件
	if err := LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
		panic(err)
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
}

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
	// Load configuration from files.
	confInit()
	// API
	s = server.NewRouter()
}

// confInit initializes configuration.
func confInit() {
	// Load environment variables from local files.
	godotenv.Load()

	// Set the log level.
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// Load translation files.
	if err := conf.LoadLocales("../conf/locales/zh-cn.yaml"); err != nil {
		util.Log().Panic("failed to load translation file: %v", err)
	}

	// Connect to the database.
	dsn, err := conf.DatabaseDSN()
	if err != nil {
		panic(err)
	}

	model.Database(dsn)
	cache.Redis()
}

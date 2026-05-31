package main

import (
	"os"
	"singo/conf"
	"singo/server"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from files.
	if err := conf.Init(); err != nil {
		panic(err)
	}

	// Load routes.
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := server.NewRouter()
	r.Run(":3000")
}

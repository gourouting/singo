package server

import (
	"os"
	"singo/api"
	"singo/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter configures routes.
func NewRouter() *gin.Engine {
	r := gin.Default()

	// Middleware. The order must not be changed.
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// User registration
		v1.POST("user/register", api.UserRegister)

		// User login
		v1.POST("user/login", api.UserLogin)

		// Routes that require login protection.
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)
		}
	}
	return r
}

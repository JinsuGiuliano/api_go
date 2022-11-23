package router

import (
	"example/api_go/controller"
	"example/api_go/server"
)

func Set(server *server.Server, cont *controller.Controller) {

	server.APIRouter.SetTrustedProxies([]string{"127.0.0.1"})

	v1 := server.APIRouter.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.GET("/", cont.HandleMain)
		auth.GET("/login-gl", cont.HandleGoogleLogin)
		auth.GET("/callback-gl", cont.CallBackFromGoogle)
		auth.GET("/login", cont.Login)
		auth.POST("/signup", cont.SignupUser)
	}

	authorized := v1.Group("/")
	authorized.Use(cont.TokenOrCookieAuth)
	{
		authorized.GET("/user", cont.ListUsers)
		authorized.POST("/user", cont.InsertUser)
		authorized.DELETE("/user/:id", cont.DeleteUSer)
	}

	server.APIRouter.Run("localhost:9090")
}

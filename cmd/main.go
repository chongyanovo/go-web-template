package main

import (
	"fmt"
	"github.com/ChongYanOvO/go-web-template/bootstrap"
	"github.com/ChongYanOvO/go-web-template/internal/repository"
	"github.com/ChongYanOvO/go-web-template/internal/repository/dao"
	"github.com/ChongYanOvO/go-web-template/internal/service"
	"github.com/ChongYanOvO/go-web-template/internal/web"
	"github.com/ChongYanOvO/go-web-template/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.InitApp()
	config := app.Config
	server := gin.Default()
	server.Use(middleware.Cors())

	store := cookie.NewStore([]byte("secret"))
	//store, err := redis.NewStore(16,
	//	"tcp",
	//	"localhost:6379",
	//	"",
	//	[]byte("secret"),
	//	[]byte("secret"))
	//if err != nil {
	//	panic(err)
	//}
	//store := memstore.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store))

	server.Use(
		middleware.
			AuthMiddleware().
			IgnorePaths("/user/login").
			IgnorePaths("/user/register").
			Builder(),
	)

	if err := dao.InitTable(app.Mysql); err != nil {
		panic("数据库初始化错误")
	}
	userDao := dao.NewUserDao(app.Mysql)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)

	userGroup := server.Group("/user")
	userGroup.POST("/register", userHandler.Register)
	userGroup.POST("/login", userHandler.Login)
	userGroup.GET("/profile", userHandler.Profile)

	server.Run(fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port))
}

package main

import (
	"fmt"
	"github.com/ChongYanOvO/go-web-template/bootstrap"
	"github.com/ChongYanOvO/go-web-template/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.InitApp()
	config := app.Config
	engine := gin.Default()
	engine.Use(middleware.Cors())
	engine.Run(fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port))
}

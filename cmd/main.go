package main

import (
	"fmt"
	"github.com/ChongYanOvO/blog/server/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.InitApp()
	config := app.Config
	gin := gin.Default()
	gin.Run(fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port))
}

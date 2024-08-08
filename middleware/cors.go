package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

// Cors 跨域配置
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "P ATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length", "Content-Type", "x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "xxx.com")
		},
		AllowCredentials: true, // 允许带cookie跨域
		MaxAge:           12 * 3600,
	})
}

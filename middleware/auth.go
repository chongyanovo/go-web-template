package middleware

import (
	"fmt"
	"github.com/ChongYanOvO/go-web-template/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	paths []string
}

func AuthMiddleware() *Auth {
	return &Auth{}
}

func (a *Auth) IgnorePaths(path string) *Auth {
	a.paths = append(a.paths, path)
	return a
}

func (a *Auth) Builder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range a.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenHeader := ctx.GetHeader("x-jwt-token")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenHeaderSplit := strings.Split(tokenHeader, " ")
		if len(tokenHeaderSplit) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := tokenHeaderSplit[1]
		userClaims := &web.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid || userClaims.UserId == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		if userClaims.ExpiresAt.Sub(now) < time.Hour*2 {
			userClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
			jwtString, err := token.SignedString([]byte("secret"))
			if err != nil {
				fmt.Println("jwt刷新失败")
			}
			ctx.Header("x-jwt-token", "Bearer "+jwtString)
		}

		ctx.Set("userClaims", userClaims)
	}
}

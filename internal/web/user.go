package web

import (
	"errors"
	"fmt"
	"github.com/ChongYanOvO/go-web-template/internal/domain"
	"github.com/ChongYanOvO/go-web-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type UserHandler struct {
	svc *service.UserService
}

type RegisterReq struct {
	UserName        string `json:"userName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64 `json:"user_id"`
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (u UserHandler) Register(ctx *gin.Context) {
	var req RegisterReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}
	err := u.svc.Register(ctx, &domain.User{
		UserName: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrorUserDuplicate) {
			ctx.String(http.StatusOK, "用户名已存在")
		} else {
			ctx.String(http.StatusOK, "系统异常")
		}
	} else {
		ctx.String(http.StatusOK, "注册成功")
	}

}

func (u UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.UserName, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrorInvalidUserNameOrPassword) {
			ctx.String(http.StatusOK, "用户名或密码错误")
		} else {
			ctx.String(http.StatusOK, "系统异常")
		}
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			&UserClaims{RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			}, UserId: user.ID})
		jwtString, err := token.SignedString([]byte("secret"))
		if err != nil {
			ctx.String(http.StatusOK, "系统异常")
		}
		ctx.Header("x-jwt-token", "Bearer "+jwtString)
		ctx.String(http.StatusOK, "登录成功")
	}

}

func (u UserHandler) Profile(ctx *gin.Context) {
	userClaims, ok := ctx.Get("userClaims")
	if !ok {
		ctx.String(http.StatusOK, "系统异常")
	}
	claims, ok := userClaims.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统异常")
	}
	fmt.Println(claims.UserId)
	ctx.String(http.StatusOK, "profile")
}

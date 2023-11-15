package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xuqil/webook/internal/web"
	"net/http"
	"strings"
	"time"
)

// LoginJTWMiddlewareBuilder JWT 登录校验
type LoginJTWMiddlewareBuilder struct {
	paths []string
}

func NewLoginJTWMiddlewareBuilder() *LoginJTWMiddlewareBuilder {
	return &LoginJTWMiddlewareBuilder{}
}

func (l *LoginJTWMiddlewareBuilder) IgnorePaths(path string) *LoginJTWMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJTWMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		// ParseWithClaims 一定要传指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("343d9040a671c45832ee5381860e2996"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 安全问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			//Secure: true,
			//HttpOnly: true,
			MaxAge: 60,
		})
		now := time.Now()
		// 说明还没有刷新，刚登陆
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		// updateTimeVal 是有的
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Minute {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}

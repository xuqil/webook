//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/xuqil/webook/internal/repository"
	"github.com/xuqil/webook/internal/repository/cache"
	"github.com/xuqil/webook/internal/repository/dao"
	"github.com/xuqil/webook/internal/service"
	"github.com/xuqil/webook/internal/web"
	"github.com/xuqil/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,

		// 初始化 DAO
		dao.NewGORMUserDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		service.NewUserService,
		service.NewCodeService,
		// 直接基于内存实现
		ioc.InitSMSService,
		web.NewUserHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

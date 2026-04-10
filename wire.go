//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
package main

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	user2 "github.com/Knowpals/Knowpals-be-go/controller/user"
	"github.com/Knowpals/Knowpals-be-go/infra/email"
	"github.com/Knowpals/Knowpals-be-go/ioc"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	"github.com/Knowpals/Knowpals-be-go/repository/cache"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	"github.com/Knowpals/Knowpals-be-go/service/user"
	"github.com/Knowpals/Knowpals-be-go/web"
	"github.com/google/wire"
)

func InitApp(conf *config.Config) *App {
	wire.Build(
		//pkg
		ijwt.NewJwtHandler,

		//middleware
		middleware.NewAuthMiddleware,
		middleware.NewLoggerMiddleware,
		//middleware.NewOtelMiddleware,

		//基础设施层
		ioc.InitDB,
		ioc.InitZapLogger,
		ioc.InitRedis,

		//infra
		email.NewEmailClient,

		//repository
		dao.NewUserDao,
		cache.NewAuthCache,

		//service
		user.NewUserService,

		//controller
		user2.NewUserController,

		//web
		web.NewGinEngine,

		NewApp,
	)

	return &App{}
}

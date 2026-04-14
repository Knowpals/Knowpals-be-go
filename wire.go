//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
package main

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	classController "github.com/Knowpals/Knowpals-be-go/controller/class"
	user2 "github.com/Knowpals/Knowpals-be-go/controller/user"
	"github.com/Knowpals/Knowpals-be-go/controller/video"
	"github.com/Knowpals/Knowpals-be-go/events"
	"github.com/Knowpals/Knowpals-be-go/events/consumer"
	"github.com/Knowpals/Knowpals-be-go/events/producer"
	"github.com/Knowpals/Knowpals-be-go/infra/cos"
	"github.com/Knowpals/Knowpals-be-go/infra/email"
	"github.com/Knowpals/Knowpals-be-go/ioc"
	"github.com/Knowpals/Knowpals-be-go/middleware"
	"github.com/Knowpals/Knowpals-be-go/pkg/ijwt"
	"github.com/Knowpals/Knowpals-be-go/repository/cache"
	"github.com/Knowpals/Knowpals-be-go/repository/dao"
	classService "github.com/Knowpals/Knowpals-be-go/service/class"
	"github.com/Knowpals/Knowpals-be-go/service/pipeline"
	"github.com/Knowpals/Knowpals-be-go/service/user"
	video2 "github.com/Knowpals/Knowpals-be-go/service/video"
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
		ioc.InitCOS,
		ioc.InitKafka,
		ioc.InitKafkaConsumerGroupID,

		//infra
		email.NewEmailClient,
		cos.NewCOSClient,

		//repository
		dao.NewUserDao,
		dao.NewClassDao,
		dao.NewVideoDao,
		dao.NewPipelineDao,
		dao.NewKnowledgeDao,
		dao.NewSegmentDao,
		dao.NewQuestionDao,
		cache.NewAuthCache,
		producer.NewSaramaProducer,
		consumer.NewSaramaConsumer,
		events.NewPipelineWorker,

		//service
		user.NewUserService,
		classService.NewClassService,
		video2.NewVideoService,
		pipeline.NewPipelineService,

		//controller
		user2.NewUserController,
		classController.NewClassController,
		video.NewVideoController,

		//web
		web.NewGinEngine,

		NewApp,
	)

	return &App{}
}

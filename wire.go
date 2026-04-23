//go:build wireinject
// +build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
package main

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/Knowpals/Knowpals-be-go/controller/agent"
	"github.com/Knowpals/Knowpals-be-go/controller/behavior"
	classController "github.com/Knowpals/Knowpals-be-go/controller/class"
	"github.com/Knowpals/Knowpals-be-go/controller/question"
	"github.com/Knowpals/Knowpals-be-go/controller/review"
	"github.com/Knowpals/Knowpals-be-go/controller/statistic"
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
	agent2 "github.com/Knowpals/Knowpals-be-go/service/agent"
	"github.com/Knowpals/Knowpals-be-go/service/agentclient"
	behaviorService "github.com/Knowpals/Knowpals-be-go/service/behavior"
	classService "github.com/Knowpals/Knowpals-be-go/service/class"
	"github.com/Knowpals/Knowpals-be-go/service/pipeline"
	question2 "github.com/Knowpals/Knowpals-be-go/service/question"
	statService "github.com/Knowpals/Knowpals-be-go/service/statistic"
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
		middleware.NewCorsMiddleware,

		//基础设施层
		ioc.InitDB,
		ioc.InitZapLogger,
		ioc.InitRedis,
		ioc.InitCOS,
		ioc.InitKafka,
		ioc.InitKafkaConsumerGroupID,
		ioc.InitGRPCConn,
		ioc.InitMemoryClient,
		ioc.InitAgentClient,

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
		dao.NewReportDao,
		dao.NewStatisticDao,
		dao.NewBehaviorDao,
		cache.NewAuthCache,
		producer.NewSaramaProducer,
		consumer.NewSaramaConsumer,
		events.NewPipelineWorker,
		dao.NewChatDao,

		//service
		user.NewUserService,
		classService.NewClassService,
		video2.NewVideoService,
		behaviorService.NewBehaviorService,
		statService.NewStatService,
		pipeline.NewPipelineService,
		agentclient.NewMemoryWriter,
		question2.NewQuestionService,
		agent2.NewAgentService,

		//controller
		user2.NewUserController,
		classController.NewClassController,
		video.NewVideoController,
		question.NewQuestionController,
		behavior.NewBehaviorController,
		statistic.NewStatController,
		agent.NewAgentController,
		review.NewReviewController,

		//web
		web.NewGinEngine,

		NewApp,
	)

	return &App{}
}

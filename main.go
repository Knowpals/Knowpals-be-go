package main

import (
	"context"
	"os"

	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/Knowpals/Knowpals-be-go/events"
	"github.com/Knowpals/Knowpals-be-go/events/consumer"
	"github.com/Knowpals/Knowpals-be-go/events/topic"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Config *config.Config

type App struct {
	r      *gin.Engine
	worker *events.PipelineWorker
	c      consumer.Consumer
}

func NewApp(r *gin.Engine, worker *events.PipelineWorker, c consumer.Consumer) *App {
	return &App{
		r:      r,
		worker: worker,
		c:      c,
	}
}

func (a *App) Run() {
	if a.worker != nil && a.c != nil {
		go func() {
			_ = a.worker.Run(context.Background(), a.c, []string{topic.RESULT_TOPIC})
		}()
	}
	a.r.Run()
}

func initViper() {
	defaultConfig := "./config/config.yaml"
	cfile := pflag.String("config", defaultConfig, "配置文件路径")
	pflag.Parse()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = *cfile
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Config = &config.Config{}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}

}

func main() {
	initViper()
	app := InitApp(Config)
	app.Run()
}

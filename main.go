package main

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Config config.Config

type App struct {
	r *gin.Engine
}

func NewApp(r *gin.Engine) *App {
	return &App{
		r: r,
	}
}

func (a *App) Run() {
	a.r.Run()
}

func initViper() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	pflag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}

}

func main() {
	initViper()

	//shutdown := ioc.InitOtel(&Config)
	//defer func() {
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//	if err := shutdown(ctx); err != nil {
	//		panic(err)
	//	}
	//}()

	app := InitApp(&Config)
	app.Run()
}

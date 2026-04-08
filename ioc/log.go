package ioc

import (
	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLogger(conf *config.Config) *zap.Logger {
	level := zap.DebugLevel

	al := zap.NewAtomicLevelAt(level)
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.Log.File,
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
		Compress:   conf.Log.Compress,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(lumberJackLogger),
		al,
	)
	//标识调用代码行
	return zap.New(core).WithOptions(zap.AddCaller())

}

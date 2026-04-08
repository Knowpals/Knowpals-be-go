package ioc

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/Knowpals/Knowpals-be-go/pkg/otelx"
)

func InitOtel(conf *config.Config) func(ctx context.Context) error {
	if conf == nil || conf.Otel == nil {
		return func(context.Context) error { return nil }
	}
	shutdown, err := otelx.Init(context.Background(), *conf.Otel)
	if err != nil {
		panic(err)
	}
	return shutdown
}

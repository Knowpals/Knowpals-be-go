package otelx

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Knowpals/Knowpals-be-go/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func Init(ctx context.Context, cfg config.OtelConf) (shutdown func(context.Context) error, err error) {
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	cfg = withDefaults(cfg)

	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	res, err := newResource(cfg.ServiceName, cfg.ServiceVersion)
	if err != nil {
		handleErr(err)
		return
	}

	otel.SetTextMapPropagator(newPropagator())

	tracerProvider, err := newTracerProvider(ctx, res, cfg)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider(ctx, res, cfg)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func withDefaults(cfg config.OtelConf) config.OtelConf {
	if cfg.ServiceName == "" {
		cfg.ServiceName = "knowpals-be-go"
	}
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = "dev"
	}
	if cfg.TraceExporter == "" {
		cfg.TraceExporter = "otlp"
	}
	if cfg.MetricsExporter == "" {
		cfg.MetricsExporter = "otlp"
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:4317"
	}
	if cfg.MetricsInterval <= 0 {
		cfg.MetricsInterval = 10
	}
	return cfg
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context, res *resource.Resource, cfg config.OtelConf) (*trace.TracerProvider, error) {
	var (
		traceExporter trace.SpanExporter
		err           error
	)

	switch {
	case strings.EqualFold(cfg.TraceExporter, "stdout"), cfg.TraceExporter == "":
		traceExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	case strings.EqualFold(cfg.TraceExporter, "otlp"):
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		traceExporter, err = otlptracegrpc.New(ctx, opts...)
	default:
		return nil, errors.New("otel trace exporter 仅支持 stdout 或 otlp")
	}
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
	), nil
}

func newMeterProvider(ctx context.Context, res *resource.Resource, cfg config.OtelConf) (*sdkmetric.MeterProvider, error) {
	interval := 10 * time.Second
	if cfg.MetricsInterval > 0 {
		interval = time.Duration(cfg.MetricsInterval) * time.Second
	}

	switch {
	case strings.EqualFold(cfg.MetricsExporter, "none"):
		return sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
		), nil
	case strings.EqualFold(cfg.MetricsExporter, "otlp"), cfg.MetricsExporter == "":
		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlpmetricgrpc.WithInsecure())
		}
		metricExporter, err := otlpmetricgrpc.New(ctx, opts...)
		if err != nil {
			return nil, err
		}

		return sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(interval))),
		), nil
	default:
		return nil, errors.New("otel metrics exporter 仅支持 otlp 或 none")
	}
}

func Tracer(name string) oteltrace.Tracer {
	return otel.Tracer(name)
}

func TraceIDFromContext(ctx context.Context) string {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ""
	}
	return spanCtx.TraceID().String()
}

func SpanIDFromContext(ctx context.Context) string {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ""
	}
	return spanCtx.SpanID().String()
}

func RecordError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	span := oteltrace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := oteltrace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	span.SetAttributes(attrs...)
}

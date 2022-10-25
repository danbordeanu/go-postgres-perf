package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"os"
	"postgres-perf/utils/logger"
)

// InitTracerJaeger Create Jaeger telemetry tracer
func InitTracerJaeger(ctx context.Context, JaegerEngine string, ServiceNameKey string, ServiceInstanceIDKey string, tenant string) (*sdktrace.TracerProvider, error) {

	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "tracer", "action", "create jaeger tracer")
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(JaegerEngine)))
	if err != nil {
		log.Errorf("Issue with the tracer:%s", err.Error())
		return nil, err
	}

	// get hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("Issue getting hostname:%s", err.Error())
	}

	// add tenant ID attribute
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		// Record information about this application in an Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ServiceNameKey),
			semconv.ServiceInstanceIDKey.String(ServiceInstanceIDKey),
			attribute.String("hostname", hostname),
			attribute.String("tenant", tenant),
			attribute.Int64("ID", int64(os.Getpid())),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

//InitTracerStdout create stdout telemetry
func InitTracerStdout(ctx context.Context) (*sdktrace.TracerProvider, error) {
	log := logger.SugaredLogger().WithContextCorrelationId(ctx).With("package", "tracer", "action", "create stdout tracer")
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Errorf("HTTP server failed to shutdown gracefully: %s", err.Error())
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

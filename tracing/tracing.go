package tracing

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/krixlion/dev_forum-lib/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// InitProvider initializes the global meter and tracer providers.
// Returns a func used to shutdown exporters running in the background or an error if the initialization failed.
// Metrics are exported to Prometheus by exposing an insecure localhost:2223/metrics endpoint.
// Traces are exported to OTEL exporter.
// Exporter's URL is read from the OTEL_EXPORTER_OTLP_ENDPOINT environment variable.
// If OTEL_EXPORTER_OTLP_ENDPOINT is unset then it is assumed the URL is 0.0.0.0:4317.
func InitProvider(ctx context.Context, serviceName string) (func(), error) {
	resource, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// The service name used to display traces in backends.
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	otelAgentAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelAgentAddr = "0.0.0.0:4317"
	}

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelAgentAddr),
	)
	if err != nil {
		return nil, err
	}

	promExporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(promExporter),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(2*time.Second),
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
	)

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)
	go func() {
		logging.Log("Serving metrics at localhost:2223/metrics")
		http.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(":2223", nil); err != nil {
			logging.Log("Failed to serve metrics", "err", err)
			return
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := metricExp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}

		if err := traceExp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
		// Pushes any last exports to the receiver.
		if err := meterProvider.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}

		if err := promExporter.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}

		if err := tracerProvider.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}, nil
}

// SetSpanErr records given error to the span and sets the span's status as Error with err.Error() as the description.
func SetSpanErr(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

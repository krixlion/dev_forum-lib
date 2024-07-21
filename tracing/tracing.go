package tracing

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// InitProvider initializes the global meter and tracer providers and returns a func used
// to close exporters running in the background or an error if the initialization failed.
// Traces and metrics are exported to OTel Collector using gRPC.
// Takes a service name to use as a label for exported resource and address to the OTel Collector.
func InitProvider(ctx context.Context, serviceName, otelColelctorAddr string) (func() error, error) {
	if otelColelctorAddr == "" {
		return nil, errors.New("failed to init providers: missing otel-collector url ")
	}

	rsc, err := resource.New(ctx,
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

	rsc, err = resource.Merge(resource.Default(), rsc)
	if err != nil {
		return nil, err
	}

	metricExp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelColelctorAddr),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(rsc),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp, sdkmetric.WithInterval(2*time.Second))),
	)

	otel.SetMeterProvider(meterProvider)

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelColelctorAddr),
	)

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(rsc),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExp)),
	)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() error {
		return errors.Join(metricExp.Shutdown(ctx), traceExp.Shutdown(ctx), meterProvider.Shutdown(ctx), tracerProvider.Shutdown(ctx))
	}, nil
}

// SetSpanErr records given error to the span and sets the span's status as Error with err.Error() as the description.
// If the span is not recording or given err is nil then this func will not do anything.
func SetSpanErr(span trace.Span, err error) {
	if span.IsRecording() && err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// ExtractMetadataFromContext extracts traceparent id from given context and
// returns it in a map following a format {"traceparent": "<id>"}.
func ExtractMetadataFromContext(ctx context.Context) map[string]string {
	metadata := map[string]string{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(metadata))
	return metadata
}

// InjectMetadataIntoContext takes in a map following a format {"traceparent": "<id>"}
// and returns a new context with the metadata appended.
func InjectMetadataIntoContext(ctx context.Context, metadata map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(metadata))
}

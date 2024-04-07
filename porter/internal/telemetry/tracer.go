package telemetry

import (
	"context"

	"github.com/honeycombio/otel-config-go/otelconfig"
)

// TracerConfig contains all config for setting up an otel tracer
type TracerConfig struct {
	// ServiceName will show as service.name in spans
	ServiceName string
	// CollectorURL is the OLTP endpoint for receiving traces
	CollectorURL string

	Debug bool
}

// Tracer is a wrapper for an otel tracer
type Tracer struct {
	config TracerConfig
	// TraceProvider *sdktrace.TracerProvider
	// Tracer        trace.Tracer
	Shutdown func()
}

// InitTracer is using the Honeycomb and Lightstep partnership launcher for setting up opentelemetry
// Make sure to run `defer tp.Shutdown(ctx)` after calling this function
// to ensure that no traces are lost on exit
func InitTracer(ctx context.Context, conf TracerConfig) (Tracer, error) {
	if conf.CollectorURL == "" {
		return Tracer{}, nil
	}

	tracer := Tracer{
		config: conf,
	}

	bsp := NewBaggageSpanProcessor()

	lnchr, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithServiceName(conf.ServiceName),
		otelconfig.WithExporterEndpoint(conf.CollectorURL),
		otelconfig.WithSpanProcessor(bsp),
		otelconfig.WithLogLevel("DEBUG"),
		otelconfig.WithMetricsEnabled(false),  // can turn this on later
		otelconfig.WithExporterInsecure(true), // TODO: disable this before production usage
		// otelconfig.WithHeaders() // TODO: add in information about runtime environment
	)
	if err != nil {
		return tracer, err
	}

	tracer.Shutdown = lnchr
	return tracer, nil
}

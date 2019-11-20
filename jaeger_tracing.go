package observability

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"
)

// traceID returns trace ID of context.
// It supports jaeger only.
// On unavailable trace ID (or unsupported tracer implementations)
// it returns empty string.
func traceID(ctx context.Context) string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ctx, ok := span.Context().(jaeger.SpanContext)
		if !ok {
			return ""
		}
		return ctx.TraceID().String()
	}
	return ""
}

// initJaegerTracer initializes tracing using jaeger for given service.
// We initialize jaeger even if OpenTracing is disabled
// (in that case it's suppressed with NewNullReporter).
// This way jaeger always generates a valid traceID.
// It would be better if we could use an externally generated traceID for jaeger
// traces, but it's not possible yet:
// https://github.com/jaegertracing/jaeger-client-go/issues/397.
func initJaegerTracer(
	serviceName string, reportSpans bool, r *prometheus.Registry,
) (io.Closer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}
	opts := []config.Option{
		config.Logger(Logger{}),
		config.Metrics(jaegerprom.New(jaegerprom.WithRegisterer(r))),
	}
	if !reportSpans {
		opts = append(opts, config.Reporter(jaeger.NewNullReporter()))
	}
	closer, err := cfg.InitGlobalTracer(serviceName, opts...)
	if err != nil {
		return nil, fmt.Errorf("could not initialize jaeger tracer: %w", err)
	}
	return closer, nil
}

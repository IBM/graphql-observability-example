package observability

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/graph-gophers/graphql-go/trace"
	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

// Observer represents a collection of services that
// instruments our application to provide runtime insights
// (metrics, logs, traces).
type Observer struct {
	Logger          Logger
	Tracer          trace.Tracer
	Closer          io.Closer
	MetricsExporter http.Handler
	TraceHeader     func(http.Handler) http.Handler
}

// NewObserver returns a new Observer for given service, configured
// using provided options.
func NewObserver(serviceName string, opts ...Opt) (*Observer, error) {
	var tracers []trace.Tracer
	var reportSpans bool
	for _, opt := range opts {
		tracers = append(tracers, opt.Tracer())
		if opt.ReportSpans() {
			reportSpans = true
		}
	}

	registry, ok := prometheus.DefaultRegisterer.(*prometheus.Registry)
	if !ok {
		return nil, errors.New("cast prometheus.DefaultRegisterer")
	}
	exporter, err := ocprom.NewExporter(ocprom.Options{Registry: registry})
	if err != nil {
		return nil, fmt.Errorf("create prometheus exporter: %w", err)
	}
	view.RegisterExporter(exporter)
	if err := registerAnalyticViews(); err != nil {
		return nil, err
	}
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		return nil, fmt.Errorf("register http views: %w", err)
	}

	closer, err := initJaegerTracer(serviceName, reportSpans, registry)
	if err != nil {
		return nil, err
	}

	return &Observer{
		Tracer:          ChainingTracer{Tracers: tracers},
		Closer:          closer,
		MetricsExporter: exporter,
		TraceHeader:     traceHeader,
	}, nil
}

// Opt enables an opt-in feature to Observer if provided at NewObserver.
type Opt struct {
	t           trace.Tracer
	reportSpans bool
}

// Tracer returns associated tracer with given Opt.
func (o Opt) Tracer() trace.Tracer { return o.t }

// ReportSpans returns whether or not to enable OpenTracing span reporting
// according to given span.
func (o Opt) ReportSpans() bool { return o.reportSpans }

// WithOpenTracing returns an Opt that enables OpenTracing.
func WithOpenTracing() Opt {
	return Opt{
		t:           trace.OpenTracingTracer{},
		reportSpans: true,
	}
}

// WithLogging returns an Opt that enables logging of resolver errors
// (with traceID).
func WithLogging() Opt { return Opt{t: LoggingTracer{}} }

// WithAnalytics returns an Opt that enables GraphQL analytics
// (exposed via Observer.MetricsExporter).
func WithAnalytics() Opt { return Opt{t: AnalyticsTracer{}} }

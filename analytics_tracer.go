package observability

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/trace"
)

// AnalyticsTracer implements graph-gophers/graphql-go/trace/Tracer interface.
// It records GraphQL analytics of queries and fields via QueryAnalyzer and
// FieldAnalyzer.
type AnalyticsTracer struct{}

// TraceQuery initializes a QueryAnalyzer at the begginning of query resolvance.
// At the end of query resolvance QueryAnalyzer.recordMeasurements is called to
// record measures of given query resolvance.
func (t AnalyticsTracer) TraceQuery(
	ctx context.Context, q string, op string,
	vars map[string]interface{}, types map[string]*introspection.Type,
) (context.Context, trace.TraceQueryFinishFunc) {
	analyzer := QueryAnalyzer{
		Ctx:   ctx,
		Start: time.Now(),
	}
	return ctx, analyzer.recordMeasurements
}

// TraceField initializes a FieldAnalyzer at the begginning of field resolvance.
// At the end of field resolvance FieldAnalyzer.recordMeasurements is called to
// record measures of given field resolvance.
func (t AnalyticsTracer) TraceField(
	ctx context.Context, label, typeName, fieldName string, trivial bool,
	args map[string]interface{},
) (context.Context, trace.TraceFieldFinishFunc) {
	analyzer := FieldAnalyzer{
		Ctx:   ctx,
		Type:  typeName,
		Field: fieldName,
		Start: time.Now(),
	}
	return ctx, analyzer.recordMeasurements
}

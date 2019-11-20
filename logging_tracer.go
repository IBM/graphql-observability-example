package observability

import (
	"context"
	"strconv"

	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/trace"
)

// LoggingTracer implements graph-gophers/graphql-go/trace/Tracer interface.
// It provides error logging of resolvers via Logger.
type LoggingTracer struct{}

// TraceQuery returns a callback function that will be called at the end of the
// query.
// This callback function logs all errors occurred in resolvers
// (aggregated, with TraceID).
func (t LoggingTracer) TraceQuery(
	ctx context.Context, q string, op string, vars map[string]interface{},
	types map[string]*introspection.Type,
) (context.Context, trace.TraceQueryFinishFunc) {
	return ctx, logErrors(ctx)
}

// TraceField is a no-op.
// It's defined to implement graph-gophers/graphql-go/trace/Tracer interface.
func (t LoggingTracer) TraceField(
	ctx context.Context, label, typeName, fieldName string, trivial bool,
	args map[string]interface{},
) (context.Context, trace.TraceFieldFinishFunc) {
	return ctx, func(err *errors.QueryError) {}
}

func logErrors(ctx context.Context) trace.TraceQueryFinishFunc {
	return func(errs []*errors.QueryError) {
		aggregated := map[string]int{}
		for _, err := range errs {
			aggregated[err.Message]++
		}
		for err, n := range aggregated {
			Logger{}.ErrorCtx(ctx, err+" (x"+strconv.Itoa(n)+")")
		}
	}
}

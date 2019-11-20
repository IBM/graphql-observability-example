package observability

import (
	"context"

	"github.com/graph-gophers/graphql-go/errors"
	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/graph-gophers/graphql-go/trace"
)

// ChainingTracer implements graph-gophers/graphql-go/trace/Tracer interface.
// It's capable of chaining the tracers provided.
// Tracers are called in the order provided.
// context.Context returned by previous tracer is passed to the next one in
// order.
type ChainingTracer struct {
	Tracers []trace.Tracer
}

// TraceQuery calls TraceQuery method of tracers in the order they are provided.
// context.Context returned by previous TraceQuery is passed to the next one in
// order.
func (t ChainingTracer) TraceQuery(
	ctx context.Context, q string, op string, vars map[string]interface{},
	types map[string]*introspection.Type,
) (context.Context, trace.TraceQueryFinishFunc) {
	var fns []trace.TraceQueryFinishFunc
	prevCtx := ctx
	for _, t := range t.Tracers {
		ctx, fn := t.TraceQuery(prevCtx, q, op, vars, types)
		fns = append(fns, fn)
		prevCtx = ctx
	}
	return prevCtx, chainTraceQueryFinishFns(fns...)
}

// TraceField calls TraceField method of tracers in the order they are provided.
// context.Context returned by previous TraceField is passed to the next one in
// order.
func (t ChainingTracer) TraceField(
	ctx context.Context, label, typeName, fieldName string, trivial bool,
	args map[string]interface{},
) (context.Context, trace.TraceFieldFinishFunc) {
	var fns []trace.TraceFieldFinishFunc
	prevCtx := ctx
	for _, t := range t.Tracers {
		ctx, fn := t.TraceField(prevCtx, label, typeName, fieldName, trivial, args)
		fns = append(fns, fn)
		prevCtx = ctx
	}
	return prevCtx, chainTraceFieldFinishFns(fns...)
}

func chainTraceQueryFinishFns(fns ...trace.TraceQueryFinishFunc) trace.TraceQueryFinishFunc {
	return func(errs []*errors.QueryError) {
		for _, fn := range fns {
			fn(errs)
		}
	}
}

func chainTraceFieldFinishFns(fns ...trace.TraceFieldFinishFunc) trace.TraceFieldFinishFunc {
	return func(err *errors.QueryError) {
		for _, fn := range fns {
			fn(err)
		}
	}
}

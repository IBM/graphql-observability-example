package observability

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Logger provides a consistent logging format.
type Logger struct{}

// Error is to implement jaeger.Logger interface.
// It can also be used if trace ID is not available.
// If trace ID is available ErrorCtx should be used instead.
func (l Logger) Error(msg string) {
	l.ErrorCtx(context.Background(), msg)
}

// ErrorCtx logs error alongside trace ID of context.
func (l Logger) ErrorCtx(ctx context.Context, msg interface{}) {
	s := "[" + l.time() + "]"
	tid := traceID(ctx)
	if tid != "" {
		s += " (trace:" + tid + ")"
	}
	fmt.Fprintf(os.Stderr, "%s error: %v\n", s, msg)
}

func (Logger) time() string {
	return time.Now().Format("02/Jan/2006:15:04:05 -0700")
}

// Infof is to implement jaeger.Logger interface. It doesn't do anything really,
// we should implement it if we want to do info level logging.
// If we do, we should have a configuration flag to turn it off.
func (Logger) Infof(string, ...interface{}) {}

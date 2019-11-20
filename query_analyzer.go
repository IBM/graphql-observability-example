package observability

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go/errors"
	"go.opencensus.io/stats"
)

// QueryAnalyzer is capable of recording measures of given GraphQL query for
// analytic purposes.
type QueryAnalyzer struct {
	Ctx   context.Context
	Start time.Time
}

// recordMeasurements is called at the end of query resolvance.
// It's responsible for recording measures of associated GraphQL
// query resolvance.
func (qa QueryAnalyzer) recordMeasurements(errs []*errors.QueryError) {
	durationMs := float64(time.Since(qa.Start)) / float64(time.Millisecond)
	m := []stats.Measurement{
		QueryResolveCount.M(1),
		QueryResolveDuration.M(durationMs),
	}
	if len(errs) != 0 {
		m = append(m, QueryResolveErrorCount.M(int64(len(errs))))
	}
	stats.Record(qa.Ctx, m...)
}

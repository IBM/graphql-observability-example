package observability

import (
	"context"
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// FieldAnalyzer is capable of recording measures of given GraphQL field for
// analytic purposes.
type FieldAnalyzer struct {
	Ctx         context.Context
	Type, Field string
	Start       time.Time
}

// recordMeasurements is called at the end of field resolvance.
// It's responsible for recording measures of associated GraphQL field
// resolvance.
func (fa FieldAnalyzer) recordMeasurements(err *errors.QueryError) {
	if strings.HasPrefix(fa.Type, "__") {
		// this was an introspection field query, we don't record it currently
		// (I don't see a reason why we would like to see these)
		return
	}
	durationMs := float64(time.Since(fa.Start)) / float64(time.Millisecond)
	m := []stats.Measurement{
		FieldResolveCount.M(1),
		FieldResolveDuration.M(durationMs),
	}
	if err != nil {
		m = append(m, FieldResolveErrorCount.M(1))
	}
	stats.RecordWithTags(fa.Ctx, []tag.Mutator{
		tag.Upsert(TagField, fa.Type+"."+fa.Field),
	}, m...)
}

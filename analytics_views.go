package observability

import (
	"fmt"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func registerAnalyticViews() error {
	err := view.Register(
		FieldResolveCountView,
		FieldResolveDurationView,
		FieldResolveErrorCountView,
		QueryResolveCountView,
		QueryResolveDurationView,
		QueryResolveErrorCountView,
	)
	if err != nil {
		return fmt.Errorf("register analytic views: %w", err)
	}
	return nil
}

// Views of GraphQL fields analytics.
var (
	FieldResolveCountView = &view.View{
		Description: "Number of times given GraphQL field was resolved",
		TagKeys:     []tag.Key{TagField},
		Measure:     FieldResolveCount,
		Aggregation: view.Count(),
	}
	FieldResolveDurationView = &view.View{
		Description: "Duration of given GraphQL field's resolvance",
		TagKeys:     []tag.Key{TagField},
		Measure:     FieldResolveDuration,
		Aggregation: ochttp.DefaultLatencyDistribution,
	}
	FieldResolveErrorCountView = &view.View{
		Description: "Number of times given GraphQL field was resolvance returned an error",
		TagKeys:     []tag.Key{TagField},
		Measure:     FieldResolveErrorCount,
		Aggregation: view.Count(),
	}
	QueryResolveCountView = &view.View{
		Description: "Number of times GraphQL queries were resolved",
		Measure:     QueryResolveCount,
		Aggregation: view.Count(),
	}
	QueryResolveDurationView = &view.View{
		Description: "Duration of GraphQL queries resolvance",
		Measure:     QueryResolveDuration,
		Aggregation: ochttp.DefaultLatencyDistribution,
	}
	QueryResolveErrorCountView = &view.View{
		Description: "Number of errors returned by GraphQL queries resolvance",
		Measure:     QueryResolveErrorCount,
		Aggregation: view.Distribution(1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200),
	}
)

// GraphQL field analytics measures.
var (
	FieldResolveCount = stats.Int64(
		"graphql/server/field_resolve_count",
		"Number of times given GraphQL field was resolved",
		stats.UnitDimensionless)
	FieldResolveDuration = stats.Float64(
		"graphql/server/field_resolve_duration",
		"Duration of given GraphQL field's resolvance",
		stats.UnitMilliseconds)
	FieldResolveErrorCount = stats.Int64(
		"graphql/server/field_resolve_error_count",
		"Number of times given GraphQL field was resolvance returned an error",
		stats.UnitDimensionless)
	QueryResolveCount = stats.Int64(
		"graphql/server/query_resolve_count",
		"Number of times GraphQL queries were resolved",
		stats.UnitDimensionless)
	QueryResolveDuration = stats.Float64(
		"graphql/server/query_resolve_duration",
		"Duration of GraphQL queries resolvance",
		stats.UnitMilliseconds)
	QueryResolveErrorCount = stats.Int64(
		"graphql/server/query_resolve_error_count",
		"Number of errors returned by GraphQL queries resolvance",
		stats.UnitDimensionless)
)

// Tags for constructing viewes of GraphQL field analytics.
var (
	TagField = tag.MustNewKey("graphql.field")
)

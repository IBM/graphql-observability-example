module github.com/IBM/graphql-observability-example

go 1.13

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/graph-gophers/graphql-go v0.0.0-20191115155744-f33e81362277
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.2.1
	github.com/uber/jaeger-client-go v2.20.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.opencensus.io v0.22.2
	go.uber.org/atomic v1.5.1 // indirect
)

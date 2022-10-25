package handlers

import (
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
	"postgres-perf/configuration"
)

// tracer init
var tracer = otel.Tracer(configuration.OTName, oteltrace.WithInstrumentationVersion(configuration.OTVersion), oteltrace.WithSchemaURL(configuration.OTSchema))

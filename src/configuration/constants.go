package configuration

// CorrelationIdKey Server related constants
const (
	CorrelationIdKey = "correlation_id"
)

// OTName Telemetry related constants
const (
	OTName          = "postgres-perf"
	OTInstanceIDKey = "postgres-perf"
	OTVersion       = "1.0"
	OTSchema        = "/v1"
	OTTenant        = "postgres-perf-id"
)

const (
	SqlQuery = `SELECT 
       time_bucket('1 minute', ts) AS minute,  min(usage), max(usage) AS min
FROM cpu_usage
WHERE host = $1
  AND ts BETWEEN $2 AND $3
GROUP BY minute
ORDER BY minute DESC;`

	StatementName = "statement"

	ExpectedHeaders = "hostname,start_time,end_time"
)

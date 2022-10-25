package configuration

import (
	"postgres-perf/utils"
)

type Configuration struct {
	Swagger CSwagger

	// jaeger
	JaegerEngine string

	// Configuration
	HttpPort int32

	// Internal settings
	CleanupTimeoutSec int32
	UseTelemetry      string
	Development       bool
	GinLogger         bool
	UseSwagger        bool
	Initialized       bool

	// Postgresql
	PostgresqlHost     string
	PostgresqlDatabase string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlPort     int32
	Workers            int32
}

var appConfig Configuration

func AppConfig() *Configuration {
	if appConfig.Initialized == false {
		loadEnvironmentVariables()
		appConfig.Initialized = true
	}
	return &appConfig
}

// loadEnvironmentVariables load env variables
func loadEnvironmentVariables() {

	// jaeger telemetry settings
	appConfig.JaegerEngine = utils.EnvOrDefault("JAEGER_ENGINE_NAME", "http://localhost:14268/api/traces")
	// postgres stuff
	appConfig.PostgresqlHost = utils.EnvOrDefault("POSTGRESQL_HOST", "localhost")
	appConfig.PostgresqlUser = utils.EnvOrDefault("POSTGRESQL_USER", "postgres")
	appConfig.PostgresqlPassword = utils.EnvOrDefault("POSTGRESQL_PASSWORD", "rdsdb")
	appConfig.PostgresqlDatabase = utils.EnvOrDefault("POSTGRESQL_DATABASE", "homework")
	appConfig.PostgresqlPort = utils.EnvOrDefaultInt32("POSTGRESQL_PORT", 5432)
	appConfig.Workers = utils.EnvOrDefaultInt32("WORKERS", 4)
}

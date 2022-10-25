package main

import (
	"context"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"postgres-perf/api"
	"postgres-perf/configuration"
	"postgres-perf/docs"
	"postgres-perf/utils/go-stats/concurrency"
	"postgres-perf/utils/logger"
	"syscall"
	"time"
)

// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @query.collection.format multi
func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	appConfig := configuration.AppConfig()

	pflag.Int32VarP(&appConfig.CleanupTimeoutSec, "timeout", "t", 60, "Time to wait for graceful shutdown on SIGTERM/SIGINT in seconds. Default: 60")
	pflag.Int32VarP(&appConfig.HttpPort, "port", "p", 8080, "TCP port for the HTTP listener to bind to. Default: 8080")
	pflag.BoolVarP(&appConfig.UseSwagger, "swagger", "s", false, "Activate swagger. Do not use this in Production!")
	pflag.BoolVarP(&appConfig.Development, "devel", "d", false, "Start in development mode. Implies --swagger. Do not use this in Production!")
	pflag.BoolVarP(&appConfig.GinLogger, "gin-logger", "g", false, "Activate Gin's logger, for debugging. Do not use this in Production!")
	pflag.StringVarP(&appConfig.UseTelemetry, "telemetry", "r", "", "Activate telemetry local or remote/jaeger")
	pflag.Parse()

	ctx = context.Background()
	ctx, cancel = context.WithCancel(ctx)
	cSignal := make(chan os.Signal)
	signal.Notify(cSignal, os.Interrupt, syscall.SIGTERM)

	logger.Init(ctx, appConfig.Development)
	logger.SetCorrelationIdFieldKey(configuration.CorrelationIdKey)
	logger.SetCorrelationIdContextKey(configuration.CorrelationIdKey)
	log := logger.SugaredLogger()
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()
	defer logger.PanicLogger()

	if !appConfig.Development {
		if appConfig.CleanupTimeoutSec < 120 {
			log.Warnf("Cleanup timeout is set to %d seconds which might be too small for production mode!", appConfig.CleanupTimeoutSec)
		}

	}

	if appConfig.Development {
		appConfig.UseSwagger = true
	}

	if appConfig.UseSwagger {
		appConfig.LoadSwaggerConf()
		docs.SwaggerInfo.Title = appConfig.Swagger.Title
		docs.SwaggerInfo.Version = appConfig.Swagger.Version
		docs.SwaggerInfo.BasePath = appConfig.Swagger.BasePath
		docs.SwaggerInfo.Description = appConfig.Swagger.Description
	}
	log.Infof(docs.SwaggerInfo.BasePath)

	go func() {
		<-cSignal
		log.Warnf("SIGTERM received, attempting graceful exit.")
		cancel()
	}()

	log.Info("Starting webapi handler")
	concurrency.GlobalWaitGroup.Add(1)
	go api.StartGin(ctx)

	<-ctx.Done()

	log.Infof("Graceful shutdown initiated. Waiting for %d seconds before forced exit.", appConfig.CleanupTimeoutSec)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(appConfig.CleanupTimeoutSec))
	go func() {
		concurrency.GlobalWaitGroup.Wait()
		log.Infof("Cleanup done.")
		cancel()
	}()
	<-ctx.Done()
	log.Info("Exiting.")
}

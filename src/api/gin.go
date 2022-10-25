package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"postgres-perf/api/handlers"
	"postgres-perf/api/middleware"
	"postgres-perf/configuration"
	"postgres-perf/utils/go-stats/concurrency"
	"postgres-perf/utils/logger"
	"postgres-perf/utils/tracer"
	"time"
)

const httpServerShutdownGracePeriodSeconds = 20

func StartGin(ctx context.Context) {
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	if conf.UseTelemetry == "remote" {

		log.Infof("Jaeger Telemetry enabled")
		// init tracer jaeger
		tp, err := tracer.InitTracerJaeger(ctx, conf.JaegerEngine, configuration.OTName, configuration.OTInstanceIDKey, configuration.OTTenant)
		if err != nil {
			log.Fatal(err)
		}

		concurrency.GlobalWaitGroup.Add(1)
		defer func() {
			defer concurrency.GlobalWaitGroup.Done()
			localCtx, localCancel := context.WithTimeout(context.Background(), httpServerShutdownGracePeriodSeconds*time.Second)
			defer localCancel()
			if err := tp.Shutdown(localCtx); err != nil {
				log.Errorf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	if conf.UseTelemetry == "local" {

		log.Infof("Stdout Telemetry enabled")
		// init tracer jaeger
		tp, err := tracer.InitTracerStdout(ctx)
		if err != nil {
			log.Fatal(err)
		}

		concurrency.GlobalWaitGroup.Add(1)
		defer func() {
			defer concurrency.GlobalWaitGroup.Done()
			localCtx, localCancel := context.WithTimeout(context.Background(), httpServerShutdownGracePeriodSeconds*time.Second)
			defer localCancel()
			if err := tp.Shutdown(localCtx); err != nil {
				log.Errorf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	// Set up gin
	log.Debugf("Setting up Gin")
	if !conf.GinLogger {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Set up the middleware
	if conf.GinLogger {
		log.Warnf("Gin's logger is active! Logs will be unstructured!")
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.Use(middleware.CorrelationId())
	router.Use(otelgin.Middleware("postgres-perf"))

	// Set up the groups
	userAPI := router.Group("/v1")
	{
		// sample sql query
		userAPI.POST("/upload", handlers.SaveFileHandler)
	}

	// Activate swagger if configured
	if conf.UseSwagger {
		log.Infof("Swagger is active, enabling endpoints")
		url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	// Set up the listener
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.HttpPort),
		Handler: router,
	}

	// Start the HTTP Server
	go func() {
		log.Infof("Listening on port %d", conf.HttpPort)
		if err := httpSrv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Unrecoverable HTTP Server failure: %s", err.Error())
			}
		}
	}()

	// Block until SIGTERM/SIGINT
	<-ctx.Done()

	// Clean up and shutdown the HTTP server
	cleanCtx, cancel := context.WithTimeout(context.Background(), httpServerShutdownGracePeriodSeconds*time.Second)
	defer cancel()
	log.Infof("Attempting to shutdown the HTTP server with a timeout of %d seconds", httpServerShutdownGracePeriodSeconds)
	if err := httpSrv.Shutdown(cleanCtx); err != nil {
		log.Errorf("HTTP server failed to shutdown gracefully: %s", err.Error())
	} else {
		log.Infof("HTTP Server was shutdown successfully")
	}
}

# PgPerf API

A go application capable of pgsql perf test

- upload_file

# Building

After cloning the repo, you can:


## Option 0: Docker compose

```shell
cd $GOPATH/src/go-postgres-perf
docker-compose up --build
```

__!!!NB!!!__ this should build all dependencies and populate the db

Open a browser and [ http://localhost:8080/swagger/index.html#/default/uploader ]

Or, curl

```shell
curl -X POST http://localhost:8080/v1/upload -F "file=@query_params.csv" -H "Content-Type: multipart/form-data"
```

This should return

```shell
{"code":200,"message":"Success","data":{"aggregate":"3.02450159s","maxtime":"60.586111ms","mintime":"3.265855ms","totalqueries":200,"totaltime":"615.07478ms"},"id":"9144d31c-5d90-4311-ac8d-bebfe96d5f44"}
```


## Option 1: Build locally

__!!!NB!!!__ First, get the swagger/swag Go application

```shell
go get -u github.com/swaggo/swag/cmd/swag
```

```shell
cd $GOPATH/src/go-postgres-perf/src
go get -d -v ./...
swag init && go build main.go && ./main -s -d -p 8080 -r local
```



## Options 2: Build the docker image

```shell
cd $GOPATH/src/go-postgres-perf
docker build -t go-postgres-perf-agent -f Dockerfile .
```

# Running it locally

If you built it locally, then execute the binary and pass necessary command-line parameters.

```shell
./go-postgres-perf [opts]
```

## Examples

### To enable swagger:

```shell
./go-postgres-perf -d -s -p 8080
```

(open browser: http://localhost:8080/swagger/index.html#/default/uploader)

### To enable telemetry:

```shell
./go-postgres-perf -r remote/local
```

(__!!!NB!!!__ user remote to print stuff at stdout)

If you built the docker image, you may do the same. Remember to expose the necessary ports.

```shell
docker run -it -d --name pgsqlrunme -p 8080:8080 go-postgres-perf-agent
```

# Command line parameters

You may specify a number of command-line parameters which change the behavior of the application

| Short | Long | Default | Usable in prod | Description |
|-----|-----|-----|-----|-----|
| -t | --timeout | 60 | Yes | Time to wait for graceful shutdown on SIGTERM/SIGINT in seconds |
| -p | --port | 8080 | Yes | TCP port for the HTTP listener to bind to |
| -s | --swagger | | No | Activate swagger. Do not use this in Production! |
| -d | --devel | | No | Start in development mode. Implies --swagger. Do not use this in Production! |
| -g | --gin-logger| | No | Activate Gin's logger, for debugging. **Warning**: This breaks structured logging. Do not use this in Production! |
| -r | --telemetry| | Yes | Enable telemetry. Values accepted: local (for local telemetry) remote(for jaeger telemetry)|


# Environment variables and options

## Workers

number of workers

```shell
export WORKERS=4 
```

## Pgsql stuff

```shell
export POSTGRESQL_HOST=localhost
export POSTGRESQL_USER=postgres
export POSTGRESQL_PASSWORD=rdsdb
export POSTGRESQL_DATABASE=homework
export POSTGRESQL_PORT=5432
```

__!!NB!!__ If not set the API will use default variables from config

## Telemetry env vars (jaeger)
For telemetry using jaeger app required jaeger endpoint (if not set, default local host will be used)

```
appConfig.JaegerEngine = utils.EnvOrDefault("JAEGER_ENGINE_NAME", "http://localhost:14268/api/traces")
```

```shell
export JAEGER_ENGINE_NAME=http://localhost:14268/api/traces
```

Starting local jaeger server

```shell
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
```

(open browser http://localhost:16686/)


# API Docs

All endpoints are documented using [swagger](http://localhost:8080/swagger/index.html)

# Swagger for development

First, get the swagger/swag Go application

```shell
go get -u github.com/swaggo/swag/cmd/swag
```

Now, every time you make a change to the Swagger headers, you will need to regenerate the docs

```shell
cd src
swag init
```

If you have a problem with your headers or mappings, you will get an error describing what's wrong. You **must** fix these before committing the code!

**Note:** The docs are regenerated automatically when building the docker image.


# API request sample

## Upload csv sample file

```shell
curl -X POST http://localhost:8080/v1/upload -F "file=@query_params.csv" -H "Content-Type: multipart/form-data"
```



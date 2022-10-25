FROM golang:1.18 as build
WORKDIR /src
ADD src /src
RUN go get -d -v ./... \
    && go build -o /postgres-perf
FROM gcr.io/distroless/base
USER 1000
EXPOSE 8080
ENTRYPOINT [ "/postgres-perf" ,"-s", "-d", "-r", "remote"]
COPY --from=build /postgres-perf /
ADD /src/swagger.yaml /swagger.yaml


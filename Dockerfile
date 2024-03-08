FROM golang:1.22-bullseye AS build
COPY . /go/src/hadoop-http-jmx-exporter/
WORKDIR /go/src/hadoop-http-jmx-exporter
RUN go build -o /bin/hadoop-http-jmx-exporter /go/src/hadoop-http-jmx-exporter/cmd/main.go

FROM ramencloud/debian:bullseye-slim
COPY --from=build /bin/hadoop-http-jmx-exporter /bin/hadoop-http-jmx-exporter
ENTRYPOINT ["/bin/hadoop-http-jmx-exporter"]

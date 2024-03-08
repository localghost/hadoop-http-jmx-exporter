MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MAKEFILE_PATH))

run:
	@go run cmd/main.go

build:
	@mkdir -p $(CURRENT_DIR)/bin
	@go build -o bin/hadoop-http-jmx-exporter cmd/main.go

docker:
	@docker build -t hadoop-http-jmx-exporter .

prom:
	@docker run \
          -p 9090:9090 \
          -v $(CURRENT_DIR)/prom_config.yml:/etc/prometheus/prometheus.yml \
          prom/prometheus
tidy:
	@go mod tidy

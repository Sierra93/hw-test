# Переменные путей
PROJECT_DIR := ./hw12_13_14_15_16_calendar
BIN_DIR := ./bin

GIT_HASH := $(shell git log --format="%h" -n 1)
# Исправлены кавычки для корректной передачи строк в Go
LDFLAGS := -X 'main.release=develop' -X 'main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X 'main.gitHash=$(GIT_HASH)'

# Сборка всех компонентов
build: build-calendar build-scheduler build-storer

build-calendar:
	go build -v -o $(BIN_DIR)/calendar -ldflags "$(LDFLAGS)" $(PROJECT_DIR)/cmd/calendar

build-scheduler:
	go build -v -o $(BIN_DIR)/scheduler -ldflags "$(LDFLAGS)" $(PROJECT_DIR)/cmd/scheduler

build-storer:
	go build -v -o $(BIN_DIR)/storer -ldflags "$(LDFLAGS)" $(PROJECT_DIR)/cmd/storer

# Тесты (исправлены пути к пакетам)
test:
	go test -race $(PROJECT_DIR)/internal/logger $(PROJECT_DIR)/internal/storage/memory $(PROJECT_DIR)/internal/server/http

# Линтер (запуск из корня с учетом всех подпапок)
lint: install-lint-deps
	golangci-lint run $(PROJECT_DIR)/...

# Запуск локально
run-calendar: build-calendar
	$(BIN_DIR)/calendar -config $(PROJECT_DIR)/configs/config.yaml

run-scheduler: build-scheduler
	$(BIN_DIR)/scheduler -config $(PROJECT_DIR)/configs/scheduler-config.yaml

run-storer: build-storer
	$(BIN_DIR)/storer -config $(PROJECT_DIR)/configs/storer-config.yaml

# Docker сборка (исправлены пути к Dockerfile и контексту)
build-calendar-img:
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t calendar:develop -f $(PROJECT_DIR)/build/calendar/Dockerfile .

build-scheduler-img:
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t scheduler:develop -f $(PROJECT_DIR)/build/scheduler/Dockerfile .

build-storer-img:
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t storer:develop -f $(PROJECT_DIR)/build/storer/Dockerfile .

# Утилиты
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.5

target:
	go generate $(PROJECT_DIR)/internal/server/http

down:
	docker-compose -f $(PROJECT_DIR)/deployments/docker-compose.yaml down

up: down
	docker-compose -f $(PROJECT_DIR)/deployments/docker-compose.yaml up -d --build

.PHONY: build build-calendar build-scheduler build-storer test lint up down target

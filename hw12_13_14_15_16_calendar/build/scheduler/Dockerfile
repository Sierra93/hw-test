FROM golang:1.23 as build

ENV BIN_FILE /opt/scheduler/scheduler-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/scheduler/*

FROM alpine:3.9

ENV BIN_FILE "/opt/scheduler/scheduler-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/scheduler/config.yaml
COPY ./configs/scheduler-config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE}
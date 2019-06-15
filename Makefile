.PHONY: build

APP = Citrix-NetScaler-Exporter
VERSION = 3.2.0
BINARY-LINUX = ${APP}_${VERSION}_amd64

BUILD = $(shell git rev-parse HEAD)

LDFLAGS = -ldflags "-X main.version=%VERSION% -X main.build=%BUILD%"

build:
    GOOS=linux GOOS=amd64 go build -o ${BINARY-LINUX} ${LDFLAGS}

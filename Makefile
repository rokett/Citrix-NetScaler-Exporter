.PHONY: build

APP = Citrix-NetScaler-Exporter
VERSION = 4.0.0
BINARY-LINUX = ${APP}_${VERSION}_Linux_amd64

BUILD_VER = $(shell git rev-parse HEAD)

LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.build=${BUILD_VER} -s -w"

build:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY-LINUX} ${LDFLAGS}

FROM golang:alpine as builder

ENV VERSION="4.1.0"

WORKDIR $GOPATH/src/github.com/rokett
RUN \
    apk add --no-cache git && \
    git clone --branch $VERSION --depth 1 https://github.com/rokett/Citrix-NetScaler-Exporter.git citrix-netscaler-exporter && \
    cd citrix-netscaler-exporter && \
    BUILD=$(git rev-list -1 HEAD) && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-X main.version=$VERSION -X main.build=$BUILD -s -w -extldflags '-static'" -o citrix-netscaler-exporter

FROM scratch
LABEL maintainer="rokett@rokett.me"
COPY --from=builder go/src/github.com/rokett/citrix-netscaler-exporter/citrix-netscaler-exporter /

EXPOSE 9280

ENTRYPOINT ["./citrix-netscaler-exporter"]

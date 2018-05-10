FROM golang:1.9rc2-alpine
MAINTAINER IPI  <ipi@ticketmaster.co.uk>

WORKDIR /go/src
RUN apk update; apk add git
RUN mkdir Citrix-NetScaler-Exporter
COPY . Citrix-NetScaler-Exporter
WORKDIR /go/src/Citrix-NetScaler-Exporter
RUN go get github.com/alecthomas/kingpin
RUN go get github.com/rokett/citrix-netscaler-exporter/netscaler
RUN go build .
RUN cp /go/src/Citrix-NetScaler-Exporter/Citrix-NetScaler-Exporter /
ENTRYPOINT ["/Citrix-NetScaler-Exporter"]

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/rokett/citrix-netscaler-exporter/collector"
)

var (
	app        = "Citrix-NetScaler-Exporter"
	version    string
	build      string
	url        = flag.String("url", "", "Base URL of the NetScaler management interface.  Normally something like https://my-netscaler.something.x")
	username   = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password   = flag.String("password", "", "Password with which to connect to the NetScaler API")
	bindPort   = flag.Int("bind_port", 9280, "Port to bind the exporter endpoint to")
	versionFlg = flag.Bool("version", false, "Display application version")
	ignoreCert = flag.Bool("ignore-cert", false, "Ignore certificate errors; use with caution")
	logger     log.Logger

	nsInstance string
)

func main() {
	flag.Parse()

	if *versionFlg {
		fmt.Printf("%s v%s build %s\n", app, version, build)
		os.Exit(0)
	}

	if *url == "" || *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	nsInstance = strings.TrimLeft(*url, "https://")
	nsInstance = strings.Trim(nsInstance, " /")

	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "app", app, "bind_port", *bindPort, "url", *url, "version", "v"+version, "build", build)

	exporter, _ := collector.NewExporter(*url, *username, *password, *ignoreCert, logger, nsInstance)
	prometheus.MustRegister(exporter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Citrix NetScaler Exporter</title></head>
			<body>
			<h1>Citrix NetScaler Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	http.Handle("/metrics", promhttp.Handler())

	listeningPort := ":" + strconv.Itoa(*bindPort)
	level.Info(logger).Log("msg", "Listening on port "+listeningPort)

	err := http.ListenAndServe(listeningPort, nil)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}

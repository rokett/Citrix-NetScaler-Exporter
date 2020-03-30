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

	"github.com/mckesson/mk-citrix-exporter/collector"
)

var (
	app        = "Citrix-NetScaler-Exporter"
	version    string
	build      string
	username   = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password   = flag.String("password", "", "Password with which to connect to the NetScaler API")
	bindPort   = flag.Int("bind_port", 9280, "Port to bind the exporter endpoint to")
	versionFlg = flag.Bool("version", false, "Display application version")
	debugFlg   = flag.Bool("debug", false, "Enable debug logging?")
	logger     log.Logger

	nsInstance string
)

func init() {
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "app", app, "bind_port", *bindPort, "version", "v"+version, "build", build)
}

func main() {
	flag.Parse()

	if *versionFlg {
		fmt.Printf("%s v%s build %s\n", app, version, build)
		os.Exit(0)
	}

	if *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>Citrix NetScaler Exporter</title></head>
				<style>
				label{
				display:inline-block;
				width:75px;
				}
				form label {
				margin: 10px;
				}
				form input {
				margin: 10px;
				}
				</style>
				<body>
				<h1>Citrix NetScaler Exporter</h1>
				<form action="/netscaler">
				<label>Target:</label> <input type="text" name="target" placeholder="https://netscaler.domain.tld"> <br>
				<p>Ignore certificate check?</p>
				<input type="radio" id="yes" name="ignore-cert" value="yes">
				<label for="yes">Yes</label>
				<input type="radio" id="no" name="ignore-cert" value="no" checked>
  				<label for="no">No</label>
				<br>
				<input type="submit" value="Submit">
				</form>
				</body>
				</html>`))
	})

	http.HandleFunc("/netscaler", handler)
	http.Handle("/metrics", promhttp.Handler())

	listeningPort := ":" + strconv.Itoa(*bindPort)
	level.Info(logger).Log("msg", "Listening on port "+listeningPort)

	err := http.ListenAndServe(listeningPort, nil)
	if err != nil {
		level.Error(logger).Log("msg", err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		return
	}

	ignoreCertCheck := false
	if strings.ToLower(r.URL.Query().Get("ignore-cert")) == "yes" {
		ignoreCertCheck = true
	}

	nsInstance = strings.TrimLeft(target, "https://")
	nsInstance = strings.TrimLeft(nsInstance, "http://")
	nsInstance = strings.Trim(nsInstance, " /")

	if *debugFlg {
		level.Debug(logger).Log("msg", "scraping target", "target", target)
	}

	exporter, err := collector.NewExporter(target, *username, *password, ignoreCertCheck, logger, nsInstance)
	if err != nil {
		http.Error(w, "Error creating exporter"+err.Error(), 400)
		level.Error(logger).Log("msg", err)
		return
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

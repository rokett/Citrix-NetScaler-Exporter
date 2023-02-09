package collector

import (
	"strconv"

	"github.com/rokett/citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	aaaAuthSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aaa_auth_success",
			Help: "Count of authentication successes",
		},
		[]string{
			"ns_instance",
		},
	)

	aaaAuthFail = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aaa_auth_fail",
			Help: "Count of authentication failures",
		},
		[]string{
			"ns_instance",
		},
	)

	aaaAuthOnlyHTTPSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aaa_auth_only_http_success",
			Help: "Count of HTTP connections that succeeded authorisation",
		},
		[]string{
			"ns_instance",
		},
	)

	aaaAuthOnlyHTTPFail = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aaa_auth_only_http_fail",
			Help: "Count of HTTP connections that failed authorisation",
		},
		[]string{
			"ns_instance",
		},
	)

	aaaCurIcaSessions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aaa_current_ica_sessions",
			Help: "Count of current Basic ICA only sessions",
		},
		[]string{
			"ns_instance",
		},
	)

	aaaCurIcaOnlyConn = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aaa_current_ica_only_connections",
			Help: "Count of current Basic ICA only connections",
		},
		[]string{
			"ns_instance",
		},
	)
)

func (e *Exporter) collectAaaAuthSuccess(ns netscaler.NSAPIResponse) {
	e.aaaAuthSuccess.Reset()

	val, _ := strconv.ParseFloat(ns.AAAStats.AuthSuccess, 64)
	e.aaaAuthSuccess.WithLabelValues(e.nsInstance).Set(val)
}

func (e *Exporter) collectAaaAuthFail(ns netscaler.NSAPIResponse) {
	e.aaaAuthFail.Reset()

	val, _ := strconv.ParseFloat(ns.AAAStats.AuthFail, 64)
	e.aaaAuthFail.WithLabelValues(e.nsInstance).Set(val)
}

func (e *Exporter) collectAaaAuthOnlyHTTPSuccess(ns netscaler.NSAPIResponse) {
	e.aaaAuthOnlyHTTPSuccess.Reset()

	val, _ := strconv.ParseFloat(ns.AAAStats.AuthOnlyHTTPSuccess, 64)
	e.aaaAuthOnlyHTTPSuccess.WithLabelValues(e.nsInstance).Set(val)
}

func (e *Exporter) collectAaaAuthOnlyHTTPFail(ns netscaler.NSAPIResponse) {
	e.aaaAuthOnlyHTTPFail.Reset()

	val, _ := strconv.ParseFloat(ns.AAAStats.AuthOnlyHTTPFail, 64)
	e.aaaAuthOnlyHTTPFail.WithLabelValues(e.nsInstance).Set(val)
}

func (e *Exporter) collectAaaCurIcaSessions(ns netscaler.NSAPIResponse) {
	e.aaaCurIcaSessions.Reset()

	val, _ := strconv.ParseFloat(ns.AAAStats.CurrentIcaSessions, 64)
	e.aaaCurIcaSessions.WithLabelValues(e.nsInstance).Set(val)
}

func (e *Exporter) collectAaaCurIcaOnlyConn(ns netscaler.NSAPIResponse) {
	e.aaaCurIcaOnlyConn.Reset()

	val, _ := strconv.ParseFloat(ns.AAAStats.CurrentIcaOnlyConnections, 64)
	e.aaaCurIcaOnlyConn.WithLabelValues(e.nsInstance).Set(val)
}

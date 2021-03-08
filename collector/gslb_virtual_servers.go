package collector

import (
	"strconv"
	"github.com/encoretechnologies/citrix-netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	gslbVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_state",
			Help: "Current state of the server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	gslbVirtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)
)

func (e *Exporter) collectGSLBVirtualServerState(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersState.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		state := 0.0

		if vs.State == "UP" {
			state = 1.0
		}

		e.gslbVirtualServersState.WithLabelValues(e.nsInstance, vs.Name).Set(state)
	}
}

func (e *Exporter) collectGSLBVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersHealth.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		health, _ := strconv.ParseFloat(vs.Health, 64)
		e.gslbVirtualServersHealth.WithLabelValues(e.nsInstance, vs.Name).Set(health)
	}
}

func (e *Exporter) collectGSLBVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersInactiveServices.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		inactiveServices, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.gslbVirtualServersInactiveServices.WithLabelValues(e.nsInstance, vs.Name).Set(inactiveServices)
	}
}

func (e *Exporter) collectGSLBVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersActiveServices.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		activeServices, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.gslbVirtualServersActiveServices.WithLabelValues(e.nsInstance, vs.Name).Set(activeServices)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalHits.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.gslbVirtualServersTotalHits.WithLabelValues(e.nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalRequests.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.gslbVirtualServersTotalRequests.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalResponses.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.gslbVirtualServersTotalResponses.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.gslbVirtualServersTotalRequestBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectGSLBVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.gslbVirtualServersTotalResponseBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.gslbVirtualServersCurrentClientConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectGSLBVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.gslbVirtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.GSLBVirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.gslbVirtualServersCurrentServerConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentServerConnections)
	}
}

package collector

import (
	"strconv"

	"github.com/rokett/citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	gslbServicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_state",
			Help: "Current state of the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_responses",
			Help: "Total number of responses received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_total_response_bytes",
			Help: "Total number of response bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_client_connections",
			Help: "Number of current client connections",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_server_connections",
			Help: "Number of current connections to the actual servers",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_established_connections",
			Help: "Number of server connections in ESTABLISHED state",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gslb_service_current_load",
			Help: "Load on the service that is calculated from the bound load based monitor",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	gslbServicesVirtualServerServiceHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gslb_service_virtual_server_service_hits",
			Help: "Number of times that the service has been provided",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)
)

func (e *Exporter) collectGSLBServicesState(ns netscaler.NSAPIResponse) {
	e.gslbServicesState.Reset()

	for _, service := range ns.GSLBServiceStats {
		state := 0.0

		if service.State == "UP" {
			state = 1.0
		}

		e.gslbServicesState.WithLabelValues(e.nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalRequests.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.gslbServicesTotalRequests.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalResponses.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.gslbServicesTotalResponses.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalRequestBytes.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.gslbServicesTotalRequestBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.gslbServicesTotalResponseBytes.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.gslbServicesTotalResponseBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentClientConns.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.gslbServicesCurrentClientConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentServerConns.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.gslbServicesCurrentServerConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.gslbServicesEstablishedConnections.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.EstablishedConnections, 64)
		e.gslbServicesEstablishedConnections.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.gslbServicesCurrentLoad.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.gslbServicesCurrentLoad.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectGSLBServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.gslbServicesVirtualServerServiceHits.Reset()

	for _, service := range ns.GSLBServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.gslbServicesVirtualServerServiceHits.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

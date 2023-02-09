package collector

import (
	"strconv"

	"github.com/rokett/citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	servicesThroughput = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_throughput",
			Help: "Number of bytes received or sent by this service (Mbps)",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server. TTFB is the time interval between sending the request packet to a service and receiving the first response from the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_state",
			Help: "Current state of the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_responses",
			Help: "Total number of responses received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_response_bytes",
			Help: "Total number of response bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_client_connections",
			Help: "Number of current client connections",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_surge_count",
			Help: "Number of requests in the surge queue",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_server_connections",
			Help: "Number of current connections to the actual servers",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_max_clients",
			Help: "Maximum open connections allowed on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_load",
			Help: "Load on the service that is calculated from the bound load based monitor",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesVirtualServerServiceHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_virtual_server_service_hits",
			Help: "Number of times that the service has been provided",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesActiveTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_active_transactions",
			Help: "Number of active transactions handled by this service. (Including those in the surge queue.) Active Transaction means number of transactions currently served by the server including those waiting in the SurgeQ",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)
)

func (e *Exporter) collectServicesThroughput(ns netscaler.NSAPIResponse) {
	e.servicesThroughput.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.Throughput, 64)
		e.servicesThroughput.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesAvgTTFB(ns netscaler.NSAPIResponse) {
	e.servicesAvgTTFB.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.AvgTimeToFirstByte, 64)
		e.servicesAvgTTFB.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesState(ns netscaler.NSAPIResponse) {
	e.servicesState.Reset()

	for _, service := range ns.ServiceStats {
		state := 0.0

		if service.State == "UP" {
			state = 1.0
		}

		e.servicesState.WithLabelValues(e.nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequests.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.servicesTotalRequests.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponses.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.servicesTotalResponses.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequestBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.servicesTotalRequestBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponseBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.servicesTotalResponseBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentClientConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.servicesCurrentClientConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesSurgeCount(ns netscaler.NSAPIResponse) {
	e.servicesSurgeCount.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.SurgeCount, 64)
		e.servicesSurgeCount.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentServerConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.servicesCurrentServerConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.servicesServerEstablishedConnections.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServerEstablishedConnections, 64)
		e.servicesServerEstablishedConnections.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentReusePool(ns netscaler.NSAPIResponse) {
	e.servicesCurrentReusePool.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentReusePool, 64)
		e.servicesCurrentReusePool.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesMaxClients(ns netscaler.NSAPIResponse) {
	e.servicesMaxClients.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.MaxClients, 64)
		e.servicesMaxClients.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.servicesCurrentLoad.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.servicesCurrentLoad.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.servicesVirtualServerServiceHits.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.servicesVirtualServerServiceHits.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesActiveTransactions(ns netscaler.NSAPIResponse) {
	e.servicesActiveTransactions.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ActiveTransactions, 64)
		e.servicesActiveTransactions.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

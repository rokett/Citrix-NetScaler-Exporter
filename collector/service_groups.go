package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rokett/citrix-netscaler-exporter/netscaler"
)

var (
	serviceGroupsState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_state",
			Help: "Current state of the server",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_responses",
			Help: "Number of responses received on this service.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_response_bytes",
			Help: "Number of response bytes received by this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_client_connections",
			Help: "Number of current client connections.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_surge_count",
			Help: "Number of requests in the surge queue.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_server_connections",
			Help: "Number of current connections to the actual servers behind the virtual server.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_max_clients",
			Help: "Maximum open connections allowed on this service.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)
)

func (e *Exporter) collectServiceGroupsState(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsState.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		state := 0.0

		if sg.State == "UP" {
			state = 1.0
		}

		e.serviceGroupsState.WithLabelValues(e.nsInstance, sgName, servername).Set(state)
	}
}

func (e *Exporter) collectServiceGroupsAvgTTFB(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsAvgTTFB.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.AvgTimeToFirstByte, 64)
		e.serviceGroupsAvgTTFB.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalRequests(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalRequests.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalRequests, 64)
		e.serviceGroupsTotalRequests.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalResponses(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalResponses.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalResponses, 64)
		e.serviceGroupsTotalResponses.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalRequestBytes(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalRequestBytes.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalRequestBytes, 64)
		e.serviceGroupsTotalRequestBytes.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalResponseBytes(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalResponseBytes.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalResponseBytes, 64)
		e.serviceGroupsTotalResponseBytes.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentClientConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentClientConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentClientConnections, 64)
		e.serviceGroupsCurrentClientConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsSurgeCount(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsSurgeCount.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.SurgeCount, 64)
		e.serviceGroupsSurgeCount.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentServerConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentServerConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentServerConnections, 64)
		e.serviceGroupsCurrentServerConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsServerEstablishedConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsServerEstablishedConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.ServerEstablishedConnections, 64)
		e.serviceGroupsServerEstablishedConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentReusePool(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentReusePool.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentReusePool, 64)
		e.serviceGroupsCurrentReusePool.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsMaxClients(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsMaxClients.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.MaxClients, 64)
		e.serviceGroupsMaxClients.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

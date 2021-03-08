package collector

import (
	"strconv"
	"github.com/encoretechnologies/citrix-netscaler-exporter/netscaler"
	"github.com/prometheus/client_golang/prometheus"
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
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
			"port",
		},
	)
)

func (e *Exporter) collectServiceGroupsState(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsState.Reset()

	state := 0.0

	if sg.State == "UP" {
		state = 1.0
	}

	port := strconv.Itoa(sg.PrimaryPort)

	e.serviceGroupsState.WithLabelValues(e.nsInstance, sgName, servername, port).Set(state)
}

func (e *Exporter) collectServiceGroupsAvgTTFB(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsAvgTTFB.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.AvgTimeToFirstByte, 64)
	e.serviceGroupsAvgTTFB.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalRequests(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalRequests.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.TotalRequests, 64)
	e.serviceGroupsTotalRequests.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponses(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalResponses.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.TotalResponses, 64)
	e.serviceGroupsTotalResponses.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalRequestBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalRequestBytes.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.TotalRequestBytes, 64)
	e.serviceGroupsTotalRequestBytes.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsTotalResponseBytes(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsTotalResponseBytes.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.TotalResponseBytes, 64)
	e.serviceGroupsTotalResponseBytes.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentClientConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsCurrentClientConnections.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.CurrentClientConnections, 64)
	e.serviceGroupsCurrentClientConnections.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsSurgeCount(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsSurgeCount.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.SurgeCount, 64)
	e.serviceGroupsSurgeCount.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentServerConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsCurrentServerConnections.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.CurrentServerConnections, 64)
	e.serviceGroupsCurrentServerConnections.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsServerEstablishedConnections(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsServerEstablishedConnections.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.ServerEstablishedConnections, 64)
	e.serviceGroupsServerEstablishedConnections.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsCurrentReusePool(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsCurrentReusePool.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.CurrentReusePool, 64)
	e.serviceGroupsCurrentReusePool.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

func (e *Exporter) collectServiceGroupsMaxClients(sg netscaler.ServiceGroupMemberStats, sgName string, servername string) {
	e.serviceGroupsMaxClients.Reset()

	port := strconv.Itoa(sg.PrimaryPort)

	val, _ := strconv.ParseFloat(sg.MaxClients, 64)
	e.serviceGroupsMaxClients.WithLabelValues(e.nsInstance, sgName, servername, port).Set(val)
}

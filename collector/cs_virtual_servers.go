package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rokett/citrix-netscaler-exporter/netscaler"
)

var (
	csVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_state",
			Help: "Current state of the server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_established_connections",
			Help: "Number of client connections in ESTABLISHED state.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalPacketsReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_packets_received",
			Help: "Total number of packets received",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalPacketsSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_packets_sent",
			Help: "Total number of packets sent.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalSpillovers = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_spillovers",
			Help: "Number of times vserver experienced spill over.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersDeferredRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_deferred_requests",
			Help: "Number of deferred request on this vserver",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersNumberInvalidRequestResponse = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_number_invalid_request_response",
			Help: "Number invalid requests/responses on this vserver",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersNumberInvalidRequestResponseDropped = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_number_invalid_request_response_dropped",
			Help: "Number invalid requests/responses dropped on this vserver",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersTotalVServerDownBackupHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cs_virtual_servers_total_vserver_down_backup_hits",
			Help: "Number of times traffic was diverted to backup vserver since primary vserver was DOWN.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentMultipathSessions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_multipath_sessions",
			Help: "Current Multipath TCP sessions",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	csVirtualServersCurrentMultipathSubflows = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cs_virtual_servers_current_multipath_subflows",
			Help: "Current Multipath TCP subflows",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)
)

func (e *Exporter) collectCSVirtualServerState(ns netscaler.NSAPIResponse) {
	e.csVirtualServersState.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		state := 0.0

		if vs.State == "UP" {
			state = 1.0
		}

		e.csVirtualServersState.WithLabelValues(e.nsInstance, vs.Name).Set(state)
	}
}

func (e *Exporter) collectCSVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalHits.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.csVirtualServersTotalHits.WithLabelValues(e.nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalRequests.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.csVirtualServersTotalRequests.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalResponses.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.csVirtualServersTotalResponses.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectCSVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.csVirtualServersTotalRequestBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectCSVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.csVirtualServersTotalResponseBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.csVirtualServersCurrentClientConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.csVirtualServersCurrentServerConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentServerConnections)
	}
}

func (e *Exporter) collectCSVirtualServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.csVirtualServersEstablishedConnections.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		EstablishedConnections, _ := strconv.ParseFloat(vs.EstablishedConnections, 64)
		e.csVirtualServersEstablishedConnections.WithLabelValues(e.nsInstance, vs.Name).Set(EstablishedConnections)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsReceived(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalPacketsReceived.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalPacketsReceived, _ := strconv.ParseFloat(vs.TotalPacketsReceived, 64)
		e.csVirtualServersTotalPacketsReceived.WithLabelValues(e.nsInstance, vs.Name).Set(totalPacketsReceived)
	}
}

func (e *Exporter) collectCSVirtualServerTotalPacketsSent(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalPacketsSent.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalPacketsSent, _ := strconv.ParseFloat(vs.TotalPacketsSent, 64)
		e.csVirtualServersTotalPacketsSent.WithLabelValues(e.nsInstance, vs.Name).Set(totalPacketsSent)
	}
}

func (e *Exporter) collectCSVirtualServerTotalSpillovers(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalSpillovers.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalSpillovers, _ := strconv.ParseFloat(vs.TotalSpillovers, 64)
		e.csVirtualServersTotalSpillovers.WithLabelValues(e.nsInstance, vs.Name).Set(totalSpillovers)
	}
}

func (e *Exporter) collectCSVirtualServerDeferredRequests(ns netscaler.NSAPIResponse) {
	e.csVirtualServersDeferredRequests.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		deferredRequests, _ := strconv.ParseFloat(vs.DeferredRequests, 64)
		e.csVirtualServersDeferredRequests.WithLabelValues(e.nsInstance, vs.Name).Set(deferredRequests)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponse(ns netscaler.NSAPIResponse) {
	e.csVirtualServersNumberInvalidRequestResponse.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		numberInvalidRequestResponse, _ := strconv.ParseFloat(vs.InvalidRequestResponse, 64)
		e.csVirtualServersNumberInvalidRequestResponse.WithLabelValues(e.nsInstance, vs.Name).Set(numberInvalidRequestResponse)
	}
}

func (e *Exporter) collectCSVirtualServerNumberInvalidRequestResponseDropped(ns netscaler.NSAPIResponse) {
	e.csVirtualServersNumberInvalidRequestResponseDropped.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		numberInvalidRequestResponseDropped, _ := strconv.ParseFloat(vs.InvalidRequestResponseDropped, 64)
		e.csVirtualServersNumberInvalidRequestResponseDropped.WithLabelValues(e.nsInstance, vs.Name).Set(numberInvalidRequestResponseDropped)
	}
}

func (e *Exporter) collectCSVirtualServerTotalVServerDownBackupHits(ns netscaler.NSAPIResponse) {
	e.csVirtualServersTotalVServerDownBackupHits.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		totalVServerDownBackupHits, _ := strconv.ParseFloat(vs.TotalVServerDownBackupHits, 64)
		e.csVirtualServersTotalVServerDownBackupHits.WithLabelValues(e.nsInstance, vs.Name).Set(totalVServerDownBackupHits)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSessions(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentMultipathSessions.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentMultipathSessions, _ := strconv.ParseFloat(vs.CurrentMultipathSessions, 64)
		e.csVirtualServersCurrentMultipathSessions.WithLabelValues(e.nsInstance, vs.Name).Set(currentMultipathSessions)
	}
}

func (e *Exporter) collectCSVirtualServerCurrentMultipathSubflows(ns netscaler.NSAPIResponse) {
	e.csVirtualServersCurrentMultipathSubflows.Reset()

	for _, vs := range ns.CSVirtualServerStats {
		currentMultipathSubflows, _ := strconv.ParseFloat(vs.CurrentMultipathSubflows, 64)
		e.csVirtualServersCurrentMultipathSubflows.WithLabelValues(e.nsInstance, vs.Name).Set(currentMultipathSubflows)
	}
}

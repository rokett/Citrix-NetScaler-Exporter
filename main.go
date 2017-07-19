package main

import (
	"Citrix-NetScaler-Exporter/netscaler"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
)

var (
	url      = flag.String("url", "", "Base URL of the NetScaler management interface.  Normally something like https://my-netscaler.something.x")
	username = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password = flag.String("password", "", "Password with which to connect to the NetScaler API")
	bindPort = flag.Int("bind_port", 9280, "Port to bind the exporter endpoint to")

	nsInstance string

	mgmtCPUUsage = prometheus.NewDesc(
		"mgmt_cpu_usage",
		"Current CPU utilisation for management",
		[]string{
			"ns_instance",
		},
		nil,
	)

	pktCPUUsage = prometheus.NewDesc(
		"pkt_cpu_usage",
		"Current CPU utilisation for packet engines, excluding management",
		[]string{
			"ns_instance",
		},
		nil,
	)

	memUsage = prometheus.NewDesc(
		"mem_usage",
		"Current memory utilisation",
		[]string{
			"ns_instance",
		},
		nil,
	)

	flashPartitionUsage = prometheus.NewDesc(
		"flash_partition_usage",
		"Used space in /flash partition of the disk, as a percentage.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	varPartitionUsage = prometheus.NewDesc(
		"var_partition_usage",
		"Used space in /var partition of the disk, as a percentage. ",
		[]string{
			"ns_instance",
		},
		nil,
	)

	rxMbPerSec = prometheus.NewDesc(
		"received_mb_per_second",
		"Number of Megabits received by the NetScaler appliance per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	txMbPerSec = prometheus.NewDesc(
		"transmit_mb_per_second",
		"Number of Megabits transmitted by the NetScaler appliance per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpRequestsRate = prometheus.NewDesc(
		"http_requests_rate",
		"HTTP requests received per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpResponsesRate = prometheus.NewDesc(
		"http_responses_rate",
		"HTTP requests sent per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentClientConnections = prometheus.NewDesc(
		"tcp_current_client_connections",
		"Client connections, including connections in the Opening, Established, and Closing state.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentClientConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_client_connections_established",
		"Current client connections in the Established state, which indicates that data transfer can occur between the NetScaler and the client.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentServerConnections = prometheus.NewDesc(
		"tcp_current_server_connections",
		"Server connections, including connections in the Opening, Established, and Closing state.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentServerConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_server_connections_established",
		"Current server connections in the Established state, which indicates that data transfer can occur between the NetScaler and the server.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	interfacesRxBytesPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_bytes_per_second",
			Help: "Number of bytes received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxBytesPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_bytes_per_second",
			Help: "Number of bytes transmitted per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesRxPacketsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_packets_per_second",
			Help: "Number of packets received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxPacketsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_packets_per_second",
			Help: "Number of packets transmitted per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsRxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_received_per_second",
			Help: "Number of bytes received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsTxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_transmitted_per_second",
			Help: "Number of jumbo packets transmitted per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesErrorPacketsRxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_error_packets_received_per_second",
			Help: "Number of error packets received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	virtualServersWaitingRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_waiting_requests",
			Help: "Number of requests waiting on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersHitsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_hits_rate",
			Help: "Number of hits/second to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersRequestsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_requests_rate",
			Help: "Number of requests/second to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersReponsesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_responses_rate",
			Help: "Number of responses/second from a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersRequestBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_request_bytes_rate",
			Help: "Number of request bytes/second to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersReponseBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_reponse_bytes_rate",
			Help: "Number of response bytes/second from a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

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

	servicesThroughputRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_throughput_rate",
			Help: "Rate (/s) counter for throughput",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server.TTFB is the time interval between sending the request packet to a service and receiving the first response from the service",
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

	servicesRequestsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_request_rate",
			Help: "Rate (/s) counter for totalrequests",
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

	servicesResponsesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_responses_rate",
			Help: "Rate (/s) counter for totalresponses",
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

	servicesRequestBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_request_bytes_rate",
			Help: "Rate (/s) counter for totalrequestbytes",
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

	servicesResponseBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_response_bytes_rate",
			Help: "Rate (/s) counter for totalresponsebytes",
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

	servicesVirtualServerServiceHitsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_virtual_server_service_hits_rate",
			Help: "Rate (/s) counter for vsvrservicehits",
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

// Exporter represents the metrics exported to Prometheus
type Exporter struct {
	mgmtCPUUsage                           *prometheus.Desc
	memUsage                               *prometheus.Desc
	pktCPUUsage                            *prometheus.Desc
	flashPartitionUsage                    *prometheus.Desc
	varPartitionUsage                      *prometheus.Desc
	rxMbPerSec                             *prometheus.Desc
	txMbPerSec                             *prometheus.Desc
	httpRequestsRate                       *prometheus.Desc
	httpResponsesRate                      *prometheus.Desc
	tcpCurrentClientConnections            *prometheus.Desc
	tcpCurrentClientConnectionsEstablished *prometheus.Desc
	tcpCurrentServerConnections            *prometheus.Desc
	tcpCurrentServerConnectionsEstablished *prometheus.Desc
	interfacesRxBytesPerSecond             *prometheus.GaugeVec
	interfacesTxBytesPerSecond             *prometheus.GaugeVec
	interfacesRxPacketsPerSecond           *prometheus.GaugeVec
	interfacesTxPacketsPerSecond           *prometheus.GaugeVec
	interfacesJumboPacketsRxPerSecond      *prometheus.GaugeVec
	interfacesJumboPacketsTxPerSecond      *prometheus.GaugeVec
	interfacesErrorPacketsRxPerSecond      *prometheus.GaugeVec
	virtualServersWaitingRequests          *prometheus.GaugeVec
	virtualServersHealth                   *prometheus.GaugeVec
	virtualServersInactiveServices         *prometheus.GaugeVec
	virtualServersActiveServices           *prometheus.GaugeVec
	virtualServersTotalHits                *prometheus.CounterVec
	virtualServersHitsRate                 *prometheus.GaugeVec
	virtualServersTotalRequests            *prometheus.CounterVec
	virtualServersRequestsRate             *prometheus.GaugeVec
	virtualServersTotalResponses           *prometheus.CounterVec
	virtualServersReponsesRate             *prometheus.GaugeVec
	virtualServersTotalRequestBytes        *prometheus.CounterVec
	virtualServersRequestBytesRate         *prometheus.GaugeVec
	virtualServersTotalResponseBytes       *prometheus.CounterVec
	virtualServersReponseBytesRate         *prometheus.GaugeVec
	virtualServersCurrentClientConnections *prometheus.GaugeVec
	virtualServersCurrentServerConnections *prometheus.GaugeVec
	servicesThroughput                     *prometheus.CounterVec
	servicesThroughputRate                 *prometheus.GaugeVec
	servicesAvgTTFB                        *prometheus.GaugeVec
	servicesState                          *prometheus.GaugeVec
	servicesTotalRequests                  *prometheus.CounterVec
	servicesRequestsRate                   *prometheus.GaugeVec
	servicesTotalResponses                 *prometheus.CounterVec
	servicesResponsesRate                  *prometheus.GaugeVec
	servicesTotalRequestBytes              *prometheus.CounterVec
	servicesRequestBytesRate               *prometheus.GaugeVec
	servicesTotalResponseBytes             *prometheus.CounterVec
	servicesResponseBytesRate              *prometheus.GaugeVec
	servicesCurrentClientConns             *prometheus.GaugeVec
	servicesSurgeCount                     *prometheus.GaugeVec
	servicesCurrentServerConns             *prometheus.GaugeVec
	servicesServerEstablishedConnections   *prometheus.GaugeVec
	servicesCurrentReusePool               *prometheus.GaugeVec
	servicesMaxClients                     *prometheus.GaugeVec
	servicesCurrentLoad                    *prometheus.GaugeVec
	servicesVirtualServerServiceHits       *prometheus.CounterVec
	servicesVirtualServerServiceHitsRate   *prometheus.GaugeVec
	servicesActiveTransactions             *prometheus.GaugeVec
}

// NewExporter initialises the exporter
func NewExporter() (*Exporter, error) {
	return &Exporter{
		mgmtCPUUsage:                           mgmtCPUUsage,
		memUsage:                               memUsage,
		pktCPUUsage:                            pktCPUUsage,
		flashPartitionUsage:                    flashPartitionUsage,
		varPartitionUsage:                      varPartitionUsage,
		rxMbPerSec:                             rxMbPerSec,
		txMbPerSec:                             txMbPerSec,
		httpRequestsRate:                       httpRequestsRate,
		httpResponsesRate:                      httpResponsesRate,
		tcpCurrentClientConnections:            tcpCurrentClientConnections,
		tcpCurrentClientConnectionsEstablished: tcpCurrentClientConnectionsEstablished,
		tcpCurrentServerConnections:            tcpCurrentServerConnections,
		tcpCurrentServerConnectionsEstablished: tcpCurrentServerConnectionsEstablished,
		interfacesRxBytesPerSecond:             interfacesRxBytesPerSecond,
		interfacesTxBytesPerSecond:             interfacesTxBytesPerSecond,
		interfacesRxPacketsPerSecond:           interfacesRxPacketsPerSecond,
		interfacesTxPacketsPerSecond:           interfacesTxPacketsPerSecond,
		interfacesJumboPacketsRxPerSecond:      interfacesJumboPacketsRxPerSecond,
		interfacesJumboPacketsTxPerSecond:      interfacesJumboPacketsTxPerSecond,
		interfacesErrorPacketsRxPerSecond:      interfacesErrorPacketsRxPerSecond,
		virtualServersWaitingRequests:          virtualServersWaitingRequests,
		virtualServersHealth:                   virtualServersHealth,
		virtualServersInactiveServices:         virtualServersInactiveServices,
		virtualServersActiveServices:           virtualServersActiveServices,
		virtualServersTotalHits:                virtualServersTotalHits,
		virtualServersHitsRate:                 virtualServersHitsRate,
		virtualServersTotalRequests:            virtualServersTotalRequests,
		virtualServersRequestsRate:             virtualServersRequestsRate,
		virtualServersTotalResponses:           virtualServersTotalResponses,
		virtualServersReponsesRate:             virtualServersReponsesRate,
		virtualServersTotalRequestBytes:        virtualServersTotalRequestBytes,
		virtualServersRequestBytesRate:         virtualServersRequestBytesRate,
		virtualServersTotalResponseBytes:       virtualServersTotalResponseBytes,
		virtualServersReponseBytesRate:         virtualServersReponseBytesRate,
		virtualServersCurrentClientConnections: virtualServersCurrentClientConnections,
		virtualServersCurrentServerConnections: virtualServersCurrentServerConnections,
		servicesThroughput:                     servicesThroughput,
		servicesThroughputRate:                 servicesThroughputRate,
		servicesAvgTTFB:                        servicesAvgTTFB,
		servicesState:                          servicesState,
		servicesTotalRequests:                  servicesTotalRequests,
		servicesRequestsRate:                   servicesRequestsRate,
		servicesTotalResponses:                 servicesTotalResponses,
		servicesResponsesRate:                  servicesResponsesRate,
		servicesTotalRequestBytes:              servicesTotalRequestBytes,
		servicesRequestBytesRate:               servicesRequestBytesRate,
		servicesTotalResponseBytes:             servicesTotalResponseBytes,
		servicesResponseBytesRate:              servicesResponseBytesRate,
		servicesCurrentClientConns:             servicesCurrentClientConns,
		servicesSurgeCount:                     servicesSurgeCount,
		servicesCurrentServerConns:             servicesCurrentServerConns,
		servicesServerEstablishedConnections:   servicesServerEstablishedConnections,
		servicesCurrentReusePool:               servicesCurrentReusePool,
		servicesMaxClients:                     servicesMaxClients,
		servicesCurrentLoad:                    servicesCurrentLoad,
		servicesVirtualServerServiceHits:       servicesVirtualServerServiceHits,
		servicesVirtualServerServiceHitsRate:   servicesVirtualServerServiceHitsRate,
		servicesActiveTransactions:             servicesActiveTransactions,
	}, nil
}

// Describe implements Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- mgmtCPUUsage
	ch <- memUsage
	ch <- pktCPUUsage
	ch <- flashPartitionUsage
	ch <- varPartitionUsage
	ch <- rxMbPerSec
	ch <- txMbPerSec
	ch <- httpRequestsRate
	ch <- httpResponsesRate
	ch <- tcpCurrentClientConnections
	ch <- tcpCurrentClientConnectionsEstablished
	ch <- tcpCurrentServerConnections
	ch <- tcpCurrentServerConnectionsEstablished

	e.interfacesRxBytesPerSecond.Describe(ch)
	e.interfacesTxBytesPerSecond.Describe(ch)
	e.interfacesRxPacketsPerSecond.Describe(ch)
	e.interfacesTxPacketsPerSecond.Describe(ch)
	e.interfacesJumboPacketsRxPerSecond.Describe(ch)
	e.interfacesJumboPacketsTxPerSecond.Describe(ch)
	e.interfacesErrorPacketsRxPerSecond.Describe(ch)

	e.virtualServersWaitingRequests.Describe(ch)
	e.virtualServersHealth.Describe(ch)
	e.virtualServersInactiveServices.Describe(ch)
	e.virtualServersActiveServices.Describe(ch)
	e.virtualServersTotalHits.Describe(ch)
	e.virtualServersHitsRate.Describe(ch)
	e.virtualServersTotalRequests.Describe(ch)
	e.virtualServersRequestsRate.Describe(ch)
	e.virtualServersTotalResponses.Describe(ch)
	e.virtualServersReponsesRate.Describe(ch)
	e.virtualServersTotalRequestBytes.Describe(ch)
	e.virtualServersRequestBytesRate.Describe(ch)
	e.virtualServersTotalResponseBytes.Describe(ch)
	e.virtualServersReponseBytesRate.Describe(ch)
	e.virtualServersCurrentClientConnections.Describe(ch)
	e.virtualServersCurrentServerConnections.Describe(ch)

	e.servicesThroughput.Describe(ch)
	e.servicesThroughputRate.Describe(ch)
	e.servicesAvgTTFB.Describe(ch)
	e.servicesState.Describe(ch)
	e.servicesTotalRequests.Describe(ch)
	e.servicesRequestsRate.Describe(ch)
	e.servicesTotalResponses.Describe(ch)
	e.servicesResponsesRate.Describe(ch)
	e.servicesTotalRequestBytes.Describe(ch)
	e.servicesRequestBytesRate.Describe(ch)
	e.servicesTotalResponseBytes.Describe(ch)
	e.servicesResponseBytesRate.Describe(ch)
	e.servicesCurrentClientConns.Describe(ch)
	e.servicesSurgeCount.Describe(ch)
	e.servicesCurrentServerConns.Describe(ch)
	e.servicesServerEstablishedConnections.Describe(ch)
	e.servicesCurrentReusePool.Describe(ch)
	e.servicesMaxClients.Describe(ch)
	e.servicesCurrentLoad.Describe(ch)
	e.servicesVirtualServerServiceHits.Describe(ch)
	e.servicesVirtualServerServiceHitsRate.Describe(ch)
	e.servicesActiveTransactions.Describe(ch)
}

func (e *Exporter) collectInterfacesRxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesRxBytesPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesRxBytesPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.ReceivedBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesTxBytesPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesTxBytesPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.TransmitBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesRxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesRxPacketsPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesRxPacketsPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.ReceivedPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesTxPacketsPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesTxPacketsPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.TransmitPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsRxPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesJumboPacketsRxPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.JumboPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsTxPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesJumboPacketsTxPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.JumboPacketsTransmittedPerSecond)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesErrorPacketsRxPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesErrorPacketsRxPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.ErrorPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersWaitingRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		waitingRequests, _ := strconv.ParseFloat(vs.WaitingRequests, 64)
		e.virtualServersWaitingRequests.WithLabelValues(nsInstance, vs.Name).Set(waitingRequests)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.virtualServersHealth.Reset()

	for _, vs := range ns.VirtualServerStats {
		health, _ := strconv.ParseFloat(vs.Health, 64)
		e.virtualServersHealth.WithLabelValues(nsInstance, vs.Name).Set(health)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersInactiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		inactiveServices, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.virtualServersInactiveServices.WithLabelValues(nsInstance, vs.Name).Set(inactiveServices)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersActiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		activeServices, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.virtualServersActiveServices.WithLabelValues(nsInstance, vs.Name).Set(activeServices)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalHits.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.virtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectVirtualServerHitsRate(ns netscaler.NSAPIResponse) {
	e.virtualServersHitsRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersHitsRate.WithLabelValues(nsInstance, vs.Name).Set(vs.HitsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.virtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectVirtualServerRequestsRate(ns netscaler.NSAPIResponse) {
	e.virtualServersRequestsRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersRequestsRate.WithLabelValues(nsInstance, vs.Name).Set(vs.RequestsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponses.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.virtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectVirtualServerResponsesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersReponsesRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersReponsesRate.WithLabelValues(nsInstance, vs.Name).Set(vs.ResponsesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.virtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectVirtualServerRequestBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersRequestBytesRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersRequestBytesRate.WithLabelValues(nsInstance, vs.Name).Set(vs.RequestBytesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.virtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectVirtualServerResponseBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersReponseBytesRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersReponseBytesRate.WithLabelValues(nsInstance, vs.Name).Set(vs.ResponseBytesRate)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.virtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.virtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(currentServerConnections)
	}
}

func (e *Exporter) collectServicesThroughput(ns netscaler.NSAPIResponse) {
	e.servicesThroughput.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.Throughput, 64)
		e.servicesThroughput.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesThroughputRate(ns netscaler.NSAPIResponse) {
	e.servicesThroughputRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesThroughputRate.WithLabelValues(nsInstance, service.Name).Set(service.ThroughputRate)
	}
}

func (e *Exporter) collectServicesAvgTTFB(ns netscaler.NSAPIResponse) {
	e.servicesAvgTTFB.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.AvgTimeToFirstByte, 64)
		e.servicesAvgTTFB.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesState(ns netscaler.NSAPIResponse) {
	e.servicesState.Reset()

	for _, service := range ns.ServiceStats {
		state := 0.0

		if service.State == "UP" {
			state = 1.0
		}

		e.servicesState.WithLabelValues(nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequests.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.servicesTotalRequests.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesRequestsRate(ns netscaler.NSAPIResponse) {
	e.servicesRequestsRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesRequestsRate.WithLabelValues(nsInstance, service.Name).Set(service.RequestsRate)
	}
}

func (e *Exporter) collectServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponses.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.servicesTotalResponses.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesResponsesRate(ns netscaler.NSAPIResponse) {
	e.servicesResponsesRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesResponsesRate.WithLabelValues(nsInstance, service.Name).Set(service.ResponsesRate)
	}
}

func (e *Exporter) collectServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequestBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.servicesTotalRequestBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesRequestBytesRate(ns netscaler.NSAPIResponse) {
	e.servicesRequestBytesRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesRequestBytesRate.WithLabelValues(nsInstance, service.Name).Set(service.RequestBytesRate)
	}
}

func (e *Exporter) collectServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponseBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.servicesTotalResponseBytes.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesResponseBytesRate(ns netscaler.NSAPIResponse) {
	e.servicesResponseBytesRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesResponseBytesRate.WithLabelValues(nsInstance, service.Name).Set(service.ResponseBytesRate)
	}
}

func (e *Exporter) collectServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentClientConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.servicesCurrentClientConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesSurgeCount(ns netscaler.NSAPIResponse) {
	e.servicesSurgeCount.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.SurgeCount, 64)
		e.servicesSurgeCount.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentServerConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.servicesCurrentServerConns.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.servicesServerEstablishedConnections.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServerEstablishedConnections, 64)
		e.servicesServerEstablishedConnections.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentReusePool(ns netscaler.NSAPIResponse) {
	e.servicesCurrentReusePool.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentReusePool, 64)
		e.servicesCurrentReusePool.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesMaxClients(ns netscaler.NSAPIResponse) {
	e.servicesMaxClients.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.MaxClients, 64)
		e.servicesMaxClients.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.servicesCurrentLoad.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.servicesCurrentLoad.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.servicesVirtualServerServiceHits.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.servicesVirtualServerServiceHits.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHitsRate(ns netscaler.NSAPIResponse) {
	e.servicesVirtualServerServiceHitsRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesVirtualServerServiceHitsRate.WithLabelValues(nsInstance, service.Name).Set(service.ServiceHitsRate)
	}
}

func (e *Exporter) collectServicesActiveTransactions(ns netscaler.NSAPIResponse) {
	e.servicesActiveTransactions.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ActiveTransactions, 64)
		e.servicesActiveTransactions.WithLabelValues(nsInstance, service.Name).Set(val)
	}
}

// Collect is initiated by the Prometheus handler and gathers the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	nsClient := netscaler.NewNitroClient(*url, *username, *password)

	ns, err := netscaler.GetNSStats(nsClient, "")
	if err != nil {
		log.Error(err)
	}

	interfaces, err := netscaler.GetInterfaceStats(nsClient, "")
	if err != nil {
		log.Error(err)
	}

	virtualServers, err := netscaler.GetVirtualServerStats(nsClient, "")
	if err != nil {
		log.Error(err)
	}

	services, err := netscaler.GetServiceStats(nsClient, "")
	if err != nil {
		log.Error(err)
	}

	fltTCPCurrentClientConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnections, 64)
	fltTCPCurrentClientConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnectionsEstablished, 64)
	fltTCPCurrentServerConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnections, 64)
	fltTCPCurrentServerConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnectionsEstablished, 64)

	ch <- prometheus.MustNewConstMetric(
		mgmtCPUUsage, prometheus.GaugeValue, ns.NSStats.MgmtCPUUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		memUsage, prometheus.GaugeValue, ns.NSStats.MemUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		pktCPUUsage, prometheus.GaugeValue, ns.NSStats.PktCPUUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		flashPartitionUsage, prometheus.GaugeValue, ns.NSStats.FlashPartitionUsage, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		varPartitionUsage, prometheus.GaugeValue, ns.NSStats.VarPartitionUsage, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		rxMbPerSec, prometheus.GaugeValue, ns.NSStats.ReceivedMbPerSecond, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		txMbPerSec, prometheus.GaugeValue, ns.NSStats.TransmitMbPerSecond, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpRequestsRate, prometheus.GaugeValue, ns.NSStats.HTTPRequestsRate, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpResponsesRate, prometheus.GaugeValue, ns.NSStats.HTTPResponsesRate, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnections, prometheus.GaugeValue, fltTCPCurrentClientConnections, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentClientConnectionsEstablished, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnections, prometheus.GaugeValue, fltTCPCurrentServerConnections, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentServerConnectionsEstablished, nsInstance,
	)

	e.collectInterfacesRxBytesPerSecond(interfaces)
	e.interfacesRxBytesPerSecond.Collect(ch)

	e.collectInterfacesTxBytesPerSecond(interfaces)
	e.interfacesTxBytesPerSecond.Collect(ch)

	e.collectInterfacesRxPacketsPerSecond(interfaces)
	e.interfacesRxPacketsPerSecond.Collect(ch)

	e.collectInterfacesTxPacketsPerSecond(interfaces)
	e.interfacesTxPacketsPerSecond.Collect(ch)

	e.collectInterfacesJumboPacketsRxPerSecond(interfaces)
	e.interfacesJumboPacketsRxPerSecond.Collect(ch)

	e.collectInterfacesJumboPacketsTxPerSecond(interfaces)
	e.interfacesJumboPacketsTxPerSecond.Collect(ch)

	e.collectInterfacesErrorPacketsRxPerSecond(interfaces)
	e.interfacesErrorPacketsRxPerSecond.Collect(ch)

	e.collectVirtualServerWaitingRequests(virtualServers)
	e.virtualServersWaitingRequests.Collect(ch)

	e.collectVirtualServerHealth(virtualServers)
	e.virtualServersHealth.Collect(ch)

	e.collectVirtualServerInactiveServices(virtualServers)
	e.virtualServersInactiveServices.Collect(ch)

	e.collectVirtualServerActiveServices(virtualServers)
	e.virtualServersActiveServices.Collect(ch)

	e.collectVirtualServerTotalHits(virtualServers)
	e.virtualServersTotalHits.Collect(ch)

	e.collectVirtualServerHitsRate(virtualServers)
	e.virtualServersHitsRate.Collect(ch)

	e.collectVirtualServerTotalRequests(virtualServers)
	e.virtualServersTotalRequests.Collect(ch)

	e.collectVirtualServerRequestsRate(virtualServers)
	e.virtualServersRequestsRate.Collect(ch)

	e.collectVirtualServerTotalResponses(virtualServers)
	e.virtualServersTotalResponses.Collect(ch)

	e.collectVirtualServerResponsesRate(virtualServers)
	e.virtualServersReponsesRate.Collect(ch)

	e.collectVirtualServerTotalRequestBytes(virtualServers)
	e.virtualServersTotalRequestBytes.Collect(ch)

	e.collectVirtualServerRequestBytesRate(virtualServers)
	e.virtualServersRequestBytesRate.Collect(ch)

	e.collectVirtualServerTotalResponseBytes(virtualServers)
	e.virtualServersTotalResponseBytes.Collect(ch)

	e.collectVirtualServerResponseBytesRate(virtualServers)
	e.virtualServersReponseBytesRate.Collect(ch)

	e.collectVirtualServerCurrentClientConnections(virtualServers)
	e.virtualServersCurrentClientConnections.Collect(ch)

	e.collectVirtualServerCurrentServerConnections(virtualServers)
	e.virtualServersCurrentServerConnections.Collect(ch)

	e.collectServicesThroughput(services)
	e.servicesThroughput.Collect(ch)

	e.collectServicesThroughputRate(services)
	e.servicesThroughputRate.Collect(ch)

	e.collectServicesAvgTTFB(services)
	e.servicesAvgTTFB.Collect(ch)

	e.collectServicesState(services)
	e.servicesState.Collect(ch)

	e.collectServicesTotalRequests(services)
	e.servicesTotalRequests.Collect(ch)

	e.collectServicesRequestsRate(services)
	e.servicesRequestsRate.Collect(ch)

	e.collectServicesTotalResponses(services)
	e.servicesTotalResponses.Collect(ch)

	e.collectServicesResponsesRate(services)
	e.servicesResponsesRate.Collect(ch)

	e.collectServicesTotalRequestBytes(services)
	e.servicesTotalRequestBytes.Collect(ch)

	e.collectServicesRequestBytesRate(services)
	e.servicesRequestBytesRate.Collect(ch)

	e.collectServicesTotalResponseBytes(services)
	e.servicesTotalResponseBytes.Collect(ch)

	e.collectServicesResponseBytesRate(services)
	e.servicesResponseBytesRate.Collect(ch)

	e.collectServicesCurrentClientConns(services)
	e.servicesCurrentClientConns.Collect(ch)

	e.collectServicesSurgeCount(services)
	e.servicesSurgeCount.Collect(ch)

	e.collectServicesCurrentServerConns(services)
	e.servicesCurrentServerConns.Collect(ch)

	e.collectServicesServerEstablishedConnections(services)
	e.servicesServerEstablishedConnections.Collect(ch)

	e.collectServicesCurrentReusePool(services)
	e.servicesCurrentReusePool.Collect(ch)

	e.collectServicesMaxClients(services)
	e.servicesMaxClients.Collect(ch)

	e.collectServicesCurrentLoad(services)
	e.servicesCurrentLoad.Collect(ch)

	e.collectServicesVirtualServerServiceHits(services)
	e.servicesVirtualServerServiceHits.Collect(ch)

	e.collectServicesVirtualServerServiceHitsRate(services)
	e.servicesVirtualServerServiceHitsRate.Collect(ch)

	e.collectServicesActiveTransactions(services)
	e.servicesActiveTransactions.Collect(ch)
}

func main() {
	flag.Parse()

	if *url == "" || *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	nsInstance = strings.TrimLeft(*url, "https://")
	nsInstance = strings.Trim(nsInstance, " /")

	if service.Interactive() != true {
		log.SetFormatter(&log.JSONFormatter{})

		logfile := nsInstance + ".log"
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}

	exporter, _ := NewExporter()
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
	log.Infof("Listening on port %s", listeningPort)
	log.Fatal(http.ListenAndServe(listeningPort, nil))
}

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

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

var (
	url      = flag.String("url", "", "Base URL of the NetScaler management interface.  Normally something like https://my-netscaler.something.x")
	username = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password = flag.String("password", "", "Password with which to connect to the NetScaler API")
	bindPort = flag.Int("bind_port", 9279, "Port to bind the exporter endpoint to")

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
}

func (e *Exporter) collectInterfacesRxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesRxBytesPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfacesRxBytesPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.ReceivedBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesTxBytesPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfacesTxBytesPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.TransmitBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesRxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesRxPacketsPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfacesRxPacketsPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.ReceivedPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesTxPacketsPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfacesTxPacketsPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.TransmitPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsRxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfacesJumboPacketsRxPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.JumboPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsTxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfacesJumboPacketsTxPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.JumboPacketsTransmittedPerSecond)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesErrorPacketsRxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfacesErrorPacketsRxPerSecond.WithLabelValues(nsInstance, iface.ID, iface.Alias).Set(iface.ErrorPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersWaitingRequests.Reset()

	for _, vs := range ns.VirtualServers {
		waitingRequests, _ := strconv.ParseFloat(vs.WaitingRequests, 64)
		e.virtualServersWaitingRequests.WithLabelValues(nsInstance, vs.Name).Set(waitingRequests)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.virtualServersHealth.Reset()

	for _, vs := range ns.VirtualServers {
		health, _ := strconv.ParseFloat(vs.Health, 64)
		e.virtualServersHealth.WithLabelValues(nsInstance, vs.Name).Set(health)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersInactiveServices.Reset()

	for _, vs := range ns.VirtualServers {
		inactiveServices, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.virtualServersInactiveServices.WithLabelValues(nsInstance, vs.Name).Set(inactiveServices)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersActiveServices.Reset()

	for _, vs := range ns.VirtualServers {
		activeServices, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.virtualServersActiveServices.WithLabelValues(nsInstance, vs.Name).Set(activeServices)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalHits.Reset()

	for _, vs := range ns.VirtualServers {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.virtualServersTotalHits.WithLabelValues(nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectVirtualServerHitsRate(ns netscaler.NSAPIResponse) {
	e.virtualServersHitsRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServersHitsRate.WithLabelValues(nsInstance, vs.Name).Set(vs.HitsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequests.Reset()

	for _, vs := range ns.VirtualServers {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.virtualServersTotalRequests.WithLabelValues(nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectVirtualServerRequestsRate(ns netscaler.NSAPIResponse) {
	e.virtualServersRequestsRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServersRequestsRate.WithLabelValues(nsInstance, vs.Name).Set(vs.RequestsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponses.Reset()

	for _, vs := range ns.VirtualServers {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.virtualServersTotalResponses.WithLabelValues(nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectVirtualServerResponsesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersReponsesRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServersReponsesRate.WithLabelValues(nsInstance, vs.Name).Set(vs.ResponsesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.VirtualServers {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.virtualServersTotalRequestBytes.WithLabelValues(nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectVirtualServerRequestBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersRequestBytesRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServersRequestBytesRate.WithLabelValues(nsInstance, vs.Name).Set(vs.RequestBytesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.VirtualServers {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.virtualServersTotalResponseBytes.WithLabelValues(nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectVirtualServerResponseBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersReponseBytesRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServersReponseBytesRate.WithLabelValues(nsInstance, vs.Name).Set(vs.ResponseBytesRate)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.VirtualServers {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.virtualServersCurrentClientConnections.WithLabelValues(nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.VirtualServers {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.virtualServersCurrentServerConnections.WithLabelValues(nsInstance, vs.Name).Set(currentServerConnections)
	}
}

// Collect is initiated by the Prometheus handler and gathers the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	nsClient := netscaler.NewNitroClient(*url, *username, *password)

	ns, err := netscaler.GetNSStats(nsClient)
	if err != nil {
		log.Error(err)
	}

	interfaces, err := netscaler.GetInterfaceStats(nsClient)
	if err != nil {
		log.Error(err)
	}

	virtualServers, err := netscaler.GetVirtualServerStats(nsClient)
	if err != nil {
		log.Error(err)
	}

	ch <- prometheus.MustNewConstMetric(
		mgmtCPUUsage, prometheus.GaugeValue, ns.NS.MgmtCPUUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		memUsage, prometheus.GaugeValue, ns.NS.MemUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		pktCPUUsage, prometheus.GaugeValue, ns.NS.PktCPUUsagePcnt, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		flashPartitionUsage, prometheus.GaugeValue, ns.NS.FlashPartitionUsage, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		varPartitionUsage, prometheus.GaugeValue, ns.NS.VarPartitionUsage, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		rxMbPerSec, prometheus.GaugeValue, ns.NS.ReceivedMbPerSecond, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		txMbPerSec, prometheus.GaugeValue, ns.NS.TransmitMbPerSecond, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpRequestsRate, prometheus.GaugeValue, ns.NS.HTTPRequestsRate, nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpResponsesRate, prometheus.GaugeValue, ns.NS.HTTPResponsesRate, nsInstance,
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
}

func main() {
	flag.Parse()

	if *url == "" || *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	nsInstance = strings.TrimLeft(*url, "https://")
	nsInstance = strings.Trim(nsInstance, " /")

	interactive, _ := svc.IsAnInteractiveSession()

	if interactive != true {
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

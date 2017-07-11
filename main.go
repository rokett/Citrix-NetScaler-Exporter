package main

import (
	"Citrix-NetScaler-Exporter/netscaler"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

//TODO: Proper commenting
//TODO: Refactor into separate files
//TODO: Label for Netscdaler name

var (
	url      = flag.String("url", "", "Base URL of the NetScaler management interface.  Normally something like https://my-netscaler.something.x")
	username = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password = flag.String("password", "", "Password with which to connect to the NetScaler API")
	bindPort = flag.Int("bind_port", 9279, "Port to bind the exporter endpoint to")

	mgmtCPUUsage = prometheus.NewDesc(
		"mgmt_cpu_usage",
		"Current CPU utilisation for management",
		nil,
		nil,
	)

	pktCPUUsage = prometheus.NewDesc(
		"pkt_cpu_usage",
		"Current CPU utilisation for packet engines, excluding management",
		nil,
		nil,
	)

	memUsage = prometheus.NewDesc(
		"mem_usage",
		"Current memory utilisation",
		nil,
		nil,
	)

	flashPartitionUsage = prometheus.NewDesc(
		"flash_partition_usage",
		"Used space in /flash partition of the disk, as a percentage.",
		nil,
		nil,
	)

	varPartitionUsage = prometheus.NewDesc(
		"var_partition_usage",
		"Used space in /var partition of the disk, as a percentage. ",
		nil,
		nil,
	)

	rxMbPerSec = prometheus.NewDesc(
		"received_mb_per_second",
		"Number of Megabits received by the NetScaler appliance per second",
		nil,
		nil,
	)

	txMbPerSec = prometheus.NewDesc(
		"transmit_mb_per_second",
		"Number of Megabits transmitted by the NetScaler appliance per second",
		nil,
		nil,
	)

	HTTPRequestsRate = prometheus.NewDesc(
		"http_requests_rate",
		"HTTP requests received per second",
		nil,
		nil,
	)

	HTTPResponsesRate = prometheus.NewDesc(
		"http_responses_rate",
		"HTTP requests sent per second",
		nil,
		nil,
	)

	interfaces_rxBytesPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_bytes_per_second",
			Help: "Number of bytes received per second by specific interfaces",
		},
		[]string{
			"interface",
			"alias",
		},
	)

	interfaces_txBytesPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_bytes_per_second",
			Help: "Number of bytes transmitted per second by specific interfaces",
		},
		[]string{
			"interface",
			"alias",
		},
	)

	interfaces_rxPacketsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_packets_per_second",
			Help: "Number of packets received per second by specific interfaces",
		},
		[]string{
			"interface",
			"alias",
		},
	)

	interfaces_txPacketsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_packets_per_second",
			Help: "Number of packets transmitted per second by specific interfaces",
		},
		[]string{
			"interface",
			"alias",
		},
	)

	interfaces_jumboPacketsRxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_received_per_second",
			Help: "Number of bytes received per second by specific interfaces",
		},
		[]string{
			"interface",
			"alias",
		},
	)

	interfaces_jumboPacketsTxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_transmitted_per_second",
			Help: "Number of jumbo packets transmitted per second by specific interfaces",
		},
		[]string{
			"interface",
			"alias",
		},
	)

	interfaces_errorPacketsRxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_error_packets_received_per_second",
			Help: "Number of error packets received per second by specific interfaces",
		},
		[]string{
			"interface",
			"alias",
		},
	)

	virtualServers_waitingRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_waiting_requests",
			Help: "Number of requests waiting on a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_health = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_inactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_activeServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_totalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_hitsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_hits_rate",
			Help: "Number of hits/second to a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_requestsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_requests_rate",
			Help: "Number of requests/second to a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_totalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_reponsesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_responses_rate",
			Help: "Number of responses/second from a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_totalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_requestBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_request_bytes_rate",
			Help: "Number of request bytes/second to a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_totalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_reponseBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_reponse_bytes_rate",
			Help: "Number of response bytes/second from a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_currentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"virtual_server",
		},
	)

	virtualServers_currentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"virtual_server",
		},
	)
)

// Exporter represents the metrics exported to Prometheus
type Exporter struct {
	mgmtCPUUsage                            *prometheus.Desc
	memUsage                                *prometheus.Desc
	pktCPUUsage                             *prometheus.Desc
	flashPartitionUsage                     *prometheus.Desc
	varPartitionUsage                       *prometheus.Desc
	rxMbPerSec                              *prometheus.Desc
	txMbPerSec                              *prometheus.Desc
	HTTPRequestsRate                        *prometheus.Desc
	HTTPResponsesRate                       *prometheus.Desc
	interfaces_rxBytesPerSecond             *prometheus.GaugeVec
	interfaces_txBytesPerSecond             *prometheus.GaugeVec
	interfaces_rxPacketsPerSecond           *prometheus.GaugeVec
	interfaces_txPacketsPerSecond           *prometheus.GaugeVec
	interfaces_jumboPacketsRxPerSecond      *prometheus.GaugeVec
	interfaces_jumboPacketsTxPerSecond      *prometheus.GaugeVec
	interfaces_errorPacketsRxPerSecond      *prometheus.GaugeVec
	virtualServers_waitingRequests          *prometheus.GaugeVec
	virtualServers_health                   *prometheus.GaugeVec
	virtualServers_inactiveServices         *prometheus.GaugeVec
	virtualServers_activeServices           *prometheus.GaugeVec
	virtualServers_totalHits                *prometheus.CounterVec
	virtualServers_hitsRate                 *prometheus.GaugeVec
	virtualServers_totalRequests            *prometheus.CounterVec
	virtualServers_requestsRate             *prometheus.GaugeVec
	virtualServers_totalResponses           *prometheus.CounterVec
	virtualServers_reponsesRate             *prometheus.GaugeVec
	virtualServers_totalRequestBytes        *prometheus.CounterVec
	virtualServers_requestBytesRate         *prometheus.GaugeVec
	virtualServers_totalResponseBytes       *prometheus.CounterVec
	virtualServers_reponseBytesRate         *prometheus.GaugeVec
	virtualServers_currentClientConnections *prometheus.GaugeVec
	virtualServers_currentServerConnections *prometheus.GaugeVec
}

// NewExporter initialises the exporter
func NewExporter() (*Exporter, error) {
	return &Exporter{
		mgmtCPUUsage:                            mgmtCPUUsage,
		memUsage:                                memUsage,
		pktCPUUsage:                             pktCPUUsage,
		flashPartitionUsage:                     flashPartitionUsage,
		varPartitionUsage:                       varPartitionUsage,
		rxMbPerSec:                              rxMbPerSec,
		txMbPerSec:                              txMbPerSec,
		HTTPRequestsRate:                        HTTPRequestsRate,
		HTTPResponsesRate:                       HTTPResponsesRate,
		interfaces_rxBytesPerSecond:             interfaces_rxBytesPerSecond,
		interfaces_txBytesPerSecond:             interfaces_txBytesPerSecond,
		interfaces_rxPacketsPerSecond:           interfaces_rxPacketsPerSecond,
		interfaces_txPacketsPerSecond:           interfaces_txPacketsPerSecond,
		interfaces_jumboPacketsRxPerSecond:      interfaces_jumboPacketsRxPerSecond,
		interfaces_jumboPacketsTxPerSecond:      interfaces_jumboPacketsTxPerSecond,
		interfaces_errorPacketsRxPerSecond:      interfaces_errorPacketsRxPerSecond,
		virtualServers_waitingRequests:          virtualServers_waitingRequests,
		virtualServers_health:                   virtualServers_health,
		virtualServers_inactiveServices:         virtualServers_inactiveServices,
		virtualServers_activeServices:           virtualServers_activeServices,
		virtualServers_totalHits:                virtualServers_totalHits,
		virtualServers_hitsRate:                 virtualServers_hitsRate,
		virtualServers_totalRequests:            virtualServers_totalRequests,
		virtualServers_requestsRate:             virtualServers_requestsRate,
		virtualServers_totalResponses:           virtualServers_totalResponses,
		virtualServers_reponsesRate:             virtualServers_reponsesRate,
		virtualServers_totalRequestBytes:        virtualServers_totalRequestBytes,
		virtualServers_requestBytesRate:         virtualServers_requestBytesRate,
		virtualServers_totalResponseBytes:       virtualServers_totalResponseBytes,
		virtualServers_reponseBytesRate:         virtualServers_reponseBytesRate,
		virtualServers_currentClientConnections: virtualServers_currentClientConnections,
		virtualServers_currentServerConnections: virtualServers_currentServerConnections,
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
	ch <- HTTPRequestsRate
	ch <- HTTPResponsesRate
	e.interfaces_rxBytesPerSecond.Describe(ch)
	e.interfaces_txBytesPerSecond.Describe(ch)
	e.interfaces_rxPacketsPerSecond.Describe(ch)
	e.interfaces_txPacketsPerSecond.Describe(ch)
	e.interfaces_jumboPacketsRxPerSecond.Describe(ch)
	e.interfaces_jumboPacketsTxPerSecond.Describe(ch)
	e.interfaces_errorPacketsRxPerSecond.Describe(ch)
	e.virtualServers_waitingRequests.Describe(ch)
	e.virtualServers_health.Describe(ch)
	e.virtualServers_inactiveServices.Describe(ch)
	e.virtualServers_activeServices.Describe(ch)
	e.virtualServers_totalHits.Describe(ch)
	e.virtualServers_hitsRate.Describe(ch)
	e.virtualServers_totalRequests.Describe(ch)
	e.virtualServers_requestsRate.Describe(ch)
	e.virtualServers_totalResponses.Describe(ch)
	e.virtualServers_reponsesRate.Describe(ch)
	e.virtualServers_totalRequestBytes.Describe(ch)
	e.virtualServers_requestBytesRate.Describe(ch)
	e.virtualServers_totalResponseBytes.Describe(ch)
	e.virtualServers_reponseBytesRate.Describe(ch)
	e.virtualServers_currentClientConnections.Describe(ch)
	e.virtualServers_currentServerConnections.Describe(ch)
}

func (e *Exporter) collectInterfacesRxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfaces_rxBytesPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_rxBytesPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.ReceivedBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfaces_rxBytesPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_txBytesPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.TransmitBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesRxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfaces_rxPacketsPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_rxPacketsPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.ReceivedPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfaces_txPacketsPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_txPacketsPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.TransmitPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfaces_jumboPacketsRxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_jumboPacketsRxPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.JumboPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfaces_jumboPacketsTxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_jumboPacketsTxPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.JumboPacketsTransmittedPerSecond)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfaces_errorPacketsRxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_errorPacketsRxPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.ErrorPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse) {
	e.virtualServers_waitingRequests.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_waitingRequests.WithLabelValues(vs.Name).Set(vs.WaitingRequests)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.virtualServers_health.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_health.WithLabelValues(vs.Name).Set(vs.Health)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServers_inactiveServices.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_inactiveServices.WithLabelValues(vs.Name).Set(vs.InactiveServices)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServers_activeServices.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_activeServices.WithLabelValues(vs.Name).Set(vs.ActiveServices)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.virtualServers_totalHits.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_totalHits.WithLabelValues(vs.Name).Set(vs.TotalHits)
	}
}

func (e *Exporter) collectVirtualServerHitsRate(ns netscaler.NSAPIResponse) {
	e.virtualServers_hitsRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_hitsRate.WithLabelValues(vs.Name).Set(vs.HitsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.virtualServers_totalRequests.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_totalRequests.WithLabelValues(vs.Name).Set(vs.TotalRequests)
	}
}

func (e *Exporter) collectVirtualServerRequestsRate(ns netscaler.NSAPIResponse) {
	e.virtualServers_requestsRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_requestsRate.WithLabelValues(vs.Name).Set(vs.RequestsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.virtualServers_totalResponses.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_totalResponses.WithLabelValues(vs.Name).Set(vs.TotalResponses)
	}
}

func (e *Exporter) collectVirtualServerResponsesRate(ns netscaler.NSAPIResponse) {
	e.virtualServers_reponsesRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_reponsesRate.WithLabelValues(vs.Name).Set(vs.ResponsesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.virtualServers_totalRequestBytes.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_totalRequestBytes.WithLabelValues(vs.Name).Set(vs.TotalRequestBytes)
	}
}

func (e *Exporter) collectVirtualServerRequestBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServers_requestBytesRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_requestBytesRate.WithLabelValues(vs.Name).Set(vs.RequestBytesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.virtualServers_totalResponseBytes.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_totalResponseBytes.WithLabelValues(vs.Name).Set(vs.TotalResponseBytes)
	}
}

func (e *Exporter) collectVirtualServerResponseBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServers_reponseBytesRate.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_reponseBytesRate.WithLabelValues(vs.Name).Set(vs.ResponseBytesRate)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.virtualServers_currentClientConnections.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_currentClientConnections.WithLabelValues(vs.Name).Set(vs.CurrentClientConnections)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.virtualServers_currentServerConnections.Reset()

	for _, vs := range ns.VirtualServers {
		e.virtualServers_currentServerConnections.WithLabelValues(vs.Name).Set(vs.CurrentServerConnections)
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
		mgmtCPUUsage, prometheus.GaugeValue, ns.NS.MgmtCPUUsagePcnt,
	)

	ch <- prometheus.MustNewConstMetric(
		memUsage, prometheus.GaugeValue, ns.NS.MemUsagePcnt,
	)

	ch <- prometheus.MustNewConstMetric(
		pktCPUUsage, prometheus.GaugeValue, ns.NS.PktCPUUsagePcnt,
	)

	ch <- prometheus.MustNewConstMetric(
		flashPartitionUsage, prometheus.GaugeValue, ns.NS.FlashPartitionUsage,
	)

	ch <- prometheus.MustNewConstMetric(
		varPartitionUsage, prometheus.GaugeValue, ns.NS.VarPartitionUsage,
	)

	ch <- prometheus.MustNewConstMetric(
		rxMbPerSec, prometheus.GaugeValue, ns.NS.ReceivedMbPerSecond,
	)

	ch <- prometheus.MustNewConstMetric(
		txMbPerSec, prometheus.GaugeValue, ns.NS.TransmitMbPerSecond,
	)

	ch <- prometheus.MustNewConstMetric(
		HTTPRequestsRate, prometheus.GaugeValue, ns.NS.HTTPRequestsRate,
	)

	ch <- prometheus.MustNewConstMetric(
		HTTPResponsesRate, prometheus.GaugeValue, ns.NS.HTTPResponsesRate,
	)

	e.collectInterfacesRxBytesPerSecond(interfaces)
	e.interfaces_rxBytesPerSecond.Collect(ch)

	e.collectInterfacesTxBytesPerSecond(interfaces)
	e.interfaces_txBytesPerSecond.Collect(ch)

	e.collectInterfacesRxPacketsPerSecond(interfaces)
	e.interfaces_rxPacketsPerSecond.Collect(ch)

	e.collectInterfacesTxPacketsPerSecond(interfaces)
	e.interfaces_txPacketsPerSecond.Collect(ch)

	e.collectInterfacesJumboPacketsRxPerSecond(interfaces)
	e.interfaces_jumboPacketsRxPerSecond.Collect(ch)

	e.collectInterfacesJumboPacketsTxPerSecond(interfaces)
	e.interfaces_jumboPacketsTxPerSecond.Collect(ch)

	e.collectInterfacesErrorPacketsRxPerSecond(interfaces)
	e.interfaces_errorPacketsRxPerSecond.Collect(ch)

	e.collectVirtualServerWaitingRequests(virtualServers)
	e.virtualServers_waitingRequests.Collect(ch)

	e.collectVirtualServerHealth(virtualServers)
	e.virtualServers_health.Collect(ch)

	e.collectVirtualServerInactiveServices(virtualServers)
	e.virtualServers_inactiveServices.Collect(ch)

	e.collectVirtualServerActiveServices(virtualServers)
	e.virtualServers_activeServices.Collect(ch)

	e.collectVirtualServerTotalHits(virtualServers)
	e.virtualServers_totalHits.Collect(ch)

	e.collectVirtualServerHitsRate(virtualServers)
	e.virtualServers_hitsRate.Collect(ch)

	e.collectVirtualServerTotalRequests(virtualServers)
	e.virtualServers_totalRequests.Collect(ch)

	e.collectVirtualServerRequestsRate(virtualServers)
	e.virtualServers_requestsRate.Collect(ch)

	e.collectVirtualServerTotalResponses(virtualServers)
	e.virtualServers_totalResponses.Collect(ch)

	e.collectVirtualServerResponsesRate(virtualServers)
	e.virtualServers_reponsesRate.Collect(ch)

	e.collectVirtualServerTotalRequestBytes(virtualServers)
	e.virtualServers_totalRequestBytes.Collect(ch)

	e.collectVirtualServerRequestBytesRate(virtualServers)
	e.virtualServers_requestBytesRate.Collect(ch)

	e.collectVirtualServerTotalResponseBytes(virtualServers)
	e.virtualServers_totalResponseBytes.Collect(ch)

	e.collectVirtualServerResponseBytesRate(virtualServers)
	e.virtualServers_reponseBytesRate.Collect(ch)

	e.collectVirtualServerCurrentClientConnections(virtualServers)
	e.virtualServers_currentClientConnections.Collect(ch)

	e.collectVirtualServerCurrentServerConnections(virtualServers)
	e.virtualServers_currentServerConnections.Collect(ch)
}

func main() {
	flag.Parse()

	if *url == "" || *username == "" || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	interactive, _ := svc.IsAnInteractiveSession()

	if interactive != true {
		log.SetFormatter(&log.JSONFormatter{})

		//TODO: Set log file
		file, err := os.OpenFile("cmdb.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}

	listeningPort := ":" + strconv.Itoa(*bindPort)

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

	log.Infof("Listening on port %s", listeningPort)
	log.Fatal(http.ListenAndServe(listeningPort, nil))
}

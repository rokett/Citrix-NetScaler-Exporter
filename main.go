package main

import (
	"Citrix-NetScaler-Exporter/collectors"
	"Citrix-NetScaler-Exporter/netscaler"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
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
)

type Exporter struct {
	mgmtCPUUsage                       *prometheus.Desc
	memUsage                           *prometheus.Desc
	pktCPUUsage                        *prometheus.Desc
	flashPartitionUsage                *prometheus.Desc
	varPartitionUsage                  *prometheus.Desc
	rxMbPerSec                         *prometheus.Desc
	txMbPerSec                         *prometheus.Desc
	HTTPRequestsRate                   *prometheus.Desc
	HTTPResponsesRate                  *prometheus.Desc
	interfaces_rxBytesPerSecond        *prometheus.GaugeVec
	interfaces_txBytesPerSecond        *prometheus.GaugeVec
	interfaces_rxPacketsPerSecond      *prometheus.GaugeVec
	interfaces_txPacketsPerSecond      *prometheus.GaugeVec
	interfaces_jumboPacketsRxPerSecond *prometheus.GaugeVec
	interfaces_jumboPacketsTxPerSecond *prometheus.GaugeVec
	interfaces_errorPacketsRxPerSecond *prometheus.GaugeVec
}

func NewExporter() (*Exporter, error) {
	return &Exporter{
		mgmtCPUUsage:                       mgmtCPUUsage,
		memUsage:                           memUsage,
		pktCPUUsage:                        pktCPUUsage,
		flashPartitionUsage:                flashPartitionUsage,
		varPartitionUsage:                  varPartitionUsage,
		rxMbPerSec:                         rxMbPerSec,
		txMbPerSec:                         txMbPerSec,
		HTTPRequestsRate:                   HTTPRequestsRate,
		HTTPResponsesRate:                  HTTPResponsesRate,
		interfaces_rxBytesPerSecond:        interfaces_rxBytesPerSecond,
		interfaces_txBytesPerSecond:        interfaces_txBytesPerSecond,
		interfaces_rxPacketsPerSecond:      interfaces_rxPacketsPerSecond,
		interfaces_txPacketsPerSecond:      interfaces_txPacketsPerSecond,
		interfaces_jumboPacketsRxPerSecond: interfaces_jumboPacketsRxPerSecond,
		interfaces_jumboPacketsTxPerSecond: interfaces_jumboPacketsTxPerSecond,
		interfaces_errorPacketsRxPerSecond: interfaces_errorPacketsRxPerSecond,
	}, nil
}

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
}

func (e *Exporter) collectInterfacesRxBytesPerSecond(ns collectors.NSAPIResponse) {
	e.interfaces_rxBytesPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_rxBytesPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.ReceivedBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxBytesPerSecond(ns collectors.NSAPIResponse) {
	e.interfaces_rxBytesPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_txBytesPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.TransmitBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesRxPacketsPerSecond(ns collectors.NSAPIResponse) {
	e.interfaces_rxPacketsPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_rxPacketsPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.ReceivedPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxPacketsPerSecond(ns collectors.NSAPIResponse) {
	e.interfaces_txPacketsPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_txPacketsPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.TransmitPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRxPerSecond(ns collectors.NSAPIResponse) {
	e.interfaces_jumboPacketsRxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_jumboPacketsRxPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.JumboPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTxPerSecond(ns collectors.NSAPIResponse) {
	e.interfaces_jumboPacketsTxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_jumboPacketsTxPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.JumboPacketsTransmittedPerSecond)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRxPerSecond(ns collectors.NSAPIResponse) {
	e.interfaces_errorPacketsRxPerSecond.Reset()

	for _, iface := range ns.Interfaces {
		e.interfaces_errorPacketsRxPerSecond.WithLabelValues(iface.ID, iface.Alias).Set(iface.ErrorPacketsReceivedPerSecond)
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	nsClient := netscaler.NewNitroClient("https://netscaler-internal-new.exe.nhs.uk", "nsroot", "rGM!WiLl!cMG=$tZZ~|Lw4q*!GNlGm")

	ns := collectors.GetNSStats(nsClient)
	interfaces := collectors.GetInterfaceStats(nsClient)

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
}

func main() {
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

	log.Fatal(http.ListenAndServe(":9279", nil))
}

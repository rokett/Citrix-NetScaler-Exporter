package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sonny-sevin-mk/citrix-netscaler-exporter/netscaler"
)

var (
	interfacesRxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_bytes",
			Help: "Number of bytes received by specific interfaces.",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_bytes",
			Help: "Number of bytes transmitted by specific interfaces.",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesRxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_packets",
			Help: "Number of packets received by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_packets",
			Help: "Number of packets transmitted by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_received",
			Help: "Number of bytes received by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsTx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_transmitted",
			Help: "Number of jumbo packets transmitted by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesErrorPacketsRx = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_error_packets_received",
			Help: "Number of error packets received by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)
)

func (e *Exporter) collectInterfacesRxBytes(ns netscaler.NSAPIResponse) {
	e.interfacesRxBytes.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedBytes, 64)
		e.interfacesRxBytes.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxBytes(ns netscaler.NSAPIResponse) {
	e.interfacesTxBytes.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitBytes, 64)
		e.interfacesTxBytes.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesRxPackets(ns netscaler.NSAPIResponse) {
	e.interfacesRxPackets.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalReceivedPackets, 64)
		e.interfacesRxPackets.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesTxPackets(ns netscaler.NSAPIResponse) {
	e.interfacesTxPackets.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.TotalTransmitPackets, 64)
		e.interfacesTxPackets.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRx(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsRx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsReceived, 64)
		e.interfacesJumboPacketsRx.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTx(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsTx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.JumboPacketsTransmitted, 64)
		e.interfacesJumboPacketsTx.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRx(ns netscaler.NSAPIResponse) {
	e.interfacesErrorPacketsRx.Reset()

	for _, iface := range ns.InterfaceStats {
		val, _ := strconv.ParseFloat(iface.ErrorPacketsReceived, 64)
		e.interfacesErrorPacketsRx.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(val)
	}
}

package collector

import (
	"strconv"

	"citrix-netscaler-exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	vpnVirtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vpn_virtual_servers_total_requests",
			Help: "Total VPN virtual server requests",
		},
		[]string{
			"ns_instance",
			"vpn_virtual_server",
		},
	)

	vpnVirtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vpn_virtual_servers_total_responses",
			Help: "Total VPN virtual server responses",
		},
		[]string{
			"ns_instance",
			"vpn_virtual_server",
		},
	)

	vpnVirtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vpn_virtual_servers_total_request_bytes",
			Help: "Total VPN virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"vpn_virtual_server",
		},
	)
	vpnVirtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vpn_virtual_servers_total_response_bytes",
			Help: "Total VPN virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"vpn_virtual_server",
		},
	)

	vpnVirtualServersState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vpn_virtual_servers_state",
			Help: "Current state of the VPN virtual server",
		},
		[]string{
			"ns_instance",
			"vpn_virtual_server",
		},
	)
)

func (e *Exporter) collectVPNVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.vpnVirtualServersTotalRequests.Reset()

	for _, vs := range ns.VPNVirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.vpnVirtualServersTotalRequests.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.vpnVirtualServersTotalResponses.Reset()

	for _, vs := range ns.VPNVirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.vpnVirtualServersTotalResponses.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.vpnVirtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.VPNVirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.vpnVirtualServersTotalRequestBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectVPNVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.vpnVirtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.VPNVirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.vpnVirtualServersTotalResponseBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectVPNVirtualServerState(ns netscaler.NSAPIResponse) {
	e.vpnVirtualServersState.Reset()

	for _, vs := range ns.VPNVirtualServerStats {
		state := 0.0

		if vs.State == "UP" {
			state = 1.0
		}

		e.vpnVirtualServersState.WithLabelValues(e.nsInstance, vs.Name).Set(state)
	}
}

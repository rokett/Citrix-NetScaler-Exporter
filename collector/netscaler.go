package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	up = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "Is the NetScaler appliance up?",
	})

	modelID = prometheus.NewDesc(
		"model_id",
		"NetScaler model - reflects the bandwidth available; for example VPX 10 would report as 10.",
		[]string{
			"ns_instance",
		},
		nil,
	)

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

	totRxMB = prometheus.NewDesc(
		"total_received_mb",
		"Total number of Megabytes received by the NetScaler appliance",
		[]string{
			"ns_instance",
		},
		nil,
	)

	totTxMB = prometheus.NewDesc(
		"total_transmit_mb",
		"Total number of Megabytes transmitted by the NetScaler appliance",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpRequests = prometheus.NewDesc(
		"http_requests",
		"Total number of HTTP requests received",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpResponses = prometheus.NewDesc(
		"http_responses",
		"Total number of HTTP responses sent",
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
)

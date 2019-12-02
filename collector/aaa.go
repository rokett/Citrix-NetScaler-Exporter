package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	aaaauthsuccess = prometheus.NewDesc(
		"aaaauthsuccess",
		"Count of authentication successes.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaaauthfail = prometheus.NewDesc(
		"aaaauthfail",
		"Count of authentication failures",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaaauthonlyhttpsuccess = prometheus.NewDesc(
		"aaaauthonlyhttpsuccess",
		"Count of HTTP connections that succeeded authorization.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaaauthonlyhttpfail = prometheus.NewDesc(
		"aaaauthonlyhttpfail",
		"Count of HTTP connections that failed authorization.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaaauthnonhttpsuccess = prometheus.NewDesc(
		"aaaauthnonhttpsuccess",
		"Count of non HTTP connections that succeeded authorization.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaaauthnonhttpfail = prometheus.NewDesc(
		"aaaauthnonhttpfail",
		"Count of non HTTP connections that failed authorization.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaacursessions = prometheus.NewDesc(
		"aaacursessions",
		"Count of current SmartAccess AAA sessions.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaatotsessions = prometheus.NewDesc(
		"aaatotsessions",
		"Count of all SmartAccess AAA sessions.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaatotsessiontimeout = prometheus.NewDesc(
		"aaatotsessiontimeout",
		"Count of AAA sessions that have timed out.",
		[]string{
			"ns_instance",
		},
		nil,
	)
	aaacuricasessions = prometheus.NewDesc(
		"aaacuricasessions",
		"Count of current ICA only sessions.",
		[]string{
			"ns_instance",
		},
		nil,
	)
	aaacuricaonlyconn = prometheus.NewDesc(
		"aaacuricaonlyconn",
		"Count of current Basic ICA only connections.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	aaacuricaconn = prometheus.NewDesc(
		"aaacuricaconn",
		"Count of current SmartAccess ICA connections.",
		[]string{
			"ns_instance",
		},
		nil,
	)
	aaacurtmsessions = prometheus.NewDesc(
		"aaacurtmsessions",
		"Count of current AAATM connections.",
		[]string{
			"ns_instance",
		},
		nil,
	)
	aaatottmsessions = prometheus.NewDesc(
		"aaatottmsessions",
		"Count of all AAATM sessions.",
		[]string{
			"ns_instance",
		},
		nil,
	)
)

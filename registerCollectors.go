package main

import (
	"fmt"

	"Citrix-NetScaler-Exporter/collectors"
	"Citrix-NetScaler-Exporter/netscaler"

	"github.com/prometheus/client_golang/prometheus"
)

func registerCollectors() {
	nsClient := netscaler.NewNitroClient("http://netscaler-test.exe.nhs.uk", "nsroot", "password")

	err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "Current CPU utilisation",
		},
		func() float64 {
			return collectors.GetCPU()
		},
	))

	if err != nil {
		fmt.Println(err)
	}

	err = prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "mem_usage",
			Help: "Current memoru utilisation",
		},
		func() float64 {
			return collectors.GetMem()
		},
	))

	if err != nil {
		fmt.Println(err)
	}

	collectors.GetNSStats(&nsclient)
}

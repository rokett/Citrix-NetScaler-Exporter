package collectors

import (
	"Citrix-NetScaler-Exporter/netscaler"
	"encoding/json"
	"fmt"
	"log"
)

// NSStats ...
type NSStats struct {
	CPUUsagePcnt        float64 `json:"cpuusagepcnt"`
	MemUsagePcnt        float64 `json:"memusagepcnt"`
	MgmtCPUUsagePcnt    float64 `json:"mgmtcpuusagepcnt"`
	PktCPUUsagePcnt     float64 `json:"pktcpuusagepcnt"`
	FlashPartitionUsage float64 `json:"disk0perusage"`
	VarPartitionUsage   float64 `json:"disk1perusage"`
	ReceivedMbPerSecond float64 `json:"rxmbitsrate"`
	TransmitMbPerSecond float64 `json:"txmbitsrate"`
	HTTPRequestsRate    float64 `json:"httprequestsrate"`
	HTTPResponsesRate   float64 `json:"httpresponsesrate"`
}

// GetNSStats ...
func GetNSStats(c *netscaler.NitroClient) NSAPIResponse {
	stats, err := c.GetStats("ns")
	if err != nil {
		log.Println(err)
	}

	var s = new(NSAPIResponse)

	err = json.Unmarshal(stats, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	return *s
}

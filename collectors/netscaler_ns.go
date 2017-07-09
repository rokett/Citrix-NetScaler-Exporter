package collectors

import (
	"Citrix-NetScaler-Exporter/netscaler"
	"log"
)

// GetNSStats ...
func GetNSStats(c *netscaler.NitroClient) float64 {
	stats, err := c.GetStats("ns")
	if err != nil {
		log.Println(err)
	}

	//parse bytes into JSON array
	return stats.cpuusagepcnt
}

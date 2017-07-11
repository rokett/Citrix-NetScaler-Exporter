package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// NSStats represents the data returned from the /stat/ns Nitro API endpoint
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

// GetNSStats queries the Nitro API for ns stats
func GetNSStats(c *NitroClient) (NSAPIResponse, error) {
	stats, err := c.GetStats("ns")
	if err != nil {
		return NSAPIResponse{}, err
	}

	var response = new(NSAPIResponse)

	err = json.Unmarshal(stats, &response)
	if err != nil {
		return NSAPIResponse{}, errors.Wrap(err, "error unmarshalling response body")
	}

	return *response, nil
}

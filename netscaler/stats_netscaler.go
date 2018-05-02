package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// NSStats represents the data returned from the /stat/ns Nitro API endpoint
type NSStats struct {
	CPUUsagePcnt                           float64 `json:"cpuusagepcnt"`
	MemUsagePcnt                           float64 `json:"memusagepcnt"`
	MgmtCPUUsagePcnt                       float64 `json:"mgmtcpuusagepcnt"`
	PktCPUUsagePcnt                        float64 `json:"pktcpuusagepcnt"`
	FlashPartitionUsage                    float64 `json:"disk0perusage"`
	VarPartitionUsage                      float64 `json:"disk1perusage"`
	TotalReceivedMB	string	`json:"totrxmbits"`
	TotalTransmitMB string `json:"tottxmbits"`
	HTTPRequests                       string `json:"httptotrequests"`
	HTTPResponses                      string `json:"httptotresponses"`
	TCPCurrentClientConnections            string  `json:"tcpcurclientconn"`
	TCPCurrentClientConnectionsEstablished string  `json:"tcpcurclientconnestablished"`
	TCPCurrentServerConnections            string  `json:"tcpcurserverconn"`
	TCPCurrentServerConnectionsEstablished string  `json:"tcpcurserverconnestablished"`
}

// GetNSStats queries the Nitro API for ns stats
func GetNSStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("ns", querystring)
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

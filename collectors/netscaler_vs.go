package collectors

import (
	"Citrix-NetScaler-Exporter/netscaler"
	"encoding/json"
	"fmt"
	"log"
)

//TODO: Proper commenting
//TODO: Proper logging

// VirtualServerStats ...
type VirtualServerStats struct {
	Name                     string  `json:"name"`
	WaitingRequests          float64 `json:"vsvrsurgecount"`
	Health                   float64 `json:"vslbhealth"`
	InactiveServices         float64 `json:"inactsvcs"`
	ActiveServices           float64 `json:"actsvcs"`
	TotalHits                float64 `json:"tothits"`
	HitsRate                 float64 `json:"hitsrate"`
	TotalRequests            float64 `json:"totalrequests"`
	RequestsRate             float64 `json:"requestsrate"`
	TotalResponses           float64 `json:"totalresponses"`
	ResponsesRate            float64 `json:"responsesrate"`
	TotalRequestBytes        float64 `json:"totalrequestsbytes"`
	RequestBytesRate         float64 `json:"requestbytesrate"`
	TotalResponseBytes       float64 `json:"totalresponsebytes"`
	ResponseBytesRate        float64 `json:"responsebytesrate"`
	CurrentClientConnections float64 `json:"curclntconnections"`
	CurrentServerConnections float64 `json:"cursrvrconnections"`
}

// GetVirtualServerStats ...
func GetVirtualServerStats(c *netscaler.NitroClient) NSAPIResponse {
	stats, err := c.GetStats("lbvserver")
	if err != nil {
		log.Println(err)
	}

	//TODO: s?  Bad name
	var s = new(NSAPIResponse)

	err = json.Unmarshal(stats, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	return *s
}

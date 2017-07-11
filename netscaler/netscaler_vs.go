package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// VirtualServerStats represents the data returned from the /stat/lbvserver Nitro API endpoint
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

// GetVirtualServerStats queries the Nitro API for virtual server stats
func GetVirtualServerStats(c *NitroClient) (NSAPIResponse, error) {
	stats, err := c.GetStats("lbvserver")
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

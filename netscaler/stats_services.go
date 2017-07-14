package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ServiceStats represents the data returned from the /stat/service Nitro API endpoint
type ServiceStats struct {
	Name                         string  `json:"name"`
	Throughput                   string  `json:"throughput"`
	ThroughputRate               float64 `json:"throughputrate"`
	AvgTimeToFirstByte           string  `json:"avgsvrttfb"`
	State                        string  `json:"state"`
	TotalRequests                string  `json:"totalrequests"`
	RequestsRate                 float64 `json:"requestsrate"`
	TotalResponses               string  `json:"totalresponses"`
	ResponsesRate                float64 `json:"responsesrate"`
	TotalRequestBytes            string  `json:"totalrequestbytes"`
	RequestBytesRate             float64 `json:"requestbytesrate"`
	TotalResponseBytes           string  `json:"totalresponsebytes"`
	ResponseBytesRate            float64 `json:"responsebytesrate"`
	CurrentClientConnections     string  `json:"curclntconnections"`
	SurgeCount                   string  `json:"surgecount"`
	CurrentServerConnections     string  `json:"cursrvrconnections"`
	ServerEstablishedConnections string  `json:"svrestablishedconn"`
	CurrentReusePool             string  `json:"curreusepool"`
	MaxClients                   string  `json:"maxclients"`
	CurrentLoad                  string  `json:"curload"`
	ServiceHits                  string  `json:"vsvrservicehits"`
	ServiceHitsRate              float64 `json:"vsvrservicehitsrate"`
	ActiveTransactions           string  `json:"activetransactions"`
}

// GetServiceStats queries the Nitro API for service stats
func GetServiceStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("service", querystring)
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

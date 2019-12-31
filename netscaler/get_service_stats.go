package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ServiceStats represents the data returned from the /stat/service Nitro API endpoint
type ServiceStats struct {
	Name                         string  `json:"name"`
	Throughput                   string  `json:"throughput"`
	AvgTimeToFirstByte           string  `json:"avgsvrttfb"`
	State                        string  `json:"state"`
	TotalRequests                string  `json:"totalrequests"`
	TotalResponses               string  `json:"totalresponses"`
	TotalRequestBytes            string  `json:"totalrequestbytes"`
	TotalResponseBytes           string  `json:"totalresponsebytes"`
	CurrentClientConnections     string  `json:"curclntconnections"`
	SurgeCount                   string  `json:"surgecount"`
	CurrentServerConnections     string  `json:"cursrvrconnections"`
	ServerEstablishedConnections string  `json:"svrestablishedconn"`
	CurrentReusePool             string  `json:"curreusepool"`
	MaxClients                   string  `json:"maxclients"`
	CurrentLoad                  string  `json:"curload"`
	ServiceHits                  string  `json:"vsvrservicehits"`
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

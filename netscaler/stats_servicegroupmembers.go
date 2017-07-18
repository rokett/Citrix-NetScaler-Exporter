package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ServiceGroupMemberStats represents the data returned from the /stat/servicegroupmember Nitro API endpoint
type ServiceGroupMemberStats struct {
	State                        string  `json:"state"`
	AvgTimeToFirstByte           string  `json:"avgsvrttfb"`
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
}

// GetServiceGroupMemberStats queries the Nitro API for servicegroup member stats
func GetServiceGroupMemberStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("servicegroupmember", querystring)
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

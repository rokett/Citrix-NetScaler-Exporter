package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GSLBServiceStats represents the data returned from the /stat/gslbservice Nitro API endpoint
type GSLBServiceStats struct {
	Name                     string `json:"servicename"`
	State                    string `json:"state"`
	TotalRequests            string `json:"totalrequests"`
	TotalResponses           string `json:"totalresponses"`
	TotalRequestBytes        string `json:"totalrequestbytes"`
	TotalResponseBytes       string `json:"totalresponsebytes"`
	CurrentClientConnections string `json:"curclntconnections"`
	CurrentServerConnections string `json:"cursrvrconnections"`
	EstablishedConnections   string `json:"establishedconn"`
	CurrentLoad              string `json:"curload"`
	ServiceHits              string `json:"vsvrservicehits"`
}

// GetGSLBServiceStats queries the Nitro API for service stats
func GetGSLBServiceStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("gslbservice", querystring)
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

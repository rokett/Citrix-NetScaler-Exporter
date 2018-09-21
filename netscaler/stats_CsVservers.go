package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// VirtualServerStats represents the data returned from the /stat/lbvserver Nitro API endpoint
type CsVserverStats struct {
	Name                     string  `json:"name"`
	State										 string  `json:"state"`
	TotalHits                string  `json:"tothits"`
	TotalRequests            string  `json:"totalrequests"`
	TotalResponses           string  `json:"totalresponses"`
	TotalRequestBytes        string  `json:"totalrequestsbytes"`
	TotalResponseBytes       string  `json:"totalresponsebytes"`
	CurrentClientConnections string  `json:"curclntconnections"`
	CurrentServerConnections string  `json:"cursrvrconnections"`
}

// GetVirtualServerStats queries the Nitro API for virtual server stats
func GetCsVserverStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("csvserver", querystring)
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

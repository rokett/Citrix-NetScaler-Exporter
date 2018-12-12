package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GSLBVirtualServerStats represents the data returned from the /stat/gslbvserver Nitro API endpoint
type GSLBVirtualServerStats struct {
	Name                     string `json:"name"`
	Health                   string `json:"vslbhealth"`
	InactiveServices         string `json:"inactsvcs"`
	ActiveServices           string `json:"actsvcs"`
	TotalHits                string `json:"tothits"`
	TotalRequests            string `json:"totalrequests"`
	TotalResponses           string `json:"totalresponses"`
	TotalRequestBytes        string `json:"totalrequestbytes"`
	TotalResponseBytes       string `json:"totalresponsebytes"`
	CurrentClientConnections string `json:"curclntconnections"`
	CurrentServerConnections string `json:"cursrvrconnections"`
}

// GetGSLBVirtualServerStats queries the Nitro API for virtual server stats
func GetGSLBVirtualServerStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("gslbvserver", querystring)
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

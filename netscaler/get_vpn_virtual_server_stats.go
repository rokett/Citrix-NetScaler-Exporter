package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// VPNVirtualServerStats represents the data returned from the /stat/vpnvserver Nitro API endpoint
type VPNVirtualServerStats struct {
	Name               string `json:"name"`
	TotalRequests      string `json:"totalrequests"`
	TotalResponses     string `json:"totalresponses"`
	TotalRequestBytes  string `json:"totalrequestbytes"`
	TotalResponseBytes string `json:"totalresponsebytes"`
	State              string `json:"state"`
}

// GetVPNVirtualServerStats queries the Nitro API for VPN virtual server stats
func GetVPNVirtualServerStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("vpnvserver", querystring)
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

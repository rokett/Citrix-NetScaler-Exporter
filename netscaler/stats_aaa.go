package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// InterfaceStats represents the data returned from the /stat/interface Nitro API endpoint
type AAAStats struct {
	ID                      string `json:"id"`
	TotalReceivedBytes      string `json:"totrxbytes"`
	TotalTransmitBytes      string `json:"tottxbytes"`
	TotalReceivedPackets    string `json:"totrxpkts"`
	TotalTransmitPackets    string `json:"tottxpkts"`
	JumboPacketsReceived    string `json:"jumbopktsreceived"`
	JumboPacketsTransmitted string `json:"jumbopktstransmitted"`
	ErrorPacketsReceived    string `json:"errpktrx"`
	Alias                   string `json:"interfacealias"`
}

// GetAAAStats queries the Nitro API for interface stats
func GetAAAStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("aaa", querystring)
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

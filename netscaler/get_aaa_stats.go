package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// AAAStats represents the data returned from the /stat/aaa Nitro API endpoint
type AAAStats struct {
	AuthSuccess               string `json:"aaaauthsuccess"`
	AuthFail                  string `json:"aaaauthfail"`
	AuthOnlyHTTPSuccess       string `json:"aaaauthonlyhttpsuccess"`
	AuthOnlyHTTPFail          string `json:"aaaauthonlyhttpfail"`
	CurrentIcaSessions        string `json:"aaacuricasessions"`
	CurrentIcaOnlyConnections string `json:"aaacuricaonlyconn"`
}

// GetAAAStats queries the Nitro API for AAA stats
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

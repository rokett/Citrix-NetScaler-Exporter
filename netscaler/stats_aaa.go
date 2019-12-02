package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// AAAStats represents the data returned from the /stat/aaa Nitro API endpoint
type AAAStats struct {
	AAAauthsuccess         string `json:"aaaauthsuccess"`
	AAAauthfail            string `json:"aaaauthfail"`
	AAAauthonlyhttpsuccess string `json:"aaaauthonlyhttpsuccess"`
	AAAauthonlyhttpfail    string `json:"aaaauthonlyhttpfail"`
	AAAauthnonhttpsuccess  string `json:"aaaauthnonhttpsuccess"`
	AAAauthnonhttpfail     string `json:"aaaauthnonhttpfail"`
	AAAcursessions         string `json:"aaacursessions"`
	AAAtotsessions         string `json:"aaatotsessions"`
	AAAtotsessiontimeout   string `json:"aaatotsessiontimeout"`
	AAAcuricasessions      string `json:"aaacuricasessions"`
	AAAcuricaonlyconn      string `json:"aaacuricaonlyconn"`
	AAAcuricaconn          string `json:"aaacuricaconn"`
	AAAcurtmsessions       string `json:"aaacurtmsessions"`
	AAAtottmsessions       string `json:"aaatottmsessions"`
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

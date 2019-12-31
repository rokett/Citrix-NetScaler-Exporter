package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// NSLicense represents the data returned from the /config/nslicense Nitro API endpoint
type NSLicense struct {
	ModelID string `json:"modelid"`
}

// GetNSLicense queries the Nitro API for license config
func GetNSLicense(c *NitroClient, querystring string) (NSAPIResponse, error) {
	cfg, err := c.GetConfig("nslicense", querystring)
	if err != nil {
		return NSAPIResponse{}, err
	}

	var response = new(NSAPIResponse)

	err = json.Unmarshal(cfg, &response)
	if err != nil {
		return NSAPIResponse{}, errors.Wrap(err, "error unmarshalling response body")
	}

	return *response, nil
}

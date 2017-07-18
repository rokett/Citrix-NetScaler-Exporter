package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ServiceGroups represents the data returned from the /config/servicegroup Nitro API endpoint
type ServiceGroups struct {
	Name string `json:"servicegroupname"`
}

// GetServiceGroups queries the Nitro API for service group config
func GetServiceGroups(c *NitroClient, querystring string) (NSAPIResponse, error) {
	cfg, err := c.GetConfig("servicegroup", querystring)
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

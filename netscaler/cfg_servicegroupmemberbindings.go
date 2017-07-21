package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ServiceGroupMemberBindings represents the data returned from the /config/servicegroup_servicegroupmember_binding Nitro API endpoint
type ServiceGroupMemberBindings struct {
	ServerName string `json:"servername"`
	Port       int64  `json:"port"`
}

// GetServiceGroupMemberBindings queries the Nitro API for service group member binding details
func GetServiceGroupMemberBindings(c *NitroClient, servicegroup string) (NSAPIResponse, error) {
	url := "servicegroup_servicegroupmember_binding/" + servicegroup
	cfg, err := c.GetConfig(url, "")
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

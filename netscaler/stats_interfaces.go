package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// InterfaceStats represents the data returned from the /stat/interface Nitro API endpoint
type InterfaceStats struct {
	ID                               string  `json:"id"`
	ReceivedBytesPerSecond           float64 `json:"rxbytesrate"`
	TransmitBytesPerSecond           float64 `json:"txbytesrate"`
	ReceivedPacketsPerSecond         float64 `json:"rxpktsrate"`
	TransmitPacketsPerSecond         float64 `json:"txpktsrate"`
	JumboPacketsReceivedPerSecond    float64 `json:"jumbopktsreceivedrate"`
	JumboPacketsTransmittedPerSecond float64 `json:"jumbopktstransmittedrate"`
	ErrorPacketsReceivedPerSecond    float64 `json:"errpktrxrate"`
	Alias                            string  `json:"interfacealias"`
}

// GetInterfaceStats queries the Nitro API for interface stats
func GetInterfaceStats(c *NitroClient) (NSAPIResponse, error) {
	stats, err := c.GetStats("interface")
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

package collectors

import (
	"Citrix-NetScaler-Exporter/netscaler"
	"encoding/json"
	"fmt"
	"log"
)

//TODO: Proper commenting
//TODO: Proper logging

// InterfaceStats ...
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

// GetInterfaceStats ...
func GetInterfaceStats(c *netscaler.NitroClient) NSAPIResponse {
	stats, err := c.GetStats("interface")
	if err != nil {
		log.Println(err)
	}

	//TODO: s?  Bad name
	var s = new(NSAPIResponse)

	err = json.Unmarshal(stats, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	return *s
}

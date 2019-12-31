package netscaler

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// CSVirtualServerStats represents the data returned from the /stat/csvserver Nitro API endpoint
type CSVirtualServerStats struct {
	Name                          string `json:"name"`
	State                         string `json:"state"`
	TotalHits                     string `json:"tothits"`
	TotalRequests                 string `json:"totalrequests"`
	TotalResponses                string `json:"totalresponses"`
	TotalRequestBytes             string `json:"totalrequestbytes"`
	TotalResponseBytes            string `json:"totalresponsebytes"`
	CurrentClientConnections      string `json:"curclntconnections"`
	CurrentServerConnections      string `json:"cursrvrconnections"`
	EstablishedConnections        string `json:"establishedconn"`
	TotalPacketsReceived          string `json:"totalpktsrecvd"`
	TotalPacketsSent              string `json:"totalpktssent"`
	TotalSpillovers               string `json:"totspillovers"`
	DeferredRequests              string `json:"deferredreq"`
	InvalidRequestResponse        string `json:"invalidrequestresponse"`
	InvalidRequestResponseDropped string `json:"invalidrequestresponsedropped"`
	TotalVServerDownBackupHits    string `json:"totvserverdownbackuphits"`
	CurrentMultipathSessions      string `json:"curmptcpsessions"`
	CurrentMultipathSubflows      string `json:"cursubflowconn"`
}

// GetCSVirtualServerStats queries the Nitro API for Content Switching virtual server stats
func GetCSVirtualServerStats(c *NitroClient, querystring string) (NSAPIResponse, error) {
	stats, err := c.GetStats("csvserver", querystring)
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

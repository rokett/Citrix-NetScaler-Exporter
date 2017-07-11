package netscaler

// NSAPIResponse represents the main portion of the Nitro API response
type NSAPIResponse struct {
	Errorcode      int64                `json:"errorcode"`
	Message        string               `json:"message"`
	Severity       string               `json:"severity"`
	NS             NSStats              `json:"ns"`
	Interfaces     []InterfaceStats     `json:"Interface"`
	VirtualServers []VirtualServerStats `json:"lbvserver"`
}

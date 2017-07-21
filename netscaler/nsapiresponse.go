package netscaler

// NSAPIResponse represents the main portion of the Nitro API response
type NSAPIResponse struct {
	Errorcode          int64                `json:"errorcode"`
	Message            string               `json:"message"`
	Severity           string               `json:"severity"`
	SessionID          string               `json:"sessionid"`
	NSLicense          NSLicense            `json:"nslicense"`
	NSStats            NSStats              `json:"ns"`
	InterfaceStats     []InterfaceStats     `json:"Interface"`
	VirtualServerStats []VirtualServerStats `json:"lbvserver"`
	ServiceStats       []ServiceStats       `json:"service"`
}

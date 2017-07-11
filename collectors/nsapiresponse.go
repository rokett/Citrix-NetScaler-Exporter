package collectors

//TODO: Proper commenting

// NSAPIResponse ...
type NSAPIResponse struct {
	Errorcode      int64                `json:"errorcode"`
	Message        string               `json:"message"`
	Severity       string               `json:"severity"`
	NS             NSStats              `json:"ns"`
	Interfaces     []InterfaceStats     `json:"Interface"`
	VirtualServers []VirtualServerStats `json:"lbvserver"`
}

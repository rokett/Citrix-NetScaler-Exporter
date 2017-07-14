package netscaler

import (
	"net/http"
	"strings"
	"time"
)

// NitroClient represents the client used to connect to the API
type NitroClient struct {
	url      string
	username string
	password string
	client   *http.Client
}

// NewNitroClient creates a new client used to interact with the Nirto API.
// URL, username and password are passed to this function to allow connections to any NetScaler endpoint.
func NewNitroClient(url string, username string, password string) *NitroClient {
	c := new(NitroClient)

	c.url = strings.Trim(url, " /") + "/nitro/v1/"

	c.username = username
	c.password = password
	c.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return c
}

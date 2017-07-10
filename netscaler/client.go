package netscaler

import (
	"net/http"
	"strings"
	"time"
)

// NitroClient ...
type NitroClient struct {
	url      string
	username string
	password string
	client   *http.Client
}

// NewNitroClient ...
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

package netscaler

import (
	"net/http"
	"strings"
	"time"
)

//TODO: Proper comments

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

	//TODO: What is the trim for?
	c.url = strings.Trim(url, " /") + "/nitro/v1/"

	c.username = username
	c.password = password
	c.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	return c
}

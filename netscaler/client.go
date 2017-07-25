package netscaler

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/pkg/errors"
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
func NewNitroClient(url string, username string, password string) (*NitroClient, error) {
	c := new(NitroClient)

	c.url = strings.Trim(url, " /") + "/nitro/v1/"

	c.username = username
	c.password = password

	jar, err := cookiejar.New(nil)
	if err != nil {
		return c, errors.Wrap(err, "error creating cookiejar")
	}

	c.client = &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}

	return c, nil
}

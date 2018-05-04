package netscaler

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
	"crypto/tls"

	"github.com/pkg/errors"
)

// NitroClient represents the client used to connect to the API
type NitroClient struct {
	url      string
	username string
	password string
	client   *http.Client
}

// NewNitroClient creates a new client used to interact with the Nitro API.
// URL, username and password are passed to this function to allow connections to any NetScaler endpoint.
// The ignoreCert parameter allows self-signed certificates to be accepted.  It should be used sparingly and only when you fully trust the endpoint.
func NewNitroClient(url string, username string, password string, ignoreCert bool) (*NitroClient, error) {
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

	if ignoreCert {
		transpCfg := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		c.client.Transport = transpCfg
	}

	return c, nil
}

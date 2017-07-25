package netscaler

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// GetConfig sends a request to the Nitro API and retrieves configuration for the given type.
func (c *NitroClient) GetConfig(configType string, querystring string) ([]byte, error) {
	url := c.url + "config/" + configType

	if querystring != "" {
		url = url + "?" + querystring
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating HTTP request")
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "error sending request")
	}

	switch resp.StatusCode {
	case 200:
		body, _ := ioutil.ReadAll(resp.Body)

		return body, nil
	default:
		body, _ := ioutil.ReadAll(resp.Body)

		return body, errors.New("read failed: " + resp.Status + " (" + string(body) + ")")
	}
}

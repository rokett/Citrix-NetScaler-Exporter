package netscaler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// DisconnectPayload is the JSON that needs to be sent to NetScaler to logout
type DisconnectPayload struct {
	Logout struct {
	} `json:"logout"`
}

// Disconnect logs out of the NetScaler
func Disconnect(c *NitroClient) error {
	url := c.url + "config/logout"

	var p DisconnectPayload

	reqBody, err := json.Marshal(p)
	if err != nil {
		return errors.Wrap(err, "error marshalling payload")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.Wrap(err, "error creating HTTP request")
	}

	req.Header.Set("Content-Type", "application/vnd.com.citrix.netscaler.logout+json")

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}

	switch resp.StatusCode {
	case 200:
		return nil
	// Although the documentation says a 200 is returned when a logout works, testing shows that it's actually a 201
	case 201:
		return nil
	default:
		body, _ := ioutil.ReadAll(resp.Body)

		return errors.New("Logout failed: " + resp.Status + " (" + string(body) + ")")
	}
}

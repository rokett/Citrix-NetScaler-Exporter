package netscaler

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// GetStats ...
func (c *NitroClient) GetStats(statsType string) ([]byte, error) {
	url := c.url + statsType

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-NITRO-USER", c.username)
	req.Header.Set("X-NITRO-PASS", c.password)

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatal(err)
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

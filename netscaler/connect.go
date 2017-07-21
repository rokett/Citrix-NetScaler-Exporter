package netscaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type LoginCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Payload struct {
	Login LoginCreds `json:"login"`
}

// Connect initiates a connection to a NetScaler and returns the session token
func Connect(c *NitroClient) (string, error) {
	url := c.url + "config/login"

	var p Payload

	p.Login = LoginCreds{
		Username: c.username,
		Password: c.password,
	}

	reqBody, err := json.Marshal(p)
	if err != nil {
		return "", errors.Wrap(err, "error marshalling payload")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", errors.Wrap(err, "error creating HTTP request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", errors.Wrap(err, "error sending request")
	}

	switch resp.StatusCode {
	case 201:
		var response = new(NSAPIResponse)

		fmt.Println(resp.Cookies())

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", errors.Wrap(err, "error unmarshalling response body")
		}

		return response.SessionID, nil
	default:
		body, _ := ioutil.ReadAll(resp.Body)

		return "", errors.New("read failed: " + resp.Status + " (" + string(body) + ")")
	}
}

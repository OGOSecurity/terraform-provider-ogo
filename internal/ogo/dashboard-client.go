package ogosecurity

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client
type Client struct {
	Endpoint    string
	HostBaseURL string
	HTTPClient  *http.Client
	Username    string
	ApiKey      string
}

func md5sum(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// NewClient
func NewClient(host *string, username *string, apikey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Check if endpoint, username and password are provided
	if host == nil {
		return &c, errors.New("Endpoint must be provided")
	} else {
		c.Endpoint = *host
	}

	if username == nil {
		return &c, errors.New("Username must be provided")
	} else {
		c.Username = *username
	}

	if apikey == nil {
		return &c, errors.New("API Key must be provided")
	} else {
		c.ApiKey = *apikey
	}

	c.HostBaseURL = *host + "/api/" + *username

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// Generate token based on URL Path
	token := md5sum(req.URL.Path + "-" + c.ApiKey)

	// Set token on query parameters
	q := req.URL.Query()
	q.Add("t", token)
	req.URL.RawQuery = q.Encode()

	//fmt.Printf("req: %+v\n", req)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("err res body: %+v\n", res)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	//fmt.Printf("res: %+v\n", res)
	//fmt.Printf("body res: %s\n", body)

	return body, err
}

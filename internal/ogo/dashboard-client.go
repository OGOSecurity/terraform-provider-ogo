// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package ogosecurity

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var sem = make(chan int, 1)

// Client
type Client struct {
	Endpoint     string
	HostBaseURL  string
	HTTPClient   *http.Client
	Organization string
	ApiKey       string
}

// NewClient
func NewClient(host *string, organization *string, apikey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Check if endpoint, organization and password are provided
	if host == nil {
		return &c, errors.New("endpoint must be provided")
	} else {
		c.Endpoint = *host
	}

	if organization == nil {
		return &c, errors.New("organization must be provided")
	} else {
		c.Organization = *organization
	}

	if apikey == nil {
		return &c, errors.New("API Key must be provided")
	} else {
		c.ApiKey = *apikey
	}

	c.HostBaseURL = *host + "/v2/organizations/" + *organization

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Ogo-Api-Key", c.ApiKey)

	// Lock access to Ogo API to restrict count of concurrent requests
	sem <- 1
	res, err := c.HTTPClient.Do(req)
	<-sem

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieved body: %s (%+v)", string(err.Error()), res)
	}

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusNoContent &&
		res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s (%+v)", res.StatusCode, body, res)
	}

	return body, err
}

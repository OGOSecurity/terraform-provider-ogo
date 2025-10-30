// Copyright (c) OGO Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package ogosecurity

import (
	"crypto/md5"
	"encoding/hex"
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
	Email        string
	ApiKey       string
	Organization string
}

func md5sum(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// NewClient
func NewClient(host *string, email *string, apikey *string, organization *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}

	// Check if endpoint, email, apikey and organization are provided
	if host == nil {
		return &c, errors.New("endpoint must be provided")
	} else {
		c.Endpoint = *host
	}

	if email == nil {
		return &c, errors.New("user email address must be provided")
	} else {
		c.Email = *email
	}

	if apikey == nil {
		return &c, errors.New("API key must be provided")
	} else {
		c.ApiKey = *apikey
	}

	if organization == nil {
		return &c, errors.New("organization must be provided")
	} else {
		c.Organization = *organization
	}

	c.HostBaseURL = *host + "/v2/organizations/" + *organization

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Generate token based on URL Path
	token := md5sum(req.URL.Path + "-" + c.ApiKey)

	// Set headers
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}
	req.Header.Set("X-Ogo-Auth", c.Email+";"+token)

	// Lock access to Ogo API to restrict concurrent requests
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

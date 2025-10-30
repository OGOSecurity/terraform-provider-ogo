// Copyright (c) OGO Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package ogosecurity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Returns all user's TLS Options.
func (c *Client) GetAllTlsOptions() ([]TlsOptions, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tls-options", c.HostBaseURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := TlsOptionsResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.TlsOptions, nil
}

// Returns a specifc TLS Options.
func (c *Client) GetTlsOptions(tlsOptionsUid string) (*TlsOptions, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tls-options/%s", c.HostBaseURL, tlsOptionsUid), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := TlsOptions{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create new TLS Options.
func (c *Client) CreateTlsOptions(tlsOptions TlsOptions) (*TlsOptions, error) {
	rb, err := json.Marshal(tlsOptions)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tls-options", c.HostBaseURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := TlsOptions{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update existing TLS Options.
func (c *Client) UpdateTlsOptions(tlsOptions TlsOptions) (*TlsOptions, error) {
	uid := tlsOptions.Uid
	tlsOptions.Uid = ""
	if uid == "" {
		return nil, fmt.Errorf("TLS options UID is required for update")
	}

	rb, err := json.Marshal(tlsOptions)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tls-options/%s", c.HostBaseURL, uid), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := TlsOptions{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Delete an existing TLS Options.
func (c *Client) DeleteTlsOptions(tlsOptionsUid string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tls-options/%s", c.HostBaseURL, tlsOptionsUid), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

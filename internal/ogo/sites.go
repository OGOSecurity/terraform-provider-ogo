// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package ogosecurity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetSite - Returns a specifc site
func (c *Client) GetSite(siteDomainName string) (*Site, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sites/%s", c.HostBaseURL, siteDomainName), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := Site{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateSite - Create new site
func (c *Client) CreateSite(site Site) (*Site, error) {
	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sites", c.HostBaseURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := Site{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateSite - Update existing site
func (c *Client) UpdateSite(site Site) (*Site, error) {
	rb, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/sites/%s", c.HostBaseURL, site.DomainName), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := Site{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteSite - Delete existing site
func (c *Client) DeleteSite(siteDomainName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/sites/%s", c.HostBaseURL, siteDomainName), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

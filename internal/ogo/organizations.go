// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package ogosecurity

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAllOrganizations - Returns all user's organization
func (c *Client) GetAllOrganizations() ([]Organization, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/organizations", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp OrganizationsResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	var organizations []Organization
	for _, o := range resp.OrganizationDetails {
		organizations = append(organizations, o.Organization)
	}

	return organizations, nil
}

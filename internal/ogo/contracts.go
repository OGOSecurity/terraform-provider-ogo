// Copyright (c) OGO Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package ogosecurity

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAllContracts - Returns all user's contract.
func (c *Client) GetAllContracts() ([]Contract, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/contracts/available", c.HostBaseURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp ContractsResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Contracts, nil
}

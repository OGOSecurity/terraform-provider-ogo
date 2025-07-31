package ogosecurity

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetAllTlsOptions - Returns all user's TLS Options
func (c *Client) GetAllTlsOptions() ([]TlsOptions, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tlsOptions", c.HostBaseURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := AllTlsOptionsResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.TlsOptions, nil
}

// GetTlsOptions - Returns a specifc TLS Options
func (c *Client) GetTlsOptions(tlsOptionsUid string) (*TlsOptions, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tlsOptions/%s", c.HostBaseURL, tlsOptionsUid), nil)
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

	return &resp.TlsOptions, nil
}

// CreateTlsOptions - Create new TLS Options
func (c *Client) CreateTlsOptions(tlsOptions TlsOptions) (*TlsOptions, error) {
	tq := TlsOptionsQuery{
		TlsOptions: tlsOptions,
	}

	rb, err := json.Marshal(tq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tlsOptions", c.HostBaseURL), strings.NewReader(string(rb)))
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

	if resp.HasError {
		return nil, errors.New("Failed to create TLS options: " + resp.Status.Message)
	}

	t, err := c.GetTlsOptions(resp.TlsOptions.Uid)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// UpdateTlsOptions - Update existing TLS Options
func (c *Client) UpdateTlsOptions(tlsOptions TlsOptions) (*TlsOptions, error) {
	uid := tlsOptions.Uid
	tlsOptions.Uid = ""

	tq := TlsOptionsQuery{
		TlsOptions: tlsOptions,
	}

	if uid == "" {
		return nil, fmt.Errorf("TLS options UID is required for update")
	}

	rb, err := json.Marshal(tq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tlsOptions/%s", c.HostBaseURL, uid), strings.NewReader(string(rb)))
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

	if resp.HasError {
		return nil, errors.New("Failed to update TLS options: " + resp.Status.Message)
	}

	t, err := c.GetTlsOptions(uid)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// DeleteTlsOptions - Delete an existing TLS Options
func (c *Client) DeleteTlsOptions(tlsOptionsUid string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tlsOptions/%s", c.HostBaseURL, tlsOptionsUid), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	resp := TlsOptionsResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	if resp.HasError {
		return errors.New("Failed to delete TLS options: " + resp.Status.Message)
	}

	return nil
}

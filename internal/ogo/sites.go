package ogosecurity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetAllSites - Returns all user's site
func (c *Client) GetAllSites() ([]Site, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sites", c.HostBaseURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := SitesResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return sr.Sites, nil
}

// GetSite - Returns a specifc site
func (c *Client) GetSite(siteName string) (*Site, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sites/%s", c.HostBaseURL, siteName), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := SiteResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return &sr.Site, nil
}

// CreateSite - Create new site
func (c *Client) CreateSite(site Site) (*Site, error) {
	sq := SiteQuery{
		Site: site,
	}

	rb, err := json.Marshal(sq)
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

	sr := SitesResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	if sr.Count == 0 {
		return nil, errors.New("Failed to create site: " + sr.Status.Message)
	}

	s, err := c.GetSite(site.Name)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// UpdateSite - Update existing site
func (c *Client) UpdateSite(site Site) (*Site, error) {
	sq := SiteQuery{
		Action: "update",
		Site:   site,
	}

	rb, err := json.Marshal(sq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sites/%s", c.HostBaseURL, site.Name), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := SitesResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	if sr.Count == 0 {
		return nil, errors.New("Failed to update site: " + sr.Status.Message)
	}

	return &sr.SitesItems[0], nil
}

// DeleteSite - Delete existing site
func (c *Client) DeleteSite(siteName string) error {
	rb := `{"action":"delete","site":{"name":"` + siteName + `"}}`

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/sites/%s", c.HostBaseURL, siteName), bytes.NewBufferString(rb))
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	sr := SitesResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return err
	}

	if sr.HasError {
		return errors.New(string(body))
	}

	return nil
}

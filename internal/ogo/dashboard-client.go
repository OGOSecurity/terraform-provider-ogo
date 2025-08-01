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

var sem = make(chan int, 1)

// Client
type Client struct {
	Endpoint    string
	HostBaseURL string
	HTTPClient  *http.Client
	Username    string
	ApiKey      string
	Clusters    map[string]string
}

func md5sum(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// NewClient
func NewClient(host *string, username *string, apikey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Clusters:   map[string]string{},
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

	c.HostBaseURL = *host + "/v2/organizations/" + *username
	clusters, err := c.GetAllClusters()

	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		c.Clusters[cluster.ClusterName] = cluster.Id
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Ogo-Api-Key", c.ApiKey)

	//fmt.Printf("req: %+v\n", req)

	// Lock access to Ogo API to restrict count of concurrent request
	sem <- 1
	res, err := c.HTTPClient.Do(req)
	<-sem

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieved body: %s (%+v)", string(err.Error()), res)
	}

	if res.StatusCode != http.StatusOK &&
		res.StatusCode != http.StatusNoContent &&
		res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s (%+v)", res.StatusCode, body, res)
	}

	//fmt.Printf("res: %+v\n", res)
	//fmt.Printf("body res: %s\n", body)

	return body, err
}

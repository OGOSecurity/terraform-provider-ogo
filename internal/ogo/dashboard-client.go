package ogosecurity

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	Clusters    map[string]int
}

func md5sum(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func printo(o any) {
	b, _ := json.MarshalIndent(o, "", "  ")
	fmt.Printf(string(b) + "\n")
}

// NewClient
func NewClient(host *string, username *string, apikey *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Clusters:   map[string]int{},
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

	// Generate token based on URL Path
	//token := md5sum(req.URL.Path + "-" + c.ApiKey)

	// Set token on query parameters
	//q := req.URL.Query()
	//q.Add("t", token)
	//req.URL.RawQuery = q.Encode()

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

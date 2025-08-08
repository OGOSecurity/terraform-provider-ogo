package ogosecurity

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAllClusters - Returns all user's cluster
func (c *Client) GetAllClusters() ([]Cluster, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/clusters", c.HostBaseURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp []ClustersResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	var clusters []Cluster
	for _, c := range resp {
		clusters = append(clusters, c.Cluster)
	}

	return clusters, nil
}

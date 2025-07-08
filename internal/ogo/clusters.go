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

	sr := ClustersResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return sr.Clusters, nil
}

// GetCluster - Returns a specifc cluster
func (c *Client) GetCluster(clusterId int) ([]Cluster, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/clusters/%d", c.HostBaseURL, clusterId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	sr := ClustersResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return sr.Clusters, nil
}

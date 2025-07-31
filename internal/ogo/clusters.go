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

	resp := ClustersResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Clusters, nil
}

// GetCluster - Returns a specifc cluster
func (c *Client) GetCluster(clusterName string) ([]Cluster, error) {
	if _, ok := c.Clusters[clusterName]; !ok {
		return nil, fmt.Errorf("Unknown cluster %s", clusterName)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/clusters/%d", c.HostBaseURL, c.Clusters[clusterName]), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	resp := ClustersResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Clusters, nil
}

package ogosecurity

// Ogo Reponse Status
type OgoResponseStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Site
type Site struct {
	Name             string `json:"name"`
	ClusterName      string `json:"clusterName,omitempty"`
	DestHost         string `json:"destHost"`
	DestHostScheme   string `json:"destHostScheme"`
	TrustSelfSigned  bool   `json:"trustSelfSigned"`
	NoCopyXForwarded bool   `json:"noCopyXForwarded"`
	ForceHttps       bool   `json:"forceHttps"`
	DryRun           bool   `json:"dryRun"`
	PanicMode        bool   `json:"panicMode"`
}

// Site Create Query
type SiteQuery struct {
	Action string `json:"action,omitempty"`
	Site   Site   `json:"site"`
}

// Site Response
type SiteResponse struct {
	Site     Site              `json:"site"`
	Status   OgoResponseStatus `json:"status"`
	HasError bool              `json:"hasError"`
	Count    int               `json:"count"`
}

// Sites Response
type SitesResponse struct {
	Sites      []Site            `json:"sites"`
	SitesItems []Site            `json:"items"`
	Status     OgoResponseStatus `json:"status"`
	HasError   bool              `json:"hasError"`
	Count      int               `json:"count"`
}

// Cluster
type Cluster struct {
	ClusterID           int      `json:"clusterId"`
	ClusterHost         string   `json:"clusterHost,omitempty"`
	ClusterName         string   `json:"clusterName"`
	SupportsCache       bool     `json:"supportsCache"`
	SupportsIpv6Origins bool     `json:"supportsIpv6Origins"`
	SupportsMtls        bool     `json:"supportsMtls"`
	SupportedCdns       []string `json:"supportedCdns"`
}

// Clusters Response
type ClustersResponse struct {
	Clusters []Cluster         `json:"clusters"`
	Status   OgoResponseStatus `json:"status"`
	HasError bool              `json:"hasError"`
	Count    int               `json:"count"`
}

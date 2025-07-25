package ogosecurity

// Sites Response
type SitesResponse struct {
	Sites      []Site            `json:"sites"`
	SitesItems []Site            `json:"items"`
	Status     OgoResponseStatus `json:"status"`
	HasError   bool              `json:"hasError"`
	Count      int               `json:"count"`
}

// Site Response
type SiteResponse struct {
	Site     Site              `json:"site"`
	Status   OgoResponseStatus `json:"status"`
	HasError bool              `json:"hasError"`
	Count    int               `json:"count"`
}

// Site
type Site struct {
	Name              string `json:"name"`
	ClusterName       string `json:"clusterName,omitempty"`
	DestHost          string `json:"destHost"`
	DestHostScheme    string `json:"destHostScheme"`
	Port              int    `json:"port,omitempty"`
	TrustSelfSigned   bool   `json:"trustSelfSigned"`
	NoCopyXForwarded  bool   `json:"noCopyXForwarded"`
	ForceHttps        bool   `json:"forceHttps"`
	DryRun            bool   `json:"dryRun"`
	PanicMode         bool   `json:"panicMode"`
	Hsts              string `json:"hsts,omitempty"`
	LogExport         bool   `json:"logExport,omitempty"`
	DestHostMtls      bool   `json:"destHostMtls,omitempty"`
	TlsOptionsUid     string `json:"tlsOptionsUid,omitempty"`
	PassTlsClientCert string `json:"passTlsClientCert,omitempty"`
	//Tags              []string `json:"tags,omitempty"`
}

// Ogo Reponse Status
type OgoResponseStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Site Create Query
type SiteQuery struct {
	Action string `json:"action,omitempty"`
	Site   Site   `json:"site"`
}

// Clusters Response
type ClustersResponse struct {
	Cluster      Cluster  `json:"cluster"`
	Role         string   `json:"role"`
	AccessRights []string `json:"accessRights"`
}

// Cluster
type Cluster struct {
	Id                  int      `json:"id"`
	ClusterName         string   `json:"name"`
	SupportsCache       bool     `json:"supportsCache"`
	SupportsIpv6Origins bool     `json:"supportsIpv6Origins"`
	SupportsMtls        bool     `json:"supportsMtls"`
	SupportedCdns       []string `json:"supportedCdns"`
}

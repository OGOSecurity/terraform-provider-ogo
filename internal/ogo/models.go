package ogosecurity

// Ogo Reponse Status
type OgoResponseStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Blacklist Countries
type BlacklistedCountry struct {
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryNameEn"`
}

// Blacklist Countries Response
type BlacklistedCountryResponse struct {
	BlacklistedCountries []BlacklistedCountry `json:"blacklistedCountries"`
	Status               OgoResponseStatus    `json:"status"`
	HasError             bool                 `json:"hasError"`
	Count                int                  `json:"count"`
}

// Site
type Site struct {
	Name                 string               `json:"name"`
	ClusterName          string               `json:"clusterName,omitempty"`
	DestHost             string               `json:"destHost"`
	DestHostScheme       string               `json:"destHostScheme"`
	TrustSelfSigned      bool                 `json:"trustSelfSigned"`
	NoCopyXForwarded     bool                 `json:"noCopyXForwarded"`
	ForceHttps           bool                 `json:"forceHttps"`
	DryRun               bool                 `json:"dryRun"`
	PanicMode            bool                 `json:"panicMode"`
	Hsts                 string               `json:"hsts,omitempty"`
	LogExport            bool                 `json:"logExport,omitempty"`
	DestHostMtls         bool                 `json:"destHostMtls,omitempty"`
	TlsOptionsUid        string               `json:"tlsOptionsUid,omitempty"`
	PassTlsClientCert    string               `json:"passTlsClientCert,omitempty"`
	BlacklistedCountries []BlacklistedCountry `json:"blacklistedCountries,omitempty"`
	//Tags              []string `json:"tags,omitempty"`
	//Port              int    `json:"port,omitempty"`
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

// All Sites Response
type AllSitesResponse struct {
	Sites      []Site            `json:"sites"`
	SitesItems []Site            `json:"items"`
	Status     OgoResponseStatus `json:"status"`
	HasError   bool              `json:"hasError"`
	Count      int               `json:"count"`
}

// TLS Options
type TlsOptions struct {
	Name              string   `json:"name"`
	ClientAuthType    string   `json:"clientAuthType,omitempty"`
	ClientAuthCaCerts []string `json:"clientAuthCaCerts,omitempty"`
	MinTlsVersion     string   `json:"minTlsVersion,omitempty"`
	MaxTlsVersion     string   `json:"maxTlsVersion,omitempty"`
	Uid               string   `json:"uid,omitempty"`
}

// TLS Options Create Query
type TlsOptionsQuery struct {
	TlsOptions TlsOptions `json:"tlsOptions"`
}

// TLS Options Response
type TlsOptionsResponse struct {
	TlsOptions TlsOptions        `json:"tlsOptions"`
	Status     OgoResponseStatus `json:"status"`
	HasError   bool              `json:"hasError"`
}

// All TLS Options Response
type AllTlsOptionsResponse struct {
	TlsOptions []TlsOptions      `json:"items"`
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

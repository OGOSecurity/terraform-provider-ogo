package ogosecurity

// Site
type Site struct {
	Name                 string          `json:"name"`
	ClusterName          string          `json:"clusterName,omitempty"`
	Cluster              Cluster         `json:"cluster,omitempty"`
	DestHost             string          `json:"destHost"`
	DestHostScheme       string          `json:"destHostScheme"`
	Port                 int             `json:"port,omitempty"`
	TrustSelfSigned      bool            `json:"trustSelfSigned"`
	NoCopyXForwarded     bool            `json:"noCopyXForwarded"`
	ForceHttps           bool            `json:"forceHttps"`
	DryRun               bool            `json:"dryRun"`
	PanicMode            bool            `json:"panicMode"`
	Hsts                 string          `json:"hsts,omitempty"`
	LogExport            bool            `json:"logExport,omitempty"`
	DestHostMtls         bool            `json:"destHostMtls,omitempty"`
	PassTlsClientCert    string          `json:"passTlsClientCert,omitempty"`
	TlsOptionsUid        string          `json:"tlsOptionsUid,omitempty"`
	TlsOptions           *TlsOptions     `json:"tlsOptions,omitempty"`
	BlacklistedCountries []string        `json:"blacklistedCountries,omitempty"`
	UrlExceptions        []UrlException  `json:"urlExceptions,omitempty"`
	RewriteRules         []RewriteRule   `json:"rewriteRules,omitempty"`
	Rules                []Rule          `json:"rules,omitempty"`
	WhitelistedIps       []WhitelistedIp `json:"whitelistedIps,omitempty"`
	Tags                 []string        `json:"tags,omitempty"`
}

// Blacklist Countries
type BlacklistedCountry struct {
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryNameEn"`
}

// Url Exceptions
type UrlException struct {
	Path    string `json:"path"`
	Comment string `json:"comment"`
}

// Rewrite Rules
type RewriteRule struct {
	Id                 int    `json:"id,omitempty"`
	Priority           int    `json:"priority"`
	Active             bool   `json:"active"`
	Comment            string `json:"comment,omitempty"`
	RewriteSource      string `json:"rewriteSource"`
	RewriteDestination string `json:"rewriteDestination"`
}

// Rule
type Rule struct {
	Id             int      `json:"id,omitempty"`
	Priority       int      `json:"priority"`
	Active         bool     `json:"active"`
	Action         string   `json:"action"`
	Cache          bool     `json:"cache"`
	Comment        string   `json:"comment,omitempty"`
	Paths          []string `json:"paths"`
	WhitelistedIps []string `json:"whitelistedIps"`
}

// Whitelisted IPs
type WhitelistedIp struct {
	Ip      string `json:"ip"`
	Comment string `json:"comment"`
}

// TLS Options
type TlsOptions struct {
	Name              string   `json:"name,omitempty"`
	ClientAuthType    string   `json:"clientAuthType,omitempty"`
	ClientAuthCaCerts []string `json:"clientAuthCaCerts,omitempty"`
	MinTlsVersion     string   `json:"minTlsVersion,omitempty"`
	MaxTlsVersion     string   `json:"maxTlsVersion,omitempty"`
	Uid               string   `json:"uid,omitempty"`
}

// All TLS Options Response
type TlsOptionsResponse struct {
	TlsOptions []TlsOptions `json:"content"`
	Count      int          `json:"totalElements"`
}

// Cluster
type Cluster struct {
	Id                  string   `json:"clusterId"`
	ClusterName         string   `json:"name"`
	SupportsCache       bool     `json:"supportsCache"`
	SupportsIpv6Origins bool     `json:"supportsIpv6Origins"`
	SupportsMtls        bool     `json:"supportsMtls"`
	SupportedCdns       []string `json:"supportedCdns"`
}

// All Clusters Response
type ClustersResponse struct {
	Cluster      Cluster  `json:"cluster"`
	Role         string   `json:"role"`
	AccessRights []string `json:"accessRights"`
}

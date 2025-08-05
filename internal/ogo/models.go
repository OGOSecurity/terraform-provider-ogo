package ogosecurity

// Site
type Site struct {
	Name                 string          `json:"name"`
	Cluster              Cluster         `json:"cluster,omitempty"`
	DestHost             string          `json:"destHost"`
	DestHostScheme       string          `json:"destHostScheme"`
	DestHostMtls         bool            `json:"destHostMtls"`
	Port                 *int32          `json:"port"`
	TrustSelfSigned      bool            `json:"trustSelfSigned"`
	NoCopyXForwarded     bool            `json:"noCopyXForwarded"`
	ForceHttps           bool            `json:"forceHttps"`
	DryRun               bool            `json:"dryRun"`
	PanicMode            bool            `json:"panicMode"`
	Hsts                 string          `json:"hsts,omitempty"`
	LogExport            bool            `json:"logExport"`
	PassTlsClientCert    string          `json:"passTlsClientCert,omitempty"`
	TlsOptions           *TlsOptions     `json:"tlsOptions"`
	BlacklistedCountries []string        `json:"blacklistedCountries"`
	UrlExceptions        []UrlException  `json:"urlExceptions"`
	RewriteRules         []RewriteRule   `json:"rewriteRules"`
	Rules                []Rule          `json:"rules"`
	WhitelistedIps       []WhitelistedIp `json:"whitelistedIps,omitempty"`
	Tags                 []string        `json:"tags"`
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
	Comment            string `json:"comment"`
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
	Comment        string   `json:"comment"`
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
	Uid                 string   `json:"clusterId"`
	Name                string   `json:"name"`
	Host4               string   `json:"ip"`
	Host6               string   `json:"ip6"`
	SupportsCache       bool     `json:"supportsCache"`
	SupportsIpv6Origins bool     `json:"supportsIpv6Origins"`
	SupportsMtls        bool     `json:"supportsMtls"`
	IpsToWhitelist      []string `json:"ipsToWhitelist"`
	SupportedCdns       []string `json:"supportedCdns"`
}

// All Clusters Response
type ClustersResponse struct {
	Cluster      Cluster  `json:"cluster"`
	Role         string   `json:"role"`
	AccessRights []string `json:"accessRights"`
}

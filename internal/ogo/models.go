package ogosecurity

// Site
type Site struct {
	DomainName           string         `json:"domainName"`
	Cluster              Cluster        `json:"cluster,omitempty"`
	OriginServer         string         `json:"originServer"`
	OriginScheme         string         `json:"originScheme"`
	OriginPort           *int32         `json:"originPort"`
	OriginSkipCertVerify bool           `json:"originSkipCertVerify"`
	OriginMtlsEnabled    bool           `json:"originMtlsEnabled"`
	RemoveXForwarded     bool           `json:"removeXForwarded"`
	LogExportEnabled     bool           `json:"logExportEnabled"`
	ForceHttps           bool           `json:"forceHttps"`
	AuditMode            bool           `json:"auditMode"`
	PassthroughMode      bool           `json:"passthroughMode"`
	Hsts                 string         `json:"hsts,omitempty"`
	PassTlsClientCert    string         `json:"passTlsClientCert,omitempty"`
	TlsOptions           *TlsOptions    `json:"tlsOptions"`
	BlacklistedCountries []string       `json:"blacklistedCountries"`
	IpExceptions         []IpException  `json:"ipExceptions,omitempty"`
	UrlExceptions        []UrlException `json:"urlExceptions"`
	RewriteRules         []RewriteRule  `json:"rewriteRules"`
	Rules                []Rule         `json:"rules"`
	Tags                 []string       `json:"tags"`
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
	Priority           int    `json:"priority"`
	Active             bool   `json:"active"`
	Comment            string `json:"comment"`
	RewriteSource      string `json:"rewriteSource"`
	RewriteDestination string `json:"rewriteDestination"`
}

// Rule
type Rule struct {
	Priority       int      `json:"priority"`
	Active         bool     `json:"active"`
	Action         string   `json:"action"`
	Cache          bool     `json:"cache"`
	Comment        string   `json:"comment"`
	Paths          []string `json:"paths"`
	WhitelistedIps []string `json:"whitelistedIps"`
}

// IP Exception
type IpException struct {
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

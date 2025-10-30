// Copyright (c) OGO Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package ogosecurity

// Site objects.
type Site struct {
	DomainName                string                     `json:"domainName"`
	Cluster                   Cluster                    `json:"cluster,omitempty"`
	Contract                  *Contract                  `json:"contract"`
	OriginServer              string                     `json:"originServer"`
	OriginScheme              string                     `json:"originScheme"`
	OriginPort                *int32                     `json:"originPort"`
	OriginSkipCertVerify      bool                       `json:"originSkipCertVerify"`
	OriginMtlsEnabled         bool                       `json:"originMtlsEnabled"`
	RemoveXForwarded          bool                       `json:"removeXForwarded"`
	LogExportEnabled          bool                       `json:"logExportEnabled"`
	CacheEnabled              bool                       `json:"cacheEnabled"`
	Status                    string                     `json:"status"`
	ActiveCustomerCertificate *ActiveCustomerCertificate `json:"activeCustomerCertificate"`
	Cdn                       *string                    `json:"cdn"`
	CdnStatus                 *string                    `json:"cdnStatus"`
	ForceHttps                bool                       `json:"forceHttps"`
	AuditMode                 bool                       `json:"auditMode"`
	PassthroughMode           bool                       `json:"passthroughMode"`
	Hsts                      string                     `json:"hsts,omitempty"`
	PassTlsClientCert         string                     `json:"passTlsClientCert,omitempty"`
	TlsOptions                *TlsOptions                `json:"tlsOptions"`
	BlacklistedCountries      []string                   `json:"blacklistedCountries"`
	BrainOverrides            map[string]float64         `json:"brainOverrides"`
	IpExceptions              []IpException              `json:"ipExceptions,omitempty"`
	UrlExceptions             []UrlException             `json:"urlExceptions"`
	RewriteRules              []RewriteRule              `json:"rewriteRules"`
	Rules                     []Rule                     `json:"rules"`
	Tags                      []string                   `json:"tags"`
}

type BlacklistedCountry struct {
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryNameEn"`
}

type UrlException struct {
	Path    string `json:"path"`
	Comment string `json:"comment"`
}

type RewriteRule struct {
	Active             bool   `json:"active"`
	Comment            string `json:"comment"`
	RewriteSource      string `json:"rewriteSource"`
	RewriteDestination string `json:"rewriteDestination"`
}

type Rule struct {
	Active         bool     `json:"active"`
	Action         string   `json:"action"`
	Cache          bool     `json:"cache"`
	Comment        string   `json:"comment"`
	Paths          []string `json:"paths"`
	WhitelistedIps []string `json:"whitelistedIps"`
}

type IpException struct {
	Ip      string `json:"ip"`
	Comment string `json:"comment"`
}

// TLS options objects.
type TlsOptions struct {
	Name              string   `json:"name,omitempty"`
	ClientAuthType    string   `json:"clientAuthType,omitempty"`
	ClientAuthCaCerts []string `json:"clientAuthCaCerts,omitempty"`
	MinTlsVersion     *string  `json:"minTlsVersion,omitempty"`
	MaxTlsVersion     *string  `json:"maxTlsVersion,omitempty"`
	Uid               string   `json:"uid,omitempty"`
}

type TlsOptionsResponse struct {
	TlsOptions []TlsOptions `json:"content"`
	Count      int          `json:"totalElements"`
}

// Cluster objects.
type Cluster struct {
	Uid                 string   `json:"clusterId"`
	Name                string   `json:"name"`
	Entrypoint4         string   `json:"ip"`
	Entrypoint6         string   `json:"ipv6"`
	EntrypointCdn       string   `json:"cdnCname"`
	SupportsCache       bool     `json:"supportsCache"`
	SupportsIpv6Origins bool     `json:"supportsIpv6Origins"`
	SupportsMtls        bool     `json:"supportsMtls"`
	IpsToWhitelist      []string `json:"ipsToWhitelist"`
	SupportedCdns       []string `json:"supportedCdns"`
}

type ClustersResponse struct {
	Cluster      Cluster  `json:"cluster"`
	Role         string   `json:"role"`
	AccessRights []string `json:"accessRights"`
}

// Organization objects.
type Organization struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type OrganizationDetails struct {
	Organization Organization `json:"organization"`
	Role         string       `json:"role"`
	Privileges   []string     `json:"privileges"`
}

type OrganizationsResponse struct {
	OrganizationDetails []OrganizationDetails `json:"content"`
	Count               int                   `json:"totalElements"`
}

// Contract objects.
type Contract struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type ContractHolder struct {
	Code        string `json:"code"`
	CompanyName string `json:"companyName"`
}

type ContractsResponse struct {
	Contracts []Contract `json:"content"`
	Count     int        `json:"totalElements"`
}

// Certificate objects.
type CertificateP12 struct {
	Data     string `json:"data"`
	Password string `json:"password"`
}

type ActiveCustomerCertificate struct {
	P12       CertificateP12 `json:"p12"`
	Cn        string         `json:"cn"`
	ExpiredAt string         `json:"expiredAt"`
	Hash      string         `json:"hash"`
}

type Certificate struct {
	Id            int32  `json:"pid"`
	Active        bool   `json:"active"`
	Cn            string `json:"cn"`
	Csr           string `json:"csr"`
	FullChainCert string `json:"fullChainCert"`
	Type          string `json:"type"`
	CreatedAt     string `json:"createdAt"`
	ExpiredAt     string `json:"expiredAt"`
	UpdatedAt     string `json:"updatedAt"`
	Error         string `json:"error"`
}

type CertificatesResponse struct {
	Certificates []Certificate `json:"content"`
	Count        int           `json:"totalElements"`
}

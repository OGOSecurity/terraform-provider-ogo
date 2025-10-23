// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	ogosecurity "terraform-provider-ogo/internal/ogo"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

// SiteResourceModel maps the resource schema data.
type SiteResourceModel struct {
	DomainName                types.String                    `tfsdk:"domain_name"`
	ClusterUid                types.String                    `tfsdk:"cluster_uid"`
	ClusterEntrypoint4        types.String                    `tfsdk:"cluster_entrypoint_4"`
	ClusterEntrypoint6        types.String                    `tfsdk:"cluster_entrypoint_6"`
	ClusterEntrypointCdn      types.String                    `tfsdk:"cluster_entrypoint_cdn"`
	ContractNumber            types.String                    `tfsdk:"contract_number"`
	OriginServer              types.String                    `tfsdk:"origin_server"`
	OriginScheme              types.String                    `tfsdk:"origin_scheme"`
	OriginPort                types.Int32                     `tfsdk:"origin_port"`
	OriginSkipCertVerify      types.Bool                      `tfsdk:"origin_skip_cert_verify"`
	OriginMtlsEnabled         types.Bool                      `tfsdk:"origin_mtls_enabled"`
	RemoveXForwarded          types.Bool                      `tfsdk:"remove_xforwarded"`
	LogExportEnabled          types.Bool                      `tfsdk:"log_export_enabled"`
	CacheEnabled              types.Bool                      `tfsdk:"cache_enabled"`
	Status                    types.String                    `tfsdk:"status"`
	Cdn                       types.String                    `tfsdk:"cdn"`
	CdnStatus                 types.String                    `tfsdk:"cdn_status"`
	ForceHttps                types.Bool                      `tfsdk:"force_https"`
	AuditMode                 types.Bool                      `tfsdk:"audit_mode"`
	PassthroughMode           types.Bool                      `tfsdk:"passthrough_mode"`
	Hsts                      types.String                    `tfsdk:"hsts"`
	PassTlsClientCert         types.String                    `tfsdk:"pass_tls_client_cert"`
	TlsOptionsUid             types.String                    `tfsdk:"tlsoptions_uid"`
	BrainOverrides            types.Map                       `tfsdk:"brain_overrides"`
	ActiveCustomerCertificate *ActiveCustomerCertificateModel `tfsdk:"active_customer_certificate"`
	BlacklistedCountries      []types.String                  `tfsdk:"blacklisted_countries"`
	IpExceptions              []IpExceptionModel              `tfsdk:"ip_exceptions"`
	UrlExceptions             []UrlExceptionModel             `tfsdk:"url_exceptions"`
	RewriteRules              []RewriteRuleModel              `tfsdk:"rewrite_rules"`
	Rules                     []RuleModel                     `tfsdk:"rules"`
	Tags                      []types.String                  `tfsdk:"tags"`
	LastUpdated               types.String                    `tfsdk:"last_updated"`
}

type ActiveCustomerCertificateModel struct {
	Cn           types.String `tfsdk:"cn"`
	ExpiredAt    types.String `tfsdk:"expired_at"`
	Hash         types.String `tfsdk:"hash"`
	P12File      types.String `tfsdk:"p12_file"`
	P12Content64 types.String `tfsdk:"p12_content64"`
	P12Password  types.String `tfsdk:"p12_password"`
}

type RewriteRuleModel struct {
	Active             types.Bool   `tfsdk:"active"`
	Comment            types.String `tfsdk:"comment"`
	RewriteSource      types.String `tfsdk:"rewrite_source"`
	RewriteDestination types.String `tfsdk:"rewrite_destination"`
}

type RuleModel struct {
	Active         types.Bool     `tfsdk:"active"`
	Action         types.String   `tfsdk:"action"`
	Cache          types.Bool     `tfsdk:"cache"`
	Comment        types.String   `tfsdk:"comment"`
	Paths          []types.String `tfsdk:"paths"`
	WhitelistedIps []types.String `tfsdk:"whitelisted_ips"`
}

type UrlExceptionModel struct {
	Path    types.String `tfsdk:"path"`
	Comment types.String `tfsdk:"comment"`
}

type IpExceptionModel struct {
	Ip      types.String `tfsdk:"ip"`
	Comment types.String `tfsdk:"comment"`
}

// NewSiteResource is a helper function to simplify the provider implementation.
func NewSiteResource() resource.Resource {
	return &siteResource{}
}

// siteResource is the resource implementation.
type siteResource struct {
	client *ogosecurity.Client
}

// Metadata returns the resource type name.
func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shield_site"
}

// Schema defines the schema for the resource.
func (r *siteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "DNS domain name of the site.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_uid": schema.StringAttribute{
				Required:    true,
				Description: "Cluster UID on which site is deployed (force site recreation if modified). List of available cluster and associated UID can be retrieved from `ogo_shield_clusters` data source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_entrypoint_4": schema.StringAttribute{
				Computed:    true,
				Description: "IPv4 cluster entrypoint to which the site DNS record can be configured.",
			},
			"cluster_entrypoint_6": schema.StringAttribute{
				Computed:    true,
				Description: "IPv6 cluster entrypoint to which site the DNS record can be configured.",
			},
			"cluster_entrypoint_cdn": schema.StringAttribute{
				Computed:    true,
				Description: "CDN entrypoint to which the site DNS record can be configured.",
			},
			"contract_number": schema.StringAttribute{
				Optional:    true,
				Description: "Contract number to which the site is attached, only required if multiple contracts exist for this organization. List of available contracts can be retrieved from `ogo_shield_contrats` data source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"origin_server": schema.StringAttribute{
				Required:    true,
				Description: "Origin server address (IP address or domain name).",
			},
			"origin_scheme": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Scheme used to access the origin server. Supported values: **https** or **http** (default: **https**).",
				Default:     stringdefault.StaticString("https"),
				Validators: []validator.String{
					stringvalidator.OneOf("http", "https"),
				},
			},
			"origin_port": schema.Int32Attribute{
				Optional:    true,
				Description: "Port to be used to access the origin server. Must be defined only if different from standard HTTP port 443 or 80, otherwise let Ogo choose the correct port.",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"origin_mtls_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable mTLS between Ogo and the origin server (default: **false**).",
				Default:     booldefault.StaticBool(false),
			},
			"origin_skip_cert_verify": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Description: "Skip origin server certificate verification if TLS is used. " +
					"If set to **true** Ogo accepts connection to the origin server even if the " +
					"certificate doesn't match site domain name, or the certificate is expired, " +
					"or the certificate is self signed (default: **false**).",
				Default: booldefault.StaticBool(false),
			},
			"remove_xforwarded": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Remove X-Forwarded-* headers. (default: **false**).",
				Default:     booldefault.StaticBool(false),
			},
			"force_https": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Redirect HTTP request to HTTPS (default: **false**).",
				Default:     booldefault.StaticBool(false),
			},
			"audit_mode": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable audit mode. Requests are analyzed by Ogo Shield but never blocked (default: **false**).",
				Default:     booldefault.StaticBool(false),
			},
			"passthrough_mode": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable passthrough mode. Requests are not analyzed by Ogo Shield and never blocked (default: **false**).",
				Default:     booldefault.StaticBool(false),
			},
			"hsts": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Enable HSTS (default: **hsts**). Supported values:\n" +
					" * **hsts**: Enable HSTS\n" +
					"  * **hstss**: Enable HSTS including subdomains\n" +
					"  * **hstsp**: Enable HSTS including subdomains and preloading\n" +
					"  * **none**: Disable HSTS.",
				Default: stringdefault.StaticString("hsts"),
				Validators: []validator.String{
					stringvalidator.OneOf("hsts", "hstss", "hstssp", "none"),
				},
			},
			"log_export_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable log export for this site (default: **false**).",
				Default:     booldefault.StaticBool(false),
			},
			"cache_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable cache for this site if supported by cluster (default: **false**).",
				Default:     booldefault.StaticBool(false),
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Get site status. Available state:\n" +
					"  * **CREATED**: site just created, waiting for DNS to be configured for site domain name to redirect to Ogo Shield cluster.\n" +
					"  * **DNS_ERROR**: Partial DNS configuration.\n" +
					"  * **ONLINE**: Site is properly provisioned, but no active TLS certificate is present (ephemeral state).\n" +
					"  * **LE_CERT**: Site is properly provisioned, DNS is configured on Ogo Shield cluster and site is protected with a Let's Encrypt certificate.\n" +
					"  * **CUST_CERT**: Site is properly provisioned, DNS is configured on Ogo Shield cluster and site is protected with a customer certificate.\n" +
					"  * **LE_EXP**: Site is properly provisioned, DNS is configured on Ogo Shield cluster and site is protected with an expired Let's Encrypt certificate.\n" +
					"  * **CUST_EXP**: Site is properly provisioned, DNS is configured on Ogo Shield cluster and site is protected with an expired customer certificate.\n" +
					"  * **OFFLINE**: Site has been properly provisioned and configured on Ogo Shield cluster, but DNS no longer redirects on Ogo Shield cluster.",
				Validators: []validator.String{
					stringvalidator.OneOf("CREATED", "DNS_ERROR", "ONLINE", "LE_CERT", "CUST_CERT", "LE_EXP", "CUST_EXP", "OFFLINE"),
				},
			},
			"cdn": schema.StringAttribute{
				Optional:    true,
				Description: "Select CDN to be used for this site if supported by cluster.",
				Validators: []validator.String{
					stringvalidator.OneOf("ORANGE"),
				},
			},
			"cdn_status": schema.StringAttribute{
				Computed: true,
				Description: "Retrieve CDN status if CDN is enabled.\n Available state:\n" +
					"  * **ACTIVE**: CDN is enabled and ready to be used.\n" +
					"  * **ACTIVATION_IN_PROGRESS**: CDN is enabled and first activation is in progress.\n" +
					"  * **SYNC_IN_PROGRESS**: New configuration is waiting to be applied on CDN.",
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "ACTIVATION_IN_PROGRESS", "SYNC_IN_PROGRESS"),
				},
			},
			"tlsoptions_uid": schema.StringAttribute{
				Optional:    true,
				Description: "UID of TLS options to be applied to this site. List of available TLS options and associated UID can be retrieved from `ogo_shield_tlsoptions` data source.",
			},
			"pass_tls_client_cert": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Client certificate informations to pass to the origin server (default: **info**). Supported values:\n" +
					"  * **all**: Send the certificate and certificate informations\n" +
					"  * **cert**: Send only the certificate\n" +
					"  * **info**: Send only the certificate information\n" +
					"  * **none***: Nothing sent.",
				Default: stringdefault.StaticString("info"),
				Validators: []validator.String{
					stringvalidator.OneOf("all", "cert", "info", "none"),
				},
			},
			"active_customer_certificate": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "P12/PFX certificate to be used for this site.",
				Attributes: map[string]schema.Attribute{
					"cn": schema.StringAttribute{
						Computed:    true,
						Description: "Common name of this certificate.",
					},
					"expired_at": schema.StringAttribute{
						Computed:    true,
						Description: "Expiration date of this certificate.",
					},
					"hash": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Hash of the certificate generated from P12 content.",
					},
					"p12_file": schema.StringAttribute{
						Optional:    true,
						Description: "P12/PFX file path containing certificate and key (conflicts with `p12_content64`).",
					},
					"p12_content64": schema.StringAttribute{
						Optional:    true,
						Description: "P12/PFX content encoded in base64 (conflicts with `p12_file`).",
					},
					"p12_password": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "Password used to decrypt P12 file.",
					},
				},
			},
			"blacklisted_countries": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List of countries to blacklist.",
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"brain_overrides": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List of brain parameters to override",
				ElementType: types.Float64Type,
			},
			"ip_exceptions": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Passthrough mode for IPs. Requests from those IPs will never be blocked.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Required:    true,
							Description: "IP address list never blocked by Ogo Shield.",
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated with this IP list.",
						},
					},
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"ip":      types.StringType,
								"comment": types.StringType,
							},
						},
						[]attr.Value{},
					),
				),
			},
			"rewrite_rules": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Rewrite a path of your website. Rewrite rules are parsed in order of declaration.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Flag to enable (**true**) or disable (**false**) rewrite rule. (default: **true**).",
							Default:     booldefault.StaticBool(true),
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated with this rewrite rule.",
						},
						"rewrite_source": schema.StringAttribute{
							Required:    true,
							Description: "Source path to be rewritten.",
						},
						"rewrite_destination": schema.StringAttribute{
							Required:    true,
							Description: "Rewritten destination path.",
						},
					},
				},
				Default: listdefault.StaticValue(
					types.ListValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"active":              types.BoolType,
								"comment":             types.StringType,
								"rewrite_source":      types.StringType,
								"rewrite_destination": types.StringType,
							},
						},
						[]attr.Value{},
					),
				),
			},
			"rules": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Restrict access to given URLs. Rules are parsed in order of declaration. The engine stops at the first URL match.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Flag to enable (**true**) or disable (**false**) rule. (default: **true**).",
							Default:     booldefault.StaticBool(true),
						},
						"action": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Description: "Action to be applied when the rule matches (default: **brain**). Supported values:\n" +
								"  * **brain**: Rule analyzed by Ogo Shield brain\n" +
								"  * **bypass**: Rule not analyzed by Ogo Shield brain.",
							Default: stringdefault.StaticString("brain"),
							Validators: []validator.String{
								stringvalidator.OneOf("brain", "bypass"),
							},
						},
						"cache": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Enable or disable caching on this rule. Option can be used only if site caching is enabled. (default: **false**).",
							Default:     booldefault.StaticBool(false),
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated with this rule.",
						},
						"paths": schema.SetAttribute{
							Required:    true,
							Description: "List of URL paths for which the rule is applied.",
							ElementType: types.StringType,
						},
						"whitelisted_ips": schema.SetAttribute{
							Required:    true,
							Description: "Authorized IP addresses list.",
							ElementType: types.StringType,
						},
					},
				},
				Default: listdefault.StaticValue(
					types.ListValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"active":          types.BoolType,
								"action":          types.StringType,
								"cache":           types.BoolType,
								"comment":         types.StringType,
								"paths":           types.SetType{ElemType: types.StringType},
								"whitelisted_ips": types.SetType{ElemType: types.StringType},
							},
						},
						[]attr.Value{},
					),
				),
			},
			"url_exceptions": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Passthrough mode on URL regular expressions. The matching requests will never be blocked.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Required:    true,
							Description: "Path of the URL never blocked by Ogo Shield.",
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated with this URL exception.",
						},
					},
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"path":    types.StringType,
								"comment": types.StringType,
							},
						},
						[]attr.Value{},
					),
				),
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List of tags.",
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Last resource updated by Terraform.",
			},
		},
		MarkdownDescription: "Resource `ogo_shield_site` can be used to create, " +
			"update or delete sites configuration in Ogo Dashboard.\n\n" +
			"This resource allow to manage all site settings.\n\n" +
			"`ogo_shield_clusters` and `ogo_shield_tlsoptions` data sources can be " +
			"used to retrieve UID needed in `ogo_shield_site` resource configuration.\n\n",
	}
}

// Configure adds the provider configured client to the resource.
func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ogosecurity.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"unexpected resource configure type",
			fmt.Sprintf("Expected *ogosecurity.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan SiteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new site
	var tlsOpt *ogosecurity.TlsOptions
	s := ogosecurity.Site{
		DomainName: string(plan.DomainName.ValueString()),
		Cluster: ogosecurity.Cluster{
			Uid: string(plan.ClusterUid.ValueString()),
		},
		OriginServer:         string(plan.OriginServer.ValueString()),
		OriginScheme:         string(plan.OriginScheme.ValueString()),
		OriginMtlsEnabled:    bool(plan.OriginMtlsEnabled.ValueBool()),
		OriginSkipCertVerify: bool(plan.OriginSkipCertVerify.ValueBool()),
		RemoveXForwarded:     bool(plan.RemoveXForwarded.ValueBool()),
		ForceHttps:           bool(plan.ForceHttps.ValueBool()),
		AuditMode:            bool(plan.AuditMode.ValueBool()),
		PassthroughMode:      bool(plan.PassthroughMode.ValueBool()),
		Hsts:                 string(plan.Hsts.ValueString()),
		LogExportEnabled:     bool(plan.LogExportEnabled.ValueBool()),
		CacheEnabled:         bool(plan.CacheEnabled.ValueBool()),
		Status:               plan.Status.ValueString(),
		Cdn:                  plan.Cdn.ValueStringPointer(),
		OriginPort:           plan.OriginPort.ValueInt32Pointer(),
		PassTlsClientCert:    string(plan.PassTlsClientCert.ValueString()),
		TlsOptions:           tlsOpt,
	}

	// Certificate
	if plan.ActiveCustomerCertificate != nil {
		if plan.ActiveCustomerCertificate.P12File.ValueString() != "" &&
			plan.ActiveCustomerCertificate.P12Content64.ValueString() != "" {
			resp.Diagnostics.AddError(
				"attribute conflicts in active_customer_certificate",
				"p12_file and p12_content64 attribute can't be used at same time",
			)
			return
		}

		p12_content64 := ""
		if plan.ActiveCustomerCertificate.P12File.ValueString() != "" {
			p12Data, err := os.ReadFile(string(plan.ActiveCustomerCertificate.P12File.ValueString()))
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to read P12/PFX file",
					"Could not read P12/PFX file, unexpected error: "+string(err.Error()),
				)
				return
			}
			p12_content64 = base64.StdEncoding.EncodeToString(p12Data)
		} else if plan.ActiveCustomerCertificate.P12Content64.ValueString() != "" {
			p12_content64 = plan.ActiveCustomerCertificate.P12Content64.ValueString()
		}

		s.ActiveCustomerCertificate = &ogosecurity.ActiveCustomerCertificate{
			P12: ogosecurity.CertificateP12{
				Data:     p12_content64,
				Password: plan.ActiveCustomerCertificate.P12Password.ValueString(),
			},
		}
	}

	// Contract
	if plan.ContractNumber.ValueString() != "" {
		s.Contract = &ogosecurity.Contract{
			Number: string(plan.ContractNumber.ValueString()),
		}
	}

	// TLS Options
	if plan.TlsOptionsUid.ValueString() != "" {
		s.TlsOptions = &ogosecurity.TlsOptions{
			Uid: string(plan.TlsOptionsUid.ValueString()),
		}
	}

	// Blacklist Countries
	s.BlacklistedCountries = []string{}
	for _, country := range plan.BlacklistedCountries {
		s.BlacklistedCountries = append(s.BlacklistedCountries, string(country.ValueString()))
	}

	// Brain parameters overrides
	s.BrainOverrides = make(map[string]float64)
	for brainParam, brainVal := range plan.BrainOverrides.Elements() {
		val, err := strconv.ParseFloat(brainVal.String(), 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating site",
				"Invalid brain parameter float value, "+brainVal.String()+" unexpected error: "+err.Error(),
			)
			return
		}
		s.BrainOverrides[brainParam] = val
	}

	// IP Exceptions
	s.IpExceptions = []ogosecurity.IpException{}
	for _, wlip := range plan.IpExceptions {
		s.IpExceptions = append(s.IpExceptions, ogosecurity.IpException{
			Ip:      string(wlip.Ip.ValueString()),
			Comment: string(wlip.Comment.ValueString()),
		})
	}

	// Rewrite Rules
	s.RewriteRules = []ogosecurity.RewriteRule{}
	for _, rewrite := range plan.RewriteRules {
		s.RewriteRules = append(s.RewriteRules, ogosecurity.RewriteRule{
			Active:             bool(rewrite.Active.ValueBool()),
			Comment:            string(rewrite.Comment.ValueString()),
			RewriteSource:      string(rewrite.RewriteSource.ValueString()),
			RewriteDestination: string(rewrite.RewriteDestination.ValueString()),
		})
	}

	// Rules access
	s.Rules = []ogosecurity.Rule{}
	for _, rule := range plan.Rules {
		r := ogosecurity.Rule{
			Active:         bool(rule.Active.ValueBool()),
			Action:         string(rule.Action.ValueString()),
			Cache:          bool(rule.Cache.ValueBool()),
			Comment:        string(rule.Comment.ValueString()),
			Paths:          []string{},
			WhitelistedIps: []string{},
		}

		for _, path := range rule.Paths {
			r.Paths = append(r.Paths, string(path.ValueString()))
		}

		for _, ip := range rule.WhitelistedIps {
			r.WhitelistedIps = append(r.WhitelistedIps, string(ip.ValueString()))
		}

		s.Rules = append(s.Rules, r)
	}

	// URL Exceptions
	s.UrlExceptions = []ogosecurity.UrlException{}
	for _, url := range plan.UrlExceptions {
		s.UrlExceptions = append(s.UrlExceptions, ogosecurity.UrlException{
			Path:    string(url.Path.ValueString()),
			Comment: string(url.Comment.ValueString()),
		})
	}

	// Tags
	s.Tags = []string{}
	for _, tag := range plan.Tags {
		s.Tags = append(s.Tags, string(tag.ValueString()))
	}

	site, err := r.client.CreateSite(s)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			"Could not create site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ClusterEntrypoint4 = types.StringValue(site.Cluster.Entrypoint4)
	plan.ClusterEntrypoint6 = types.StringValue(site.Cluster.Entrypoint6)
	if plan.ActiveCustomerCertificate != nil {
		plan.ActiveCustomerCertificate.Cn = types.StringValue(site.ActiveCustomerCertificate.Cn)
		plan.ActiveCustomerCertificate.ExpiredAt = types.StringValue(site.ActiveCustomerCertificate.ExpiredAt)
		plan.ActiveCustomerCertificate.Hash = types.StringValue(site.ActiveCustomerCertificate.Hash)
	}
	if plan.Cdn.String() != "" {
		plan.ClusterEntrypointCdn = types.StringValue(site.Cluster.EntrypointCdn)
		plan.CdnStatus = types.StringPointerValue(site.CdnStatus)
	}
	plan.Status = types.StringValue(site.Status)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state SiteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed site value from Ogo
	site, err := r.client.GetSite(state.DomainName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Ogo site",
			"Could not read Ogo site domain name "+state.DomainName.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite properties with refreshed state
	state.ClusterUid = types.StringValue(site.Cluster.Uid)
	state.ClusterEntrypoint4 = types.StringValue(site.Cluster.Entrypoint4)
	state.ClusterEntrypoint6 = types.StringValue(site.Cluster.Entrypoint6)
	state.ClusterEntrypointCdn = types.StringValue(site.Cluster.EntrypointCdn)
	state.OriginServer = types.StringValue(site.OriginServer)
	state.OriginScheme = types.StringValue(site.OriginScheme)
	state.OriginMtlsEnabled = types.BoolValue(site.OriginMtlsEnabled)
	state.OriginPort = types.Int32PointerValue(site.OriginPort)
	state.OriginSkipCertVerify = types.BoolValue(site.OriginSkipCertVerify)
	state.RemoveXForwarded = types.BoolValue(site.RemoveXForwarded)
	state.ForceHttps = types.BoolValue(site.ForceHttps)
	state.AuditMode = types.BoolValue(site.AuditMode)
	state.PassthroughMode = types.BoolValue(site.PassthroughMode)
	state.Hsts = types.StringValue(site.Hsts)
	state.LogExportEnabled = types.BoolValue(site.LogExportEnabled)
	state.CacheEnabled = types.BoolValue(site.CacheEnabled)
	state.Status = types.StringValue(site.Status)
	state.LogExportEnabled = types.BoolValue(site.LogExportEnabled)
	state.PassTlsClientCert = types.StringValue(site.PassTlsClientCert)

	// Activate certificate
	if state.ActiveCustomerCertificate != nil {
		if state.ActiveCustomerCertificate.Hash != types.StringValue(site.ActiveCustomerCertificate.Hash) {
			if state.ActiveCustomerCertificate.P12File.ValueString() != "" {
				state.ActiveCustomerCertificate.P12File = types.StringValue("")
			} else if state.ActiveCustomerCertificate.P12Content64.ValueString() != "" {
				state.ActiveCustomerCertificate.P12Content64 = types.StringValue("")
			}
		}
	}

	// Brain parameters overrides
	state.BrainOverrides, diags = types.MapValueFrom(ctx, types.Float64Type, site.BrainOverrides)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// CDN
	if site.Cdn != nil {
		state.Cdn = types.StringPointerValue(site.Cdn)
		state.CdnStatus = types.StringPointerValue(site.CdnStatus)
		state.ClusterEntrypointCdn = types.StringValue(site.Cluster.EntrypointCdn)
	}

	// Contract
	if site.Contract != nil {
		if site.Contract.Number != "" {
			state.ContractNumber = types.StringValue(site.Contract.Number)
		}
	}

	// TLS Options
	if site.TlsOptions != nil {
		state.TlsOptionsUid = types.StringValue(site.TlsOptions.Uid)
	}

	// Blacklist Countries
	state.BlacklistedCountries = []types.String{}
	for _, country := range site.BlacklistedCountries {
		state.BlacklistedCountries = append(state.BlacklistedCountries, types.StringValue(country))
	}

	// IP Exceptions
	state.IpExceptions = []IpExceptionModel{}
	for _, wlip := range site.IpExceptions {
		state.IpExceptions = append(state.IpExceptions, IpExceptionModel{
			Ip:      types.StringValue(wlip.Ip),
			Comment: types.StringValue(wlip.Comment),
		})
	}

	// Rewrite rules
	state.RewriteRules = []RewriteRuleModel{}
	for _, rewrite := range site.RewriteRules {
		state.RewriteRules = append(state.RewriteRules, RewriteRuleModel{
			Active:             types.BoolValue(rewrite.Active),
			Comment:            types.StringValue(rewrite.Comment),
			RewriteSource:      types.StringValue(rewrite.RewriteSource),
			RewriteDestination: types.StringValue(rewrite.RewriteDestination),
		})
	}

	// Rules access
	state.Rules = []RuleModel{}
	for _, rule := range site.Rules {
		r := RuleModel{
			Active:         types.BoolValue(rule.Active),
			Action:         types.StringValue(rule.Action),
			Cache:          types.BoolValue(rule.Cache),
			Comment:        types.StringValue(rule.Comment),
			Paths:          []types.String{},
			WhitelistedIps: []types.String{},
		}

		for _, path := range rule.Paths {
			r.Paths = append(r.Paths, types.StringValue(path))
		}

		for _, ip := range rule.WhitelistedIps {
			r.WhitelistedIps = append(r.WhitelistedIps, types.StringValue(ip))
		}

		state.Rules = append(state.Rules, r)
	}

	// URL Exceptions
	state.UrlExceptions = []UrlExceptionModel{}
	for _, url := range site.UrlExceptions {
		state.UrlExceptions = append(state.UrlExceptions, UrlExceptionModel{
			Path:    types.StringValue(url.Path),
			Comment: types.StringValue(url.Comment),
		})
	}

	// Tags
	state.Tags = []types.String{}
	for _, tag := range site.Tags {
		state.Tags = append(state.Tags, types.StringValue(tag))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan SiteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new site
	var tlsOpt *ogosecurity.TlsOptions
	s := ogosecurity.Site{
		DomainName: string(plan.DomainName.ValueString()),
		Cluster: ogosecurity.Cluster{
			Uid: string(plan.ClusterUid.ValueString()),
		},
		OriginServer:         string(plan.OriginServer.ValueString()),
		OriginScheme:         string(plan.OriginScheme.ValueString()),
		OriginMtlsEnabled:    bool(plan.OriginMtlsEnabled.ValueBool()),
		OriginSkipCertVerify: bool(plan.OriginSkipCertVerify.ValueBool()),
		RemoveXForwarded:     bool(plan.RemoveXForwarded.ValueBool()),
		ForceHttps:           bool(plan.ForceHttps.ValueBool()),
		AuditMode:            bool(plan.AuditMode.ValueBool()),
		PassthroughMode:      bool(plan.PassthroughMode.ValueBool()),
		Hsts:                 string(plan.Hsts.ValueString()),
		LogExportEnabled:     bool(plan.LogExportEnabled.ValueBool()),
		CacheEnabled:         bool(plan.CacheEnabled.ValueBool()),
		Status:               plan.Status.ValueString(),
		Cdn:                  plan.Cdn.ValueStringPointer(),
		OriginPort:           plan.OriginPort.ValueInt32Pointer(),
		PassTlsClientCert:    string(plan.PassTlsClientCert.ValueString()),
		TlsOptions:           tlsOpt,
	}

	// Certificate
	if plan.ActiveCustomerCertificate != nil {
		if plan.ActiveCustomerCertificate.P12File.ValueString() != "" &&
			plan.ActiveCustomerCertificate.P12Content64.ValueString() != "" {
			resp.Diagnostics.AddError(
				"attribute conflicts in active_customer_certificate",
				"p12_file and p12_content64 attribute can't be used at same time",
			)
			return
		}

		p12_content64 := ""
		if plan.ActiveCustomerCertificate.P12File.ValueString() != "" {
			p12Data, err := os.ReadFile(string(plan.ActiveCustomerCertificate.P12File.ValueString()))
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to read certificate",
					"Could not read certificate file, unexpected error: "+string(err.Error()),
				)
				return
			}
			p12_content64 = base64.StdEncoding.EncodeToString(p12Data)
		} else if plan.ActiveCustomerCertificate.P12Content64.ValueString() != "" {
			p12_content64 = plan.ActiveCustomerCertificate.P12Content64.ValueString()
		}
		s.ActiveCustomerCertificate = &ogosecurity.ActiveCustomerCertificate{
			P12: ogosecurity.CertificateP12{
				Data:     p12_content64,
				Password: plan.ActiveCustomerCertificate.P12Password.ValueString(),
			},
		}
	}

	// Contract
	if plan.ContractNumber.ValueString() != "" {
		s.Contract = &ogosecurity.Contract{
			Number: string(plan.ContractNumber.ValueString()),
		}
	}

	// TLS Options
	if plan.TlsOptionsUid.ValueString() != "" {
		s.TlsOptions = &ogosecurity.TlsOptions{
			Uid: string(plan.TlsOptionsUid.ValueString()),
		}
	}

	// Blacklist Countries
	s.BlacklistedCountries = []string{}
	for _, country := range plan.BlacklistedCountries {
		s.BlacklistedCountries = append(s.BlacklistedCountries, string(country.ValueString()))
	}

	// Brain parameters overrides
	s.BrainOverrides = make(map[string]float64)
	for brainParam, brainVal := range plan.BrainOverrides.Elements() {
		val, err := strconv.ParseFloat(brainVal.String(), 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating site",
				"Invalid brain parameter float value, "+brainVal.String()+" unexpected error: "+err.Error(),
			)
			return
		}
		s.BrainOverrides[brainParam] = val
	}

	// IP Exceptions
	s.IpExceptions = []ogosecurity.IpException{}
	for _, wlip := range plan.IpExceptions {
		s.IpExceptions = append(s.IpExceptions, ogosecurity.IpException{
			Ip:      string(wlip.Ip.ValueString()),
			Comment: string(wlip.Comment.ValueString()),
		})
	}

	// Rewrite Rules
	s.RewriteRules = []ogosecurity.RewriteRule{}
	for _, rewrite := range plan.RewriteRules {
		s.RewriteRules = append(s.RewriteRules, ogosecurity.RewriteRule{
			Active:             bool(rewrite.Active.ValueBool()),
			Comment:            string(rewrite.Comment.ValueString()),
			RewriteSource:      string(rewrite.RewriteSource.ValueString()),
			RewriteDestination: string(rewrite.RewriteDestination.ValueString()),
		})
	}

	// Rules access
	s.Rules = []ogosecurity.Rule{}
	for _, rule := range plan.Rules {
		r := ogosecurity.Rule{
			Active:         bool(rule.Active.ValueBool()),
			Action:         string(rule.Action.ValueString()),
			Cache:          bool(rule.Cache.ValueBool()),
			Comment:        string(rule.Comment.ValueString()),
			Paths:          []string{},
			WhitelistedIps: []string{},
		}

		for _, path := range rule.Paths {
			r.Paths = append(r.Paths, string(path.ValueString()))
		}

		for _, ip := range rule.WhitelistedIps {
			r.WhitelistedIps = append(r.WhitelistedIps, string(ip.ValueString()))
		}

		s.Rules = append(s.Rules, r)
	}

	// URL Exceptions
	s.UrlExceptions = []ogosecurity.UrlException{}
	for _, url := range plan.UrlExceptions {
		s.UrlExceptions = append(s.UrlExceptions, ogosecurity.UrlException{
			Path:    string(url.Path.ValueString()),
			Comment: string(url.Comment.ValueString()),
		})
	}

	// Tags
	s.Tags = []string{}
	for _, tag := range plan.Tags {
		s.Tags = append(s.Tags, string(tag.ValueString()))
	}

	site, err := r.client.UpdateSite(s)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating site",
			"Could not update site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ClusterEntrypoint4 = types.StringValue(site.Cluster.Entrypoint4)
	plan.ClusterEntrypoint6 = types.StringValue(site.Cluster.Entrypoint6)
	if plan.ActiveCustomerCertificate != nil {
		plan.ActiveCustomerCertificate.Cn = types.StringValue(site.ActiveCustomerCertificate.Cn)
		plan.ActiveCustomerCertificate.ExpiredAt = types.StringValue(site.ActiveCustomerCertificate.ExpiredAt)
		plan.ActiveCustomerCertificate.Hash = types.StringValue(site.ActiveCustomerCertificate.Hash)
	}
	if plan.Cdn.String() != "" {
		plan.ClusterEntrypointCdn = types.StringValue(site.Cluster.EntrypointCdn)
		plan.CdnStatus = types.StringPointerValue(site.CdnStatus)
	}
	plan.Status = types.StringValue(site.Status)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state SiteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing site
	err := r.client.DeleteSite(state.DomainName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Ogo Site",
			"Could not delete site, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import domain name and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("domain_name"), req, resp)
}

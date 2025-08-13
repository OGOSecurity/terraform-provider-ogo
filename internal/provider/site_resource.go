package provider

import (
	"context"
	"fmt"
	"time"

	ogosecurity "terraform-provider-ogo/internal/ogo"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	DomainName           types.String        `tfsdk:"domain_name"`
	ClusterUid           types.String        `tfsdk:"cluster_uid"`
	OriginServer         types.String        `tfsdk:"origin_server"`
	OriginScheme         types.String        `tfsdk:"origin_scheme"`
	OriginPort           types.Int32         `tfsdk:"origin_port"`
	OriginSkipCertVerify types.Bool          `tfsdk:"origin_skip_cert_verify"`
	OriginMtlsEnabled    types.Bool          `tfsdk:"origin_mtls_enabled"`
	RemoveXForwarded     types.Bool          `tfsdk:"remove_xforwarded"`
	LogExportEnabled     types.Bool          `tfsdk:"log_export_enabled"`
	ForceHttps           types.Bool          `tfsdk:"force_https"`
	AuditMode            types.Bool          `tfsdk:"audit_mode"`
	PassthroughMode      types.Bool          `tfsdk:"passthrough_mode"`
	Hsts                 types.String        `tfsdk:"hsts"`
	PassTlsClientCert    types.String        `tfsdk:"pass_tls_client_cert"`
	TlsOptionsUid        types.String        `tfsdk:"tls_options_uid"`
	Tags                 []types.String      `tfsdk:"tags"`
	BlacklistedCountries []types.String      `tfsdk:"blacklisted_countries"`
	IpExceptions         []IpExceptionModel  `tfsdk:"ip_exceptions"`
	UrlExceptions        []UrlExceptionModel `tfsdk:"url_exceptions"`
	RewriteRules         []RewriteRuleModel  `tfsdk:"rewrite_rules"`
	Rules                []RuleModel         `tfsdk:"rules"`
	LastUpdated          types.String        `tfsdk:"last_updated"`
}

type RewriteRuleModel struct {
	Priority           types.Int32  `tfsdk:"priority"`
	Active             types.Bool   `tfsdk:"active"`
	Comment            types.String `tfsdk:"comment"`
	RewriteSource      types.String `tfsdk:"rewrite_source"`
	RewriteDestination types.String `tfsdk:"rewrite_destination"`
}

type RuleModel struct {
	Priority       types.Int32    `tfsdk:"priority"`
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
				Description: "DNS domain name of site",
			},
			"cluster_uid": schema.StringAttribute{
				Required:    true,
				Description: "Cluster UID on which site is deployed (force site recreation if modified)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"origin_server": schema.StringAttribute{
				Required:    true,
				Description: "Origin server address (IP address or domain name)",
			},
			"origin_scheme": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Scheme to used to access origin server. Supported values: **https** or **http** (Default: **https**)",
				Default:     stringdefault.StaticString("https"),
				Validators: []validator.String{
					stringvalidator.OneOf("http", "https"),
				},
			},
			"origin_port": schema.Int32Attribute{
				Optional:    true,
				Description: "Port to be used to access origin server. Must be defined only if different to standard HTTP port 443 or 80, otherwise let Ogo choose the correct port",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
			},
			"origin_mtls_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable mTLS between Ogo and origin server (Default: **false**)",
				Default:     booldefault.StaticBool(false),
			},
			"origin_skip_cert_verify": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Skip origin server certificate verification if TLS is used. If set to **true** Ogo accept connection to origin server even if certificate doesn't match site domain name, or certificate is expired, or certificate is self signed (Default: **false**)",
				Default:     booldefault.StaticBool(false),
			},
			"remove_xforwarded": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Remove X-Forwarded-* headers. (Default: **false** )",
				Default:     booldefault.StaticBool(false),
			},
			"force_https": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Redirect HTTP request to HTTPS (Default: **true** )",
				Default:     booldefault.StaticBool(true),
			},
			"audit_mode": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable audit mode. Requests are analysed by Ogo Shield but never blocked (Default: **false**)",
				Default:     booldefault.StaticBool(false),
			},
			"passthrough_mode": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable passthrough mode. Requests are not analysed by Ogo Shield and never blocked (Default: **false**)",
				Default:     booldefault.StaticBool(false),
			},
			"hsts": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable HSTS (Default: **hsts**). Supported values:\n  * **hsts**: Enable HSTS\n  * **hstss**: Enable HSTS including subdomains\n  * **hstsp**: Enable HSTS including subdomains and preloading\n  * **none**: Disable HSTS",
				Default:     stringdefault.StaticString("hsts"),
				Validators: []validator.String{
					stringvalidator.OneOf("hsts", "hstss", "hstssp", "none"),
				},
			},
			"log_export_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable log export for this site (Default: **false**)",
				Default:     booldefault.StaticBool(false),
			},
			"tls_options_uid": schema.StringAttribute{
				Optional:    true,
				Description: "UID of TLS options to be applied to this site",
			},
			"pass_tls_client_cert": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Client certificate information to pass to origin server (Default: **info**). Supported values:\n  * **all**: Send certificate and certificate information\n  * **cert**: Send only certificate\n  * **info**: Send only certificate information\n  * **none***: Nothing send",
				Default:     stringdefault.StaticString("info"),
				Validators: []validator.String{
					stringvalidator.OneOf("all", "cert", "info", "none"),
				},
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List of tags",
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"blacklisted_countries": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List of countries to blacklist",
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"ip_exceptions": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Required:    true,
							Description: "IP address list never blocked by Ogo Shield",
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated to this IP list",
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
			"rewrite_rules": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"priority": schema.Int32Attribute{
							Required:    true,
							Description: "Rewrite rule priority",
						},
						"active": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Flag to enabled (**true**) or disabled (**false**) rewrite rule. (Default: **true**)",
							Default:     booldefault.StaticBool(true),
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated to this rewrite rule",
						},
						"rewrite_source": schema.StringAttribute{
							Required:    true,
							Description: "Source path to be rewrited",
						},
						"rewrite_destination": schema.StringAttribute{
							Required:    true,
							Description: "Rewrited destination path",
						},
					},
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"priority":            types.Int32Type,
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
			"rules": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"priority": schema.Int32Attribute{
							Required:    true,
							Description: "Rule priority",
						},
						"active": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Flag to enabled (**true**) or disabled (**false**) rule. (Default: **true**)",
							Default:     booldefault.StaticBool(true),
						},
						"action": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Action to applied when rule match (Default: **brain**). Supported values:\n  * **brain**: Rule analysed by ogo Shield brain\n  * **bypass**: Rule not analysed by Ogo Shield brain",
							Default:     stringdefault.StaticString("brain"),
							Validators: []validator.String{
								stringvalidator.OneOf("brain", "bypass"),
							},
						},
						"cache": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "",
							Default:     booldefault.StaticBool(false),
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated to this rule",
						},
						"paths": schema.SetAttribute{
							Required:    true,
							Description: "List of URL path for which rule is applied",
							ElementType: types.StringType,
						},
						"whitelisted_ips": schema.SetAttribute{
							Required:    true,
							Description: "Authorized IP addresses list",
							ElementType: types.StringType,
						},
					},
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"priority":        types.Int32Type,
								"active":          types.BoolType,
								"action":          types.StringType,
								"cache":           types.BoolType,
								"comment":         types.StringType,
								"paths":           types.SetType{types.StringType},
								"whitelisted_ips": types.SetType{types.StringType},
							},
						},
						[]attr.Value{},
					),
				),
			},
			"url_exceptions": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Required:    true,
							Description: "Path of URL never blocked by Ogo Shield",
						},
						"comment": schema.StringAttribute{
							Optional:    true,
							Description: "Description associated to this URL exception",
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
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Last resource update by terraform",
			},
		},
		MarkdownDescription: "Resource *site* can be used to create, update or delete site " +
			"configuration in Ogo Dashboard.\n\n" +
			"This resource allowed to managed all site settings.\n\n",
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
			"Unexpected Data Source Configure Type",
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
		OriginPort:           plan.OriginPort.ValueInt32Pointer(),
		PassTlsClientCert:    string(plan.PassTlsClientCert.ValueString()),
		TlsOptions:           tlsOpt,
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
			Priority:           int(rewrite.Priority.ValueInt32()),
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
			Priority:       int(rule.Priority.ValueInt32()),
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

	_, err := r.client.CreateSite(s)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			"Could not create site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
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
	state.PassTlsClientCert = types.StringValue(site.PassTlsClientCert)

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
			Priority:           types.Int32Value(int32(rewrite.Priority)),
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
			Priority:       types.Int32Value(int32(rule.Priority)),
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
		OriginPort:           plan.OriginPort.ValueInt32Pointer(),
		PassTlsClientCert:    string(plan.PassTlsClientCert.ValueString()),
		TlsOptions:           tlsOpt,
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
			Priority:           int(rewrite.Priority.ValueInt32()),
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
			Priority:       int(rule.Priority.ValueInt32()),
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

	_, err := r.client.UpdateSite(s)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating site",
			"Could not update site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
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

package provider

import (
	"context"
	"fmt"
	"time"

	ogosecurity "terraform-provider-ogo/internal/ogo"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

// siteResourceModel maps the resource schema data.
type siteResourceModel struct {
	Name                 types.String         `tfsdk:"name"`
	ClusterUid           types.String         `tfsdk:"cluster_uid"`
	DestHost             types.String         `tfsdk:"dest_host"`
	DestHostScheme       types.String         `tfsdk:"dest_host_scheme"`
	DestHostMtls         types.Bool           `tfsdk:"dest_host_mtls"`
	Port                 types.Int32          `tfsdk:"port"`
	LogExport            types.Bool           `tfsdk:"log_export"`
	TrustSelfSigned      types.Bool           `tfsdk:"trust_selfsigned"`
	NoCopyXForwarded     types.Bool           `tfsdk:"no_copy_xforwarded"`
	ForceHttps           types.Bool           `tfsdk:"force_https"`
	DryRun               types.Bool           `tfsdk:"dry_run"`
	PanicMode            types.Bool           `tfsdk:"panic_mode"`
	Hsts                 types.String         `tfsdk:"hsts"`
	PassTlsClientCert    types.String         `tfsdk:"pass_tls_client_cert"`
	TlsOptionsUid        types.String         `tfsdk:"tls_options_uid"`
	Tags                 []types.String       `tfsdk:"tags"`
	BlacklistedCountries []types.String       `tfsdk:"blacklisted_countries"`
	WhitelistedIps       []whitelistedIpModel `tfsdk:"whitelisted_ips"`
	RewriteRules         []rewriteRuleModel   `tfsdk:"rewrite_rules"`
	Rules                []ruleModel          `tfsdk:"rules"`
	UrlExceptions        []urlExceptionModel  `tfsdk:"url_exceptions"`
	LastUpdated          types.String         `tfsdk:"last_updated"`
}

type rewriteRuleModel struct {
	Id                 types.Int32  `tfsdk:"id"`
	Priority           types.Int32  `tfsdk:"priority"`
	Active             types.Bool   `tfsdk:"active"`
	Comment            types.String `tfsdk:"comment"`
	RewriteSource      types.String `tfsdk:"rewrite_source"`
	RewriteDestination types.String `tfsdk:"rewrite_destination"`
}

type ruleModel struct {
	Id             types.Int32    `tfsdk:"id"`
	Priority       types.Int32    `tfsdk:"priority"`
	Active         types.Bool     `tfsdk:"active"`
	Action         types.String   `tfsdk:"action"`
	Cache          types.Bool     `tfsdk:"cache"`
	Comment        types.String   `tfsdk:"comment"`
	Paths          []types.String `tfsdk:"paths"`
	WhitelistedIps []types.String `tfsdk:"whitelisted_ips"`
}

type urlExceptionModel struct {
	Path    types.String `tfsdk:"path"`
	Comment types.String `tfsdk:"comment"`
}

type whitelistedIpModel struct {
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
			"name": schema.StringAttribute{
				Required: true,
			},
			"cluster_uid": schema.StringAttribute{
				Required: true,
			},
			"dest_host": schema.StringAttribute{
				Required: true,
			},
			"dest_host_scheme": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("https"),
			},
			"port": schema.Int32Attribute{
				Optional: true,
			},
			"trust_selfsigned": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"no_copy_xforwarded": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"force_https": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"dry_run": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"panic_mode": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"hsts": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("hsts"),
			},
			"log_export": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"dest_host_mtls": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"tls_options_uid": schema.StringAttribute{
				Optional: true,
			},
			"pass_tls_client_cert": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("info"),
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
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
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"whitelisted_ips": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Required: true,
						},
						"comment": schema.StringAttribute{
							Optional: true,
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
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							Optional: true,
							Computed: true,
						},
						"priority": schema.Int32Attribute{
							Required: true,
						},
						"active": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(true),
						},
						"comment": schema.StringAttribute{
							Optional: true,
						},
						"rewrite_source": schema.StringAttribute{
							Required: true,
						},
						"rewrite_destination": schema.StringAttribute{
							Required: true,
						},
					},
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"id":                  types.Int32Type,
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
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int32Attribute{
							Optional: true,
							Computed: true,
						},
						"priority": schema.Int32Attribute{
							Required: true,
						},
						"active": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(true),
						},
						"action": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("brain"),
						},
						"cache": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
						"comment": schema.StringAttribute{
							Optional: true,
						},
						"paths": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
						"whitelisted_ips": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"id":              types.Int32Type,
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
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Required: true,
						},
						"comment": schema.StringAttribute{
							Optional: true,
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
				Computed: true,
			},
		},
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
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new site
	var tlsOpt *ogosecurity.TlsOptions
	s := ogosecurity.Site{
		Name: string(plan.Name.ValueString()),
		Cluster: ogosecurity.Cluster{
			Uid: string(plan.ClusterUid.ValueString()),
		},
		DestHost:          string(plan.DestHost.ValueString()),
		DestHostScheme:    string(plan.DestHostScheme.ValueString()),
		DestHostMtls:      bool(plan.DestHostMtls.ValueBool()),
		TrustSelfSigned:   bool(plan.TrustSelfSigned.ValueBool()),
		NoCopyXForwarded:  bool(plan.NoCopyXForwarded.ValueBool()),
		ForceHttps:        bool(plan.ForceHttps.ValueBool()),
		DryRun:            bool(plan.DryRun.ValueBool()),
		PanicMode:         bool(plan.PanicMode.ValueBool()),
		Hsts:              string(plan.Hsts.ValueString()),
		LogExport:         bool(plan.LogExport.ValueBool()),
		Port:              plan.Port.ValueInt32Pointer(),
		PassTlsClientCert: string(plan.PassTlsClientCert.ValueString()),
		TlsOptions:        tlsOpt,
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

	// Whitelist IP
	s.WhitelistedIps = []ogosecurity.WhitelistedIp{}
	for _, wlip := range plan.WhitelistedIps {
		s.WhitelistedIps = append(s.WhitelistedIps, ogosecurity.WhitelistedIp{
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
	var state siteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed site value from Ogo
	site, err := r.client.GetSite(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ogo Site",
			"Could not read Ogo site Name "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite properties with refreshed state
	state.ClusterUid = types.StringValue(site.Cluster.Uid)
	state.DestHost = types.StringValue(site.DestHost)
	state.DestHostScheme = types.StringValue(site.DestHostScheme)
	state.DestHostMtls = types.BoolValue(site.DestHostMtls)
	state.Port = types.Int32PointerValue(site.Port)
	state.TrustSelfSigned = types.BoolValue(site.TrustSelfSigned)
	state.NoCopyXForwarded = types.BoolValue(site.NoCopyXForwarded)
	state.ForceHttps = types.BoolValue(site.ForceHttps)
	state.DryRun = types.BoolValue(site.DryRun)
	state.PanicMode = types.BoolValue(site.PanicMode)
	state.Hsts = types.StringValue(site.Hsts)
	state.LogExport = types.BoolValue(site.LogExport)
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

	// Whitelist IP
	state.WhitelistedIps = []whitelistedIpModel{}
	for _, wlip := range site.WhitelistedIps {
		state.WhitelistedIps = append(state.WhitelistedIps, whitelistedIpModel{
			Ip:      types.StringValue(wlip.Ip),
			Comment: types.StringValue(wlip.Comment),
		})
	}

	// Rewrite rules
	state.RewriteRules = []rewriteRuleModel{}
	for _, rewrite := range site.RewriteRules {
		state.RewriteRules = append(state.RewriteRules, rewriteRuleModel{
			Id:                 types.Int32Value(int32(rewrite.Id)),
			Priority:           types.Int32Value(int32(rewrite.Priority)),
			Active:             types.BoolValue(rewrite.Active),
			Comment:            types.StringValue(rewrite.Comment),
			RewriteSource:      types.StringValue(rewrite.RewriteSource),
			RewriteDestination: types.StringValue(rewrite.RewriteDestination),
		})
	}

	// Rules access
	state.Rules = []ruleModel{}
	for _, rule := range site.Rules {
		r := ruleModel{
			Id:             types.Int32Value(int32(rule.Id)),
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
	state.UrlExceptions = []urlExceptionModel{}
	for _, url := range site.UrlExceptions {
		state.UrlExceptions = append(state.UrlExceptions, urlExceptionModel{
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
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new site
	var tlsOpt *ogosecurity.TlsOptions
	s := ogosecurity.Site{
		Name: string(plan.Name.ValueString()),
		Cluster: ogosecurity.Cluster{
			Uid: string(plan.ClusterUid.ValueString()),
		},
		DestHost:          string(plan.DestHost.ValueString()),
		DestHostScheme:    string(plan.DestHostScheme.ValueString()),
		DestHostMtls:      bool(plan.DestHostMtls.ValueBool()),
		TrustSelfSigned:   bool(plan.TrustSelfSigned.ValueBool()),
		NoCopyXForwarded:  bool(plan.NoCopyXForwarded.ValueBool()),
		ForceHttps:        bool(plan.ForceHttps.ValueBool()),
		DryRun:            bool(plan.DryRun.ValueBool()),
		PanicMode:         bool(plan.PanicMode.ValueBool()),
		Hsts:              string(plan.Hsts.ValueString()),
		LogExport:         bool(plan.LogExport.ValueBool()),
		Port:              plan.Port.ValueInt32Pointer(),
		PassTlsClientCert: string(plan.PassTlsClientCert.ValueString()),
		TlsOptions:        tlsOpt,
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

	// Whitelist IP
	s.WhitelistedIps = []ogosecurity.WhitelistedIp{}
	for _, wlip := range plan.WhitelistedIps {
		s.WhitelistedIps = append(s.WhitelistedIps, ogosecurity.WhitelistedIp{
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state siteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing site
	err := r.client.DeleteSite(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Ogo Site",
			"Could not delete site, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Name and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

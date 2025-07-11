package provider

import (
	"context"
	"fmt"
	"time"

	ogosecurity "terraform-provider-ogo/internal/ogo"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	Name             types.String `tfsdk:"name"`
	ClusterID        types.Int64  `tfsdk:"cluster_id"`
	ClusterName      types.String `tfsdk:"cluster_name"`
	DestHost         types.String `tfsdk:"dest_host"`
	DestHostScheme   types.String `tfsdk:"dest_host_scheme"`
	TrustSelfSigned  types.Bool   `tfsdk:"trust_selfsigned"`
	NoCopyXForwarded types.Bool   `tfsdk:"no_copy_xforwarded"`
	ForceHttps       types.Bool   `tfsdk:"force_https"`
	DryRun           types.Bool   `tfsdk:"dry_run"`
	PanicMode        types.Bool   `tfsdk:"panic_mode"`
	LastUpdated      types.String `tfsdk:"last_updated"`
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
			"cluster_id": schema.Int64Attribute{
				Required: true,
			},
			"cluster_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dest_host": schema.StringAttribute{
				Required: true,
			},
			"dest_host_scheme": schema.StringAttribute{
				Required: true,
			},
			"trust_selfsigned": schema.BoolAttribute{
				Required: true,
			},
			"no_copy_xforwarded": schema.BoolAttribute{
				Required: true,
			},
			"force_https": schema.BoolAttribute{
				Required: true,
			},
			"dry_run": schema.BoolAttribute{
				Required: true,
			},
			"panic_mode": schema.BoolAttribute{
				Required: true,
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
	s := ogosecurity.Site{
		Name:             string(plan.Name.ValueString()),
		ClusterID:        int(plan.ClusterID.ValueInt64()),
		DestHost:         string(plan.DestHost.ValueString()),
		DestHostScheme:   string(plan.DestHostScheme.ValueString()),
		TrustSelfSigned:  bool(plan.TrustSelfSigned.ValueBool()),
		NoCopyXForwarded: bool(plan.NoCopyXForwarded.ValueBool()),
		ForceHttps:       bool(plan.ForceHttps.ValueBool()),
		DryRun:           bool(plan.DryRun.ValueBool()),
		PanicMode:        bool(plan.PanicMode.ValueBool()),
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
	plan.ClusterName = types.StringValue(string(site.ClusterName))
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
	state.ClusterID = types.Int64Value(int64(site.ClusterID))
	state.ClusterName = types.StringValue(site.ClusterName)
	state.DestHost = types.StringValue(site.DestHost)
	state.DestHostScheme = types.StringValue(site.DestHostScheme)
	state.TrustSelfSigned = types.BoolValue(site.TrustSelfSigned)
	state.NoCopyXForwarded = types.BoolValue(site.NoCopyXForwarded)
	state.ForceHttps = types.BoolValue(site.ForceHttps)
	state.DryRun = types.BoolValue(site.DryRun)
	state.PanicMode = types.BoolValue(site.PanicMode)

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
	s := ogosecurity.Site{
		Name:             string(plan.Name.ValueString()),
		ClusterID:        int(plan.ClusterID.ValueInt64()),
		DestHost:         string(plan.DestHost.ValueString()),
		DestHostScheme:   string(plan.DestHostScheme.ValueString()),
		TrustSelfSigned:  bool(plan.TrustSelfSigned.ValueBool()),
		NoCopyXForwarded: bool(plan.NoCopyXForwarded.ValueBool()),
		ForceHttps:       bool(plan.ForceHttps.ValueBool()),
		DryRun:           bool(plan.DryRun.ValueBool()),
		PanicMode:        bool(plan.PanicMode.ValueBool()),
	}

	site, err := r.client.UpdateSite(s)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			"Could not create site, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ClusterName = types.StringValue(string(site.ClusterName))
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

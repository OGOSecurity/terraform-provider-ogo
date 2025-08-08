package provider

import (
	"context"
	"fmt"
	"time"

	ogosecurity "terraform-provider-ogo/internal/ogo"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &tlsOptionsResource{}
	_ resource.ResourceWithConfigure   = &tlsOptionsResource{}
	_ resource.ResourceWithImportState = &tlsOptionsResource{}
)

// TlsOptionsResourceModel maps the resource schema data.
type TlsOptionsResourceModel struct {
	Name              types.String   `tfsdk:"name"`
	Uid               types.String   `tfsdk:"uid"`
	ClientAuthType    types.String   `tfsdk:"client_auth_type"`
	ClientAuthCaCerts []types.String `tfsdk:"client_auth_ca_certs"`
	MinTlsVersion     types.String   `tfsdk:"min_tls_version"`
	MaxTlsVersion     types.String   `tfsdk:"max_tls_version"`
	LastUpdated       types.String   `tfsdk:"last_updated"`
}

// NewTlsOptionsResource is a helper function to simplify the provider implementation.
func NewTlsOptionsResource() resource.Resource {
	return &tlsOptionsResource{}
}

// tlsOptionsResource is the resource implementation.
type tlsOptionsResource struct {
	client *ogosecurity.Client
}

// Metadata returns the resource type name.
func (r *tlsOptionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shield_tls_options"
}

// Schema defines the schema for the resource.
func (r *tlsOptionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"uid": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_auth_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("VerifyClientCertIfGiven"),
				Validators: []validator.String{
					stringvalidator.OneOf("VerifyClientCertIfGiven", "RequireAndVerifyClientCert"),
				},
			},
			"client_auth_ca_certs": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"min_tls_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("TLS_12"),
				Validators: []validator.String{
					stringvalidator.OneOf("TLS_10", "TLS_11", "TLS_12", "TLS_13"),
				},
			},
			"max_tls_version": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("TLS_10", "TLS_11", "TLS_12", "TLS_13"),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *tlsOptionsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *tlsOptionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan TlsOptionsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new TLS options
	t := ogosecurity.TlsOptions{
		Name:           string(plan.Name.ValueString()),
		ClientAuthType: string(plan.ClientAuthType.ValueString()),
		MinTlsVersion:  string(plan.MinTlsVersion.ValueString()),
		MaxTlsVersion:  string(plan.MaxTlsVersion.ValueString()),
	}

	// CA certificates
	t.ClientAuthCaCerts = []string{}
	for _, cert := range plan.ClientAuthCaCerts {
		t.ClientAuthCaCerts = append(t.ClientAuthCaCerts, string(cert.ValueString()))
	}

	// Create TLS options
	tlsOpt, err := r.client.CreateTlsOptions(t)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating TLS options",
			"Could not create TLS options, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Uid = types.StringValue(tlsOpt.Uid)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *tlsOptionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state TlsOptionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed TLS options value from Ogo
	tlsOptions, err := r.client.GetTlsOptions(state.Uid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ogo TLS options",
			"Could not read Ogo TLS options "+state.Uid.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite properties with refreshed state
	state.Name = types.StringValue(tlsOptions.Name)
	state.ClientAuthType = types.StringValue(tlsOptions.ClientAuthType)
	state.MinTlsVersion = types.StringValue(tlsOptions.MinTlsVersion)
	state.MaxTlsVersion = types.StringValue(tlsOptions.MaxTlsVersion)

	// CA certificates
	state.ClientAuthCaCerts = []types.String{}
	for _, cert := range tlsOptions.ClientAuthCaCerts {
		state.ClientAuthCaCerts = append(state.ClientAuthCaCerts, types.StringValue(cert))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tlsOptionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan TlsOptionsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set TLS options
	t := ogosecurity.TlsOptions{
		Name:           string(plan.Name.ValueString()),
		Uid:            string(plan.Uid.ValueString()),
		ClientAuthType: string(plan.ClientAuthType.ValueString()),
		MinTlsVersion:  string(plan.MinTlsVersion.ValueString()),
		MaxTlsVersion:  string(plan.MaxTlsVersion.ValueString()),
	}

	// CA certificate
	t.ClientAuthCaCerts = []string{}
	for _, cert := range plan.ClientAuthCaCerts {
		t.ClientAuthCaCerts = append(t.ClientAuthCaCerts, string(cert.ValueString()))
	}

	// Update TLS options
	_, err := r.client.UpdateTlsOptions(t)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating TLS options",
			"Could not update TLS options, unexpected error: "+err.Error(),
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
func (r *tlsOptionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state TlsOptionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing TLS options
	err := r.client.DeleteTlsOptions(state.Uid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TLS options",
			"Could not delete TLS options, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *tlsOptionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import domain name and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uid"), req, resp)
}

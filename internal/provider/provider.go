package provider

import (
	"context"
	"os"

	ogosecurity "terraform-provider-ogo/internal/ogo"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ogoProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &ogoProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ogoProvider{
			version: version,
		}
	}
}

// ogoProvider defines the provider implementation.
type ogoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ogoProviderModel describes the provider data model.
type ogoProviderModel struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	Organization types.String `tfsdk:"organization"`
	ApiKey       types.String `tfsdk:"apikey"`
}

func (p *ogoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ogo"
	resp.Version = p.version
}

func (p *ogoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Ogo API endpoint",
				Required:            true,
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "Organization used to authenticate to Ogo Dashboard (usualy user email address)",
				Required:            true,
			},
			"apikey": schema.StringAttribute{
				MarkdownDescription: "API Key",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ogoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Ogo client")
	var config ogoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	endpoint := os.Getenv("OGO_ENDPOINT")
	organization := os.Getenv("OGO_ORGANIZATION")
	apikey := os.Getenv("OGO_APIKEY")

	// Configuration values are now available.
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.Organization.IsNull() {
		organization = config.Organization.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return error
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Ogo API endpoint",
			"The provider cannot create the Ogo API client as there is a missing or empty value for the Ogo API endpoint. "+
				"Set the endpoint value in the configuration or use the OGO_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if organization == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Missing Ogo API organization",
			"The provider cannot create the Ogo API client as there is a missing or empty value for the Ogo API organization. "+
				"Set the organization value in the configuration or use the OGO_ORGANIZATION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing Ogo API apikey",
			"The provider cannot create the Ogo API client as there is a missing or empty value for the Ogo API apikey. "+
				"Set the apikey value in the configuration or use the OGO_APIKEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ogo_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "ogo_organization", organization)
	ctx = tflog.SetField(ctx, "ogo_apikey", apikey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ogo_apikey")

	tflog.Debug(ctx, "Creating Ogo client")

	// Create a new HashiCups client using the configuration values
	client, err := ogosecurity.NewClient(&endpoint, &organization, &apikey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OgoSecurity Dashboard API Client",
			"An unexpected error occurred when creating the OgoSecurity Dashboard API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OgoSecurity Dashboard Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Ogo client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *ogoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClustersDataSource,
		NewTlsOptionsDataSource,
	}
}

func (p *ogoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSiteResource,
		NewTlsOptionsResource,
	}
}

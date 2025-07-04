package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ogoProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &ogoProvider{}
)

// ogoProvider defines the provider implementation.
type ogoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ogoProviderModel describes the provider data model.
type ogoProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ogoProvider{
			version: version,
		}
	}
}

func (p *ogoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ogo"
	resp.Version = p.version
}

func (p *ogoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *ogoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ogoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ogoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	//return []func() datasource.DataSource{
	//	NewExampleDataSource,
	//}
	return nil
}

func (p *ogoProvider) Resources(ctx context.Context) []func() resource.Resource {
	//return []func() resource.Resource{
	//	NewExampleResource,
	//}
	return nil
}

package provider

import (
	"context"
	"fmt"

	ogosecurity "terraform-provider-ogo/internal/ogo"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clustersDataSource{}
	_ datasource.DataSourceWithConfigure = &clustersDataSource{}
)

// clustersDataSourceModel maps the data source schema data
type clustersDataSourceModel struct {
	Clusters []clustersModel `tfsdk:"clusters"`
}

// clusterModel maps cluster schema data
type clustersModel struct {
	Uid                 types.String `tfsdk:"uid"`
	Name                types.String `tfsdk:"name"`
	Host4               types.String `tfsdk:"host4"`
	Host6               types.String `tfsdk:"host6"`
	SupportsCache       types.Bool   `tfsdk:"supports_cache"`
	SupportsIpv6Origins types.Bool   `tfsdk:"supports_ipv6_origins"`
	SupportsMtls        types.Bool   `tfsdk:"supports_mtls"`
}

func NewClustersDataSource() datasource.DataSource {
	return &clustersDataSource{}
}

type clustersDataSource struct {
	client *ogosecurity.Client
}

func (d *clustersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shield_clusters"
}

func (d *clustersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"clusters": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uid": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"host4": schema.StringAttribute{
							Computed: true,
						},
						"host6": schema.StringAttribute{
							Computed: true,
						},
						"supports_cache": schema.BoolAttribute{
							Computed: true,
						},
						"supports_ipv6_origins": schema.BoolAttribute{
							Computed: true,
						},
						"supports_mtls": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *clustersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *clustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clustersDataSourceModel

	clusters, err := d.client.GetAllClusters()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Ogo Clusters",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, c := range clusters {
		clusterState := clustersModel{
			Uid:                 types.StringValue(c.Uid),
			Name:                types.StringValue(c.Name),
			Host4:               types.StringValue(c.Host4),
			Host6:               types.StringValue(c.Host6),
			SupportsCache:       types.BoolValue(c.SupportsCache),
			SupportsIpv6Origins: types.BoolValue(c.SupportsIpv6Origins),
			SupportsMtls:        types.BoolValue(c.SupportsMtls),
		}

		state.Clusters = append(state.Clusters, clusterState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

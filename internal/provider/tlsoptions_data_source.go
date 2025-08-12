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
	_ datasource.DataSource              = &tlsoptionsDataSource{}
	_ datasource.DataSourceWithConfigure = &tlsoptionsDataSource{}
)

// tlsoptionsDataSourceModel maps the data source schema data
type tlsoptionsDataSourceModel struct {
	TlsOptions []tlsoptionsModel `tfsdk:"tlsoptions"`
}

// tlsoptionsModel maps TLS Options schema data
type tlsoptionsModel struct {
	Uid               types.String   `tfsdk:"uid"`
	Name              types.String   `tfsdk:"name"`
	ClientAuthType    types.String   `tfsdk:"client_auth_type"`
	ClientAuthCaCerts []types.String `tfsdk:"client_auth_ca_certs"`
	MinTlsVersion     types.String   `tfsdk:"min_tls_version"` //TLS_1.2 => enum
	MaxTlsVersion     types.String   `tfsdk:"max_tls_version"` //TLS_1.3 => enum
}

func NewTlsOptionsDataSource() datasource.DataSource {
	return &tlsoptionsDataSource{}
}

type tlsoptionsDataSource struct {
	client *ogosecurity.Client
}

func (d *tlsoptionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shield_tlsoptions"
}

func (d *tlsoptionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tlsoptions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uid": schema.StringAttribute{
							Computed:    true,
							Description: "UID used to reference this TLS Options",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the TLS Options",
						},
						"client_auth_type": schema.StringAttribute{
							Computed:    true,
							Description: "Authentication type needed to authenticate client.\n  * *VerifyClientCertIfGiven*: if a certificate is provided, verifies if it is signed by a CA listed in `client_auth_ca_certs`. Otherwise proceeds without any certificate.\n  * *RequireAndVerifyClientCert*: requires a certificate, which must be signed by a CA listed in `client_auth_ca_certs`.",
						},
						"client_auth_ca_certs": schema.ListAttribute{
							Computed:    true,
							Description: "List of certificate authority used to verify client certificate",
							ElementType: types.StringType,
						},
						"min_tls_version": schema.StringAttribute{
							Computed:    true,
							Description: "Minimum TLS version accepted",
						},
						"max_tls_version": schema.StringAttribute{
							Computed:    true,
							Description: "Maximum TLS version accepted",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *tlsoptionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tlsoptionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tlsoptionsDataSourceModel

	tlsoptions, err := d.client.GetAllTlsOptions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Ogo TLS Options",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, t := range tlsoptions {
		tlsoptionsState := tlsoptionsModel{
			Uid:            types.StringValue(t.Uid),
			Name:           types.StringValue(t.Name),
			ClientAuthType: types.StringValue(t.ClientAuthType),
			MinTlsVersion:  types.StringValue(t.MinTlsVersion),
			MaxTlsVersion:  types.StringValue(t.MaxTlsVersion),
		}

		for _, cert := range t.ClientAuthCaCerts {
			tlsoptionsState.ClientAuthCaCerts = append(tlsoptionsState.ClientAuthCaCerts, types.StringValue(cert))
		}

		state.TlsOptions = append(state.TlsOptions, tlsoptionsState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

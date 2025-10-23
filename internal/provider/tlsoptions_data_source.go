// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	MinTlsVersion     types.String   `tfsdk:"min_tls_version"`
	MaxTlsVersion     types.String   `tfsdk:"max_tls_version"`
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
							Description: "UID used to reference this TLS Options.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the TLS Options.",
						},
						"client_auth_type": schema.StringAttribute{
							Computed: true,
							Description: "Authentication type needed to authenticate clients.\n" +
								"  * **VerifyClientCertIfGiven**: If a certificate is provided, verify if it is signed by a CA listed in `client_auth_ca_certs`. Otherwise, proceed without any certificate.\n" +
								"  * **RequireAndVerifyClientCert**: Require a certificate, which must be signed by a CA listed in `client_auth_ca_certs`.",
						},
						"client_auth_ca_certs": schema.ListAttribute{
							Computed:    true,
							Description: "List of certificate authorities used to verify client certificates.",
							ElementType: types.StringType,
						},
						"min_tls_version": schema.StringAttribute{
							Computed:    true,
							Description: "Minimum TLS version accepted.",
						},
						"max_tls_version": schema.StringAttribute{
							Computed:    true,
							Description: "Maximum TLS version accepted.",
						},
					},
				},
			},
		},
		MarkdownDescription: "Get a list of organization TLS options and associated configurations.\n\n" +
			"Use this data source to retrieve information, in particular TLS options UID, to be used " +
			"in `ogo_shield_site` resource configuration for which TLS default settings need to be overridden.",
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
			"unexpected data source configure type",
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
			"Unable to read Ogo TLS Options",
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
			MinTlsVersion:  types.StringPointerValue(t.MinTlsVersion),
			MaxTlsVersion:  types.StringPointerValue(t.MaxTlsVersion),
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

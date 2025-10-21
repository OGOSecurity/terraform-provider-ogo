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
	_ datasource.DataSource              = &contractsDataSource{}
	_ datasource.DataSourceWithConfigure = &contractsDataSource{}
)

// contractsDataSourceModel maps the data source schema data
type contractsDataSourceModel struct {
	Contracts []contractsModel `tfsdk:"contracts"`
}

// contractModel maps contract schema data
type contractsModel struct {
	Number types.String `tfsdk:"number"`
	Name   types.String `tfsdk:"name"`
}

func NewContractsDataSource() datasource.DataSource {
	return &contractsDataSource{}
}

type contractsDataSource struct {
	client *ogosecurity.Client
}

func (d *contractsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shield_contracts"
}

func (d *contractsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"contracts": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"number": schema.StringAttribute{
							Computed:    true,
							Description: "Number used to reference this contract.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the contract.",
						},
					},
				},
			},
		},
		MarkdownDescription: "Get list of contract and associated informations.\n\n" +
			"Use this data source to retrieve list of available contract and related informations, " +
			"in particular contract Number needed to create new site.",
	}
}

// Configure adds the provider configured client to the data source.
func (d *contractsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *contractsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state contractsDataSourceModel

	contracts, err := d.client.GetAllContracts()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Ogo Contracts",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, c := range contracts {
		contractState := contractsModel{
			Number: types.StringValue(c.Number),
			Name:   types.StringValue(c.Name),
		}

		state.Contracts = append(state.Contracts, contractState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

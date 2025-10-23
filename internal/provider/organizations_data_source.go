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
	_ datasource.DataSource              = &organizationsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationsDataSource{}
)

// organizationsDataSourceModel maps the data source schema data
type organizationsDataSourceModel struct {
	Organizations []organizationsModel `tfsdk:"organizations"`
}

// organizationModel maps organization schema data
type organizationsModel struct {
	Code        types.String `tfsdk:"code"`
	CompanyName types.String `tfsdk:"company_name"`
}

func NewOrganizationsDataSource() datasource.DataSource {
	return &organizationsDataSource{}
}

type organizationsDataSource struct {
	client *ogosecurity.Client
}

func (d *organizationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shield_organizations"
}

func (d *organizationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							Computed:    true,
							Description: "Code identifier used to reference this organization.",
						},
						"company_name": schema.StringAttribute{
							Computed:    true,
							Description: "Company name of this organization.",
						},
					},
				},
			},
		},
		MarkdownDescription: "Get a list of organizations.\n\n" +
			"Use this data source to retrieve the list of available organizations, " +
			"in particular the organization code needed to create a new site.",
	}
}

// Configure adds the provider configured client to the data source.
func (d *organizationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organizationsDataSourceModel

	organizations, err := d.client.GetAllOrganizations()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Ogo Organizations",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, c := range organizations {
		organizationState := organizationsModel{
			Code:        types.StringValue(c.Code),
			CompanyName: types.StringValue(c.CompagnyName),
		}

		state.Organizations = append(state.Organizations, organizationState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

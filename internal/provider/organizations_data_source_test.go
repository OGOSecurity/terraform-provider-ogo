// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationsDataSource(t *testing.T) {
	providerConfig := testAccProviderConfig()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "ogo_shield_organizations" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ogo_shield_organizations.test", "organizations.0.code", "unit1896"),
					resource.TestCheckResourceAttr("data.ogo_shield_organizations.test", "organizations.0.company_name", "UnitTest-Terraform"),
				),
			},
		},
	})
}

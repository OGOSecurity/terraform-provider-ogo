// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContractsDataSource(t *testing.T) {
	providerConfig := testAccProviderConfig()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "ogo_shield_contracts" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.number", "unitt-40466"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.name", "UnitTest-Terraform"),
				),
			},
		},
	})
}

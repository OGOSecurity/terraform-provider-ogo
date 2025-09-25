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
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.type", "POC"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.bandwidth_per_month", "500"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.million_requests_per_month", "5"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.nb_sites_advanced", "15"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.nb_sites_expert", "10"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.cdn_enabled", "true"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.start_date", "2025-09-01"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.end_date", "2042-09-01"),
					resource.TestCheckResourceAttr("data.ogo_shield_contracts.test", "contracts.0.renewal_date", "2026-09-01"),
				),
			},
		},
	})
}

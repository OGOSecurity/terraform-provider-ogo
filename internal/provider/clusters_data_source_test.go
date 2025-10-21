// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClustersDataSource(t *testing.T) {
	providerConfig := testAccProviderConfig()

	clusterUid := os.Getenv("OGO_CLUSTER_UID")
	if clusterUid == "" {
		t.Errorf("OGO_CLUSTER_UID must be set")
	}

	clusterName := os.Getenv("OGO_CLUSTER_NAME")
	if clusterName == "" {
		t.Errorf("OGO_CLUSTER_NAME must be set")
	}

	clusterEntrypoint4 := os.Getenv("OGO_CLUSTER_HOST4")
	if clusterEntrypoint4 == "" {
		t.Errorf("OGO_CLUSTER_HOST4 must be set")
	}

	clusterEntrypoint6 := os.Getenv("OGO_CLUSTER_HOST6")
	if clusterEntrypoint6 == "" {
		t.Errorf("OGO_CLUSTER_HOST6 must be set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "ogo_shield_clusters" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.entrypoint4", clusterEntrypoint4),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.entrypoint6", clusterEntrypoint6),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.ips_to_whitelist.#", "2"),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.name", clusterName),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.supported_cdns.#", "1"),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.supports_cache", "true"),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.supports_ipv6_origins", "true"),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.supports_mtls", "true"),
					resource.TestCheckResourceAttr("data.ogo_shield_clusters.test", "clusters.0.uid", clusterUid),
				),
			},
		},
	})
}

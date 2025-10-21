// Copyright (c) OgoSecurity, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteResource(t *testing.T) {
	providerConfig := testAccProviderConfig()

	clusterUid := os.Getenv("OGO_CLUSTER_UID")
	if clusterUid == "" {
		t.Errorf("OGO_CLUSTER_UID must be set")
	}

	tlsOptionsUid := os.Getenv("OGO_TLSOPTIONS_UID")
	if tlsOptionsUid == "" {
		t.Errorf("OGO_TLSOPTIONS_UID must be set")
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
			// Create and Read testing
			{
				Config: providerConfig + `
resource "ogo_shield_site" "test" {
  domain_name             = "bar.example.com"
  cluster_uid             = "` + clusterUid + `"
  origin_server           = "172.18.1.11"
  origin_scheme           = "http"
  origin_skip_cert_verify = true
  origin_mtls_enabled     = false
  remove_xforwarded       = true
  log_export_enabled      = false
  force_https             = false
  audit_mode              = true
  cache_enabled           = true
  cdn                     = "ORANGE"
  passthrough_mode        = true
  hsts                    = "hsts"
  tags                    = ["app", "dev"]
  blacklisted_countries   = ["CN"]
  active_customer_certificate = {
    p12_content64 = <<-EOT
      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
      EOT
	p12_password = "Pr@t3ctMe!"
  }
  brain_overrides         = {
    "/BRAIN/DRIVE_43E5A99D_5A09_47FC_A9D9_C4FF0248B6C1_Priority": 0.0,
    "/ACTOR/DRIVE_493EE2EC_4776_4A98_8D56_75C2DDD28215_BELIEF": 0.8,
    "/ACTOR/DRIVE_E550246A_FA9A_4EB9_AFA5_C4C3D7C3FBA8_MAX_CONSUMED_RESPONSE_TIME": 200000.0,
    "/ACTOR/DRIVE_01234567_BELIEF": 0.7
  }
  rewrite_rules = [
    {
      active              = true
      comment             = "Rewrite old to new"
      rewrite_source      = "^/old"
      rewrite_destination = "/new"
    },
    {
      active              = true
      comment             = "Rewrite from to"
      rewrite_source      = "^/from"
      rewrite_destination = "/to"
    }
  ]
  rules = [
    {
      comment         = "Admin from office"
      paths           = ["/admin", "/wp-admin"]
      whitelisted_ips = ["fded:b552:6f7e:fc6f::/64", "10.10.9.0/24"]
    },
    {
      comment         = "Monitoring from office"
      paths           = ["/monitor"]
      whitelisted_ips = ["10.10.10.1/32"]
    }
  ]
  url_exceptions = [
    {
      path    = "/monitoring"
      comment = "Supervision"
    }
  ]
  ip_exceptions = [
    {
      ip      = "131.220.78.219/32"
      comment = "Home IPv4"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first site
					resource.TestCheckResourceAttr("ogo_shield_site.test", "domain_name", "bar.example.com"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_server", "172.18.1.11"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_scheme", "http"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_skip_cert_verify", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_mtls_enabled", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "remove_xforwarded", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "log_export_enabled", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "force_https", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "audit_mode", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cache_enabled", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cdn", "ORANGE"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cdn_status", "ACTIVATION_IN_PROGRESS"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cluster_entrypoint_4", clusterEntrypoint4),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cluster_entrypoint_6", clusterEntrypoint6),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cluster_entrypoint_cdn", "cl-gla36e56b1.maps.cdn.orange.com"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "hsts", "hsts"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "status", "CREATED"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "tags.*", "app"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "tags.*", "dev"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "active_customer_certificate.hash", "6b0fe950fa7935cf8c55c790398e3093fad331183b59c47c8b63b3d02f7c9b5a"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "blacklisted_countries.#", "1"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "blacklisted_countries.*", "CN"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "brain_overrides.%", "4"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "brain_overrides./ACTOR/DRIVE_01234567_BELIEF", "0.7"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "brain_overrides./ACTOR/DRIVE_493EE2EC_4776_4A98_8D56_75C2DDD28215_BELIEF", "0.8"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "brain_overrides./ACTOR/DRIVE_E550246A_FA9A_4EB9_AFA5_C4C3D7C3FBA8_MAX_CONSUMED_RESPONSE_TIME", "200000"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "brain_overrides./BRAIN/DRIVE_43E5A99D_5A09_47FC_A9D9_C4FF0248B6C1_Priority", "0"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.active", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.comment", "Rewrite old to new"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.rewrite_source", "^/old"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.rewrite_destination", "/new"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.active", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.comment", "Rewrite from to"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.rewrite_source", "^/from"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.rewrite_destination", "/to"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.comment", "Admin from office"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.paths.0", "/admin"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.paths.1", "/wp-admin"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.whitelisted_ips.#", "2"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.1.comment", "Monitoring from office"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.1.paths.0", "/monitor"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.1.whitelisted_ips.0", "10.10.10.1/32"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "url_exceptions.0.path", "/monitoring"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "url_exceptions.0.comment", "Supervision"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "ip_exceptions.0.ip", "131.220.78.219/32"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "ip_exceptions.0.comment", "Home IPv4"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ogo_shield_site.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "ogo_shield_site.test",
				ImportStateId:                        "bar.example.com",
				ImportStateVerifyIdentifierAttribute: "domain_name",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"last_updated", "active_customer_certificate"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "ogo_shield_site" "test" {
  domain_name             = "bar.example.com"
  cluster_uid             = "` + clusterUid + `"
  origin_server           = "172.18.1.12"
  origin_scheme           = "https"
  origin_skip_cert_verify = true
  origin_mtls_enabled     = false
  remove_xforwarded       = false
  log_export_enabled      = false
  force_https             = true
  audit_mode              = false
  passthrough_mode        = false
  hsts                    = "hstss"
  tlsoptions_uid          = "` + tlsOptionsUid + `"
  pass_tls_client_cert    = "all"
  tags                    = ["app", "prod", "platinium"]
  blacklisted_countries   = ["DE", "CN", "IT"]
  brain_overrides         = {
    "/ACTOR/DRIVE_493EE2EC_4776_4A98_8D56_75C2DDD28215_BELIEF": 0.5,
  }
  rewrite_rules = [
    {
      active              = false
      comment             = "Rewrite old to new"
      rewrite_source      = "^/old"
      rewrite_destination = "/new"
    },
    {
      active              = true
      comment             = "Rewrite informations"
      rewrite_source      = "^/informations"
      rewrite_destination = "/contacts"
    }
  ]
  rules = [
    {
      comment         = "Admin from office"
      paths           = ["/admin", "/wp-admin"]
      whitelisted_ips = ["fded:b552:6f7e:fc6f::/64", "10.10.9.0/24"]
    },
    {
      comment         = "Monitoring from internal"
      paths           = ["/health"]
      whitelisted_ips = ["10.10.10.10/32"]
    }
  ]
  url_exceptions = [
    {
      path    = "/health"
      comment = "Health check endpoint"
    }
  ]
  ip_exceptions = [
    {
      ip      = "fda1:a9bb:d292:ada6::/64"
      comment = "Home IPv6"
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first site
					resource.TestCheckResourceAttr("ogo_shield_site.test", "domain_name", "bar.example.com"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_server", "172.18.1.12"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_scheme", "https"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_skip_cert_verify", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "origin_mtls_enabled", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "remove_xforwarded", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "log_export_enabled", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "force_https", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "audit_mode", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cache_enabled", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cluster_entrypoint_4", clusterEntrypoint4),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "cluster_entrypoint_6", clusterEntrypoint6),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "passthrough_mode", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "hsts", "hstss"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "tlsoptions_uid", tlsOptionsUid),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "pass_tls_client_cert", "all"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "tags.#", "3"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "tags.*", "app"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "tags.*", "prod"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "tags.*", "platinium"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "blacklisted_countries.#", "3"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "blacklisted_countries.*", "DE"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "blacklisted_countries.*", "CN"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "blacklisted_countries.*", "IT"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "brain_overrides.%", "1"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "brain_overrides./ACTOR/DRIVE_493EE2EC_4776_4A98_8D56_75C2DDD28215_BELIEF", "0.5"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.active", "false"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.comment", "Rewrite old to new"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.rewrite_source", "^/old"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.0.rewrite_destination", "/new"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.active", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.comment", "Rewrite informations"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.rewrite_source", "^/informations"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rewrite_rules.1.rewrite_destination", "/contacts"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.comment", "Admin from office"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.paths.0", "/admin"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.paths.1", "/wp-admin"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.0.whitelisted_ips.#", "2"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.1.comment", "Monitoring from internal"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.1.paths.0", "/health"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "rules.1.whitelisted_ips.0", "10.10.10.10/32"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "url_exceptions.0.path", "/health"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "url_exceptions.0.comment", "Health check endpoint"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "ip_exceptions.0.ip", "fda1:a9bb:d292:ada6::/64"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "ip_exceptions.0.comment", "Home IPv6"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

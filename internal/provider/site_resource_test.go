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
  passthrough_mode        = true
  hsts                    = "hsts"
  tlsoptions_uid          = "` + tlsOptionsUid + `"
  pass_tls_client_cert    = "info"
  tags                    = ["app", "dev"]
  blacklisted_countries   = ["CN"]
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
					resource.TestCheckResourceAttr("ogo_shield_site.test", "passthrough_mode", "true"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "hsts", "hsts"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "tlsoptions_uid", tlsOptionsUid),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "pass_tls_client_cert", "info"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "tags.*", "app"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "tags.*", "dev"),
					resource.TestCheckResourceAttr("ogo_shield_site.test", "blacklisted_countries.#", "1"),
					resource.TestCheckTypeSetElemAttr("ogo_shield_site.test", "blacklisted_countries.*", "CN"),
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
				ImportStateVerifyIgnore:              []string{"last_updated"},
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

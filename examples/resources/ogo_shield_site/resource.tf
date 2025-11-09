# Define variable with default cluster UID to used in sites definition
variable "cluster_uid" {
  type        = string
  description = "Default Cluster ID where sites will be provisioned"
  default     = "802448cf-e2f9-40eb-b0d8-2983e018a0f4"
}

# Simple example with only required attributes
resource "ogo_shield_site" "foo_example_com" {
  domain_name   = "foo.example.com"
  cluster_uid   = var.cluster_uid
  origin_server = "172.18.1.10"
}

# Complete example with all attributes
resource "ogo_shield_site" "bar_example_com" {
  domain_name             = "bar.example.com"
  cluster_uid             = var.cluster_uid
  origin_server           = "172.18.1.11"
  origin_scheme           = "https"
  origin_skip_cert_verify = true
  origin_mtls_enabled     = false
  remove_xforwarded       = false
  log_export_enabled      = false
  force_https             = true
  audit_mode              = false
  passthrough_mode        = false
  hsts                    = "hstss"
  tlsoptions_uid          = "example00812-f4d2574e-d85e-5dg7-ad11-1edd0489jmp1"
  pass_tls_client_cert    = "info"
  tags                    = ["app", "dev"]
  blacklisted_countries   = ["DE", "ES", "CN"]
  brain_overrides = {
    "/BRAIN/DRIVE_43E5A99D_5A09_47FC_A9D9_C4FF0248B6C1_Priority" : 0.0,
    "/ACTOR/DRIVE_493EE2EC_4776_4A98_8D56_75C2DDD28215_BELIEF" : 0.8,
    "/ACTOR/DRIVE_E550246A_FA9A_4EB9_AFA5_C4C3D7C3FBA8_MAX_CONSUMED_RESPONSE_TIME" : 200000.0,
    "/ACTOR/DRIVE_01234567_BELIEF" : 0.7
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
      path    = "/demo"
      comment = "demo"
    },
    {
      path    = "/monitoring"
      comment = "Supervision"
    },
  ]
  ip_exceptions = [
    {
      ip      = "131.220.78.219/32"
      comment = "Home IPv4"
    },
    {
      ip      = "fda1:a9bb:d292:ada6::/64"
      comment = "Home IPv6"
    },
  ]
}

# Define variable with default cluster UID to used in sites definition
variable "cluster_uid" {
  type        = string
  description = "Default Cluster ID where sites will be provisioned"
  default     = "802448cf-e2f9-40eb-b0d8-2983e018a0f4"
}

# Simple example with only required attributes
resource "ogo_shield_site" "foo_example_com" {
  domain_name          = "foo.example.com"
  cluster_uid          = var.cluster_uid
  origin_server        = "172.18.1.10"
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
  tlsoptions_uid         = "example00812-f4d2574e-d85e-5dg7-ad11-1edd0489jmp1"
  pass_tls_client_cert    = "info"
  tags                    = ["app", "dev"]
  blacklisted_countries   = ["DE", "EN", "FR", "IT"]
  rewrite_rules = [
    {
      active              = true
      comment             = "Rewrite old to new"
      priority            = 1
      rewrite_destination = "/new"
      rewrite_source      = "^/old"
    },
    {
      active              = true
      comment             = "Rewrite from to"
      priority            = 2
      rewrite_destination = "/to"
      rewrite_source      = "^/from"
    }
  ]
  rules = [
    {
      priority        = 1
      comment         = "Admin from office"
      paths           = ["/admin", "/wp-admin"]
      whitelisted_ips = ["fded:b552:6f7e:fc6f::/64", "10.10.9.0/24"]
    },
    {
      priority        = 2
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

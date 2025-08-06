terraform {
  required_providers {
    ogo = {
      source = "ogosecurity.com/ogosecurity/ogo"
    }
  }
}

# Provider
## Provider variables
variable "ogo_endpoint" {
  type        = string
  default     = "https://api-stg.ogosecurity.com"
  description = "Ogo Dashboard API endpoint"
}

variable "ogo_username" {
  type        = string
  description = "Username to access Ogo Dashboard API"
}

variable "ogo_apikey" {
  type        = string
  description = "API Key used to authenticate to Ogo Dashboard API"
}

## Provider configuration
provider "ogo" {
  endpoint = var.ogo_endpoint
  username = var.ogo_username
  apikey   = var.ogo_apikey
}




# Datasources
## Datasource variables
variable "cluster_uid" {
  type        = string
  description = "Default Cluster ID to be used"
  default     = "67f276d4-d71b-4941-95e5-e86f07647a41" # OGO GYS
}

## Datasource clusters
data "ogo_shield_clusters" "shield" {}
//output "shield_clusters" {
//  value = data.ogo_shield_clusters.shield
//}

## Datasource tlsoptions
data "ogo_shield_tlsoptions" "tlsoptions" {}

## Datasource outputs
output "default_cluster_uid" {
  value = var.cluster_uid
}


# Resources
## Resources site
#resource "ogo_shield_site" "gys_tf_ogosecurity_com" {
#  domain_name             = "gys-tf.ogosecurity.com"
#  cluster_uid             = var.cluster_uid
#  origin_server           = "192.168.122.13"
#  origin_scheme           = "https"
#  origin_skip_cert_verify = true
#  origin_mtls_enabled     = false
#  remove_xforwarded       = false
#  log_export_enabled      = false
#  force_https             = true
#  audit_mode              = false
#  passthrough_mode        = false
#  hsts                    = "hstss"
#  tls_options_uid         = "ogo00795-f4c2670b-d75b-4cd8-ad11-1edd8409bfc0"
#  pass_tls_client_cert    = "info"
#  tags                    = ["test", "dev"]
#  blacklisted_countries   = ["IT", "FR"]
#  rewrite_rules = [
#    {
#      active              = true
#      comment             = "Rewrite old to new"
#      priority            = 1
#      rewrite_destination = "/new"
#      rewrite_source      = "^/old"
#    },
#    {
#      active              = true
#      comment             = "Rewrite from to"
#      priority            = 2
#      rewrite_destination = "/to"
#      rewrite_source      = "^/from"
#    }
#  ]
#  rules = [
#    {
#      priority        = 1
#      comment         = "Admin from home"
#      paths           = ["/admin", "/backoffice"]
#      whitelisted_ips = ["2a01:e0a:4:4e30:b001:6e70:cf1e:56df/64", "82.65.146.54/32"]
#    },
#    {
#      priority        = 2
#      comment         = "Admin from office"
#      paths           = ["/wp-admin"]
#      whitelisted_ips = ["10.10.10.1/32"]
#    }
#  ]
#  url_exceptions = [
#    {
#      path    = "/demo"
#      comment = "demo"
#    },
#    {
#      path    = "/monitoring"
#      comment = "Supervision"
#    },
#  ]
#  ip_exceptions = [
#    {
#      ip      = "82.65.146.54/32"
#      comment = "Home IPv4"
#    },
#    {
#      ip      = "2a01:e0a:4:4e30:5201:9e56:7417:38c9/64"
#      comment = "Home IPv6"
#    },
#  ]
#}
resource "ogo_shield_site" "gys_tf_ogosecurity_com" {
  domain_name          = "gys-tf.ogosecurity.com"
  cluster_uid          = var.cluster_uid
  origin_server        = "192.168.122.13"
}

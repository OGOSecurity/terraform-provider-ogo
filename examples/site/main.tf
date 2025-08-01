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
resource "ogo_shield_site" "gys_tf_ogosecurity_com" {
  name                 = "gys-tf.ogosecurity.com"
  cluster_uid          = var.cluster_uid
  dest_host            = "192.168.122.13"
  dest_host_scheme     = "https"
  dest_host_mtls       = false
  log_export           = false
  trust_selfsigned     = false
  no_copy_xforwarded   = false
  force_https          = true
  dry_run              = true
  panic_mode           = false
  hsts                 = "hstss"
  pass_tls_client_cert = "info"
  tls_options_uid      = "ogo00795-f4c2670b-d75b-4cd8-ad11-1edd8409bfc0"
  tags                 = ["staging", "demo"]
}

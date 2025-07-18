terraform {
  required_providers {
    ogo = {
      source = "ogosecurity.com/ogosecurity/ogo"
    }
  }
}

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

provider "ogo" {
  endpoint = var.ogo_endpoint
  username = var.ogo_username
  apikey   = var.ogo_apikey
}

data "ogo_shield_clusters" "shield" {}
//output "shield_clusters" {
//  value = data.ogo_shield_clusters.shield
//}

resource "ogo_shield_site" "gys_webapp_ogosecurity_com" {
  name               = "gys-webapp.ogosecurity.com"
  cluster_name       = "OGO GYS"
  dest_host          = "192.168.122.13"
  dest_host_scheme   = "https"
  trust_selfsigned   = false
  no_copy_xforwarded = false
  force_https        = true
  dry_run            = true
  panic_mode         = false
}

variable "cluster_name" {
  type        = string
  description = "Default Cluster ID to be used"
  default     = "c1-stg3"
}

output "default_cluster_name" {
  value = var.cluster_name
}

output "origin_stub" {
  value = ogo_shield_site.gys_webapp_ogosecurity_com.dest_host
}

output "gys_webapp_site" {
  value = ogo_shield_site.gys_webapp_ogosecurity_com
}


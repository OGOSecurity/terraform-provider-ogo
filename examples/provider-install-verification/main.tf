terraform {
  required_providers {
    ogo = {
      source = "ogosecurity.com/ogosecurity/ogo"
    }
  }
}

provider "ogo" {
  endpoint = "https://api-stg.ogosecurity.com"
}

data "ogo_clusters" "shield" {}

output "shield_clusters" {
  value = data.ogo_clusters.shield
}

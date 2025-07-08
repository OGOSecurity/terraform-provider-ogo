terraform {
  required_providers {
    ogo = {
      source = "ogosecurity.com/ogosecurity/ogo"
    }
  }
}

provider "ogo" {
  endpoint = "https://api-stg.ogosecurity.com"
  username = "changeme"
  apikey   = "changeme"
}

data "ogo_clusters" "shield" {}
output "shield_clusters" {
  value = data.ogo_clusters.shield
}

resource "ogo_site" "gys_tf_ogosecurity_com" {
  name               = "gys-tf.ogosecurity.com"
  cluster_id         = 495
  dest_host          = "192.168.122.15"
  dest_host_scheme   = "http"
  trust_selfsigned   = true
  no_copy_xforwarded = false
  force_https        = false
  dry_run            = true
  panic_mode         = false
}

resource "ogo_site" "stubapache_ogosecurity_com" {
  name               = "stubapache.ogosecurity.com"
  cluster_id         = 5
  dest_host          = "195.154.168.43"
  dest_host_scheme   = "https"
  trust_selfsigned   = false
  no_copy_xforwarded = false
  force_https        = true
  dry_run            = true
  panic_mode         = false
}

//output "gys_site" {
//  value = ogo_site.gys_tf_ogosecurity_com
//}

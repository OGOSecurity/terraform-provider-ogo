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

variable "ogo_organization" {
  type        = string
  description = "Organization to access Ogo Dashboard API"
}

variable "ogo_apikey" {
  type        = string
  description = "API Key used to authenticate to Ogo Dashboard API"
}

## Provider configuration
provider "ogo" {
  endpoint     = var.ogo_endpoint
  organization = var.ogo_organization
  apikey       = var.ogo_apikey
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
#resource "ogo_shield_site" "min_tf_ogosecurity_com" {
#  domain_name          = "min-tf.ogosecurity.com"
#  cluster_uid          = var.cluster_uid
#  origin_server        = "192.168.122.13"
#  tlsoptions_uid      = "ogo00795-f4c2670b-d75b-4cd8-ad11-1edd8409bfc0"
#  pass_tls_client_cert = "none"
#}

resource "ogo_shield_site" "gys_tf_ogosecurity_com" {
  domain_name             = "gys-tf.ogosecurity.com"
  cluster_uid             = var.cluster_uid
  origin_server           = "192.168.122.13"
  origin_scheme           = "https"
  origin_skip_cert_verify = true
  origin_mtls_enabled     = false
  remove_xforwarded       = false
  log_export_enabled      = false
  force_https             = true
  audit_mode              = false
  passthrough_mode        = false
  hsts                    = "hstss"
  tlsoptions_uid          = "ogo00795-f4c2670b-d75b-4cd8-ad11-1edd8409bfc0"
  pass_tls_client_cert    = "info"
  tags                    = ["test", "dev"]
  blacklisted_countries   = ["IT", "FR"]
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
      action          = "brain"
      comment         = "Admin from office"
      paths           = ["/wp-admin"]
      whitelisted_ips = ["10.10.10.1/32"]
    },
    {
      comment         = "Admin from home"
      paths           = ["/admin", "/backoffice"]
      whitelisted_ips = ["2a01:e0a:4:4e30:b001:6e70:cf1e:56df/64", "82.65.146.54/32"]
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
      ip      = "82.65.146.54/32"
      comment = "Home IPv4"
    },
    {
      ip      = "2a01:e0a:4:4e30:5201:9e56:7417:38c9/64"
      comment = "Home IPv6"
    },
  ]
}

resource "ogo_shield_tlsoptions" "gys_tf" {
  name             = "gys-tf"
  min_tls_version  = "TLS_12"
  client_auth_type = "VerifyClientCertIfGiven"
  client_auth_ca_certs = [
    <<-EOT
-----BEGIN CERTIFICATE-----
MIIE6jCCAtKgAwIBAgIUalcKYItsBi+0uncayhnAEHKeZiQwDQYJKoZIhvcNAQEL
BQAwEzERMA8GA1UEAwwIdG90by5jb20wHhcNMjUwNDE3MTMxMTU4WhcNMzcxMjMx
MTMxMTU4WjATMREwDwYDVQQDDAh0b3RvLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAJPxW2fieWNzQrOgvEPV4iPZgiJcLGfwtMVR4SAL2zLcNhVL
p1peOqlsbs2YCx3i3lhaeXQYMSmRs/wT04ls/QkKfqkyRphhSyYtM6UFJZfdLnDl
gI5tUlys9V7XpHIVgsPZNqQ9WZTkboay7oO/1FX4JY+mvwtY9B0J+VQ0T2CcOQkP
utbIRURDcSFaZIoAtlLH7KS5oaTl18PcaKpK8nV5ohYLyB8EtxNRNW3/c6swPx7+
36zREbW+kG+MsEA4lKNPzxYUuZGIpmfq06TYoJtz3IRtEbgMb05Q+cj57alqJHpr
wOhKVpxkKcVStj0yNT5tG69SD/9WFCSfsIWfG3lb9ewJD9gJbqHOALf6CYuAQLbz
gQgrWUJVE1saMet3f+7BW+C1ywpHvNGpWJ5qNklEIneRXEE1+DKOQBKDAEetIVSp
R8AK2Y8m6aGq4UD2ohtX255SdVk2I5AprXpo2KrqINn4K9AZBD3V2FzSTmgMScR3
1X3DDm2KGldjZQubF47AwThRb8v1186Onx4WnTsBXCZE2o3cII5U30heEs8o8N6S
ciejBYjqLdGI0v0RpGbrV0ReJugkHL/pQYUyyX8vdR53KLq9rN39dze3lNfbsAKT
XItUDIXMLsz9pj3UU6FCk5VsKkFxcLdQiSbzEfTvF/9xq9NNoI5t9cISF51JAgMB
AAGjNjA0MBMGA1UdEQQMMAqCCHRvdG8uY29tMB0GA1UdDgQWBBRTldEdG8cQSU7+
kYoqzBuKtZZQiDANBgkqhkiG9w0BAQsFAAOCAgEAe3Yv9lxb5jJw0e1mvxhm1yzA
VDzwMOiohOKepd+hS9d2GDrelU+WYkrUk7sztKL/Z1ES4XM0BBnmi3PDeAdpiyka
qvFGw54tRGscMNuzKJIzsykcEcs3PCHbl3rEa5e0uaPQBNTaJS384o8/n3KeTsOB
UY4mfqyqzxYZ9mOnpYB0pOv8+0qM5WJEznsqavqFi1xA7trtjL/vjVJrQtqJrc3j
1rRHe9StX2KdQgcFH4hNVrERmJXnkweH7uon9cyNsWrsUWNU8Yxc5/jxHihKHf/G
Kl0/ciRO7K50P7lKhBD3sXGk2tCzxBDurNVXOA7iL3vjcrqfibYNfeEZTEazwaCk
chiFmyyLmZny/ifJT826/Zmz0S2HHMwcLc5/n79d9yW3g/0CgfqjBjEU0PnAkb4j
ev8sSILccM+UMNGEsWe8eA8sas4nX3t9tPWTzHHGRIY4JlyDleqFsLHpnl6Bu4Oe
90CHXyjPNa04c2+0HbnV5dA/MTcyDMCUMTPK89Bcsiy/8PG+Ppdu4bt+E6QI6Tum
BHBrNTm3yxF4l7kkCg6zMGQGCLCssRTkjlNYudE1Al6jQzSMdqu6QgP6rzv1/3dj
WZSsQcTqb56y5FJHFK4tYHoM4m4EzqYlvKjyv+1tsZ6gCH3JcemxW2R6eohD/yba
zCESYTnxjtaeZU7HMn4=
-----END CERTIFICATE-----
EOT
    ,
    <<-EOT
-----BEGIN CERTIFICATE-----
MIIElDCCAvygAwIBAgICAXswDQYJKoZIhvcNAQELBQAwLjEYMBYGA1UECgwPT0dP
U0VDVVJJVFkuQ09NMRIwEAYDVQQDDAlPR08tQURNSU4wHhcNMjUwNDI3MjIxNzUy
WhcNMjcwNDI4MjIxNzUyWjAuMRwwGgYDVQQKDBNBRE0uT0dPU0VDVVJJVFkuTEFO
MQ4wDAYDVQQDDAVzZWdpcjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AMLaNpRz6Mln8Ip7rPiwsR1cbztTJ0/8PIA/WlD0UOjBPdLHUG+ZrXvwvvnYuoqX
S4icuRq80Y0hKlOtUwF1+WeOndvsMaOM9cDtsxE8D8gaajlFy5G/i1GOLN4aPXYm
hciUe2PQgv2ryjheCS8HykrQuw4StSWdGbUzspXj0ucSHRxijcLXu1nuBuTngLH5
lKw0SdmgeLD3iZBtCJ+yHO15MnXa9KBSvsOA+hFqoA0l4UuqoJHwotoK9kKpSlCH
9a/1ggYbjX3pnsiLVaGYVWN/oUi+82m/zuKeGyPo3RdKcgHIT6LjYEFBE/fBZsQD
nFz0jr0MELqxN0ZExvgEXBUCAwEAAaOCATowggE2MB8GA1UdIwQYMBaAFIbHkmfH
ub6ZMmZ9UHiOAgQPjdmFMEUGCCsGAQUFBwEBBDkwNzA1BggrBgEFBQcwAYYpaHR0
cDovL2lwYS1jYS5hZG0ub2dvc2VjdXJpdHkubGFuL2NhL29jc3AwDgYDVR0PAQH/
BAQDAgTwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjB+BgNVHR8EdzB1
MHOgO6A5hjdodHRwOi8vaXBhLWNhLmFkbS5vZ29zZWN1cml0eS5sYW4vaXBhL2Ny
bC9NYXN0ZXJDUkwuYmluojSkMjAwMQ4wDAYDVQQKDAVpcGFjYTEeMBwGA1UEAwwV
Q2VydGlmaWNhdGUgQXV0aG9yaXR5MB0GA1UdDgQWBBTAVA36cQC0qwu7vUc1a4J9
MPKugDANBgkqhkiG9w0BAQsFAAOCAYEAX6Xj6Wo96h3hkEZmTWqerEE/TH16W7Pp
82daqpgQiWhFTaXiQcVZtBw1/ou6yOGFueLdjTGeC7Lr9XOzG4ZX5A+7bxfw8sLq
ab8y4xw0+M/n6dS5D/SRx+WuUK/hWciPcAdY1fCoLrfghp9UVtoFJBCTeUr3RnfM
SDqdfN1V9kIjHf5jfDryyYBcdYtGdjQO2izWbXPcUWDHpiRZ+kaX1X5vyEUZfKSw
LN2saw6cZohXfNu9bfQykD9Qx/UiN5bEcs7TZzolK+xPSQZUkCNdd1f3IQCTgHQE
yubOaofMxN8lJE089HbK5zljsJ8fYkbYWGL4juUI0CiCJ4CY94XoJOTMJdfwcdWc
SFU8utLNQd+QfRcH9/4DPpO1oHvmw3Y3xjl8uzO/z/x5ssNf98C73Mp40phF16/J
NYx6TMbY7rF14TMTT3FHTTpPKfIKTrcj/xk6dRDJjJTKIGBv23eoqps4bC8eJ/yn
PAV3CIWULm1Xcp4yhwdyimWWk8nKqDbp
-----END CERTIFICATE-----
EOT
    ,
  ]
}

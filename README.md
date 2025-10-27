# OgoSecurity Terraform Provider

The OgoSecurity provider allows Terraform to manage your Site resources configuration through Ogo Dashboard API.

⚠️  Provider is under development and may lead to changes, please feel free to give us feedback! ⚠️ 

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0

## Usage

Create `main.tf` terraform configuration file with the following content:

```hcl
terraform {
  required_providers {
    ogo = {
      source = "ogosecurity/ogo"

    }
  }
}

# Provider configuration
provider "ogo" {
  # Ogo Dashboard API endpoint (or use env variable: OGO_ENDPOINT)
  endpoint     = "https://api.ogosecurity.com"
  # Organization to access Ogo Dashboard API (or use env variable: OGO_ORGANIZATION)
  organization = "orga04242"
  # Organization to access Ogo Dashboard API (or use env variable: OGO_APIKEY)
  apikey       = "cd583abf-a02f-49e7-949e-424bf42419ea"
}

# Datasources
## Clusters datasource
data "ogo_shield_clusters" "shield" {}

## TLS options datasource
data "ogo_shield_tlsoptions" "tlsoptions" {}

# Resources
## Define variable with default cluster UID to used in site resources definition
variable "cluster_uid" {
  type        = string
  description = "Default Cluster ID where sites will be provisioned"
  default     = "802448cf-e2f9-40eb-b0d8-2983e018a0f4"
}

## Simple example with only required attributes
resource "ogo_shield_site" "foo_example_com" {
  domain_name   = "foo.example.com"
  cluster_uid   = var.cluster_uid
  origin_server = "172.18.1.10"
}
```

Use `terraform init` command to initialize your project.


## Documentation

Full OgoSecurity provider documentation is available on the official Hashicorp Terraform provider registry [the Terraform Registry](https://registry.terraform.io/providers/ogosecurity/ogo/latest/docs).

Some examples can also be found in this [./examples](./examples) project directory.


## Contacts

If you believe you have found a security issue in the Terraform OgoSecurity Provider, please responsibly disclose it by contacting us at security@ogosecurity.com.

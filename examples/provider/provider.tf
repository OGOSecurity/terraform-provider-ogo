terraform {
  required_providers {
    ogo = {
      source  = "ogosecurity/ogo"
      version = "~> 0.1"
    }
  }
}

# Provider
## Provider variables
variable "ogo_endpoint" {
  type        = string
  default     = "https://api.ogosecurity.com"
  description = "Ogo Dashboard API endpoint (or use env variable OGO_ENDPOINT)"
}

variable "ogo_email" {
  type        = string
  description = "Email address used to authenticate to Ogo Dashboard API (or use env variable OGO_EMAIL)"
}

variable "ogo_apikey" {
  type        = string
  description = "API Key used to authenticate to Ogo Dashboard API (or use env variable OGO_APIKEY)"
}

variable "ogo_organization" {
  type        = string
  description = "Organization used to access Ogo Dashboard API (or use env variable OGO_ORGANIZATION)"
}


## Provider configuration
provider "ogo" {
  endpoint     = var.ogo_endpoint
  email        = var.ogo_email
  apikey       = var.ogo_apikey
  organization = var.ogo_organization
}

// Copyright (c) OGO Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "ogo" {
  endpoint     = "%s"
  email        = "%s"
  organization = "%s"
  apikey       = "%s"
}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ogo": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccProviderConfig() string {
	endpoint := os.Getenv("OGO_ENDPOINT")
	email := os.Getenv("OGO_EMAIL")
	organization := os.Getenv("OGO_ORGANIZATION")
	apikey := os.Getenv("OGO_APIKEY")

	return fmt.Sprintf(providerConfig, endpoint, email, organization, apikey)
}

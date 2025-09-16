package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTlsOptionsDataSource(t *testing.T) {
	providerConfig := testAccProviderConfig()

	tlsOptionsUid := os.Getenv("OGO_TLSOPTIONS_UID")
	if tlsOptionsUid == "" {
		t.Errorf("OGO_TLSOPTIONS_UID must be set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "ogo_shield_tlsoptions" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ogo_shield_tlsoptions.test", "tlsoptions.0.name", "UnitTest"),
					resource.TestCheckResourceAttr("data.ogo_shield_tlsoptions.test", "tlsoptions.0.client_auth_type", "VerifyClientCertIfGiven"),
					resource.TestCheckResourceAttr("data.ogo_shield_tlsoptions.test", "tlsoptions.0.min_tls_version", "TLS_12"),
					resource.TestCheckResourceAttr("data.ogo_shield_tlsoptions.test", "tlsoptions.0.uid", tlsOptionsUid),
					resource.TestCheckResourceAttr("data.ogo_shield_tlsoptions.test", "tlsoptions.0.client_auth_ca_certs.0", "-----BEGIN CERTIFICATE-----\nMIIDrzCCApegAwIBAgIUbKqK408DCxCOjSmmUQ08qlu3ptEwDQYJKoZIhvcNAQEL\nBQAwZzELMAkGA1UEBhMCRlIxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy\naXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MSEwHwYDVQQDDBh1bml0dGVzdC5vZ29z\nZWN1cml0eS5jb20wHhcNMjUwOTE2MDgzNTExWhcNMzUwOTE0MDgzNTExWjBnMQsw\nCQYDVQQGEwJGUjEPMA0GA1UECAwGRnJhbmNlMQ4wDAYDVQQHDAVQYXJpczEUMBIG\nA1UECgwLT2dvU2VjdXJpdHkxITAfBgNVBAMMGHVuaXR0ZXN0Lm9nb3NlY3VyaXR5\nLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKO5fFdD99piLXha\nF/GvRVcurOdz/I9XxlfsYJ/82WLAK8eiekcymod3fbHOZa7oDjXXBRqBBHx456H3\nVAlqv7TlyJ5I90rJlBk9ot40D69bTLuQU93jm02xxXRll3S+v36joVFpwE2rf0av\n1KjJNnDiR3uLvA9+hLqdryVGh1Mj+y4Chmc+zlJ1BWOi02dq9saOKmXd0MkTiHi6\nMBX91oBq6O/+VLl13keHpNwPsP5XeJvsFHcnj4vbuGHQQ74NJxh2DQBksPlmPatS\nr5km0z5HIm3ei7M6SZO7xMAG3PWyrc3WTHJx2hgqORJmWv4Fp88NM1CdWolzc6jv\nPcnCvg8CAwEAAaNTMFEwHQYDVR0OBBYEFMAI6xPvV3x3B48FYX44d5qP3EQxMB8G\nA1UdIwQYMBaAFMAI6xPvV3x3B48FYX44d5qP3EQxMA8GA1UdEwEB/wQFMAMBAf8w\nDQYJKoZIhvcNAQELBQADggEBAFn2SsMoGWfvxFQf1a56x+qTLm8gWQAyGXJh0lRp\nKN/DERkdUMMZ7Fa1h/rJZ3EkO/rHtPsiSrSn6Hl2tcajcEvIvh9nVeUAh0XF8leO\nH2VM9hbW7bJ7BJZ2lvPVDBcxnu/goFN6oNIUsSV7qy8uUhXIVFlWtes/P+1jedNd\nP4F1sQFpridU2lx4UTeP+Stq0SWq5ONNKC+TCIOZaWnqKS5NcjbSYdZv+BCZ7wT4\nQLSVInDJO5ddDcBBh8LQgyMHXs6iR0EOYR4cc/pIIRH1y4V0Yw/TY9dlxQhno7KX\nXLQPB27Jo1/0o8qQLwvC1D1Rcf0lNZBVVY/AtlEiCwO+TpI=\n-----END CERTIFICATE-----"),
				),
			},
		},
	})
}

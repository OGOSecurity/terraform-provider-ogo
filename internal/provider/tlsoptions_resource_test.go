package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccImportStateIdFromAttribute(resourceName string, attributeName string) resource.ImportStateIdFunc {
	importStateIdFunc := func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		id := rs.Primary.Attributes[attributeName]
		if id == "" {
			return "", fmt.Errorf("%s: Attribute '%s' not found", resourceName, attributeName)
		}

		return id, nil
	}

	return importStateIdFunc
}

func TestAccTlsOptionsResource(t *testing.T) {
	providerConfig := testAccProviderConfig()

	tlsOptionsUid := os.Getenv("OGO_TLSOPTIONS_UID")
	if tlsOptionsUid == "" {
		t.Errorf("OGO_TLSOPTIONS_UID must be set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "ogo_shield_tlsoptions" "test" {
  name            = "mTLS foo bar"
  min_tls_version = "TLS_12"
  client_auth_ca_certs = [
    <<-EOT
-----BEGIN CERTIFICATE-----
MIIDnTCCAoWgAwIBAgIUU5bIl5SJavP6YWPL/RUPLCbGu9owDQYJKoZIhvcNAQEL
BQAwXjELMAkGA1UEBhMCRlIxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy
aXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MRgwFgYDVQQDDA9mb28uZXhhbXBsZS5j
b20wHhcNMjUwODEzMDYyODM5WhcNMzUwODExMDYyODM5WjBeMQswCQYDVQQGEwJG
UjEPMA0GA1UECAwGRnJhbmNlMQ4wDAYDVQQHDAVQYXJpczEUMBIGA1UECgwLT2dv
U2VjdXJpdHkxGDAWBgNVBAMMD2Zvby5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBANZ0zrEH22IMXp8tQ5PbwLFHmCQRc1/T1ge5z7ho
p6zdyFn5GEFrMv1ZOywPBPlCz+Lb/5sWWj9qhcMw6JkPogKKVx9PQZDwfpc9ov+M
mujh/SM1Ms07AFt286h9e0yZzQfP9t6B9+Dns4Lgn6/+Ua8g7VW+Hrq3V9Ait0bx
kDOZUj0djOp9H3tShtgl8p9Z+dcqYIAPtkjSTt/U7jUDtR9PH6qz4/gXE/mCKq4e
Lf+63nLqGfZ3S1mIwjysRhPsJwy4g9v+E6fHO4Emfk4KF6EvFj3GVXyckxLbxKNa
1yYRUhSLZvNCNyDVvySVUta7yOhdzyC53YvS/Emtrh/7I6kCAwEAAaNTMFEwHQYD
VR0OBBYEFF3MC9L6J3lESlcdryXXFzC1uJGvMB8GA1UdIwQYMBaAFF3MC9L6J3lE
SlcdryXXFzC1uJGvMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB
ABUa9KqjfCB5lf9G7rpVTqhrg5LIKzuYzH4c7MZ76R4GZyH475yV5jQYCj2Qr3Pq
2m50UpEzVKICjfBSCbJulv5ZofSn8DWTpEBoLZA2pVM9yutI3wOQW350HX01nY82
9j9im1yMVtdf1uAPd1O3pm+RUcSICI5YBFQ1/LAAEoSqrmSoVUPwH6pt9Gr+4E/w
PpUcdAju8piy48Nx9ZD9vwCVjD67oRNnF00wEDJgrl8RpI9i0zOzflBxXyllGD8L
XT6wvPmUpso+jn04qnizfMWaYy9P2ip8RgOslrH6WIe6GyXSy6VjAu9JSuVE5OYX
XpYou5FLSGMhNaPTuaukAgY=
-----END CERTIFICATE-----
EOT
    ,
    <<-EOT
-----BEGIN CERTIFICATE-----
MIIDnTCCAoWgAwIBAgIUHvOpeMH+4Lk1ewQZKMwOygGBQe0wDQYJKoZIhvcNAQEL
BQAwXjELMAkGA1UEBhMCRlIxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy
aXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MRgwFgYDVQQDDA9iYXIuZXhhbXBsZS5j
b20wHhcNMjUwODEzMDYyODU1WhcNMzUwODExMDYyODU1WjBeMQswCQYDVQQGEwJG
UjEPMA0GA1UECAwGRnJhbmNlMQ4wDAYDVQQHDAVQYXJpczEUMBIGA1UECgwLT2dv
U2VjdXJpdHkxGDAWBgNVBAMMD2Jhci5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAO4dBU9DGbgBzjIYy/Qls0IglivSHyughVRa4nfZ
b2b3iGP1rEa+xNlnmOlgxp8ihjxF4yBz/DMGVDEDnErwITUOxEG4fJ5gdX7a5Iyd
OgYYyoh1RJKRkyWSGQoU4RmbVidTCyxq15j+yRBJDt3fll+Y9rlL+Ejl9QJCe+Zt
kSab7pBn9SmUzX8IeHyX1IpEMA4nNtFI8ysNSZNxPJa1hB3tXtVGZrkhpecCZvx4
IBpuRrjBSY3MaRE5YW51l7nC7jExC+IeNGe3mfKYUu0Re7fkK7n1auGmAJhTlzIR
4126rTJDbZlKyDFSfoaDFsyYeNe2t2W6KlhG4d0dSiFwIucCAwEAAaNTMFEwHQYD
VR0OBBYEFKBppFca57l7wutRyaIRZ3fwzZP7MB8GA1UdIwQYMBaAFKBppFca57l7
wutRyaIRZ3fwzZP7MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB
AGwIPOoZqd/3Uu3W2dUcd8HWw/VmjokjKrC811KUfhiijpFpQjGMcGQrjti3rIkk
5ZQyItkw91/IaUPOnyO8H5O/I/4RmTPaqbhmZ2gn8Ekw3/TO79tBB3bQWcfaSkK9
b+4+ryk2fCe3Um6Q/NCeSRwYe3Z8Xe5ByqJfjGmrXLyU//folGAtnx4uaAeJ98ze
jUXT17x8AbdEt2JIpYoJI7xFC8mOr0s3LvA/gFmpNkuRNbCNQF2v5Qt9L2AYT0Fv
B5uT42VuHQvRRNReAxa5oNGp/zcCjspaouPia03Tf5ZNZEd5LUFANLHPtsJg4jBB
kjHKjCnt0/9fttE1u/gMW7k=
-----END CERTIFICATE-----
EOT
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first tlsoptions
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "name", "mTLS foo bar"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "min_tls_version", "TLS_12"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "client_auth_ca_certs.#", "2"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "client_auth_ca_certs.0", "-----BEGIN CERTIFICATE-----\nMIIDnTCCAoWgAwIBAgIUHvOpeMH+4Lk1ewQZKMwOygGBQe0wDQYJKoZIhvcNAQEL\nBQAwXjELMAkGA1UEBhMCRlIxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy\naXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MRgwFgYDVQQDDA9iYXIuZXhhbXBsZS5j\nb20wHhcNMjUwODEzMDYyODU1WhcNMzUwODExMDYyODU1WjBeMQswCQYDVQQGEwJG\nUjEPMA0GA1UECAwGRnJhbmNlMQ4wDAYDVQQHDAVQYXJpczEUMBIGA1UECgwLT2dv\nU2VjdXJpdHkxGDAWBgNVBAMMD2Jhci5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAO4dBU9DGbgBzjIYy/Qls0IglivSHyughVRa4nfZ\nb2b3iGP1rEa+xNlnmOlgxp8ihjxF4yBz/DMGVDEDnErwITUOxEG4fJ5gdX7a5Iyd\nOgYYyoh1RJKRkyWSGQoU4RmbVidTCyxq15j+yRBJDt3fll+Y9rlL+Ejl9QJCe+Zt\nkSab7pBn9SmUzX8IeHyX1IpEMA4nNtFI8ysNSZNxPJa1hB3tXtVGZrkhpecCZvx4\nIBpuRrjBSY3MaRE5YW51l7nC7jExC+IeNGe3mfKYUu0Re7fkK7n1auGmAJhTlzIR\n4126rTJDbZlKyDFSfoaDFsyYeNe2t2W6KlhG4d0dSiFwIucCAwEAAaNTMFEwHQYD\nVR0OBBYEFKBppFca57l7wutRyaIRZ3fwzZP7MB8GA1UdIwQYMBaAFKBppFca57l7\nwutRyaIRZ3fwzZP7MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB\nAGwIPOoZqd/3Uu3W2dUcd8HWw/VmjokjKrC811KUfhiijpFpQjGMcGQrjti3rIkk\n5ZQyItkw91/IaUPOnyO8H5O/I/4RmTPaqbhmZ2gn8Ekw3/TO79tBB3bQWcfaSkK9\nb+4+ryk2fCe3Um6Q/NCeSRwYe3Z8Xe5ByqJfjGmrXLyU//folGAtnx4uaAeJ98ze\njUXT17x8AbdEt2JIpYoJI7xFC8mOr0s3LvA/gFmpNkuRNbCNQF2v5Qt9L2AYT0Fv\nB5uT42VuHQvRRNReAxa5oNGp/zcCjspaouPia03Tf5ZNZEd5LUFANLHPtsJg4jBB\nkjHKjCnt0/9fttE1u/gMW7k=\n-----END CERTIFICATE-----\n"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "client_auth_ca_certs.1", "-----BEGIN CERTIFICATE-----\nMIIDnTCCAoWgAwIBAgIUU5bIl5SJavP6YWPL/RUPLCbGu9owDQYJKoZIhvcNAQEL\nBQAwXjELMAkGA1UEBhMCRlIxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy\naXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MRgwFgYDVQQDDA9mb28uZXhhbXBsZS5j\nb20wHhcNMjUwODEzMDYyODM5WhcNMzUwODExMDYyODM5WjBeMQswCQYDVQQGEwJG\nUjEPMA0GA1UECAwGRnJhbmNlMQ4wDAYDVQQHDAVQYXJpczEUMBIGA1UECgwLT2dv\nU2VjdXJpdHkxGDAWBgNVBAMMD2Zvby5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBANZ0zrEH22IMXp8tQ5PbwLFHmCQRc1/T1ge5z7ho\np6zdyFn5GEFrMv1ZOywPBPlCz+Lb/5sWWj9qhcMw6JkPogKKVx9PQZDwfpc9ov+M\nmujh/SM1Ms07AFt286h9e0yZzQfP9t6B9+Dns4Lgn6/+Ua8g7VW+Hrq3V9Ait0bx\nkDOZUj0djOp9H3tShtgl8p9Z+dcqYIAPtkjSTt/U7jUDtR9PH6qz4/gXE/mCKq4e\nLf+63nLqGfZ3S1mIwjysRhPsJwy4g9v+E6fHO4Emfk4KF6EvFj3GVXyckxLbxKNa\n1yYRUhSLZvNCNyDVvySVUta7yOhdzyC53YvS/Emtrh/7I6kCAwEAAaNTMFEwHQYD\nVR0OBBYEFF3MC9L6J3lESlcdryXXFzC1uJGvMB8GA1UdIwQYMBaAFF3MC9L6J3lE\nSlcdryXXFzC1uJGvMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB\nABUa9KqjfCB5lf9G7rpVTqhrg5LIKzuYzH4c7MZ76R4GZyH475yV5jQYCj2Qr3Pq\n2m50UpEzVKICjfBSCbJulv5ZofSn8DWTpEBoLZA2pVM9yutI3wOQW350HX01nY82\n9j9im1yMVtdf1uAPd1O3pm+RUcSICI5YBFQ1/LAAEoSqrmSoVUPwH6pt9Gr+4E/w\nPpUcdAju8piy48Nx9ZD9vwCVjD67oRNnF00wEDJgrl8RpI9i0zOzflBxXyllGD8L\nXT6wvPmUpso+jn04qnizfMWaYy9P2ip8RgOslrH6WIe6GyXSy6VjAu9JSuVE5OYX\nXpYou5FLSGMhNaPTuaukAgY=\n-----END CERTIFICATE-----\n"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ogo_shield_tlsoptions.test", "uid"),
					resource.TestCheckResourceAttrSet("ogo_shield_tlsoptions.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "ogo_shield_tlsoptions.test",
				ImportStateIdFunc:                    testAccImportStateIdFromAttribute("ogo_shield_tlsoptions.test", "uid"),
				ImportStateVerifyIdentifierAttribute: "uid",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "ogo_shield_tlsoptions" "test" {
  name            = "mTLS foo bar"
  min_tls_version = "TLS_11"
  max_tls_version = "TLS_13"
  client_auth_ca_certs = [
    <<-EOT
-----BEGIN CERTIFICATE-----
MIIDnTCCAoWgAwIBAgIUHvOpeMH+4Lk1ewQZKMwOygGBQe0wDQYJKoZIhvcNAQEL
BQAwXjELMAkGA1UEBhMCRlIxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy
aXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MRgwFgYDVQQDDA9iYXIuZXhhbXBsZS5j
b20wHhcNMjUwODEzMDYyODU1WhcNMzUwODExMDYyODU1WjBeMQswCQYDVQQGEwJG
UjEPMA0GA1UECAwGRnJhbmNlMQ4wDAYDVQQHDAVQYXJpczEUMBIGA1UECgwLT2dv
U2VjdXJpdHkxGDAWBgNVBAMMD2Jhci5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAO4dBU9DGbgBzjIYy/Qls0IglivSHyughVRa4nfZ
b2b3iGP1rEa+xNlnmOlgxp8ihjxF4yBz/DMGVDEDnErwITUOxEG4fJ5gdX7a5Iyd
OgYYyoh1RJKRkyWSGQoU4RmbVidTCyxq15j+yRBJDt3fll+Y9rlL+Ejl9QJCe+Zt
kSab7pBn9SmUzX8IeHyX1IpEMA4nNtFI8ysNSZNxPJa1hB3tXtVGZrkhpecCZvx4
IBpuRrjBSY3MaRE5YW51l7nC7jExC+IeNGe3mfKYUu0Re7fkK7n1auGmAJhTlzIR
4126rTJDbZlKyDFSfoaDFsyYeNe2t2W6KlhG4d0dSiFwIucCAwEAAaNTMFEwHQYD
VR0OBBYEFKBppFca57l7wutRyaIRZ3fwzZP7MB8GA1UdIwQYMBaAFKBppFca57l7
wutRyaIRZ3fwzZP7MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB
AGwIPOoZqd/3Uu3W2dUcd8HWw/VmjokjKrC811KUfhiijpFpQjGMcGQrjti3rIkk
5ZQyItkw91/IaUPOnyO8H5O/I/4RmTPaqbhmZ2gn8Ekw3/TO79tBB3bQWcfaSkK9
b+4+ryk2fCe3Um6Q/NCeSRwYe3Z8Xe5ByqJfjGmrXLyU//folGAtnx4uaAeJ98ze
jUXT17x8AbdEt2JIpYoJI7xFC8mOr0s3LvA/gFmpNkuRNbCNQF2v5Qt9L2AYT0Fv
B5uT42VuHQvRRNReAxa5oNGp/zcCjspaouPia03Tf5ZNZEd5LUFANLHPtsJg4jBB
kjHKjCnt0/9fttE1u/gMW7k=
-----END CERTIFICATE-----
EOT
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first tlsoptions
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "name", "mTLS foo bar"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "min_tls_version", "TLS_11"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "max_tls_version", "TLS_13"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "client_auth_ca_certs.#", "1"),
					resource.TestCheckResourceAttr("ogo_shield_tlsoptions.test", "client_auth_ca_certs.0", "-----BEGIN CERTIFICATE-----\nMIIDnTCCAoWgAwIBAgIUHvOpeMH+4Lk1ewQZKMwOygGBQe0wDQYJKoZIhvcNAQEL\nBQAwXjELMAkGA1UEBhMCRlIxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy\naXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MRgwFgYDVQQDDA9iYXIuZXhhbXBsZS5j\nb20wHhcNMjUwODEzMDYyODU1WhcNMzUwODExMDYyODU1WjBeMQswCQYDVQQGEwJG\nUjEPMA0GA1UECAwGRnJhbmNlMQ4wDAYDVQQHDAVQYXJpczEUMBIGA1UECgwLT2dv\nU2VjdXJpdHkxGDAWBgNVBAMMD2Jhci5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN\nAQEBBQADggEPADCCAQoCggEBAO4dBU9DGbgBzjIYy/Qls0IglivSHyughVRa4nfZ\nb2b3iGP1rEa+xNlnmOlgxp8ihjxF4yBz/DMGVDEDnErwITUOxEG4fJ5gdX7a5Iyd\nOgYYyoh1RJKRkyWSGQoU4RmbVidTCyxq15j+yRBJDt3fll+Y9rlL+Ejl9QJCe+Zt\nkSab7pBn9SmUzX8IeHyX1IpEMA4nNtFI8ysNSZNxPJa1hB3tXtVGZrkhpecCZvx4\nIBpuRrjBSY3MaRE5YW51l7nC7jExC+IeNGe3mfKYUu0Re7fkK7n1auGmAJhTlzIR\n4126rTJDbZlKyDFSfoaDFsyYeNe2t2W6KlhG4d0dSiFwIucCAwEAAaNTMFEwHQYD\nVR0OBBYEFKBppFca57l7wutRyaIRZ3fwzZP7MB8GA1UdIwQYMBaAFKBppFca57l7\nwutRyaIRZ3fwzZP7MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB\nAGwIPOoZqd/3Uu3W2dUcd8HWw/VmjokjKrC811KUfhiijpFpQjGMcGQrjti3rIkk\n5ZQyItkw91/IaUPOnyO8H5O/I/4RmTPaqbhmZ2gn8Ekw3/TO79tBB3bQWcfaSkK9\nb+4+ryk2fCe3Um6Q/NCeSRwYe3Z8Xe5ByqJfjGmrXLyU//folGAtnx4uaAeJ98ze\njUXT17x8AbdEt2JIpYoJI7xFC8mOr0s3LvA/gFmpNkuRNbCNQF2v5Qt9L2AYT0Fv\nB5uT42VuHQvRRNReAxa5oNGp/zcCjspaouPia03Tf5ZNZEd5LUFANLHPtsJg4jBB\nkjHKjCnt0/9fttE1u/gMW7k=\n-----END CERTIFICATE-----\n"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

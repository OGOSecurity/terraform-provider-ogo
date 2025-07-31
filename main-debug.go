package main

import (
	"encoding/json"
	"fmt"
	"log"
	ogosecurity "terraform-provider-ogo/internal/ogo"
)

func printo(o any) {
	b, _ := json.MarshalIndent(o, "", "  ")
	fmt.Printf(string(b) + "\n")
}

func main() {
	host := "https://api-stg.ogosecurity.com"
	username := "sebastien.giraud@ogosecurity.com"
	//username := "ogo00795"
	apikey := "fAlJLNrjJk4i8dfBCESEXaRgjCgApcLg"
	sitename := "gys-tf.ogosecurity.com"

	client, err := ogosecurity.NewClient(&host, &username, &apikey)
	if err != nil {
		log.Fatal("Failed to initialize OgoSecurity Dashboard Client Error: " + err.Error())
	}

	fmt.Printf("=> Client\n")
	printo(client)

	//	fmt.Printf("=> Get all clusters:\n")
	//	clusters, err := client.GetAllClusters()
	//	if err != nil {
	//		log.Fatal("Failed to get all clusters: " + string(err.Error()))
	//	}
	//	printo(clusters)
	//
	//	fmt.Printf("=> Get cluster:\n")
	//	cluster, err := client.GetCluster("OGO GYS")
	//	if err != nil {
	//		log.Fatal("Failed to get cluster: " + string(err.Error()))
	//	}
	//	printo(cluster)
	//
	//	fmt.Printf("=> Get all TLS Options:\n")
	//	tlsOptions, err := client.GetAllTlsOptions()
	//	if err != nil {
	//		log.Fatal("Failed to get all TLS Options: " + string(err.Error()))
	//	}
	//	printo(tlsOptions)
	//
	//	t := ogosecurity.TlsOptions{
	//		Name:              "tf-test",
	//		ClientAuthType:    "VerifyClientCertIfGiven",
	//		ClientAuthCaCerts: []string{"-----BEGIN CERTIFICATE-----\nMIIFqzCCA5OgAwIBAgIUS2cmTFXmpoo1ZFDOAugp9TDyN4IwDQYJKoZIhvcNAQEL\nBQAwZTELMAkGA1UEBhMCQ0kxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFy\naXMxFDASBgNVBAoMC09nb1NlY3VyaXR5MR8wHQYDVQQDDBZzYW11ZWwub2dvc2Vj\ndXJpdHkuY29tMB4XDTIyMTEwMjE3MTY1NVoXDTMyMTAzMDE3MTY1NVowZTELMAkG\nA1UEBhMCQ0kxDzANBgNVBAgMBkZyYW5jZTEOMAwGA1UEBwwFUGFyaXMxFDASBgNV\nBAoMC09nb1NlY3VyaXR5MR8wHQYDVQQDDBZzYW11ZWwub2dvc2VjdXJpdHkuY29t\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEArWJiXEJG+ec3tS3ooeSR\nDaMZqonAqe3qdxPhaez+QioSXNL5diD/gvWZtl46Oxn7uKCkcdxQSSXPzy6bk3PK\n9j9OSN3Vlqku+vXCrpBDXW7psLusWmG+wAbXiSuB8CfB8LLioeQENShvHjaXYVDN\nymm7wQjNvK8uCFwnF6dRSCEVqkjiyFBXl7LALuK+dsjHoqJP8jeF04Z8MPehMslU\nLTtCvWYBycGfo5NSUntZ9mYakH35sIiga2YDHYDfGg4w/JtLqUvRng3Eo2vf8FDt\ngnDJR+p1zklulrAexGeVXYs3af1mGKSzLpkDyzne7fhWwEAd/PhAyppyj6NgDXwE\nWldACHwkYenI8Rs4fGIucJVIlNF2rrrhvb7ZjAFmqQnTghcOKR46VyoL1E+LpnMO\nPUkUdDVDOlBHRLTbBZo25HA6XLUUyxog+EVk0Me7VGAVNXKKa9ZICOY1UapsfmtW\nmsZaOHWdfS/avaWkhptAJvPv/oY7YQflPQN6labn6TSCKWpnW3PjiVWYky39BMO5\n8vdkRbmfdbM2SIl4maSOEI5e+5HnCzpwahtlfcd60Glnn56UV+4uIU0KMwCuBtvK\nxDRxKIIQu8Fz4tyqqE4xlTZYWeYFvmZ+Tt0bJyuG6/SZ5NPhOWdSQCkCI29l+Vxk\nqYBxafhxhR7sYLzLwsIoHLMCAwEAAaNTMFEwHQYDVR0OBBYEFHtt8ZYuEVDQUNGi\natjOpAuTCAY3MB8GA1UdIwQYMBaAFHtt8ZYuEVDQUNGiatjOpAuTCAY3MA8GA1Ud\nEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBAHq4qlfVHPpYL9xgd1oe/3yh\nLdzDnOerhOop+R+w7Y/Dxv9pz61aK5thBaroR6jwL0Sb5iTvltm1YtzHf9EHxf2V\nEfGXrSfta8tLXxZ1OApw2sTLPcKJL8te94ZgVZoEy9xCLgaovPtafiJgH+VKUj7r\n6v1vD/KiE0EkW36jTGO29bZaNuIVnNm67t9LqTH25Qdf4G3BPk1mfC8OPuO89MuG\n9/0H3R1g65kFOCwmT99eaNxzphn2S2RLNKto1dTA1Pb+BkIMxwC5L2yeHFXO2iME\nmKTAxxlhN19F4nsnVUc8yCN4A06bBhYHrYZOF2bqOfjeX5deCZMdDE98t42txOwV\nsYSMkR90jCBNl2tiH2GWUk5Mb8dZeNXjFPnEAKGcPkf/l/PNMuR187tb4UBDK2+e\nOKcurFt1SfF7kUUnFTtTtpc/imQGoNaoyTCvnm4Ivz8A/YOgJUm1sx+7DRh4G9Az\npRozjW8kl9/54HnVaDA/tVIfwVhtwv3WNT0F5RZgxqnC0F5Rux2yf1WDvHz7vL/V\nzBSGJFGd9CQ6qDAfrloUdeN8KWiqU2RiQ3qmyrubv+rZEEGj66PKqnMR3OqGqAHs\nWv166DxaUFUJtzKkx9e0j+VS9ABIp/w83tQskc6/jCzsp1+BqOJorhmWHSFklW8Q\nsWCOJOL6eRhzRNg5n3JL\n-----END CERTIFICATE-----"},
	//		MinTlsVersion:     "TLS_10",
	//		MaxTlsVersion:     "TLS_13",
	//	}
	//
	//	fmt.Printf("=> Create TLS Options:\n")
	//	tlsOpt, err := client.CreateTlsOptions(t)
	//	if err != nil {
	//		//log.Fatal("Failed to create TLS Options: " + string(err.Error()))
	//		fmt.Println("Failed to create TLS Options: " + string(err.Error()))
	//		t.Uid = "ogo00795-77c0456b-5b2a-4b29-9c64-aaec802f8b2f"
	//	} else {
	//		t.Uid = tlsOpt.Uid
	//	}
	//
	//	fmt.Printf("=> Get TLS Options:\n")
	//	tlsOpt, err = client.GetTlsOptions(t.Uid)
	//	if err != nil {
	//		log.Fatal("Failed to get TLS Options: " + string(err.Error()))
	//	}
	//	printo(tlsOpt)
	//
	//	fmt.Printf("=> Update TLS Options:\n")
	//	t.MinTlsVersion = "TLS_12"
	//	tlsOpt, err = client.UpdateTlsOptions(t)
	//	if err != nil {
	//		log.Fatal("Failed to update TLS Options: " + string(err.Error()))
	//	}
	//	printo(tlsOpt)
	//
	//	fmt.Printf("=> Delete TLS Options:\n")
	//	err = client.DeleteTlsOptions(t.Uid)
	//	if err != nil {
	//		log.Fatal("Failed to delete TLS Options: " + string(err.Error()))
	//	}
	//
	//	fmt.Printf("=> Get all sites:\n")
	//	sites, err := client.GetAllSites()
	//	if err != nil {
	//		log.Fatal("Failed to get all sites: " + string(err.Error()))
	//	}
	//	printo(sites)

	fmt.Printf("\n=> Get site " + sitename + ":\n")
	site, err := client.GetSite("gys-webapp.ogosecurity.com")
	if err != nil {
		log.Fatal("Failed to get site: " + string(err.Error()))
	}
	printo(site)

	return

	fmt.Printf("\n=> Create site " + sitename + ":\n")
	s := ogosecurity.Site{
		Name:             sitename,
		ClusterName:      "OGO GYS",
		DestHost:         "192.168.122.13",
		DestHostScheme:   "http",
		TrustSelfSigned:  true,
		NoCopyXForwarded: false,
		ForceHttps:       true,
		DryRun:           false,
		PanicMode:        false,
	}
	site, err = client.CreateSite(s)
	if err != nil {
		log.Fatal("Failed to create site: " + string(err.Error()))
	}
	printo(site)

	fmt.Printf("\n=> Get site " + sitename + ":\n")
	site, err = client.GetSite(sitename)
	if err != nil {
		log.Fatal("Failed to get site: " + string(err.Error()))
	}
	printo(site)

	fmt.Printf("\n=> Update site " + sitename + ":\n")
	s.DryRun = true
	site, err = client.UpdateSite(s)
	if err != nil {
		log.Fatal("Failed to update site: " + string(err.Error()))
	}
	printo(site)

	fmt.Printf("\n=> Get site " + sitename + ":\n")
	site, err = client.GetSite(sitename)
	if err != nil {
		log.Fatal("Failed to get site: " + string(err.Error()))
	}
	printo(site)

	fmt.Printf("\n=> Delete site " + sitename + ": \n")
	err = client.DeleteSite(sitename)
	if err != nil {
		log.Fatal("Failed to delete site: " + string(err.Error()))
	}
}

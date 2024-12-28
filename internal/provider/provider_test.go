package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "slash" {
 access_token = "foobarbaz"
 host         = "http://localhost:5231"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"slash": providerserver.NewProtocol6WithError(New("test")()),
	}
)

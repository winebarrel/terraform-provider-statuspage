package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"statuspage": providerserver.NewProtocol6WithError(New("test")()),
}

var testAccOrganizationID = "test-org-id"

func init() {
	os.Setenv("STATUSPAGE_API_KEY", "test-api-key")
	apiclient.RateLimitInterval = 0
}

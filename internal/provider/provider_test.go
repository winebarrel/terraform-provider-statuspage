package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"statuspage": providerserver.NewProtocol6WithError(New("test", apiclient.WithRateLimitInterval(0))()),
}

func init() {
	_ = os.Setenv("STATUSPAGE_API_KEY", "test-api-key")
}

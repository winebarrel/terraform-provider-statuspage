package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"statuspage": providerserver.NewProtocol6WithError(New("test")()),
}

var testAccPageID string
var testAccOrganizationID string

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if v := os.Getenv("STATUSPAGE_API_KEY"); v == "" {
		t.Fatal("STATUSPAGE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("STATUSPAGE_PAGE_ID"); v == "" {
		t.Fatal("STATUSPAGE_PAGE_ID must be set for acceptance tests")
	}
	testAccPageID = os.Getenv("STATUSPAGE_PAGE_ID")
	testAccOrganizationID = os.Getenv("STATUSPAGE_ORGANIZATION_ID")
}

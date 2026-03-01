package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestPage_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	page := &apiclient.Page{
		ID:        "test-page-id",
		Name:      "existing-page",
		Subdomain: "test-subdomain",
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, page)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PageRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			if body.Page.Name != "" {
				page.Name = body.Page.Name
			}
			return httpmock.NewJsonResponse(200, page)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import existing page
			{
				Config:             testPageConfig("tf-test-page"),
				ResourceName:       "statuspage_page.test",
				ImportState:        true,
				ImportStateId:      "test-page-id",
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			// Update page name
			{
				Config: testPageConfig("tf-test-page-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page.test", "name", "tf-test-page-updated"),
					resource.TestCheckResourceAttr("statuspage_page.test", "id", "test-page-id"),
					resource.TestCheckResourceAttrSet("statuspage_page.test", "subdomain"),
				),
			},
		},
	})
}

func testPageConfig(name string) string {
	return fmt.Sprintf(`
resource "statuspage_page" "test" {
  id   = %q
  name = %q
}
`, "test-page-id", name)
}

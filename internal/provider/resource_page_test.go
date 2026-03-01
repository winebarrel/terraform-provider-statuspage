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

func TestAccPage_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	page := &apiclient.Page{
		ID:        testAccPageID,
		Name:      "existing-page",
		Subdomain: "test-subdomain",
	}

	httpmock.RegisterResponder("GET", testBaseURL+"/pages/"+testAccPageID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, page)
		})

	httpmock.RegisterResponder("PATCH", testBaseURL+"/pages/"+testAccPageID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PageRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.Page.Name != "" {
				page.Name = body.Page.Name
			}
			return httpmock.NewJsonResponse(200, page)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import existing page
			{
				Config:             testAccPageConfig("tf-test-page"),
				ResourceName:       "statuspage_page.test",
				ImportState:        true,
				ImportStateId:      testAccPageID,
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			// Update page name
			{
				Config: testAccPageConfig("tf-test-page-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page.test", "name", "tf-test-page-updated"),
					resource.TestCheckResourceAttr("statuspage_page.test", "id", testAccPageID),
					resource.TestCheckResourceAttrSet("statuspage_page.test", "subdomain"),
				),
			},
		},
	})
}

func testAccPageConfig(name string) string {
	return fmt.Sprintf(`
resource "statuspage_page" "test" {
  id   = %q
  name = %q
}
`, testAccPageID, name)
}

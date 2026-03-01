package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestAccPageAccessGroup_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	groupID := "page-access-group-id-1"
	group := &apiclient.PageAccessGroup{
		ID:     groupID,
		PageID: "test-page-id",
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/page_access_groups",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PageAccessGroupRequest
			json.NewDecoder(req.Body).Decode(&body)
			group.Name = body.PageAccessGroup.Name
			return httpmock.NewJsonResponse(201, group)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/page_access_groups/"+groupID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, group)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/page_access_groups/"+groupID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PageAccessGroupRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.PageAccessGroup.Name != "" {
				group.Name = body.PageAccessGroup.Name
			}
			return httpmock.NewJsonResponse(200, group)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/page_access_groups/"+groupID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccPageAccessGroupConfig("tf-test-pag"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_group.test", "name", "tf-test-pag"),
					resource.TestCheckResourceAttrSet("statuspage_page_access_group.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_page_access_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncPageAccessGroup("statuspage_page_access_group.test"),
			},
			// Update
			{
				Config: testAccPageAccessGroupConfig("tf-test-pag-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_group.test", "name", "tf-test-pag-updated"),
				),
			},
		},
	})
}

func testAccPageAccessGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "statuspage_page_access_group" "test" {
  page_id = %q
  name    = %q
}
`, "test-page-id", name)
}

func importStateIDFuncPageAccessGroup(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		pageID := rs.Primary.Attributes["page_id"]
		id := rs.Primary.ID
		return fmt.Sprintf("%s/%s", pageID, id), nil
	}
}

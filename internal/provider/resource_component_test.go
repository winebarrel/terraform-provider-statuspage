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

func TestAccComponent_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	componentID := "component-id-1"
	component := &apiclient.Component{
		ID:              componentID,
		PageID:          testAccPageID,
		AutomationEmail: "component+test@notifications.statuspage.io",
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/components",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.ComponentRequest
			json.NewDecoder(req.Body).Decode(&body)
			component.Name = body.Component.Name
			component.Status = body.Component.Status
			return httpmock.NewJsonResponse(201, component)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/components/"+componentID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, component)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/components/"+componentID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.ComponentRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.Component.Name != "" {
				component.Name = body.Component.Name
			}
			if body.Component.Status != "" {
				component.Status = body.Component.Status
			}
			return httpmock.NewJsonResponse(200, component)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/components/"+componentID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentConfig("tf-test-component", "operational"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component.test", "name", "tf-test-component"),
					resource.TestCheckResourceAttr("statuspage_component.test", "status", "operational"),
					resource.TestCheckResourceAttrSet("statuspage_component.test", "id"),
					resource.TestCheckResourceAttrSet("statuspage_component.test", "automation_email"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_component.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFunc("statuspage_component.test"),
			},
			// Update
			{
				Config: testAccComponentConfig("tf-test-component-updated", "degraded_performance"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component.test", "name", "tf-test-component-updated"),
					resource.TestCheckResourceAttr("statuspage_component.test", "status", "degraded_performance"),
				),
			},
		},
	})
}

func testAccComponentConfig(name, status string) string {
	return fmt.Sprintf(`
resource "statuspage_component" "test" {
  page_id = %q
  name    = %q
  status  = %q
}
`, testAccPageID, name, status)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
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

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

func TestAccIncidentTemplate_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	templateID := "incident-template-id-1"
	tmpl := &apiclient.IncidentTemplate{
		ID:     templateID,
		PageID: "test-page-id",
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/incident_templates",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.IncidentTemplateRequest
			json.NewDecoder(req.Body).Decode(&body)
			tmpl.Name = body.Template.Name
			tmpl.Title = body.Template.Title
			tmpl.Body = body.Template.Body
			tmpl.UpdateStatus = body.Template.UpdateStatus
			return httpmock.NewJsonResponse(201, tmpl)
		})

	// GET uses list endpoint - returns array
	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/incident_templates",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []apiclient.IncidentTemplate{*tmpl})
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/incident_templates/"+templateID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.IncidentTemplateRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.Template.Name != "" {
				tmpl.Name = body.Template.Name
			}
			if body.Template.Title != "" {
				tmpl.Title = body.Template.Title
			}
			if body.Template.Body != "" {
				tmpl.Body = body.Template.Body
			}
			if body.Template.UpdateStatus != "" {
				tmpl.UpdateStatus = body.Template.UpdateStatus
			}
			return httpmock.NewJsonResponse(200, tmpl)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/incident_templates/"+templateID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccIncidentTemplateConfig("tf-test-template", "Test Incident", "investigating"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "name", "tf-test-template"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "title", "Test Incident"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "update_status", "investigating"),
					resource.TestCheckResourceAttrSet("statuspage_incident_template.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_incident_template.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncIncidentTemplate("statuspage_incident_template.test"),
			},
			// Update
			{
				Config: testAccIncidentTemplateConfig("tf-test-template-updated", "Updated Incident", "identified"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "name", "tf-test-template-updated"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "title", "Updated Incident"),
					resource.TestCheckResourceAttr("statuspage_incident_template.test", "update_status", "identified"),
				),
			},
		},
	})
}

func testAccIncidentTemplateConfig(name, title, updateStatus string) string {
	return fmt.Sprintf(`
resource "statuspage_incident_template" "test" {
  page_id       = %q
  name          = %q
  title         = %q
  body          = "This is a template body."
  update_status = %q
}
`, "test-page-id", name, title, updateStatus)
}

func importStateIDFuncIncidentTemplate(resourceName string) resource.ImportStateIdFunc {
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

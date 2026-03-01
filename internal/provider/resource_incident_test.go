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

func TestIncident_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	incident := &apiclient.Incident{
		ID:                   "incident-id-1",
		PageID:               "test-page-id",
		Shortlink:            "https://stspg.io/test123",
		DeliverNotifications: true,
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/incidents",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.IncidentRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			incident.Name = body.Incident.Name
			incident.Status = body.Incident.Status
			incident.Body = body.Incident.Body
			return httpmock.NewJsonResponse(201, incident)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/incidents/incident-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, incident)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/incidents/incident-id-1",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.IncidentRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			if body.Incident.Name != "" {
				incident.Name = body.Incident.Name
			}
			if body.Incident.Status != "" {
				incident.Status = body.Incident.Status
			}
			if body.Incident.Body != "" {
				incident.Body = body.Incident.Body
			}
			return httpmock.NewJsonResponse(200, incident)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/incidents/incident-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testIncidentConfig("tf-test-incident", "investigating", "Initial investigation"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident.test", "name", "tf-test-incident"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "status", "investigating"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "body", "Initial investigation"),
					resource.TestCheckResourceAttrSet("statuspage_incident.test", "id"),
					resource.TestCheckResourceAttrSet("statuspage_incident.test", "shortlink"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_incident.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncIncident("statuspage_incident.test"),
			},
			// Update
			{
				Config: testIncidentConfig("tf-test-incident-updated", "identified", "Issue has been identified"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident.test", "name", "tf-test-incident-updated"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "status", "identified"),
					resource.TestCheckResourceAttr("statuspage_incident.test", "body", "Issue has been identified"),
				),
			},
		},
	})
}

func testIncidentConfig(name, status, body string) string {
	return fmt.Sprintf(`
resource "statuspage_incident" "test" {
  page_id = %q
  name    = %q
  status  = %q
  body    = %q
}
`, "test-page-id", name, status, body)
}

func importStateIDFuncIncident(resourceName string) resource.ImportStateIdFunc {
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

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

func TestAccIncidentPostmortem_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	incident := &apiclient.Incident{
		ID:                   "postmortem-incident-id-1",
		PageID:               "test-page-id",
		Name:                 "tf-test-postmortem-incident",
		Status:               "resolved",
		Body:                 "Incident for postmortem testing.",
		DeliverNotifications: true,
	}

	postmortem := &apiclient.Postmortem{}

	// Incident responders
	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/incidents",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.IncidentRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			incident.Name = body.Incident.Name
			incident.Status = body.Incident.Status
			incident.Body = body.Incident.Body
			return httpmock.NewJsonResponse(201, incident)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/incidents/postmortem-incident-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, incident)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/incidents/postmortem-incident-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Postmortem responders
	httpmock.RegisterResponder("PUT", "https://api.statuspage.io/v1/pages/test-page-id/incidents/postmortem-incident-id-1/postmortem",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PostmortemRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			if body.Postmortem.Body != "" {
				postmortem.Body = body.Postmortem.Body
			}
			return httpmock.NewJsonResponse(200, postmortem)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/incidents/postmortem-incident-id-1/postmortem",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, postmortem)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/incidents/postmortem-incident-id-1/postmortem",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read (creates an incident first, then attaches postmortem)
			{
				Config: testAccIncidentPostmortemConfig("Initial postmortem body"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_postmortem.test", "body", "Initial postmortem body"),
					resource.TestCheckResourceAttrSet("statuspage_incident_postmortem.test", "incident_id"),
				),
			},
			// Import
			{
				ResourceName:                        "statuspage_incident_postmortem.test",
				ImportState:                         true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "incident_id",
				ImportStateIdFunc:                    importStateIDFuncPostmortem("statuspage_incident_postmortem.test"),
			},
			// Update
			{
				Config: testAccIncidentPostmortemConfig("Updated postmortem body"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_incident_postmortem.test", "body", "Updated postmortem body"),
				),
			},
		},
	})
}

func testAccIncidentPostmortemConfig(body string) string {
	return fmt.Sprintf(`
resource "statuspage_incident" "postmortem_incident" {
  page_id = %[1]q
  name    = "tf-test-postmortem-incident"
  status  = "resolved"
  body    = "Incident for postmortem testing."
}

resource "statuspage_incident_postmortem" "test" {
  page_id            = %[1]q
  incident_id        = statuspage_incident.postmortem_incident.id
  body               = %[2]q
  notify_subscribers = false
  notify_twitter     = false
}
`, "test-page-id", body)
}

func importStateIDFuncPostmortem(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		pageID := rs.Primary.Attributes["page_id"]
		incidentID := rs.Primary.Attributes["incident_id"]
		return fmt.Sprintf("%s/%s", pageID, incidentID), nil
	}
}

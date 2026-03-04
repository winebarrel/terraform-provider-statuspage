package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestIncidentTemplateDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tmpl := &apiclient.IncidentTemplate{
		ID:                      "template-id-1",
		PageID:                  "test-page-id",
		Name:                    "Test Template",
		Title:                   "Incident Title",
		Body:                    "Incident body text",
		GroupID:                 "group-1",
		UpdateStatus:            "investigating",
		ShouldTweet:             false,
		ShouldSendNotifications: true,
		ComponentIDs:            []string{"comp-1"},
	}

	// GetIncidentTemplate lists all templates and filters by ID
	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/incident_templates",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []apiclient.IncidentTemplate{*tmpl})
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_incident_template" "test" {
  page_id = "test-page-id"
  id      = "template-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "id", "template-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "name", "Test Template"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "title", "Incident Title"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "body", "Incident body text"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "update_status", "investigating"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "should_send_notifications", "true"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "component_ids.#", "1"),
					resource.TestCheckResourceAttr("data.statuspage_incident_template.test", "component_ids.0", "comp-1"),
				),
			},
		},
	})
}

package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestIncidentDataSource_withComponents(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	incident := &apiclient.Incident{
		ID:                   "incident-id-1",
		PageID:               "test-page-id",
		Name:                 "Test Incident",
		Status:               "investigating",
		Shortlink:            "https://stspg.io/test",
		DeliverNotifications: true,
		Components: []apiclient.Component{
			{ID: "comp-1", PageID: "test-page-id", Name: "API", Status: "major_outage"},
			{ID: "comp-2", PageID: "test-page-id", Name: "Web App", Status: "degraded_performance"},
		},
		ComponentIDs: []string{"comp-1", "comp-2"},
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/incidents/incident-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, incident)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_incident" "test" {
  page_id = "test-page-id"
  id      = "incident-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "components.#", "2"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "components.0.id", "comp-1"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "components.0.name", "API"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "components.0.status", "major_outage"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "components.1.id", "comp-2"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "components.1.name", "Web App"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "components.1.status", "degraded_performance"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "component_ids.#", "2"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "component_ids.0", "comp-1"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "component_ids.1", "comp-2"),
				),
			},
		},
	})
}

func TestIncidentDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	incident := &apiclient.Incident{
		ID:                   "incident-id-1",
		PageID:               "test-page-id",
		Name:                 "Test Incident",
		Status:               "investigating",
		ImpactOverride:       "minor",
		Body:                 "We are investigating.",
		Shortlink:            "https://stspg.io/test",
		DeliverNotifications: true,
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/incidents/incident-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, incident)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_incident" "test" {
  page_id = "test-page-id"
  id      = "incident-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "id", "incident-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "name", "Test Incident"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "status", "investigating"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "impact_override", "minor"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "body", "We are investigating."),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "shortlink", "https://stspg.io/test"),
					resource.TestCheckResourceAttr("data.statuspage_incident.test", "deliver_notifications", "true"),
				),
			},
		},
	})
}

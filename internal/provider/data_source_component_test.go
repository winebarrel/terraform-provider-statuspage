package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestComponentDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	component := &apiclient.Component{
		ID:                 "component-id-1",
		PageID:             "test-page-id",
		Name:               "My Component",
		Description:        "A test component",
		Position:           1,
		Status:             "operational",
		Showcase:           true,
		OnlyShowIfDegraded: false,
		AutomationEmail:    "component+test@notifications.statuspage.io",
		StartDate:          "2024-01-01",
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/components/component-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, component)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_component" "test" {
  page_id = "test-page-id"
  id      = "component-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_component.test", "id", "component-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_component.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_component.test", "name", "My Component"),
					resource.TestCheckResourceAttr("data.statuspage_component.test", "description", "A test component"),
					resource.TestCheckResourceAttr("data.statuspage_component.test", "status", "operational"),
					resource.TestCheckResourceAttr("data.statuspage_component.test", "showcase", "true"),
					resource.TestCheckResourceAttr("data.statuspage_component.test", "automation_email", "component+test@notifications.statuspage.io"),
					resource.TestCheckResourceAttr("data.statuspage_component.test", "start_date", "2024-01-01"),
				),
			},
		},
	})
}

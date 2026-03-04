package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestComponentGroupDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	group := &apiclient.ComponentGroup{
		ID:          "group-id-1",
		PageID:      "test-page-id",
		Name:        "My Group",
		Description: "A test group",
		Components:  []string{"comp-1", "comp-2"},
		Position:    1,
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/component-groups/group-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, group)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_component_group" "test" {
  page_id = "test-page-id"
  id      = "group-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "id", "group-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "name", "My Group"),
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "description", "A test group"),
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "components.#", "2"),
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "components.0", "comp-1"),
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "components.1", "comp-2"),
					resource.TestCheckResourceAttr("data.statuspage_component_group.test", "position", "1"),
				),
			},
		},
	})
}

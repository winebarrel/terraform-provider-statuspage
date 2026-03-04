package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestPageAccessGroupDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	group := &apiclient.PageAccessGroup{
		ID:                 "group-id-1",
		PageID:             "test-page-id",
		Name:               "Test Group",
		ExternalIdentifier: "ext-group-1",
		ComponentIDs:       []string{"comp-1", "comp-2"},
		MetricIDs:          []string{"metric-1"},
		PageAccessUserIDs:  []string{"user-1", "user-2"},
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/page_access_groups/group-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, group)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_page_access_group" "test" {
  page_id = "test-page-id"
  id      = "group-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "id", "group-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "name", "Test Group"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "external_identifier", "ext-group-1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "component_ids.#", "2"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "component_ids.0", "comp-1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "component_ids.1", "comp-2"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "metric_ids.#", "1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "metric_ids.0", "metric-1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "page_access_user_ids.#", "2"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "page_access_user_ids.0", "user-1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_group.test", "page_access_user_ids.1", "user-2"),
				),
			},
		},
	})
}

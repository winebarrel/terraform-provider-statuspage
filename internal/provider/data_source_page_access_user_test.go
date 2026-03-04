package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestPageAccessUserDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	user := &apiclient.PageAccessUser{
		ID:            "user-id-1",
		PageID:        "test-page-id",
		ExternalLogin: "testuser",
		ExternalEmail: "testuser@example.com",
		ComponentIDs:  []string{"comp-1"},
		MetricIDs:     []string{"metric-1"},
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/page_access_users/user-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, user)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_page_access_user" "test" {
  page_id = "test-page-id"
  id      = "user-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "id", "user-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "external_login", "testuser"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "external_email", "testuser@example.com"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "component_ids.#", "1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "component_ids.0", "comp-1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "metric_ids.#", "1"),
					resource.TestCheckResourceAttr("data.statuspage_page_access_user.test", "metric_ids.0", "metric-1"),
				),
			},
		},
	})
}

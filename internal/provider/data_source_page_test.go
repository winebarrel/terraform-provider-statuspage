package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestPageDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	page := &apiclient.Page{
		ID:                   "test-page-id",
		Name:                 "My Status Page",
		Subdomain:            "mypage",
		Domain:               "status.example.com",
		URL:                  "https://status.example.com",
		Branding:             "premium",
		TimeZone:             "Asia/Tokyo",
		AllowPageSubscribers: true,
		AllowEmailSubscribers: true,
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, page)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_page" "test" {
  id = "test-page-id"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_page.test", "id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "name", "My Status Page"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "subdomain", "mypage"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "domain", "status.example.com"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "url", "https://status.example.com"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "branding", "premium"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "time_zone", "Asia/Tokyo"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "allow_page_subscribers", "true"),
					resource.TestCheckResourceAttr("data.statuspage_page.test", "allow_email_subscribers", "true"),
				),
			},
		},
	})
}

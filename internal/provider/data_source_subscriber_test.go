package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestSubscriberDataSource_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	subscriber := &apiclient.Subscriber{
		ID:           "subscriber-id-1",
		PageID:       "test-page-id",
		Email:        "subscriber@example.com",
		Mode:         "email",
		ComponentIDs: []string{"comp-1"},
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/subscribers/subscriber-id-1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, subscriber)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "statuspage_subscriber" "test" {
  page_id = "test-page-id"
  id      = "subscriber-id-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspage_subscriber.test", "id", "subscriber-id-1"),
					resource.TestCheckResourceAttr("data.statuspage_subscriber.test", "page_id", "test-page-id"),
					resource.TestCheckResourceAttr("data.statuspage_subscriber.test", "email", "subscriber@example.com"),
					resource.TestCheckResourceAttr("data.statuspage_subscriber.test", "mode", "email"),
					resource.TestCheckResourceAttr("data.statuspage_subscriber.test", "component_ids.#", "1"),
					resource.TestCheckResourceAttr("data.statuspage_subscriber.test", "component_ids.0", "comp-1"),
				),
			},
		},
	})
}

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

func TestAccSubscriber_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	subscriberID := "subscriber-id-1"
	subscriber := &apiclient.Subscriber{
		ID:     subscriberID,
		PageID: "test-page-id",
		Mode:   "email",
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/subscribers",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.SubscriberRequest
			json.NewDecoder(req.Body).Decode(&body)
			subscriber.Email = body.Subscriber.Email
			return httpmock.NewJsonResponse(201, subscriber)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/subscribers/"+subscriberID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, subscriber)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/subscribers/"+subscriberID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.SubscriberRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.Subscriber.Email != "" {
				subscriber.Email = body.Subscriber.Email
			}
			return httpmock.NewJsonResponse(200, subscriber)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/subscribers/"+subscriberID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read (email subscriber)
			{
				Config: testAccSubscriberConfig("tf-test-subscriber@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_subscriber.test", "email", "tf-test-subscriber@example.com"),
					resource.TestCheckResourceAttr("statuspage_subscriber.test", "mode", "email"),
					resource.TestCheckResourceAttrSet("statuspage_subscriber.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_subscriber.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncSubscriber("statuspage_subscriber.test"),
				// skip_confirmation_notification is write-only and not returned by the API
				ImportStateVerifyIgnore: []string{"skip_confirmation_notification"},
			},
			// Update email
			{
				Config: testAccSubscriberConfig("tf-test-subscriber-updated@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_subscriber.test", "email", "tf-test-subscriber-updated@example.com"),
				),
			},
		},
	})
}

func testAccSubscriberConfig(email string) string {
	return fmt.Sprintf(`
resource "statuspage_subscriber" "test" {
  page_id                      = %q
  email                        = %q
  skip_confirmation_notification = true
}
`, "test-page-id", email)
}

func importStateIDFuncSubscriber(resourceName string) resource.ImportStateIdFunc {
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

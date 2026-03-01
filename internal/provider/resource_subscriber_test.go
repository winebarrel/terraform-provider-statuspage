package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSubscriber_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
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
`, testAccPageID, email)
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

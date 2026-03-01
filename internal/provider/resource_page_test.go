package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Import existing page
			{
				Config:             testAccPageConfig("tf-test-page"),
				ResourceName:       "statuspage_page.test",
				ImportState:        true,
				ImportStateId:      testAccPageID,
				ImportStateVerify:  false, // Name may differ from what we set initially
				ImportStatePersist: true,
			},
			// Update page name
			{
				Config: testAccPageConfig("tf-test-page-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page.test", "name", "tf-test-page-updated"),
					resource.TestCheckResourceAttr("statuspage_page.test", "id", testAccPageID),
					resource.TestCheckResourceAttrSet("statuspage_page.test", "subdomain"),
				),
			},
		},
	})
}

func testAccPageConfig(name string) string {
	return fmt.Sprintf(`
resource "statuspage_page" "test" {
  id   = %q
  name = %q
}
`, testAccPageID, name)
}

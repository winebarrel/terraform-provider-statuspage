package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusEmbedConfig_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Import existing status embed config
			{
				Config:             testAccStatusEmbedConfigConfig("bl"),
				ResourceName:       "statuspage_status_embed_config.test",
				ImportState:        true,
				ImportStateId:      testAccPageID,
				ImportStateVerify:  false,
				ImportStatePersist: true,
			},
			// Update
			{
				Config: testAccStatusEmbedConfigConfig("br"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_status_embed_config.test", "position", "br"),
					resource.TestCheckResourceAttrSet("statuspage_status_embed_config.test", "incident_background_color"),
					resource.TestCheckResourceAttrSet("statuspage_status_embed_config.test", "incident_text_color"),
				),
			},
		},
	})
}

func testAccStatusEmbedConfigConfig(position string) string {
	return fmt.Sprintf(`
resource "statuspage_status_embed_config" "test" {
  page_id  = %q
  position = %q
}
`, testAccPageID, position)
}

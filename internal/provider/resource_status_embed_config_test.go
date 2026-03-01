package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestAccStatusEmbedConfig_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &apiclient.StatusEmbedConfig{
		PageID:                  "test-page-id",
		Position:                "bl",
		IncidentBackgroundColor: "#FF0000",
		IncidentTextColor:       "#FFFFFF",
	}

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/status_embed_config",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, config)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/status_embed_config",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.StatusEmbedConfigRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			if body.StatusEmbedConfig.Position != "" {
				config.Position = body.StatusEmbedConfig.Position
			}
			return httpmock.NewJsonResponse(200, config)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import existing status embed config
			{
				Config:             testAccStatusEmbedConfigConfig("bl"),
				ResourceName:       "statuspage_status_embed_config.test",
				ImportState:        true,
				ImportStateId:      "test-page-id",
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
`, "test-page-id", position)
}

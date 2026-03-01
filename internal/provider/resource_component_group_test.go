package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/winebarrel/terraform-provider-statuspage/internal/apiclient"
)

func TestAccComponentGroup_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var mu sync.Mutex
	compIDs := []string{"comp-in-group-1", "comp-in-group-2"}
	compIdx := 0
	components := map[string]*apiclient.Component{}

	groupID := "comp-group-id-1"
	group := &apiclient.ComponentGroup{
		ID:     groupID,
		PageID: "test-page-id",
	}

	// POST /pages/{pageID}/components
	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/components",
		func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			defer mu.Unlock()
			var body apiclient.ComponentRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			id := compIDs[compIdx]
			compIdx++
			comp := &apiclient.Component{
				ID:     id,
				PageID: "test-page-id",
				Name:   body.Component.Name,
				Status: body.Component.Status,
			}
			components[id] = comp
			return httpmock.NewJsonResponse(201, comp)
		})

	// GET/DELETE for each component
	for _, id := range compIDs {
		id := id
		httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/components/"+id,
			func(req *http.Request) (*http.Response, error) {
				mu.Lock()
				defer mu.Unlock()
				if comp, ok := components[id]; ok {
					return httpmock.NewJsonResponse(200, comp)
				}
				return httpmock.NewJsonResponse(404, map[string]string{"error": "not found"})
			})
		httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/components/"+id,
			func(req *http.Request) (*http.Response, error) {
				mu.Lock()
				defer mu.Unlock()
				var body apiclient.ComponentRequest
				_ = json.NewDecoder(req.Body).Decode(&body)
				if comp, ok := components[id]; ok {
					if body.Component.Name != "" {
						comp.Name = body.Component.Name
					}
					if body.Component.Status != "" {
						comp.Status = body.Component.Status
					}
					return httpmock.NewJsonResponse(200, comp)
				}
				return httpmock.NewJsonResponse(404, map[string]string{"error": "not found"})
			})
		httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/components/"+id,
			func(req *http.Request) (*http.Response, error) {
				mu.Lock()
				defer mu.Unlock()
				delete(components, id)
				return httpmock.NewStringResponse(204, ""), nil
			})
	}

	// POST /pages/{pageID}/component-groups
	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/test-page-id/component-groups",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.ComponentGroupRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			group.Name = body.ComponentGroup.Name
			group.Components = body.ComponentGroup.Components
			return httpmock.NewJsonResponse(201, group)
		})

	// GET /pages/{pageID}/component-groups/{id}
	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/test-page-id/component-groups/"+groupID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, group)
		})

	// PATCH /pages/{pageID}/component-groups/{id}
	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/test-page-id/component-groups/"+groupID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.ComponentGroupRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			if body.ComponentGroup.Name != "" {
				group.Name = body.ComponentGroup.Name
			}
			if body.ComponentGroup.Components != nil {
				group.Components = body.ComponentGroup.Components
			}
			return httpmock.NewJsonResponse(200, group)
		})

	// DELETE /pages/{pageID}/component-groups/{id}
	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/test-page-id/component-groups/"+groupID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccComponentGroupConfig("tf-test-group", "tf-test-group-comp1", "tf-test-group-comp2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component_group.test", "name", "tf-test-group"),
					resource.TestCheckResourceAttr("statuspage_component_group.test", "components.#", "2"),
					resource.TestCheckResourceAttrSet("statuspage_component_group.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_component_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncComponentGroup("statuspage_component_group.test"),
			},
			// Update
			{
				Config: testAccComponentGroupConfig("tf-test-group-updated", "tf-test-group-comp1", "tf-test-group-comp2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_component_group.test", "name", "tf-test-group-updated"),
				),
			},
		},
	})
}

func testAccComponentGroupConfig(groupName, comp1Name, comp2Name string) string {
	return fmt.Sprintf(`
resource "statuspage_component" "group_comp1" {
  page_id = %[1]q
  name    = %[2]q
  status  = "operational"
}

resource "statuspage_component" "group_comp2" {
  page_id = %[1]q
  name    = %[3]q
  status  = "operational"
}

resource "statuspage_component_group" "test" {
  page_id    = %[1]q
  name       = %[4]q
  components = [statuspage_component.group_comp1.id, statuspage_component.group_comp2.id]
}
`, "test-page-id", comp1Name, comp2Name, groupName)
}

func importStateIDFuncComponentGroup(resourceName string) resource.ImportStateIdFunc {
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

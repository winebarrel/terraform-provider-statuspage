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

func TestAccUserPermissions_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	userID := "user-perm-id-1"
	user := &apiclient.User{
		ID:             userID,
		OrganizationID: "test-org-id",
	}

	permissions := &apiclient.Permissions{
		UserID: userID,
		Pages:  map[string]string{},
	}

	// User responders
	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/organizations/test-org-id/users",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.UserRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			user.Email = body.User.Email
			user.FirstName = body.User.FirstName
			user.LastName = body.User.LastName
			return httpmock.NewJsonResponse(201, user)
		})

	// GET users uses list endpoint
	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/organizations/test-org-id/users",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []apiclient.User{*user})
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/organizations/test-org-id/users/"+userID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Permissions responders
	httpmock.RegisterResponder("PUT", "https://api.statuspage.io/v1/organizations/test-org-id/permissions/"+userID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PermissionsRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			pages := make(map[string]string)
			for k, v := range body.Pages {
				if len(v) > 0 {
					pages[k] = v[0]
				}
			}
			permissions.Pages = pages
			return httpmock.NewJsonResponse(200, permissions)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/organizations/test-org-id/permissions/"+userID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, permissions)
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read (creates a user first, then sets permissions)
			{
				Config: testAccUserPermissionsConfig("page_configuration"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_user_permissions.test", "user_id"),
					resource.TestCheckResourceAttrSet("statuspage_user_permissions.test", "organization_id"),
				),
			},
			// Import
			{
				ResourceName:                         "statuspage_user_permissions.test",
				ImportState:                          true,
				ImportStateVerify:                     true,
				ImportStateVerifyIdentifierAttribute:  "user_id",
				ImportStateIdFunc:                     importStateIDFuncUserPermissions("statuspage_user_permissions.test"),
			},
			// Update permissions
			{
				Config: testAccUserPermissionsConfig("incident_manager"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("statuspage_user_permissions.test", "user_id"),
				),
			},
		},
	})
}

func testAccUserPermissionsConfig(permission string) string {
	return fmt.Sprintf(`
resource "statuspage_user" "perm_user" {
  organization_id = %[1]q
  email           = "tf-test-perms@example.com"
  password        = "TestP@ssw0rd456!"
  first_name      = "PermTest"
  last_name       = "User"
}

resource "statuspage_user_permissions" "test" {
  organization_id = %[1]q
  user_id         = statuspage_user.perm_user.id
  pages = {
    %[2]q = %[3]q
  }
}
`, "test-org-id", "test-page-id", permission)
}

func importStateIDFuncUserPermissions(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["organization_id"]
		userID := rs.Primary.Attributes["user_id"]
		return fmt.Sprintf("%s/%s", orgID, userID), nil
	}
}

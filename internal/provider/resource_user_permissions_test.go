package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserPermissions_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck: func() {
			testAccPreCheck(t)
			if testAccOrganizationID == "" {
				t.Skip("STATUSPAGE_ORGANIZATION_ID must be set for user permissions acceptance tests")
			}
		},
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
				ResourceName:      "statuspage_user_permissions.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncUserPermissions("statuspage_user_permissions.test"),
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
`, testAccOrganizationID, testAccPageID, permission)
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

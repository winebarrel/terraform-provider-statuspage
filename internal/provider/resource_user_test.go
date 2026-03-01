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

func TestAccUser_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	userID := "user-id-1"
	user := &apiclient.User{
		ID:             userID,
		OrganizationID: testAccOrganizationID,
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/organizations/"+testAccOrganizationID+"/users",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.UserRequest
			_ = json.NewDecoder(req.Body).Decode(&body)
			user.Email = body.User.Email
			user.FirstName = body.User.FirstName
			user.LastName = body.User.LastName
			return httpmock.NewJsonResponse(201, user)
		})

	// GET uses list endpoint - returns array
	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/organizations/"+testAccOrganizationID+"/users",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []apiclient.User{*user})
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/organizations/"+testAccOrganizationID+"/users/"+userID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserConfig("tf-test-user@example.com", "TestFirst", "TestLast"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_user.test", "email", "tf-test-user@example.com"),
					resource.TestCheckResourceAttr("statuspage_user.test", "first_name", "TestFirst"),
					resource.TestCheckResourceAttr("statuspage_user.test", "last_name", "TestLast"),
					resource.TestCheckResourceAttrSet("statuspage_user.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncUser("statuspage_user.test"),
				// password is sensitive and not returned by the API
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUserConfig(email, firstName, lastName string) string {
	return fmt.Sprintf(`
resource "statuspage_user" "test" {
  organization_id = %q
  email           = %q
  password        = "TestP@ssw0rd123!"
  first_name      = %q
  last_name       = %q
}
`, testAccOrganizationID, email, firstName, lastName)
}

func importStateIDFuncUser(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["organization_id"]
		id := rs.Primary.ID
		return fmt.Sprintf("%s/%s", orgID, id), nil
	}
}

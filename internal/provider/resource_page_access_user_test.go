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

func TestAccPageAccessUser_basic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	userID := "page-access-user-id-1"
	user := &apiclient.PageAccessUser{
		ID:     userID,
		PageID: testAccPageID,
	}

	httpmock.RegisterResponder("POST", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/page_access_users",
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PageAccessUserRequest
			json.NewDecoder(req.Body).Decode(&body)
			user.ExternalLogin = body.PageAccessUser.ExternalLogin
			user.ExternalEmail = body.PageAccessUser.ExternalEmail
			return httpmock.NewJsonResponse(201, user)
		})

	httpmock.RegisterResponder("GET", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/page_access_users/"+userID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, user)
		})

	httpmock.RegisterResponder("PATCH", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/page_access_users/"+userID,
		func(req *http.Request) (*http.Response, error) {
			var body apiclient.PageAccessUserRequest
			json.NewDecoder(req.Body).Decode(&body)
			if body.PageAccessUser.ExternalLogin != "" {
				user.ExternalLogin = body.PageAccessUser.ExternalLogin
			}
			if body.PageAccessUser.ExternalEmail != "" {
				user.ExternalEmail = body.PageAccessUser.ExternalEmail
			}
			return httpmock.NewJsonResponse(200, user)
		})

	httpmock.RegisterResponder("DELETE", "https://api.statuspage.io/v1/pages/"+testAccPageID+"/page_access_users/"+userID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccPageAccessUserConfig("tf-test-pau", "tf-test-pau@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_login", "tf-test-pau"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_email", "tf-test-pau@example.com"),
					resource.TestCheckResourceAttrSet("statuspage_page_access_user.test", "id"),
				),
			},
			// Import
			{
				ResourceName:      "statuspage_page_access_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFuncPageAccessUser("statuspage_page_access_user.test"),
			},
			// Update
			{
				Config: testAccPageAccessUserConfig("tf-test-pau-updated", "tf-test-pau-updated@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_login", "tf-test-pau-updated"),
					resource.TestCheckResourceAttr("statuspage_page_access_user.test", "external_email", "tf-test-pau-updated@example.com"),
				),
			},
		},
	})
}

func testAccPageAccessUserConfig(login, email string) string {
	return fmt.Sprintf(`
resource "statuspage_page_access_user" "test" {
  page_id        = %q
  external_login = %q
  external_email = %q
}
`, testAccPageID, login, email)
}

func importStateIDFuncPageAccessUser(resourceName string) resource.ImportStateIdFunc {
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

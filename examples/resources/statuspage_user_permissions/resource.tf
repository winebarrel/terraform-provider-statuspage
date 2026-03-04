resource "statuspage_user_permissions" "example" {
  organization_id = "org123"
  user_id         = statuspage_user.example.id

  pages = {
    "page_id_1" = "manage"
    "page_id_2" = "view"
  }
}

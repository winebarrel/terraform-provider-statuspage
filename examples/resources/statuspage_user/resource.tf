import {
  to = statuspage_user.example
  id = "org123/user-id-1"
}

resource "statuspage_user" "example" {
  organization_id = "org123"
  email           = "newuser@example.com"
  password        = "securepassword123"
  first_name      = "John"
  last_name       = "Doe"
}

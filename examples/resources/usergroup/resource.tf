resource "slack_usergroup" "example_1" {
  name   = "example"
  handle = "example"
}

resource "slack_usergroup" "example_2" {
  name              = "example"
  handle            = "example"
  users             = ["U1234567890"]
  prevent_conflicts = True
}

---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "slack_user Data Source - slack"
subcategory: ""
description: |-
  Retrieve Slack user information. Either id or email must be specified, but not both.
---

# slack_user (Data Source)

Retrieve Slack user information. Either `id` or `email` must be specified, but not both.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `email` (String) Email of the user to look up.
- `id` (String) Slack user ID to look up.

### Read-Only

- `name` (String) User's name.
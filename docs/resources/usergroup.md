---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "slack_usergroup Resource - slack"
subcategory: ""
description: |-
  Manages a Slack user group.
  This resource requires the following scopes:
  usergroups:writeusergroups:read
---

# slack_usergroup (Resource)

Manages a Slack user group.

This resource requires the following scopes:

- usergroups:write
- usergroups:read



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `handle` (String)
- `name` (String)

### Optional

- `channels` (List of String) Channels shared by the user group.
- `description` (String)
- `prevent_conflicts` (Boolean) If true, the plan fails if there's an enabled user group with the same name or handle.
- `users` (Set of String) List of user IDs in the user group.

### Read-Only

- `id` (String) The ID of this resource.
